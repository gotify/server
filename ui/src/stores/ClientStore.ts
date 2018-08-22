import {EventEmitter} from 'events';
import {default as dispatcher, IEvent} from './dispatcher';

class ClientStore extends EventEmitter {
    private clients: IClient[] = [];

    public get(): IClient[] {
        return this.clients;
    }

    public getById(id: number): IClient {
        const client = this.clients.find((c) => c.id === id);
        if (!client) {
            throw new Error('client is required to exist');
        }
        return client;
    }

    public getIdByToken(token: string): number {
        const client = this.clients.find((c) => c.token === token);
        return client !== undefined ? client.id : -1;
    }

    public handle(data: IEvent): void {
        if (data.type === 'UPDATE_CLIENTS') {
            this.clients = data.payload;
            this.emit('change');
        }
    }
}

const store = new ClientStore();
dispatcher.register(store.handle.bind(store));
export default store;
