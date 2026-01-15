import {action, observable, makeObservable} from 'mobx';

interface HasID {
    id: number;
}

export interface IClearable {
    clear(): void;
}

/**
 * Base implementation for handling items with ids.
 */
export abstract class BaseStore<T extends HasID> implements IClearable {
    protected items: T[] = [];

    protected abstract requestItems(): Promise<T[]>;

    protected abstract requestDelete(id: number): Promise<void>;

    public remove = async (id: number): Promise<void> => {
        await this.requestDelete(id);
        await this.refresh();
    };

    public refresh = (): Promise<void> =>
        this.requestItems().then(
            action((items) => {
                this.items = items || [];
            })
        );

    public refreshIfMissing = async (id: number): Promise<void> => {
        if (this.getByIDOrUndefined(id) === undefined) {
            await this.refresh();
        }
    };

    public getByID = (id: number): T => {
        const item = this.getByIDOrUndefined(id);
        if (item === undefined) {
            throw new Error('cannot find item with id ' + id);
        }
        return item;
    };

    public getByIDOrUndefined = (id: number): T | undefined =>
        this.items.find((hasId: HasID) => hasId.id === id);

    public getItems = (): T[] => this.items;

    public clear = (): void => {
        this.items = [];
    };

    constructor() {
        // eslint-disable-next-line
        makeObservable<BaseStore<any>, 'items'>(this, {
            items: observable,
            remove: action,
            refresh: action,
            refreshIfMissing: action,
            clear: action,
        });
    }
}
