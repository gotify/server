import {EventEmitter} from 'events';
import dispatcher from './dispatcher';

class MessageStore extends EventEmitter {
    constructor() {
        super();
        this.messages = [];
        this.messagesForApp = [];
    }

    get() {
        return this.messages;
    }

    getForAppId(id) {
        if (id === -1) {
            return this.messages;
        }
        if (this.messagesForApp[id]) {
            return this.messagesForApp[id];
        }
        return [];
    }

    handle(data) {
        if (data.type === 'UPDATE_MESSAGES') {
            this.messages = data.payload;
            this.messagesForApp = [];
            this.messages.forEach(function(message) {
                this.createIfNotExist(message.appid);
                this.messagesForApp[message.appid].push(message);
            }.bind(this));
            this.emit('change');
        } else if (data.type === 'ONE_MESSAGE') {
            const {payload} = data;
            this.createIfNotExist(payload.appid);
            this.messagesForApp[payload.appid].unshift(payload);
            this.messages.unshift(payload);
            this.emit('change');
        }
    }

    createIfNotExist(id) {
        if (!(id in this.messagesForApp)) {
            this.messagesForApp[id] = [];
        }
    }
}

const store = new MessageStore();
dispatcher.register(store.handle.bind(store));
export default store;
