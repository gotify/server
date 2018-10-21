import axios, {AxiosResponse} from 'axios';
import * as config from '../config';
import {detect} from 'detect-browser';
import * as GlobalAction from '../actions/GlobalAction';
import SnackManager, {SnackReporter} from './SnackManager';
import {observable} from 'mobx';

const tokenKey = 'gotify-login-key';

class CurrentUser {
    private tokenCache: string | null = null;
    @observable
    public loggedIn = false;
    @observable
    public authenticating = false;
    @observable
    public user: IUser = {name: 'unknown', admin: false, id: -1};

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
        this.authenticating = true;
        const browser = detect();
        const name = (browser && browser.name + ' ' + browser.version) || 'unknown browser';
        axios
            .create()
            .request({
                url: config.get('url') + 'client',
                method: 'POST',
                data: {name},
                auth: {username, password},
            })
            .then((resp: AxiosResponse<IClient>) => {
                this.snack(`A client named '${name}' was created for your session.`);
                this.setToken(resp.data.token);
                this.tryAuthenticate()
                    .then((user) => {
                        this.authenticating = false;
                        this.loggedIn = true;
                        GlobalAction.initialLoad(user);
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
                return passThrough;
            })
            .catch((error) => {
                this.logout();
                return Promise.reject(error);
            });
    };

    public logout = () => {
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

export const currentUser = new CurrentUser(SnackManager.snack);
