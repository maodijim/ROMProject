import { api } from '../api.js';
import ConfigModal from './ConfigModal.js';
import LogModal from './LogModal.js';

export default {
    components: {
        ConfigModal,
        LogModal
    },
    template: `
    <div>
        <div class="container">
            <h1>仙境传说守护永恒的爱自动化管理</h1>

            <div v-if="alert.show" :class="'alert alert-' + alert.type">
                {{ alert.message }}
            </div>

            <div v-if="loading" class="loading">
                Loading...
            </div>

            <div v-else-if="users.length === 0" class="empty-state">
                <h3>没有用户</h3>
                <p>开始添加您的第一个用户。</p>
                <button @click="$emit('navigate', 'addUser')" class="btn">添加用户</button>
            </div>

            <div v-else>
                <div v-for="user in users" :key="user.username" class="card" style="margin-bottom: 20px;">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px;">
                        <h2 style="margin: 0;">{{ user.username }} (角色 {{ user.roleNum }})</h2>
                        <div>
                            <button @click="editUser(user)" class="btn btn-info btn-sm">编辑账户</button>
                        </div>
                    </div>

                    <table>
                        <thead>
                            <tr>
                                <th>功能名</th>
                                <th>当前任务</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="feature in features" :key="feature.name">
                                <td>{{ feature.name }}</td>
                                <td>
                                    <span v-if="getRunningTask(user.username, feature.name)" class="status-badge status-running">
                                        {{ getRunningTask(user.username, feature.name) }}
                                    </span>
                                    <span v-else class="status-badge status-idle">空闲</span>
                                </td>
                                <td>
                                    <template v-for="action in feature.actions" :key="action">
                                        <button
                                            v-if="shouldShowAction(user.username, feature.name, action)"
                                            @click="performAction(user.username, feature.name, action, feature.functionName)"
                                            :class="getActionButtonClass(action)"
                                        >
                                            {{ getActionText(action) }}
                                        </button>
                                    </template>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <config-modal
            :show="showConfigModal"
            :username="configModal.username"
            :feature-name="configModal.featureName"
            :function-name="configModal.functionName"
            @close="closeConfigModal"
            @success="handleConfigSuccess"
            @error="handleConfigError"
        />
        
        <log-modal
            :show="showLogModal"
            :username="logModal.username"
            :feature-name="logModal.featureName"
            @close="closeLogModal"
        />
    </div>
`,
    data() {
        return {
            users: [],
            features: [],
            runningTasks: {},
            loading: false,
            showConfigModal: false,
            showLogModal: false,
            logModal: {
                username: '',
                featureName: ''
            },
            configModal: {
                username: '',
                featureName: '',
                functionName: ''
            },
            alert: {
                show: false,
                type: 'success',
                message: ''
            },
            actionMapping: {
                'start': '启动',
                'stop': '停止',
                'configure': '配置',
                'log': '日志'
            }
        };
    },
    async mounted() {
        await this.loadData();
    },
    methods: {
        async loadData() {
            this.loading = true;
            this.alert.show = false;

            try {
                const [usersData, featuresResponse, runningTasksResponse] = await Promise.all([
                    api.getUsers(),
                    fetch('http://localhost:8081/api/feature'),
                    fetch('http://localhost:8081/api/feature/running')
                ]);

                const featuresData = await featuresResponse.json();
                const runningTasksData = await runningTasksResponse.json();

                if (usersData.success) {
                    this.users = usersData.data || [];
                } else {
                    this.showAlert('error', usersData.message || 'Failed to load users');
                }

                if (featuresData.success) {
                    this.features = featuresData.data || [];
                } else {
                    this.showAlert('error', featuresData.message || 'Failed to load features');
                }

                if (runningTasksData.success) {
                    const tasks = runningTasksData.data || [];
                    this.runningTasks = {};
                    tasks.forEach(task => {
                        const key = `${task.username}:${task.featureName}`;
                        this.runningTasks[key] = task.status;
                    });
                }
            } catch (error) {
                this.showAlert('error', 'Network error: ' + error.message);
            } finally {
                this.loading = false;
            }
        },
        getActionText(action) {
            return this.actionMapping[action.toLowerCase()] || action;
        },
        getActionButtonClass(action) {
            const actionLower = action.toLowerCase();
            const classMap = {
                'start': 'btn btn-success btn-sm',
                'stop': 'btn btn-danger btn-sm',
                'configure': 'btn btn-warning btn-sm',
                'log': 'btn btn-info btn-sm'
            };
            return classMap[actionLower] || 'btn btn-sm';
        },
        getRunningTask(username, featureName) {
            const key = `${username}:${featureName}`;
            return this.runningTasks[key] || null;
        },
        shouldShowAction(username, featureName, action) {
            const isRunning = !!this.getRunningTask(username, featureName);
            const actionLower = action.toLowerCase();

            if (actionLower === 'start') {
                return !isRunning;
            } else if (actionLower === 'stop') {
                return isRunning;
            }
            return true;
        },
        editUser(user) {
            this.$emit('edit-user', user);
        },
        async performAction(username, featureName, action, functionName) {
            if (action.toLowerCase() === 'configure') {
                this.openConfigModal(username, featureName, functionName);
                return;
            }
            if (action.toLowerCase() === 'log') {
                this.openLogModal(username, featureName);
                return;
            }
            if (action.toLowerCase() === 'start') {
                try {
                    const response = await fetch('http://localhost:8081/api/feature/start', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            username: username,
                            featureName: featureName
                        })
                    });
                    const data = await response.json();

                    if (data.success) {
                        this.showAlert('success', `成功启动 ${featureName}`);
                        await this.loadData();
                    } else {
                        this.showAlert('error', data.message || '启动失败');
                    }
                } catch (error) {
                    this.showAlert('error', 'Network error: ' + error.message);
                }
            } else if (action.toLowerCase() === 'stop') {
                try {
                    const response = await fetch('http://localhost:8081/api/feature/stop', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            username: username
                        })
                    });
                    const data = await response.json();

                    if (data.success) {
                        this.showAlert('success', `成功停止 ${featureName}`);
                        await this.loadData();
                    } else {
                        this.showAlert('error', data.message || '停止失败');
                    }
                } catch (error) {
                    this.showAlert('error', 'Network error: ' + error.message);
                }
            } else {
                this.showAlert('info', `执行 ${action} 操作于用户 ${username} 的 ${featureName}...`);
            }
        },
        showAlert(type, message) {
            this.alert = { show: true, type, message };
            setTimeout(() => {
                this.alert.show = false;
            }, 10000);
        },
        openConfigModal(username, featureName, functionName) {
            this.showConfigModal = true;
            this.configModal = {
                username,
                featureName,
                functionName
            };
        },
        closeConfigModal() {
            this.showConfigModal = false;
            this.configModal = {
                username: '',
                featureName: '',
                functionName: ''
            };
        },
        handleConfigSuccess(message) {
            this.showAlert('success', message);
        },
        handleConfigError(message) {
            this.showAlert('error', message);
        },
        openLogModal(username, featureName) {
            this.showLogModal = true;
            this.logModal = {
                username,
                featureName
            };
        },
        closeLogModal() {
            this.showLogModal = false;
            this.logModal = {
                username: '',
                featureName: ''
            };
        }
    }
};
