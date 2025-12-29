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
                                <div class="chat-messages" ref="chatMessages" style="flex:1; overflow:auto; padding:8px 4px;">
                                    <div
                                        v-for="(msg, idx) in chatMessages"
                                        :key="idx"
                                        :class="['chat-message', { sent: msg.isSent }]"
                                        style="margin-bottom:8px;"
                                    >
                                        <span class="chat-timestamp" style="font-size:12px; color:#6b7280;">
                                            {{ msg.timestamp }}
                                            <span v-if="msg.sender"> · {{ msg.sender }}</span>
                                            <span v-if="msg.channel" class="chat-channel">[{{ msg.channel }}]</span>
                                        </span>
                                        <div class="chat-text" style="background:#fff; padding:6px 8px; border-radius:4px; margin-top:4px;">{{ msg.text }}</div>
                                    </div>
                                </div>
                                <div class="chat-input-container" style="margin-top:8px; display:flex; align-items:center; gap:8px;">
                                    <select v-model="chatChannel" class="chat-channel-select" style="width:110px; padding:8px; border-radius:6px; border:1px solid #d1d5db; background:#fff;">
                                        <option v-for="ch in chatChannels" :key="ch" :value="ch">{{ ch }}</option>
                                    </select>

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
                                        style="flex:1; padding:8px; border-radius:6px; border:1px solid #d1d5db;"
                                    />
                                    <button @click="sendChatMessage" class="btn btn-primary btn-sm" style="padding:6px 10px;">发送</button>
                                </div>
                            </div>

                            <div class="log-section" style="position: relative; flex:1; display:flex; flex-direction:column; min-width:260px;">
                                <div class="log-header" style="display:flex; justify-content:space-between; align-items:center; padding-bottom:8px;">
                                    <h3 style="margin:0; font-size:14px;">日志</h3>
                                    <button @click="clearLogs" class="btn btn-secondary btn-xs" style="padding:6px 8px;">清空</button>
                                </div>

                                <!-- replaced textarea with styled div that supports per-line coloring -->
                                <div
                                    ref="logArea"
                                    class="log-textarea"
                                    style="flex:1; overflow:auto; padding:8px; background:#fff; border:1px solid #d1d5db; border-radius:4px; white-space:pre-wrap; font-family: monospace;"
                                >
                                    <div v-for="(line, idx) in logLines" :key="idx" :style="{ color: logColor(line) }">
                                        {{ line }}
                                    </div>
                                </div>

                                <!-- floating bottom-right button -->
                                <button
                                    @click="scrollLogsToBottom"
                                    class="btn btn-primary btn-xs"
                                    title="滚到底"
                                    style="position: absolute; right: 8px; bottom: 8px; z-index: 3;"
                                >到底⬇️</button>
                            </div>
                        </div>
                    </div>

                    <div class="log-modal-footer" style="padding:12px 16px; border-top:1px solid #e5e7eb; display:flex; justify-content:flex-end;">
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
                1: "附近",
                2: "组队",
                3: "公会",
                4: "朋友",
                5: "世界"
            },
            senderIds: {},
            ws: null,
            chatPollInterval: null // interval id for polling chat history
        };
    },
    computed: {
        // split logs into lines for per-line rendering
        logLines() {
            if (!this.logs) return [];
            // keep empty lines (so splitting preserves them)
            return this.logs.split('\n');
        }
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
        formatTimestamp(ts) {
            if (ts == null || ts === '') return '';
            if (ts instanceof Date) {
                return isNaN(ts.getTime()) ? '' : ts.toLocaleString();
            }
            const asNum = (typeof ts === 'string') ? Number(ts) : ts;
            if (typeof asNum === 'number' && !isNaN(asNum)) {
                // treat epoch seconds (small numbers) as seconds, else milliseconds
                let ms = asNum;
                if (Math.abs(asNum) < 1e12) ms = asNum * 1000;
                const d = new Date(ms);
                if (!isNaN(d.getTime())) return d.toLocaleString();
            }
            // fallback: try parsing string timestamps like "2023-01-01T12:34:56Z"
            const parsed = new Date(String(ts));
            if (!isNaN(parsed.getTime())) return parsed.toLocaleString();
            return String(ts);
        },

        scrollLogsToBottom() {
            const logArea = this.$refs.logArea;
            if (!logArea) return;
            // immediately jump to bottom
            logArea.scrollTop = logArea.scrollHeight;
            // ensure any Vue DOM updates are applied then keep bottom (safe fallback)
            this.$nextTick(() => {
                if (this.$refs.logArea) {
                    this.$refs.logArea.scrollTop = this.$refs.logArea.scrollHeight;
                }
            });
        },
        applyBodyScrollLock(disable) {
            document.body.style.overflow = disable ? 'hidden' : '';
        },

        // helper: map incoming channel id/string to friendly label
        mapChannel(channel) {
            if (channel == null) return '公会';
            if (typeof channel === 'number') {
                return this.channelMap[channel] || String(channel);
            }
            const asNum = Number(channel);
            if (!isNaN(asNum) && this.channelMap[asNum]) {
                return this.channelMap[asNum];
            }
            if (this.channelMap.hasOwnProperty(channel)) {
                return this.channelMap[channel];
            }
            if (this.chatChannels.includes(channel)) {
                return channel;
            }
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
                    let timestamp = this.formatTimestamp(tsVal);
                    this.senderIds[item.senderName] = item.senderId;
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

            const wsUrl = `${api.WS_BASE_URL}/logs?username=${this.username}&feature=${this.featureName}`;
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

        /*
         * Append a log line but only auto-scroll if the user was already at the bottom.
         * This prevents interrupting the user when they've scrolled up to read older logs.
         */
        addLog(message) {
            const logArea = this.$refs.logArea;
            // consider within 20px of bottom as "at bottom"
            const THRESHOLD_PX = 20;
            let wasAtBottom = true;
            if (logArea) {
                wasAtBottom = (logArea.scrollTop + logArea.clientHeight >= logArea.scrollHeight - THRESHOLD_PX);
            }

            const timestamp = new Date().toLocaleTimeString();
            this.logs += `[${timestamp}] ${message}\n`;

            this.$nextTick(() => {
                if (!logArea) return;
                if (wasAtBottom) {
                    // keep following new logs
                    logArea.scrollTop = logArea.scrollHeight;
                }
                // otherwise leave user's scroll position untouched
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
            this.$nextTick(() => {
                if (this.$refs.logArea) {
                    this.$refs.logArea.scrollTop = 0;
                }
            });
        },

        // determine color based on log line content
        logColor(line) {
            if (!line) return 'inherit';
            const l = String(line).toLowerCase();
            if (/\berror\b/.test(l)) return 'red';
            if (/\bwarn(?:ing)?\b/.test(l)) return '#b59f00';
            if (/\binfo\b/.test(l)) return 'green';
            // also allow explicit prefixes like "[ERROR]" or "ERROR:"
            if (/\[?error\]?[:\s]/i.test(line)) return 'red';
            if (/\[?warn(?:ing)?\]?[:\s]/i.test(line)) return '#b59f00';
            if (/\[?info\]?[:\s]/i.test(line)) return 'green';
            return 'inherit';
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
