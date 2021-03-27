import axios from 'axios';
import {action} from 'mobx';
import {BaseStore} from '../common/BaseStore';
import * as config from '../config';
import {SnackReporter} from '../snack/SnackManager';
import {IPlugin} from '../types';

export class PluginStore extends BaseStore<IPlugin> {
    public onDelete: () => void = () => {};

    public constructor(private readonly snack: SnackReporter) {
        super();
    }

    public requestConfig = (id: number): Promise<string> =>
        axios.get(`${config.get('url')}plugin/${id}/config`).then((response) => response.data);

    public requestDisplay = (id: number): Promise<string> =>
        axios.get(`${config.get('url')}plugin/${id}/display`).then((response) => response.data);

    protected requestItems = (): Promise<IPlugin[]> =>
        axios.get<IPlugin[]>(`${config.get('url')}plugin`).then((response) => response.data);

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
        await axios.post(`${config.get('url')}plugin/${id}/config`, newConfig, {
            headers: {'content-type': 'application/x-yaml'},
        });
        this.snack(`Plugin config updated`);
        await this.refresh();
    };

    @action
    public changeEnabledState = async (id: number, enabled: boolean): Promise<void> => {
        await axios.post(`${config.get('url')}plugin/${id}/${enabled ? 'enable' : 'disable'}`);
        this.snack(`Plugin ${enabled ? 'enabled' : 'disabled'}`);
        await this.refresh();
    };
}
