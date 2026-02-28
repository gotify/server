import {BaseStore} from '../common/BaseStore';
import * as config from '../config';
import {action} from 'mobx';
import {SnackReporter} from '../snack/SnackManager';
import {IUser} from '../types';
import {CurrentUser} from '../CurrentUser';
import {identityTransform, jsonBody, jsonTransform} from '../fetchUtils';

export class UserStore extends BaseStore<IUser> {
    constructor(
        private readonly currentUser: CurrentUser,
        private readonly snack: SnackReporter
    ) {
        super();
    }

    protected requestItems = (): Promise<IUser[]> =>
        this.currentUser.authenticatedFetch(config.get('url') + 'user', {}, jsonTransform<IUser[]>);

    protected requestDelete(id: number): Promise<void> {
        return this.currentUser
            .authenticatedFetch(
                config.get('url') + 'user/' + id,
                {
                    method: 'DELETE',
                },
                identityTransform
            )
            .then(() => this.snack('User deleted'));
    }

    @action
    public create = async (name: string, pass: string, admin: boolean) => {
        await this.currentUser
            .authenticatedFetch(
                config.get('url') + 'user',
                jsonBody({name, pass, admin}),
                identityTransform
            )
            .then(() => this.snack('User created'));
        await this.refresh();
        this.snack('User created');
    };

    @action
    public update = async (id: number, name: string, pass: string | null, admin: boolean) => {
        await this.currentUser
            .authenticatedFetch(
                config.get('url') + 'user/' + id,
                jsonBody({name, pass, admin}),
                identityTransform
            )
            .then(() => this.snack('User updated'));
        await this.refresh();
        this.snack('User updated');
    };
}
