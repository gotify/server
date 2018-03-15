import {EventEmitter} from 'events';
import dispatcher from './dispatcher';

class UserStore extends EventEmitter {
    constructor() {
        super();
        this.users = [];
    }

    get() {
        return this.users;
    }

    getById(id) {
        return this.users.find((app) => app.id === id);
    }

    handle(data) {
        if (data.type === 'UPDATE_USERS') {
            this.users = data.payload;
            this.emit('change');
        }
    }
}

const store = new UserStore();
dispatcher.register(store.handle.bind(store));
export default store;
