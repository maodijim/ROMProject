// File: `tools/webApp/static/js/api.js`
const PORT = 8081;
const API_BASE_URL = `http://localhost:${PORT}/api`;

export const api = {
    PORT,
    API_BASE_URL,
    async getUsers() {
        const response = await fetch(`${API_BASE_URL}/user/get`);
        return await response.json();
    },

    async createUser(userData) {
        const response = await fetch(`${API_BASE_URL}/user`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });
        return await response.json();
    },

    async updateUser(userData) {
        const response = await fetch(`${API_BASE_URL}/user`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });
        return await response.json();
    },

    async deleteUser(userData) {
        const response = await fetch(`${API_BASE_URL}/user`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });
        return await response.json();
    },

    async getUserRunningTask(username) {
        const response = await fetch(`${API_BASE_URL}/feature/running/user?username=${username}`);
        return await response.json();
    },

    async fetchLogs(username, featureName) {
        const url = `${API_BASE_URL}/feature/log?username=${encodeURIComponent(username)}&featureName=${encodeURIComponent(featureName)}`;
        const response = await fetch(url);
        return await response.json();
    },

    // New: fetch chat history for a user's feature
    async fetchChatHistory(username, featureName) {
        const url = `${API_BASE_URL}/feature/chat?username=${encodeURIComponent(username)}&featureName=${encodeURIComponent(featureName)}`;
        const response = await fetch(url);
        return await response.json();
    },

    async sendChatMessage(username, featureName, message, channelId, destId = null) {
        const url = `${API_BASE_URL}/feature/chat/send`;
        const payload = {
            username,
            featureName,
            message
        };
        if (channelId != null) payload.channelId = channelId;
        if (destId != null) payload.destId = destId;

        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        return await response.json();
    },

    async fetchFeatures() {
        const response = await fetch(`${API_BASE_URL}/feature`);
        return await response.json();
    },

    async fetchFeatureConfig(username, functionName) {
        const response = await fetch(`${API_BASE_URL}/feature/config?username=${username}&functionName=${encodeURIComponent(functionName)}`);
        return await response.json();
    },

    async updateFeatureConfig(username, functionName, configData) {
        const response = await fetch(`${API_BASE_URL}/feature/config/update?username=${username}&functionName=${encodeURIComponent(functionName)}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(configData)
        });
        return await response.json();
    },

    async fetchRunningTasks() {
        const response = await fetch(`${API_BASE_URL}/feature/running`);
        return await response.json();
    },

    async startFeatureTask(taskData) {
        const response = await fetch(`${API_BASE_URL}/feature/start`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(taskData)
        });
        return await response.json();
    },

    async stopFeatureTask(taskData) {
        const response = await fetch(`${API_BASE_URL}/feature/stop`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(taskData)
        });
        return await response.json();
    },
};
