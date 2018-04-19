import {EventEmitter} from 'events';
import dispatcher, {IEvent} from './dispatcher';

class UserStore extends EventEmitter {
    private users: IUser[] = [];

    public get(): IUser[] {
        return this.users;
    }

    public getById(id: number): IUser {
        const user = this.users.find((u) => u.id === id);
        if (!user) {
            throw new Error('user must exist');
        }
        return user;
    }

    public handle(data: IEvent): void {
        if (data.type === 'UPDATE_USERS') {
            this.users = data.payload;
            this.emit('change');
        }
    }
}

const store = new UserStore();
dispatcher.register(store.handle.bind(store));
export default store;
