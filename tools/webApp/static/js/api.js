const API_BASE_URL = 'http://localhost:8081/api';

export const api = {
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
        const response = await fetch(`http://localhost:8081/api/feature/running/user?username=${username}`);
        return await response.json();
    }
};
