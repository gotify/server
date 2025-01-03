import axios from 'axios';
import {getAuthToken} from './common/Auth.ts';
import store from './store';
import {tryAuthenticate} from './store/auth-actions.ts';
import {uiActions} from './store/ui-slice.ts';

export const initAxios = () => {
    axios.interceptors.request.use((config) => {
        config.headers['X-Gotify-Key'] = getAuthToken();
        return config;
    });

    axios.interceptors.response.use(undefined, (error) => {
        if (!error.response) {
            store.dispatch(uiActions.addSnackMessage('Gotify server is not reachable, try refreshing the page.'));
            return Promise.reject(error);
        }

        const status = error.response.status;

        if (status === 401) {
            store.dispatch(tryAuthenticate()).then(() => store.dispatch(uiActions.addSnackMessage('Could not complete request')));
        }

        if (status === 400 || status === 403 || status === 500) {
            store.dispatch(uiActions.addSnackMessage(error.response.data.error + ': ' + error.response.data.errorDescription));
        }

        return Promise.reject(error);
    });
};
