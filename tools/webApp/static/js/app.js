import Home from './components/Home.js';
import UserList from './components/UserList.js';
import AddUser from './components/AddUser.js';
import EditUser from './components/EditUser.js';
import Settings from './components/Settings.js';
import FeatureList from './components/FeatureList.js';

const { createApp } = Vue;

createApp({
    components: {
        Home,
        UserList,
        AddUser,
        EditUser,
        Settings,
        FeatureList
    },
    data() {
        return {
            currentPage: 'home',
            editUserData: null
        };
    },
    methods: {
        navigate(page) {
            this.currentPage = page;
            this.editUserData = null;
        },
        handleEditUser(user) {
            this.editUserData = user;
            this.currentPage = 'editUser';
        }
    }
}).mount('#app');
