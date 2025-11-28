import {api} from '../api.js';

export default {
    template: `
        <teleport to="body">
            <div v-if="show" class="log-modal-overlay" @click="handleOverlayClick">
                <div class="log-modal-content" @click.stop>
                    <div class="log-modal-header">
                        <h2>{{ featureName }} - {{ username }}</h2>
                        <button @click="close" class="log-btn-close">&times;</button>
                    </div>

                    <div class="log-modal-body">
                        <div class="log-chat-container">
                            <div class="chat-section">
                                <div class="chat-header">
                                    <h3>聊天窗口</h3>
                                </div>
                                <div class="chat-messages" ref="chatMessages">
                                    <div
                                        v-for="(msg, idx) in chatMessages"
                                        :key="idx"
                                        :class="['chat-message', { sent: msg.isSent }]"
                                    >
                                        <span class="chat-timestamp">
                                            {{ msg.timestamp }}
                                            <span v-if="msg.sender"> · {{ msg.sender }}</span>
                                            <span v-if="msg.channel" class="chat-channel">[{{ msg.channel }}]</span>
                                        </span>
                                        <div class="chat-text">{{ msg.text }}</div>
                                    </div>
                                </div>
                                <div class="chat-input-container">
                                    <select v-model="chatChannel" class="chat-channel-select" style="width:110px; padding:8px; border-radius:6px; border:1px solid #d1d5db; background:#fff;">
                                        <option v-for="ch in chatChannels" :key="ch" :value="ch">{{ ch }}</option>
                                    </select>

                                    <!-- new: recipient dropdown shown when chatChannel is 朋友 -->
                                    <select
                                        v-if="chatChannel === '朋友'"
                                        v-model="selectedRecipient"
                                        class="chat-recipient-select"
                                        style="width:140px; padding:8px; border-radius:6px; border:1px solid #d1d5db; background:#fff; margin-left:8px;"
                                    >
                                        <option value="" disabled selected>选择朋友</option>
                                        <option v-for="(id, name) in senderIds" :key="name" :value="name">
                                            {{ name }}
                                        </option>
                                    </select>

                                    <input
                                        v-model="chatInput"
                                        @keyup.enter="sendChatMessage"
                                        type="text"
                                        class="chat-input"
                                        placeholder="输入消息..."
                                    />
                                    <button @click="sendChatMessage" class="btn btn-primary btn-sm">发送</button>
                                </div>
                            </div>

                            <div class="log-section">
                                <div class="log-header">
                                    <h3>日志</h3>
                                    <button @click="clearLogs" class="btn btn-secondary btn-xs">清空</button>
                                </div>
                                <textarea
                                    ref="logArea"
                                    v-model="logs"
                                    class="log-textarea"
                                    readonly
                                ></textarea>
                            </div>
                        </div>
                    </div>

                    <div class="log-modal-footer">
                        <button @click="close" class="btn btn-secondary">关闭</button>
                    </div>
                </div>
            </div>
        </teleport>
    `,
    props: {
        show: Boolean,
        username: String,
        featureName: String
    },
    data() {
        return {
            logs: '',
            chatMessages: [],
            chatInput: '',
            selectedRecipient: null,
            chatChannel: '组队',
            chatChannels: ['公会', '组队', '世界', '附近', '朋友'],
            channelMap: {
                2: "组队",
                3: "公会",
                4: "朋友",
                5: "世界",
                6: "附近"
            },
            senderIds: {},
            ws: null,
            chatPollInterval: null // interval id for polling chat history
        };
    },
    watch: {
        show(newVal) {
            this.applyBodyScrollLock(newVal);
            if (newVal) {
                this.fetchLogs();
                this.fetchChatHistory();
                // start polling every 10 seconds
                if (!this.chatPollInterval) {
                    this.chatPollInterval = window.setInterval(() => {
                        this.fetchChatHistory();
                    }, 10000);
                }
                this.connectWebSocket();
            } else {
                // stop polling
                if (this.chatPollInterval) {
                    clearInterval(this.chatPollInterval);
                    this.chatPollInterval = null;
                }
                this.disconnectWebSocket();
            }
        }
    },
    methods: {
        applyBodyScrollLock(disable) {
            document.body.style.overflow = disable ? 'hidden' : '';
        },

        // helper: map incoming channel id/string to friendly label
        mapChannel(channel) {
            if (channel == null) return '公会';
            // if numeric
            if (typeof channel === 'number') {
                return this.channelMap[channel] || String(channel);
            }
            // numeric-string like "3"
            const asNum = Number(channel);
            if (!isNaN(asNum) && this.channelMap[asNum]) {
                return this.channelMap[asNum];
            }
            // direct key in channelMap (object keys are strings)
            if (this.channelMap.hasOwnProperty(channel)) {
                return this.channelMap[channel];
            }
            // known human-readable channel already in list
            if (this.chatChannels.includes(channel)) {
                return channel;
            }
            // fallback to channel string
            return String(channel);
        },

        async fetchLogs() {
            if (!this.username || !this.featureName) return;
            try {
                const data = await api.fetchLogs(this.username, this.featureName);
                if (data.success && data.data && Array.isArray(data.data.logs)) {
                    this.logs = data.data.logs.join('\n');
                    if (data.data.logs.length) {
                        this.logs += '\n';
                    }
                } else {
                    this.logs = '';
                }
            } catch (error) {
                this.logs = 'Failed to fetch logs: ' + error.message;
            }
        },

        // load chat history and normalize entries to { timestamp, text, channel }
        async fetchChatHistory() {
            if (!this.username || !this.featureName) return;
            try {
                const res = await api.fetchChatHistory(this.username, this.featureName);
                if (!res || !res.success) return;

                const raw = res.data;
                if (!Array.isArray(raw)) return;

                const normalized = raw.map((item) => {
                    if (typeof item === 'string') {
                        return {
                            timestamp: new Date().toLocaleTimeString(),
                            text: item,
                            channel: '公会',
                            sender: '系统'
                        };
                    }
                    const text = item.message || item.text || item.msg || item.content || JSON.stringify(item);
                    const rawChannel = item.msgChannel || null;
                    const channelLabel = this.mapChannel(rawChannel);
                    let tsVal = item.timestamp || null;
                    let timestamp = new Date().toLocaleTimeString();
                    this.senderIds[item.senderName] = item.senderId;
                    if (tsVal) {
                        try {
                            const parsed = new Date(tsVal);
                            if (!isNaN(parsed.getTime())) {
                                timestamp = parsed.toLocaleTimeString();
                            } else {
                                timestamp = String(tsVal);
                            }
                        } catch (e) {
                            timestamp = String(tsVal);
                        }
                    }
                    return {
                        timestamp,
                        text,
                        channel: channelLabel,
                        sender: item.senderName || '未知'
                    };
                });

                // replace chat messages with fresh history
                this.chatMessages = normalized;
                this.$nextTick(() => {
                    if (this.$refs.chatMessages) {
                        this.$refs.chatMessages.scrollTop = this.$refs.chatMessages.scrollHeight;
                    }
                });
            } catch (error) {
                this.addLog('Failed to fetch chat history: ' + (error && error.message ? error.message : error));
            }
        },

        connectWebSocket() {
            if (this.ws) {
                this.ws.close();
            }

            const wsUrl = `ws://localhost:8081/ws/logs?username=${this.username}&feature=${this.featureName}`;
            this.ws = new WebSocket(wsUrl);

            this.ws.onopen = () => {
                this.addLog('WebSocket connected');
            };

            this.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    if (data.type === 'log') {
                        this.addLog(data.message);
                    } else if (data.type === 'chat') {
                        const channelLabel = this.mapChannel(data.channel || data.ch || null);
                        this.addChatMessage(data.message, channelLabel);
                    }
                } catch (e) {
                    this.addLog(event.data);
                }
            };

            this.ws.onerror = (error) => {
                this.addLog('WebSocket error: ' + error);
            };

            this.ws.onclose = () => {
                this.addLog('WebSocket disconnected');
            };
        },
        disconnectWebSocket() {
            if (this.ws) {
                this.ws.close();
                this.ws = null;
            }
        },
        addLog(message) {
            const timestamp = new Date().toLocaleTimeString();
            this.logs += `[${timestamp}] ${message}\n`;
            this.$nextTick(() => {
                if (this.$refs.logArea) {
                    this.$refs.logArea.scrollTop = this.$refs.logArea.scrollHeight;
                }
            });
        },
        addChatMessage(text, channel) {
            const timestamp = new Date().toLocaleTimeString();
            // ensure channel label is human-readable
            const channelLabel = this.mapChannel(channel);
            this.chatMessages.push({ timestamp, text, channel: channelLabel });
            this.$nextTick(() => {
                if (this.$refs.chatMessages) {
                    this.$refs.chatMessages.scrollTop = this.$refs.chatMessages.scrollHeight;
                }
            });
        },
        getChannelId(channel) {
            if (channel == null) return null;
            if (typeof channel === 'number') return channel;
            const asNum = Number(channel);
            if (!isNaN(asNum) && this.channelMap[asNum]) return asNum;
            for (const k in this.channelMap) {
                if (this.channelMap.hasOwnProperty(k) && this.channelMap[k] === channel) {
                    return Number(k);
                }
            }
            return null;
        },

        async sendChatMessage() {
            if (!this.chatInput.trim()) return;

            const channelId = this.getChannelId(this.chatChannel);
            if (!channelId) {
                this.addLog('Cannot send message: unknown channel id for "' + this.chatChannel + '"');
                return;
            }

            const destId = this.chatChannel === '朋友' && this.selectedRecipient
                ? this.senderIds[this.selectedRecipient] || null
                : null;

            try {
                const res = await api.sendChatMessage(this.username, this.featureName, this.chatInput, channelId, destId);
                if (res && res.success) {
                    // clear input on success
                    this.chatInput = '';
                } else {
                    this.addLog('Send chat failed: ' + (res && res.message ? res.message : 'unknown'));
                }
            } catch (err) {
                this.addLog('Send chat error: ' + (err && err.message ? err.message : err));
            }
        },
        clearLogs() {
            this.logs = '';
        },
        handleOverlayClick() {
            this.close();
        },
        close() {
            this.$emit('close');
        }
    },
    beforeUnmount() {
        this.applyBodyScrollLock(false);
        if (this.chatPollInterval) {
            clearInterval(this.chatPollInterval);
            this.chatPollInterval = null;
        }
        this.disconnectWebSocket();
    }
};
