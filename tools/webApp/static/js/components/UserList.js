import { api } from '../api.js';

export default {
    template: `
        <div class="container">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
                <h1 style="margin-bottom: 0;">账户列表</h1>
                <button @click="$emit('navigate', 'addUser')" class="btn">添加账户</button>
            </div>

            <div v-if="alert.show" :class="'alert alert-' + alert.type">
                {{ alert.message }}
            </div>

            <div v-if="loading" class="loading">
                Loading users...
            </div>

            <div v-else-if="users.length === 0" class="empty-state">
                <h3>No Users Found</h3>
                <p>Get started by adding your first user.</p>
                <button @click="$emit('navigate', 'addUser')" class="btn">Add User</button>
            </div>

            <div v-else>
                <table>
                    <thead>
                        <tr>
                            <th>用户名</th>
                            <th>角色编号</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in users" :key="user.username">
                            <td>{{ user.username }}</td>
                            <td>{{ user.roleNum }}</td>
                            <td>
                                <button @click="editUser(user)" class="btn btn-warning btn-sm">修改</button>
                                <button @click="deleteUser(user.username)" class="btn btn-danger btn-sm">删除</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    `,
    data() {
        return {
            users: [],
            loading: false,
            alert: {
                show: false,
                type: 'success',
                message: ''
            }
        };
    },
    async mounted() {
        await this.loadUsers();
    },
    methods: {
        async loadUsers() {
            this.loading = true;
            this.alert.show = false;

            try {
                const data = await api.getUsers();
                if (data.success) {
                    this.users = data.data || [];
                } else {
                    this.showAlert('error', data.message || 'Failed to load users');
                }
            } catch (error) {
                this.showAlert('error', 'Network error: ' + error.message);
            } finally {
                this.loading = false;
            }
        },
        editUser(user) {
            this.$emit('edit-user', user);
        },
        async deleteUser(username) {
            if (!confirm(`Are you sure you want to delete user "${username}"?`)) {
                return;
            }

            this.alert.show = false;

            try {
                const data = await api.deleteUser({"username": username});
                if (data.success) {
                    this.showAlert('success', 'User deleted successfully');
                    await this.loadUsers();
                } else {
                    this.showAlert('error', data.message || 'Failed to delete user');
                }
            } catch (error) {
                this.showAlert('error', 'Network error: ' + error.message);
            }
        },
        showAlert(type, message) {
            this.alert = { show: true, type, message };
            setTimeout(() => {
                this.alert.show = false;
            }, 10000);
        }
    }
};
