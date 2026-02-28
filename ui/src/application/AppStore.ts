import {generateKeyBetween} from 'fractional-indexing';
import {action, runInAction} from 'mobx';
import {BaseStore} from '../common/BaseStore';
import * as config from '../config';
import {SnackReporter} from '../snack/SnackManager';
import {IApplication} from '../types';
import {arrayMove} from '@dnd-kit/sortable';
import {CurrentUser} from '../CurrentUser';
import {identityTransform, jsonBody, jsonTransform, multipartBody} from '../fetchUtils';

export class AppStore extends BaseStore<IApplication> {
    public onDelete: () => void = () => {};

    public constructor(
        private readonly currentUser: CurrentUser,
        private readonly snack: SnackReporter
    ) {
        super();
    }

    protected requestItems = (): Promise<IApplication[]> =>
        this.currentUser.authenticatedFetch(
            config.get('url') + 'application',
            {},
            jsonTransform<IApplication[]>
        );

    protected requestDelete = (id: number): Promise<void> =>
        this.currentUser
            .authenticatedFetch(
                config.get('url') + 'application/' + id,
                {
                    method: 'DELETE',
                },
                identityTransform
            )
            .then(() => {
                this.onDelete();
                return this.snack('Application deleted');
            });

    @action
    public uploadImage = async (id: number, file: Blob): Promise<void> => {
        const formData = new FormData();
        formData.append('file', file);
        await this.currentUser.authenticatedFetch(
            config.get('url') + 'application/' + id + '/image',
            multipartBody(formData),
            jsonTransform
        );
        await this.refresh();
        this.snack('Application image updated');
    };

    public async deleteImage(id: number): Promise<void> {
        try {
            await this.currentUser.authenticatedFetch(
                config.get('url') + 'application/' + id + '/image',
                {
                    method: 'DELETE',
                },
                identityTransform
            );
            await this.refresh();
            this.snack('Application image deleted');
        } catch (error) {
            console.error('Error deleting application image:', error);
            throw error;
        }
    }

    @action
    public reorder = async (fromId: number, toId: number): Promise<void> => {
        const fromIndex = this.items.findIndex((app) => app.id === fromId);
        const toIndex = this.items.findIndex((app) => app.id === toId);
        if (fromIndex === -1 || toIndex === -1) {
            throw Error('unknown apps');
        }

        const toUpdate = this.items[fromIndex];

        const normalizedIndex =
            toUpdate.sortKey > this.items[toIndex].sortKey ? toIndex - 1 : toIndex;

        const newSortKey = generateKeyBetween(
            this.items[normalizedIndex]?.sortKey,
            this.items[normalizedIndex + 1]?.sortKey
        );

        runInAction(() => (this.items = arrayMove(this.items, fromIndex, toIndex)));

        await this.update({...toUpdate, sortKey: newSortKey});
    };

    @action
    public update = async ({
        id,
        ...app
    }: Pick<
        IApplication,
        'id' | 'name' | 'description' | 'defaultPriority' | 'sortKey'
    >): Promise<void> => {
        await this.currentUser.authenticatedFetch(
            config.get('url') + 'application/' + id,
            {...jsonBody(app), method: 'PUT'},
            jsonTransform
        );
        await this.refresh();
        this.snack('Application updated');
    };

    @action
    public create = async (
        name: string,
        description: string,
        defaultPriority: number
    ): Promise<void> => {
        await this.currentUser.authenticatedFetch(
            config.get('url') + 'application',
            jsonBody({name, description, defaultPriority}),
            jsonTransform
        );
        await this.refresh();
        this.snack('Application created');
    };

    public getName = (id: number): string => {
        const app = this.getByIDOrUndefined(id);
        return id === -1 ? 'All Messages' : app !== undefined ? app.name : 'unknown';
    };
}
