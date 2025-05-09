import axios, {AxiosError, AxiosResponse} from 'axios';
import {detect} from 'detect-browser';
import {Base64} from 'js-base64';
import {getAuthToken} from '../common/Auth.ts';
import * as config from '../config.ts';
import {IClient, IUser} from '../types.ts';
import {appActions} from '../application/app-slice.ts';
import {authActions} from './auth-slice.ts';
import {clientActions} from '../client/client-slice.ts';
import {AppDispatch} from '../store/index.ts';
import {messageActions} from '../message/message-slice.ts';
import {pluginActions} from '../plugin/plugin-slice.ts';
import {uiActions} from '../store/ui-slice.ts';
import {userActions} from './user-slice.ts';

export const tokenKey = 'gotify-login-key';
let reconnectTimeoutId: number | null = null;
let reconnectTime = 7500;

export const register = (username: string, password: string) => {
    return async (dispatch: AppDispatch) => {
        try {
            await axios.create().post(config.get('url') + 'user', {name: username, pass: password});
            dispatch(uiActions.addSnackMessage('User Created. Logging in...'));
            await dispatch(login(username, password));
            return true;
        } catch (error: unknown) {
            if (error instanceof AxiosError) {
                if (!error || !error.response) {
                    dispatch(uiActions.addSnackMessage('No network connection or server unavailable.'));
                    return false;
                }
                const {data} = error.response;
                dispatch(
                    uiActions.addSnackMessage(
                        `Register failed: ${data?.error ?? 'unknown'}: ${data?.errorDescription ?? ''}`
                    )
                );
            }
            return Promise.reject(error);
        }
    };
};

export const login = (username: string, password: string) => {
    return async (dispatch: AppDispatch) => {
        dispatch(authActions.logout());
        dispatch(authActions.isAuthenticating(true));
        const browser = detect();
        const name = (browser && browser.name + ' ' + browser.version) || 'unknown browser';

        try {
            const response = await axios.create().request({
                url: config.get('url') + 'client',
                method: 'POST',
                data: {name},
                // eslint-disable-next-line @typescript-eslint/naming-convention
                headers: {Authorization: 'Basic ' + Base64.encode(username + ':' + password)},
            });
            dispatch(
                uiActions.addSnackMessage(`A client named '${name}' was created for your session.`)
            );
            localStorage.setItem(tokenKey, response.data.token);
            return await dispatch(tryAuthenticate());
        } catch (error) {
            dispatch(uiActions.addSnackMessage('Login failed'));
            return Promise.reject(new Error('Login failed.'));
        } finally {
            dispatch(authActions.isAuthenticating(false));
        }
    };
};

export const tryAuthenticate = () => {
    return async (dispatch: AppDispatch) => {
        if (!getAuthToken()) {
            return Promise.reject(new Error('No token provided'));
        }

        try {
            const response = await axios
                .create()
                // eslint-disable-next-line @typescript-eslint/naming-convention
                .get<IUser>(config.get('url') + 'current/user', {
                    headers: {'X-Gotify-Key': getAuthToken()},
                });
            dispatch(authActions.login(response.data));
            dispatch(uiActions.setConnectionErrorMessage(null));
            reconnectTime = 7500;
        } catch (error: unknown) {
            if (error instanceof AxiosError) {
                if (!error || !error.response) {
                    dispatch(connectionError('No network connection or server unavailable.'));
                    return Promise.reject(error);
                }

                if (error.response.status >= 500) {
                    dispatch(
                        connectionError(
                            `${error.response.statusText} (code: ${error.response.status}).`
                        )
                    );
                    return Promise.reject(new Error(error.message));
                }

                dispatch(uiActions.setConnectionErrorMessage(null));

                if (error.response.status >= 400 && error.response.status < 500) {
                    dispatch(logout());
                }
            }
            return Promise.reject(error);
        }
    };
};

export const logout = () => {
    return async (dispatch: AppDispatch) => {
        await axios
            .get(config.get('url') + 'client')
            .then((resp: AxiosResponse<IClient[]>) => {
                resp.data
                    .filter((client) => client.token === getAuthToken())
                    .forEach((client) => axios.delete(config.get('url') + 'client/' + client.id));
            })
            .catch(() => Promise.resolve());

        localStorage.removeItem(tokenKey);
        dispatch(authActions.logout());
        dispatch(userActions.clear());
        dispatch(messageActions.clear());
        dispatch(appActions.clear());
        dispatch(clientActions.clear());
        dispatch(pluginActions.clear());
    };
};

export const changePassword = (pass: string) => {
    return async (dispatch: AppDispatch) => {
        await axios.post(config.get('url') + 'current/user/password', {pass});
        dispatch(uiActions.addSnackMessage('Password changed.'));
    };
};

export const tryReconnect = (quiet = false) => {
    return async (dispatch: AppDispatch) => {
        try {
            dispatch(tryAuthenticate());
        } catch (error) {
            if (!quiet) {
                dispatch(uiActions.addSnackMessage('Reconnect failed'));
            }
        }
    };
};

const connectionError = (message: string) => {
    return async (dispatch: AppDispatch) => {
        dispatch(uiActions.setConnectionErrorMessage(message));
        if (reconnectTimeoutId !== null) {
            window.clearTimeout(reconnectTimeoutId);
        }
        reconnectTimeoutId = window.setTimeout(() => dispatch(tryReconnect(true)), reconnectTime);
        reconnectTime = Math.min(reconnectTime * 2, 120000);
    };
};
