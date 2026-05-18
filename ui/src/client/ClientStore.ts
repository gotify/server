import {BaseStore} from '../common/BaseStore';
import axios from 'axios';
import * as config from '../config';
import {action} from 'mobx';
import {SnackReporter} from '../snack/SnackManager';
import {IClient} from '../types';

export class ClientStore extends BaseStore<IClient> {
    public constructor(private readonly snack: SnackReporter) {
        super();
    }

    protected requestItems = (): Promise<IClient[]> =>
        axios.get<IClient[]>(`${config.get('url')}client`).then((response) => response.data);

    protected requestDelete(id: number): Promise<void> {
        return axios
            .delete(`${config.get('url')}client/${id}`)
            .then(() => this.snack('Client deleted'));
    }

    @action
    public update = async (
        id: number,
        name: string,
        expiresAfterInactivitySeconds: number
    ): Promise<void> => {
        await axios.put(`${config.get('url')}client/${id}`, {
            name,
            expiresAfterInactivitySeconds,
        });
        await this.refresh();
        this.snack('Client updated');
    };

    @action
    public createNoNotifcation = async (
        name: string,
        expiresAfterInactivitySeconds = 0
    ): Promise<IClient> => {
        const client = await axios.post(`${config.get('url')}client`, {
            name,
            expiresAfterInactivitySeconds,
        });
        await this.refresh();
        return client.data;
    };

    @action
    public create = async (name: string, expiresAfterInactivitySeconds = 0): Promise<void> => {
        await this.createNoNotifcation(name, expiresAfterInactivitySeconds);
        this.snack('Client added');
    };

    @action
    public elevate = async (id: number, durationSeconds: number): Promise<void> => {
        await axios.post(`${config.get('url')}client/${id}/elevate`, {durationSeconds});
        await this.refresh();
        if (durationSeconds < 0) {
            this.snack('Canceled client elevation');
        } else {
            this.snack('Client elevated');
        }
    };
}
