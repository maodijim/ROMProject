import { api } from '../api.js';

export default {
    template: `
        <div class="container">
            <h1>Add New User</h1>

            <div v-if="alert.show" :class="'alert alert-' + alert.type">
                {{ alert.message }}
            </div>

            <form @submit.prevent="submitUser">
                <div class="form-group">
                    <label for="username">Username:</label>
                    <input
                        type="text"
                        id="username"
                        v-model="form.username"
                        required
                        :disabled="loading"
                    >
                </div>
                <div class="form-group">
                    <label for="password">Password:</label>
                    <input
                        type="password"
                        id="password"
                        v-model="form.password"
                        required
                        :disabled="loading"
                    >
                </div>
                <div class="form-group">
                    <label for="rolenum">Role Number (1-3):</label>
                    <input
                        type="number"
                        id="rolenum"
                        v-model.number="form.roleNum"
                        min="1"
                        max="3"
                        required
                        :disabled="loading"
                    >
                </div>
                <button type="submit" class="btn" :disabled="loading">
                    {{ loading ? 'Adding...' : 'Add User' }}
                </button>
                <button type="button" class="btn btn-secondary" @click="$emit('navigate', 'listUsers')" :disabled="loading">
                    Cancel
                </button>
            </form>
        </div>
    `,
    data() {
        return {
            loading: false,
            form: {
                username: '',
                password: '',
                roleNum: 1
            },
            alert: {
                show: false,
                type: 'success',
                message: ''
            }
        };
    },
    methods: {
        async submitUser() {
            this.loading = true;
            this.alert.show = false;

            try {
                const data = await api.createUser(this.form);

                if (data.success) {
                    this.showAlert('success', 'User added successfully!');
                    this.resetForm();

                    setTimeout(() => {
                        this.$emit('navigate', 'listUsers');
                    }, 1500);
                } else {
                    this.showAlert('error', data.message || 'Failed to add user');
                }
            } catch (error) {
                this.showAlert('error', 'Network error: ' + error.message);
            } finally {
                this.loading = false;
            }
        },
        showAlert(type, message) {
            this.alert = { show: true, type, message };
            setTimeout(() => {
                this.alert.show = false;
            }, 5000);
        },
        resetForm() {
            this.form = {
                username: '',
                password: '',
                roleNum: 1
            };
        }
    }
};
