import * as config from './config';
import {detect} from 'detect-browser';
import {SnackReporter} from './snack/SnackManager';
import {observable, runInAction, action} from 'mobx';
import {IClient, IUser} from './types';
import {identityTransform, jsonBody, jsonTransform, ResponseTransformer} from './fetchUtils';

const tokenKey = 'gotify-login-key';

export class CurrentUser {
    private tokenCache: string | null = null;
    private reconnectTimeoutId: number | null = null;
    private reconnectTime = 7500;
    @observable accessor loggedIn = false;
    @observable accessor refreshKey = 0;
    @observable accessor authenticating = true;
    @observable accessor user: IUser = {name: 'unknown', admin: false, id: -1};
    @observable accessor connectionErrorMessage: string | null = null;

    public constructor(private readonly snack: SnackReporter) {}

    public token = (): string => {
        if (this.tokenCache !== null) {
            return this.tokenCache;
        }

        const localStorageToken = window.localStorage.getItem(tokenKey);
        if (localStorageToken) {
            this.tokenCache = localStorageToken;
            return localStorageToken;
        }

        return '';
    };

    public authenticatedFetch = async <T>(
        url: string,
        init: RequestInit,
        xform: ResponseTransformer<T>
    ): Promise<T> => {
        const headers = new Headers(init?.headers);
        if (this.loggedIn && !headers.has('X-Gotify-Key'))
            headers.set('X-Gotify-Key', this.token());
        let response;
        try {
            response = await fetch(url, {...init, headers});
        } catch (error) {
            this.snack('Gotify server is not reachable, try refreshing the page.');
            throw error;
        }
        if (response.ok) {
            try {
                return xform(response);
            } catch (error) {
                this.snack('Response transformation failed: ' + error);
                throw error;
            }
        }
        if (response.status === 401) {
            this.tryAuthenticate().then(() => this.snack('Could not complete request.'));
        }

        let error = 'Unexpected status code: ' + response.status;
        if (response.status === 400 || response.status === 403 || response.status === 500) {
            if (response.headers.get('content-type')?.includes('application/json')) {
                const data = await response.json();
                error = data.error + ': ' + data.errorDescription;
            } else {
                const text = await response.text();
                error = 'Unexpected response: ' + text;
            }
        }
        this.snack(error);
        throw new Error(error);
    };

    private readonly setToken = (token: string) => {
        this.tokenCache = token;
        window.localStorage.setItem(tokenKey, token);
    };

    public register = async (name: string, pass: string): Promise<boolean> => {
        runInAction(() => {
            this.loggedIn = false;
        });
        return this.authenticatedFetch(
            config.get('url') + 'user',
            jsonBody({name, pass}),
            identityTransform
        )
            .then(() => {
                this.snack('User Created. Logging in...');
                this.login(name, pass);
                return true;
            })
            .catch((error) => {
                if (error instanceof TypeError) {
                    this.snack('No network connection or server unavailable.');
                    return false;
                }
                this.snack(`Register failed: ${error?.message ?? error}`);
                return false;
            });
    };

    public login = async (username: string, password: string) => {
        runInAction(() => {
            this.loggedIn = false;
            this.authenticating = true;
        });
        const browser = detect();
        const name = (browser && browser.name + ' ' + browser.version) || 'unknown browser';
        const fetchInit = jsonBody({name});
        fetchInit.headers = new Headers(fetchInit.headers);
        fetchInit.headers.set('Authorization', 'Basic ' + btoa(username + ':' + password));
        return this.authenticatedFetch(
            config.get('url') + 'client',
            fetchInit,
            jsonTransform<IClient>
        )
            .then((resp) => {
                this.snack(`A client named '${name}' was created for your session.`);
                this.setToken(resp.token);
                this.tryAuthenticate().catch(() => {
                    console.log(
                        'create client succeeded, but authenticated with given token failed'
                    );
                });
            })
            .catch(
                action(() => {
                    this.authenticating = false;
                    return this.snack('Login failed');
                })
            );
    };

    public tryAuthenticate = async (): Promise<IUser> => {
        if (this.token() === '') {
            runInAction(() => {
                this.authenticating = false;
            });
            return Promise.reject();
        }

        return fetch(config.get('url') + 'current/user', {headers: {'X-Gotify-Key': this.token()}})
            .then(async (response) => {
                if (response.ok) {
                    const user = await response.json();
                    runInAction(() => {
                        this.user = user;
                        this.loggedIn = true;
                        this.authenticating = false;
                        this.connectionErrorMessage = null;
                        this.reconnectTime = 7500;
                    });
                    return user;
                }
                if (response.status >= 500) {
                    this.connectionError(`${response.statusText} (code: ${response.status}).`);
                    return Promise.reject(new Error('Server error'));
                }

                this.connectionErrorMessage = null;

                if (response.status >= 400 && response.status < 500) {
                    this.logout();
                }
                throw new Error('Unexpected status code: ' + response.status);
            })
            .catch(
                action((error) => {
                    this.authenticating = false;
                    this.connectionError('No network connection or server unavailable.');
                    return Promise.reject(error);
                })
            );
    };

    public logout = async () => {
        await this.authenticatedFetch(config.get('url') + 'client', {}, jsonTransform<IClient[]>)
            .then((resp) => {
                resp.filter((client) => client.token === this.tokenCache).forEach((client) =>
                    this.authenticatedFetch(
                        config.get('url') + 'client/' + client.id,
                        {},
                        jsonTransform
                    )
                );
            })
            .catch(() => Promise.resolve());
        window.localStorage.removeItem(tokenKey);
        this.tokenCache = null;
        runInAction(() => {
            this.loggedIn = false;
        });
    };

    public changePassword = (pass: string) => {
        this.authenticatedFetch(
            config.get('url') + 'current/user/password',
            jsonBody({pass}),
            identityTransform
        )
            .then(() => this.snack('Password changed'))
            .catch((error) => {
                this.snack(`Change password failed: ${error?.message ?? error}`);
            });
    };

    public tryReconnect = (quiet = false) => {
        this.tryAuthenticate().catch(() => {
            if (!quiet) {
                this.snack('Reconnect failed');
            }
        });
    };

    private readonly connectionError = (message: string) => {
        this.connectionErrorMessage = message;
        if (this.reconnectTimeoutId !== null) {
            window.clearTimeout(this.reconnectTimeoutId);
        }
        this.reconnectTimeoutId = window.setTimeout(
            () => this.tryReconnect(true),
            this.reconnectTime
        );
        this.reconnectTime = Math.min(this.reconnectTime * 2, 120000);
    };
}
