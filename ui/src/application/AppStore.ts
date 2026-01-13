import {BaseStore} from '../common/BaseStore';
import axios from 'axios';
import * as config from '../config';
import {action, makeObservable} from 'mobx';
import {SnackReporter} from '../snack/SnackManager';
import {IApplication} from '../types';

export class AppStore extends BaseStore<IApplication> {
    public onDelete: () => void = () => {};

    public constructor(private readonly snack: SnackReporter) {
        super();

        makeObservable(this, {
            uploadImage: action,
            update: action,
            create: action,
        });
    }

    protected requestItems = (): Promise<IApplication[]> =>
        axios
            .get<IApplication[]>(`${config.get('url')}application`)
            .then((response) => response.data);

    protected requestDelete = (id: number): Promise<void> =>
        axios.delete(`${config.get('url')}application/${id}`).then(() => {
            this.onDelete();
            this.snack('Application deleted');
        });

    public uploadImage = async (id: number, file: Blob): Promise<void> => {
        const formData = new FormData();
        formData.append('file', file);
        await axios.post(`${config.get('url')}application/${id}/image`, formData, {
            headers: {'content-type': 'multipart/form-data'},
        });
        await this.refresh();
        this.snack('Application image updated');
    };

    public async deleteImage(id: number): Promise<void> {
        try {
            await axios.delete(`${config.get('url')}application/${id}/image`);
            await this.refresh();
            this.snack('Application image deleted');
        } catch (error) {
            console.error('Error deleting application image:', error);
            throw error;
        }
    }

    public update = async (
        id: number,
        name: string,
        description: string,
        defaultPriority: number
    ): Promise<void> => {
        await axios.put(`${config.get('url')}application/${id}`, {
            name,
            description,
            defaultPriority,
        });
        await this.refresh();
        this.snack('Application updated');
    };

    public create = async (
        name: string,
        description: string,
        defaultPriority: number
    ): Promise<void> => {
        await axios.post(`${config.get('url')}application`, {
            name,
            description,
            defaultPriority,
        });
        await this.refresh();
        this.snack('Application created');
    };

    public getName = (id: number): string => {
        const app = this.getByIDOrUndefined(id);
        return id === -1 ? 'All Messages' : app !== undefined ? app.name : 'unknown';
    };
}
