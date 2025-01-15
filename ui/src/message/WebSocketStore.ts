import {Middleware} from '@reduxjs/toolkit';
import {AxiosError} from 'axios';
import {getAuthToken} from '../common/Auth.ts';
import * as config from '../config';
import * as Notifications from '../snack/browserNotification.ts';
import {tryAuthenticate} from '../user/auth-actions.ts';
import {uiActions} from '../store/ui-slice.ts';
import {IMessage} from '../types';
import state, {RootState} from '../store/index.ts';
import {messageActions} from './message-slice.ts';

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

const ws = new WebSocketStore();

export const handleWebsocketMiddleware: Middleware<{}, RootState> =
    (store) => (next) => (action) => {
        const prevLoginStatus = store.getState().auth.loggedIn;
        const prevReloadRequired = store.getState().ui.reloadRequired;
        const result = next(action);
        const nextLoginStatus = store.getState().auth.loggedIn;
        const nextReloadRequired = store.getState().ui.reloadRequired;

        // Open websocket connection on login and if a reload is required
        if ((action.type.startsWith('auth/login') && prevLoginStatus !== nextLoginStatus)
            || (action.type.startsWith('ui/setReloadRequired') && prevReloadRequired !== nextReloadRequired && nextReloadRequired)
        ) {
            ws.listen((message) => {
                store.dispatch(messageActions.loading(true));
                store.dispatch(messageActions.add(message));
                Notifications.notifyNewMessage(message);
                if (message.priority >= 4) {
                    const src = 'static/notification.ogg';
                    const audio = new Audio(src);
                    audio.play();
                }
            });
            window.onbeforeunload = () => {
                ws.close();
            };
        }
        // Close websocket if the user logs out
        if (action.type.startsWith('auth/logout') && prevLoginStatus !== nextLoginStatus && !nextLoginStatus) {
            ws.close();
        }
        return result;
    };
