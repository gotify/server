import {BaseStore} from '../common/BaseStore';
import axios from 'axios';
import * as config from '../config';
import {action} from 'mobx';
import {SnackReporter} from '../snack/SnackManager';
import {IUser} from '../types';

export class UserStore extends BaseStore<IUser> {
    private readonly snack: SnackReporter;

    constructor(snack: SnackReporter) {
        super();
        this.snack = snack;
        this.create = action(this.create);
        this.update = action(this.update);
    }

    protected requestItems = (): Promise<IUser[]> =>
        axios.get<IUser[]>(`${config.get('url')}user`).then((response) => response.data);

    protected requestDelete(id: number): Promise<void> {
        return axios
            .delete(`${config.get('url')}user/${id}`)
            .then(() => this.snack('User deleted'));
    }

    public create = async (name: string, pass: string, admin: boolean) => {
        await axios.post(`${config.get('url')}user`, {name, pass, admin});
        await this.refresh();
        this.snack('User created');
    };

    public update = async (id: number, name: string, pass: string | null, admin: boolean) => {
        await axios.post(config.get('url') + 'user/' + id, {name, pass, admin});
        await this.refresh();
        this.snack('User updated');
    };
}
