import {EventEmitter} from 'events';
import dispatcher from './dispatcher';

class ClientStore extends EventEmitter {
    constructor() {
        super();
        this.clients = [];
    }

    get() {
        return this.clients;
    }

    getById(id) {
        return this.clients.find((client) => client.id === id);
    }

    getIdByToken(token) {
        const client = this.clients.find((client) => client.token === token);
        return client !== undefined ? client.id : '';
    }

    handle(data) {
        if (data.type === 'UPDATE_CLIENTS') {
            this.clients = data.payload;
            this.emit('change');
        }
    }
}


const store = new ClientStore();
dispatcher.register(store.handle.bind(store));
export default store;
