import { api } from '../api.js';

export default {
    template: `
        <div class="container">
            <h1>Edit User</h1>

            <div v-if="alert.show" :class="'alert alert-' + alert.type">
                {{ alert.message }}
            </div>

            <div v-if="loading" class="loading">
                Loading user data...
            </div>

            <form v-else @submit.prevent="updateUser">
                <div class="form-group">
                    <label for="edit-username">Username:</label>
                    <input
                        type="text"
                        id="edit-username"
                        v-model="form.username"
                        required
                        :disabled="submitting"
                    >
                </div>
                <div class="form-group">
                    <label for="edit-password">Password:</label>
                    <input
                        type="password"
                        id="edit-password"
                        v-model="form.password"
                        :disabled="submitting"
                        placeholder="Enter new password leave blank to keep current"
                    >
                </div>
                <div class="form-group">
                    <label for="edit-rolenum">Role Number (1-3):</label>
                    <input
                        type="number"
                        id="edit-rolenum"
                        v-model.number="form.roleNum"
                        min="1"
                        max="3"
                        required
                        :disabled="submitting"
                    >
                </div>
                <button type="submit" class="btn" :disabled="submitting">
                    {{ submitting ? 'Updating...' : 'Update User' }}
                </button>
                <button type="button" class="btn btn-secondary" @click="$emit('navigate', 'listUsers')" :disabled="submitting">
                    Cancel
                </button>
            </form>
        </div>
    `,
    props: {
        userData: {
            type: Object,
            default: null
        }
    },
    data() {
        return {
            loading: false,
            submitting: false,
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
    mounted() {
        if (this.userData) {
            this.form.username = this.userData.username;
            this.form.roleNum = this.userData.roleNum;
        }
    },
    methods: {
        async updateUser() {
            this.submitting = true;
            this.alert.show = false;

            try {
                const data = await api.updateUser(this.form);

                if (data.success) {
                    this.showAlert('success', 'User updated successfully!');

                    setTimeout(() => {
                        this.$emit('navigate', 'listUsers');
                    }, 1500);
                } else {
                    this.showAlert('error', data.message || 'Failed to update user');
                }
            } catch (error) {
                this.showAlert('error', 'Network error: ' + error.message);
            } finally {
                this.submitting = false;
            }
        },
        showAlert(type, message) {
            this.alert = { show: true, type, message };
            setTimeout(() => {
                this.alert.show = false;
            }, 5000);
        }
    }
};
