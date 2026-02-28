import {BaseStore} from '../common/BaseStore';
import * as config from '../config';
import {action} from 'mobx';
import {SnackReporter} from '../snack/SnackManager';
import {IClient} from '../types';
import {CurrentUser} from '../CurrentUser';
import {identityTransform, jsonBody, jsonTransform} from '../fetchUtils';

export class ClientStore extends BaseStore<IClient> {
    public constructor(
        private readonly currentUser: CurrentUser,
        private readonly snack: SnackReporter
    ) {
        super();
    }

    protected requestItems = (): Promise<IClient[]> =>
        this.currentUser.authenticatedFetch(
            config.get('url') + 'client',
            {},
            jsonTransform<IClient[]>
        );

    protected requestDelete(id: number): Promise<void> {
        return this.currentUser
            .authenticatedFetch(
                config.get('url') + 'client/' + id,
                {
                    method: 'DELETE',
                },
                identityTransform
            )
            .then(() => this.snack('Client deleted'));
    }

    @action
    public update = async (id: number, name: string): Promise<void> => {
        await this.currentUser
            .authenticatedFetch(
                config.get('url') + 'client/' + id,
                {...jsonBody({name}), method: 'PUT'},
                jsonTransform
            )
            .then(() => this.snack('Client updated'));
        await this.refresh();
        this.snack('Client updated');
    };

    @action
    public createNoNotifcation = async (name: string): Promise<IClient> => {
        const client = await this.currentUser.authenticatedFetch(
            config.get('url') + 'client',
            jsonBody({name}),
            jsonTransform<IClient>
        );
        await this.refresh();
        return client;
    };

    @action
    public create = async (name: string): Promise<void> => {
        await this.createNoNotifcation(name);
        this.snack('Client added');
    };
}
