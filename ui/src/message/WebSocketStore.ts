import {SnackReporter} from '../snack/SnackManager';
import {CurrentUser} from '../CurrentUser';
import * as config from '../config';
import {AxiosError} from 'axios';
import {IMessage} from '../types';

export class WebSocketStore {
    private wsActive = false;
    private ws: WebSocket | null = null;

    public constructor(
        private readonly snack: SnackReporter,
        private readonly currentUser: CurrentUser
    ) {}

    public listen = (callback: (msg: IMessage) => void) => {
        if (!this.currentUser.loggedIn || this.wsActive) {
            return;
        }
        this.wsActive = true;

        const wsUrl = config.get('url').replace('http', 'ws').replace('https', 'wss');
        const ws = new WebSocket(wsUrl + 'stream');

        ws.onerror = (e) => {
            this.wsActive = false;
            console.log('WebSocket connection errored', e);
        };

        ws.onmessage = (data) => callback(JSON.parse(data.data));

        ws.onclose = () => {
            this.wsActive = false;
            if (!this.currentUser.loggedIn) {
                return;
            }
            this.currentUser
                .tryAuthenticate()
                .then(() => {
                    this.snack('WebSocket connection closed, trying again in 30 seconds.');
                    setTimeout(() => this.listen(callback), 30000);
                })
                .catch((error: AxiosError) => {
                    if (error?.response?.status === 401) {
                        this.snack('Could not authenticate with client token, logging out.');
                    }
                });
        };

        this.ws = ws;
    };

    public close = () => this.ws?.close(1000, 'WebSocketStore#close');
}
