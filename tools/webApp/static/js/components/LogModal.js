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
                                    <div v-for="(msg, index) in chatMessages" :key="index" class="chat-message">
                                        <span class="chat-timestamp">{{ msg.timestamp }}</span>
                                        <span class="chat-text">{{ msg.text }}</span>
                                    </div>
                                </div>
                                <div class="chat-input-container">
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
            ws: null
        };
    },
    watch: {
        show(newVal) {
            if (newVal) {
                this.connectWebSocket();
            } else {
                this.disconnectWebSocket();
            }
        }
    },
    methods: {
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
                        this.addChatMessage(data.message);
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
        addChatMessage(text) {
            const timestamp = new Date().toLocaleTimeString();
            this.chatMessages.push({ timestamp, text });
            this.$nextTick(() => {
                if (this.$refs.chatMessages) {
                    this.$refs.chatMessages.scrollTop = this.$refs.chatMessages.scrollHeight;
                }
            });
        },
        sendChatMessage() {
            if (!this.chatInput.trim()) return;

            if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({
                    type: 'chat',
                    message: this.chatInput
                }));
                this.chatInput = '';
            } else {
                this.addLog('Cannot send message: WebSocket not connected');
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
        this.disconnectWebSocket();
    }
};
