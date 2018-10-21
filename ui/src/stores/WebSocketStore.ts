import {SnackReporter} from './SnackManager';
import {CurrentUser} from './CurrentUser';
import * as config from '../config';
import {MessagesStore} from './MessagesStore';

export class WebSocketStore {
    private wsActive = false;
    private ws: WebSocket | null = null;

    public constructor(
        private readonly snack: SnackReporter,
        private readonly currentUser: CurrentUser,
        private readonly messagesStore: MessagesStore
    ) {}

    public listen = () => {
        if (!this.currentUser.token() || this.wsActive) {
            return;
        }
        this.wsActive = true;

        const wsUrl = config
            .get('url')
            .replace('http', 'ws')
            .replace('https', 'wss');
        const ws = new WebSocket(wsUrl + 'stream?token=' + this.currentUser.token());

        ws.onerror = (e) => {
            this.wsActive = false;
            console.log('WebSocket connection errored', e);
        };

        ws.onmessage = (data) => this.messagesStore.publishSingleMessage(JSON.parse(data.data));

        ws.onclose = () => {
            this.wsActive = false;
            this.currentUser.tryAuthenticate().then(() => {
                this.snack('WebSocket connection closed, trying again in 30 seconds.');
                setTimeout(this.listen, 30000);
            });
        };

        this.ws = ws;
    };

    public close = () => this.ws && this.ws.close();
}
