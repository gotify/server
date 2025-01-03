import {AxiosError} from 'axios';
import {getAuthToken} from '../common/Auth.ts';
import * as config from '../config';
import {tryAuthenticate} from '../store/auth-actions.ts';
import {uiActions} from '../store/ui-slice.ts';
import {IMessage} from '../types';
import state from '../store/index.ts';

export class WebSocketStore {
    private wsActive = false;
    private ws: WebSocket | null = null;

    public constructor() {}

    public listen = (callback: (msg: IMessage) => void) => {
        if (!getAuthToken() || this.wsActive) {
            return;
        }
        this.wsActive = true;

        const wsUrl = config.get('url').replace('http', 'ws');
        const ws = new WebSocket(wsUrl + 'stream?token=' + getAuthToken());

        ws.onerror = (e) => {
            this.wsActive = false;
            console.log('WebSocket connection errored', e);
        };

        ws.onmessage = (data) => callback(JSON.parse(data.data));

        ws.onclose = () => {
            this.wsActive = false;
            // try to reopen websocket or completely logout the user by a side effect
            state
                .dispatch(tryAuthenticate())
                .then(() => {
                    state.dispatch(
                        uiActions.addSnackMessage(
                            'WebSocket connection closed, trying again in 30 seconds.'
                        )
                    );
                    setTimeout(() => this.listen(callback), 30000);
                })
                .catch((error: AxiosError) => {
                    if (error?.response?.status === 401) {
                        state.dispatch(
                            uiActions.addSnackMessage(
                                'WebSocket connection closed, trying again in 30 seconds.'
                            )
                        );
                    }
                });
        };

        this.ws = ws;
    };

    public close = () => this.ws?.close(1000, 'WebSocketStore#close');
}
