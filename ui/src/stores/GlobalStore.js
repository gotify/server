import {EventEmitter} from 'events';
import dispatcher from './dispatcher';

class GlobalStore extends EventEmitter {
    constructor() {
        super();
        this.currentUser = null;
        this.isAuthenticating = true;
    }

    authenticating() {
        return this.isAuthenticating;
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
        this.isAuthenticating = false;
        this.currentUser = user;
        this.emit('change');
    }

    handle(data) {
        if (data.type === 'NO_AUTHENTICATION') {
            this.set(null);
        } else if (data.type === 'AUTHENTICATED') {
            this.set(data.payload);
        } else if (data.type === 'AUTHENTICATING') {
            this.isAuthenticating = true;
            this.emit('change');
        }
    }
}

const store = new GlobalStore();
dispatcher.register(store.handle.bind(store));
export default store;
