import axios from 'axios';
import {getAuthToken} from '../common/Auth.ts';
import * as config from '../config';
import {IApplication} from '../types.ts';
import {appActions} from './app-slice.ts';
import {AppDispatch, RootState} from './index.ts';
import {uiActions} from './ui-slice.ts';

export const fetchApps = () => {
    return async (dispatch: AppDispatch) => {
        if (!getAuthToken()) {
            return;
        }
        const sendRequest = async () => {
            const response = await axios.get<IApplication[]>(`${config.get('url')}application`);

            return response.data;
        };

        // TODO: handle error case
        const appData = await sendRequest();
        dispatch(appActions.set(appData));
    };
};

export const deleteApp = (id: number) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.delete(`${config.get('url')}application/${id}`);
            return response.data;
        };

        // TODO: handle error case
        await sendRequest();
        dispatch(appActions.remove(id));
        dispatch(uiActions.addSnackMessage('Application deleted'));
    };
};

export const uploadImage = (id: number, file: Blob) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.post(
                `${config.get('url')}application/${id}/image`,
                formData,
                {
                    headers: {'content-type': 'multipart/form-data'},
                }
            );
            return response.data;
        };

        const formData = new FormData();
        formData.append('file', file);
        const data = await sendRequest();
        dispatch(appActions.replace(data));
        dispatch(uiActions.addSnackMessage('Application image updated'));
    };
};

export const updateApp = (
    id: number,
    name: string,
    description: string,
    defaultPriority: number
) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.put(`${config.get('url')}application/${id}`, {
                name,
                description,
                defaultPriority,
            });
            return response.data;
        };
        const data = await sendRequest();
        dispatch(appActions.replace(data));
        dispatch(uiActions.addSnackMessage('Application updated'));
    };
};

export const createApp = (name: string, description: string, defaultPriority: number) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.post(`${config.get('url')}application`, {
                name,
                description,
                defaultPriority,
            });
            return response.data;
        };
        const data = await sendRequest();
        dispatch(appActions.add(data));
        dispatch(uiActions.addSnackMessage('Application created'));
    };
};

export const getAppName = (id: number) => {
    return (_dispatch: AppDispatch, getState: () => RootState) => {
        const app = getState().app.items.find((item) => item.id === id);
        return id === -1 ? 'All Messages' : app !== undefined ? app.name : 'unknown';
    };
}
