import {SnackReporter} from './SnackManager';
import {currentUser} from './CurrentUser';
import * as config from '../config';
import NewMessagesStore from './MessagesStore';

export class WebSocketStore {
    private wsActive = false;
    private ws: WebSocket | null = null;

    public constructor(private readonly snack: SnackReporter) {}

    public listen = () => {
        if (!currentUser.token() || this.wsActive) {
            return;
        }
        this.wsActive = true;

        const wsUrl = config
            .get('url')
            .replace('http', 'ws')
            .replace('https', 'wss');
        const ws = new WebSocket(wsUrl + 'stream?token=' + currentUser.token());

        ws.onerror = (e) => {
            this.wsActive = false;
            console.log('WebSocket connection errored', e);
        };

        ws.onmessage = (data) => NewMessagesStore.publishSingleMessage(JSON.parse(data.data));

        ws.onclose = () => {
            this.wsActive = false;
            currentUser.tryAuthenticate().then(() => {
                this.snack('WebSocket connection closed, trying again in 30 seconds.');
                setTimeout(this.listen, 30000);
            });
        };

        this.ws = ws;
    };

    public close = () => this.ws && this.ws.close();
}
