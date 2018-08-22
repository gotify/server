import {EventEmitter} from 'events';
import * as MessageAction from '../actions/MessageAction';
import AppStore from './AppStore';
import dispatcher, {IEvent} from './dispatcher';

class MessageStore extends EventEmitter {
    private appToMessages: {[appId: number]: IAppMessages} = {};
    private reset: false | number = false;
    private resetOnAll: false | number = false;
    private loading = false;

    constructor() {
        super();
        AppStore.on('change', () => {
            this.updateApps();
            this.emit('change');
        });
    }

    public shouldReset(appId: number): false | number {
        const reset = appId === -1 ? this.resetOnAll : this.reset;
        if (reset !== false) {
            this.reset = false;
            this.resetOnAll = false;
        }
        return reset;
    }

    public loadNext(id: number): void {
        if (this.loading || !this.get(id).hasMore) {
            return;
        }
        this.loading = true;
        MessageAction.fetchMessagesApp(id, this.get(id).nextSince).catch(
            () => (this.loading = false)
        );
    }

    public get(id: number): IAppMessages {
        if (this.exists(id)) {
            return this.appToMessages[id];
        } else {
            return {messages: [], nextSince: 0, hasMore: true};
        }
    }

    public exists(id: number): boolean {
        return this.appToMessages[id] !== undefined;
    }

    public handle(data: IEvent): void {
        const {payload} = data;
        if (data.type === 'UPDATE_MESSAGES') {
            if (this.exists(payload.id)) {
                payload.messages = this.get(payload.id).messages.concat(payload.messages);
            }
            this.appToMessages[payload.id] = payload;
            this.updateApps();
            this.loading = false;
            this.emit('change');
        } else if (data.type === 'ONE_MESSAGE') {
            if (this.exists(payload.appid)) {
                this.appToMessages[payload.appid].messages.unshift(payload);
                this.reset = 0;
            }
            if (this.exists(-1)) {
                this.appToMessages[-1].messages.unshift(payload);
                this.resetOnAll = 0;
            }
            this.updateApps();
            this.emit('change');
        } else if (data.type === 'DELETE_MESSAGE') {
            this.resetOnAll = this.removeFromList(this.appToMessages[-1], payload);
            this.reset = this.removeFromList(this.appToMessages[payload.appid], payload);
            this.emit('change');
        } else if (data.type === 'DELETE_MESSAGES') {
            const id = payload;
            if (id === -1) {
                this.appToMessages = {};
            } else {
                delete this.appToMessages[-1];
                delete this.appToMessages[id];
            }
            this.reset = 0;
            this.emit('change');
        }
    }

    private removeFromList(messages: IAppMessages, messageToDelete: IMessage): false | number {
        if (messages) {
            const index = messages.messages.findIndex(
                (message) => message.id === messageToDelete.id
            );
            if (index !== -1) {
                messages.messages.splice(index, 1);
                return index;
            }
        }
        return false;
    }

    private updateApps = (): void => {
        const appToUrl: {[appId: number]: string} = {};
        AppStore.get().forEach((app) => (appToUrl[app.id] = app.image));
        Object.keys(this.appToMessages).forEach((key) => {
            const appMessages: IAppMessages = this.appToMessages[key];
            appMessages.messages.forEach((message) => (message.image = appToUrl[message.appid]));
        });
    };
}

const store = new MessageStore();
dispatcher.register(store.handle.bind(store));
export default store;
