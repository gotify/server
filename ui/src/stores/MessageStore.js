import {EventEmitter} from 'events';
import dispatcher from './dispatcher';
import AppStore from './AppStore';
import * as MessageAction from '../actions/MessageAction';

class MessageStore extends EventEmitter {
    constructor() {
        super();
        this.appToMessages = {};
        this.loading = false;
        AppStore.on('change', this.updateApps);
    }

    loadNext(id) {
        if (this.loading || !this.get(id).hasMore) {
            return;
        }
        this.loading = true;
        MessageAction.fetchMessagesApp(id, this.get(id).nextSince).catch(() => this.loading = false);
    }

    get(id) {
        if (this.exists(id)) {
            return this.appToMessages[id];
        } else {
            return {messages: [], nextSince: 0, hasMore: true};
        }
    }

    exists(id) {
        return this.appToMessages[id] !== undefined;
    }

    handle(data) {
        if (data.type === 'UPDATE_MESSAGES') {
            const payload = data.payload;
            if (this.exists(payload.id)) {
                payload.messages = this.get(payload.id).messages.concat(payload.messages);
            }
            this.appToMessages[payload.id] = payload;
            this.updateApps();
            this.loading = false;
            this.emit('change');
        } else if (data.type === 'ONE_MESSAGE') {
            const {payload} = data;
            this.createIfNotExist(payload.appid);
            this.createIfNotExist(-1);
            this.appToMessages[payload.appid].messages.unshift(payload);
            this.appToMessages[-1].messages.unshift(payload);
            this.updateApps();
            this.emit('change');
        } else if (data.type === 'DELETE_MESSAGE') {
            Object.keys(this.appToMessages).forEach((key) => {
                const appMessages = this.appToMessages[key];
                const index = appMessages.messages.indexOf(data.payload);
                if (index !== -1) {
                    appMessages.messages.splice(index, 1);
                }
            });
            this.emit('change');
        } else if (data.type === 'DELETE_MESSAGES') {
            const id = data.payload;
            if (id === -1) {
                this.appToMessages = {};
            } else {
                delete this.appToMessages[-1];
                delete this.appToMessages[id];
            }
            this.emit('change');
        }
    }

    updateApps = () => {
        const appToUrl = {};
        AppStore.get().forEach((app) => appToUrl[app.id] = app.image);
        Object.keys(this.appToMessages).forEach((key) => {
            const appMessages = this.appToMessages[key];
            appMessages.messages.forEach((message) => message.image = appToUrl[message.appid]);
        });
    };

    createIfNotExist(id) {
        if (!(id in this.appToMessages)) {
            this.appToMessages[id] = this.get(id);
        }
    }
}

const store = new MessageStore();
dispatcher.register(store.handle.bind(store));
export default store;
