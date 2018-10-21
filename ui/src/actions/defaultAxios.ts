import axios from 'axios';
import {currentUser} from '../stores/CurrentUser';
import SnackManager from '../stores/SnackManager';

axios.interceptors.request.use((config) => {
    config.headers['X-Gotify-Key'] = currentUser.token();
    return config;
});

axios.interceptors.response.use(undefined, (error) => {
    if (!error.response) {
        SnackManager.snack('Gotify server is not reachable, try refreshing the page.');
        return Promise.reject(error);
    }

    const status = error.response.status;

    if (status === 401) {
        currentUser.tryAuthenticate().then(() => SnackManager.snack('Could not complete request.'));
    }

    if (status === 400) {
        SnackManager.snack(error.response.data.error + ': ' + error.response.data.errorDescription);
    }

    return Promise.reject(error);
});
