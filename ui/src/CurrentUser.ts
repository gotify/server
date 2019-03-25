import axios, {AxiosError, AxiosResponse} from 'axios';
import * as config from './config';
import {Base64} from 'js-base64';
import {detect} from 'detect-browser';
import {SnackReporter} from './snack/SnackManager';
import {observable} from 'mobx';

const tokenKey = 'gotify-login-key';

export class CurrentUser {
    private tokenCache: string | null = null;
    @observable
    public loggedIn = false;
    @observable
    public authenticating = false;
    @observable
    public user: IUser = {name: 'unknown', admin: false, id: -1};
    @observable
    public hasNetwork = true;

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

    private setToken = (token: string) => {
        this.tokenCache = token;
        window.localStorage.setItem(tokenKey, token);
    };

    public login = async (username: string, password: string) => {
        this.loggedIn = false;
        this.authenticating = true;
        const browser = detect();
        const name = (browser && browser.name + ' ' + browser.version) || 'unknown browser';
        axios
            .create()
            .request({
                url: config.get('url') + 'client',
                method: 'POST',
                data: {name},
                headers: {Authorization: 'Basic ' + Base64.encode(username + ':' + password)},
            })
            .then((resp: AxiosResponse<IClient>) => {
                this.snack(`A client named '${name}' was created for your session.`);
                this.setToken(resp.data.token);
                this.tryAuthenticate()
                    .then(() => {
                        this.authenticating = false;
                        this.loggedIn = true;
                    })
                    .catch(() => {
                        this.authenticating = false;
                        console.log(
                            'create client succeeded, but authenticated with given token failed'
                        );
                    });
            })
            .catch(() => {
                this.authenticating = false;
                return this.snack('Login failed');
            });
    };

    public tryAuthenticate = async (): Promise<AxiosResponse<IUser>> => {
        if (this.token() === '') {
            return Promise.reject();
        }

        return axios
            .create()
            .get(config.get('url') + 'current/user', {headers: {'X-Gotify-Key': this.token()}})
            .then((passThrough) => {
                this.user = passThrough.data;
                this.loggedIn = true;
                this.hasNetwork = true;
                return passThrough;
            })
            .catch((error: AxiosError) => {
                if (!error || !error.response) {
                    this.hasNetwork = false;
                    return Promise.reject(error);
                }

                this.hasNetwork = true;

                if (error.response.status >= 400 && error.response.status < 500) {
                    this.logout();
                }
                return Promise.reject(error);
            });
    };

    public logout = async () => {
        await axios
            .get(config.get('url') + 'client')
            .then((resp: AxiosResponse<IClient[]>) => {
                resp.data.filter((client) => client.token === this.tokenCache).forEach((client) => {
                    return axios.delete(config.get('url') + 'client/' + client.id);
                });
            })
            .catch(() => Promise.resolve());
        window.localStorage.removeItem(tokenKey);
        this.tokenCache = null;
        this.loggedIn = false;
    };

    public changePassword = (pass: string) => {
        axios
            .post(config.get('url') + 'current/user/password', {pass})
            .then(() => this.snack('Password changed'));
    };
}
