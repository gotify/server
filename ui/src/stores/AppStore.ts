import {EventEmitter} from 'events';
import dispatcher, {IEvent} from './dispatcher';

class AppStore extends EventEmitter {
    private apps: IApplication[] = [];

    public get(): IApplication[] {
        return this.apps;
    }

    public getById(id: number): IApplication {
        const app = this.getByIdOrUndefined(id);
        if (!app) {
            throw new Error('app is required to exist')
        }
        return app;
    }

    public getName(id: number): string {
        const app = this.getByIdOrUndefined(id);
        return id === -1 ? 'All Messages' : app !== undefined ? app.name : 'unknown';
    }

    public handle(data: IEvent): void {
        if (data.type === 'UPDATE_APPS') {
            this.apps = data.payload;
            this.emit('change');
        }
    }

    private getByIdOrUndefined(id: number): IApplication | undefined {
        return this.apps.find((a) => a.id === id);
    }
}

const store = new AppStore();
dispatcher.register(store.handle.bind(store));
export default store;
