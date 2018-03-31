import {EventEmitter} from 'events';
import dispatcher from './dispatcher';

class SnackBarStore extends EventEmitter {
    messages = [];

    next() {
        return this.messages.shift();
    }

    hasNext() {
        return this.messages.length !== 0;
    }

    handle(data) {
        if (data.type === 'SNACK') {
            this.messages.push(data.payload);
            this.emit('change');
        }
    }
}

const store = new SnackBarStore();
dispatcher.register(store.handle.bind(store));
export default store;
