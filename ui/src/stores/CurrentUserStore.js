import {EventEmitter} from 'events';
import dispatcher from './dispatcher';

class CurrentUserStore extends EventEmitter {
    constructor() {
        super();
        this.currentUser = null;
    }

    get() {
        return this.currentUser || {name: 'unknown', admin: false};
    }

    isAdmin() {
        return this.get().admin;
    }

    getName() {
        return this.get().name;
    }

    isLoggedIn() {
        return this.currentUser != null;
    }

    set(user) {
        this.currentUser = user;
        this.emit('change');
    }

    handle(data) {
        if (data.type === 'REMOVE_CURRENT_USER') {
            this.set(null);
        } else if (data.type === 'SET_CURRENT_USER') {
            this.set(data.payload);
        }
    }
}

const store = new CurrentUserStore();
dispatcher.register(store.handle.bind(store));
export default store;
