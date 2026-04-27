import { api } from '../api.js';

export default {
    template: `
        <div class="container">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
                <h1 style="margin-bottom: 0;">功能列表</h1>
            </div>

            <div v-if="alert.show" :class="'alert alert-' + alert.type">
                {{ alert.message }}
            </div>

            <div v-if="loading" class="loading">
                Loading features...
            </div>

            <div v-else-if="features.length === 0" class="empty-state">
                <h3>No Features Found</h3>
                <p>Get started by adding your first feature.</p>
            </div>

            <div v-else>
                <table>
                    <thead>
                        <tr>
                            <th>功能名</th>
                            <th>描述</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="feature in features" :key="feature.name">
                            <td>{{ feature.name }}</td>
                            <td>{{ feature.desc }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    `,
    data() {
        return {
            features: [],
            loading: false,
            alert: {
                show: false,
                type: 'success',
                message: ''
            }
        };
    },
    async mounted() {
        await this.loadFeatures();
    },
    methods: {
        async loadFeatures() {
            this.loading = true;
            this.alert.show = false;
            try {
                const response = await api.fetchFeatures();
                const data = await response;
                if (data.success) {
                    this.features = data.data || [];
                } else {
                    this.showAlert('error', data.message || 'Failed to load features');
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
            }, 10000);
        }
    }
};
