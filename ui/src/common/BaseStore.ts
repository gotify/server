import {action, observable} from 'mobx';

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
    @observable
    protected items: T[] = [];

    protected abstract requestItems(): Promise<T[]>;

    protected abstract requestDelete(id: number): Promise<void>;

    @action
    public remove = async (id: number): Promise<void> => {
        await this.requestDelete(id);
        await this.refresh();
    };

    @action
    public refresh = async (): Promise<void> => {
        this.items = await this.requestItems().then((items) => items || []);
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

    @action
    public clear = (): void => {
        this.items = [];
    };
}
