import {EventEmitter} from 'events';
import dispatcher from './dispatcher';

class AppStore extends EventEmitter {
    constructor() {
        super();
        this.apps = [];
    }

    get() {
        return this.apps;
    }

    getById(id) {
        return this.apps.find((app) => app.id === id);
    }

    getName(id) {
        const app = this.getById(id);
        return id === -1 ? 'All Messages' : app !== undefined ? app.name : 'unknown';
    }

    handle(data) {
        if (data.type === 'UPDATE_APPS') {
            this.apps = data.payload;
            this.emit('change');
        }
    }
}

const store = new AppStore();
dispatcher.register(store.handle.bind(store));
export default store;
