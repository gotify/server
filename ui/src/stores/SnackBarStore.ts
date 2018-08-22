import {EventEmitter} from 'events';
import dispatcher, {IEvent} from './dispatcher';

class SnackBarStore extends EventEmitter {
    public messages: string[] = [];

    public next(): string {
        if (!this.hasNext()) {
            throw new Error('no such element');
        }
        return this.messages.shift() as string;
    }

    public hasNext(): boolean {
        return this.messages.length !== 0;
    }

    public handle(data: IEvent): void {
        if (data.type === 'SNACK') {
            this.messages.push(data.payload);
            this.emit('change');
        }
    }
}

const store = new SnackBarStore();
dispatcher.register(store.handle.bind(store));
export default store;
