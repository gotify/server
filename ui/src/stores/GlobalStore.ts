import {EventEmitter} from 'events';
import dispatcher, {IEvent} from './dispatcher';

class GlobalStore extends EventEmitter {
    private currentUser: IUser | null = null;
    private isAuthenticating = true;

    public authenticating(): boolean {
        return this.isAuthenticating;
    }

    public get(): IUser {
        return this.currentUser || {name: 'unknown', admin: false, id: -1};
    }

    public isAdmin(): boolean {
        return this.get().admin;
    }

    public getName(): string {
        return this.get().name;
    }

    public isLoggedIn(): boolean {
        return this.currentUser != null;
    }

    public handle(data: IEvent): void {
        if (data.type === 'NO_AUTHENTICATION') {
            this.set(null);
        } else if (data.type === 'AUTHENTICATED') {
            this.set(data.payload);
        } else if (data.type === 'AUTHENTICATING') {
            this.isAuthenticating = true;
            this.emit('change');
        }
    }

    private set(user: IUser | null): void {
        this.isAuthenticating = false;
        this.currentUser = user;
        this.emit('change');
    }
}

const store = new GlobalStore();
dispatcher.register(store.handle.bind(store));
export default store;
