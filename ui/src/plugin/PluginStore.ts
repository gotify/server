import {action} from 'mobx';
import {BaseStore} from '../common/BaseStore';
import * as config from '../config';
import {SnackReporter} from '../snack/SnackManager';
import {IPlugin} from '../types';
import {CurrentUser} from '../CurrentUser';
import {identityTransform, yamlBody, jsonTransform, textTransform} from '../fetchUtils';

export class PluginStore extends BaseStore<IPlugin> {
    public onDelete: () => void = () => {};

    public constructor(
        private readonly currentUser: CurrentUser,
        private readonly snack: SnackReporter
    ) {
        super();
    }

    public requestConfig = (id: number): Promise<string> =>
        this.currentUser.authenticatedFetch(
            config.get('url') + 'plugin/' + id + '/config',
            {},
            textTransform
        );

    public requestDisplay = (id: number): Promise<string> =>
        this.currentUser.authenticatedFetch(
            config.get('url') + 'plugin/' + id + '/display',
            {},
            jsonTransform<string>
        );

    protected requestItems = (): Promise<IPlugin[]> =>
        this.currentUser.authenticatedFetch(
            config.get('url') + 'plugin',
            {},
            jsonTransform<IPlugin[]>
        );

    protected requestDelete = (): Promise<void> => {
        this.snack('Cannot delete plugin');
        throw new Error('Cannot delete plugin');
    };

    public getName = (id: number): string => {
        const plugin = this.getByIDOrUndefined(id);
        return id === -1 ? 'All Plugins' : plugin !== undefined ? plugin.name : 'unknown';
    };

    @action
    public changeConfig = async (id: number, newConfig: string): Promise<void> => {
        await this.currentUser
            .authenticatedFetch(
                config.get('url') + 'plugin/' + id + '/config',
                yamlBody(newConfig),
                identityTransform
            )
            .then(() => this.snack('Plugin config updated'));
        await this.refresh();
    };

    @action
    public changeEnabledState = async (id: number, enabled: boolean): Promise<void> => {
        await this.currentUser
            .authenticatedFetch(
                config.get('url') + 'plugin/' + id + '/' + (enabled ? 'enable' : 'disable'),
                {method: 'POST'},
                identityTransform
            )
            .then(() => this.snack('Plugin ' + (enabled ? 'enabled' : 'disabled')));
        await this.refresh();
    };
}
