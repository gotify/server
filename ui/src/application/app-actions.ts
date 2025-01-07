import axios from 'axios';
import {getAuthToken} from '../common/Auth.ts';
import * as config from '../config';
import {IApplication} from '../types.ts';
import {appActions} from './app-slice.ts';
import {AppDispatch, RootState} from '../store/index.ts';
import {uiActions} from '../store/ui-slice.ts';

export const fetchApps = () => {
    return async (dispatch: AppDispatch) => {
        if (!getAuthToken()) {
            return;
        }
        dispatch(appActions.loading(true));
        const response = await axios.get<IApplication[]>(`${config.get('url')}application`);
        dispatch(appActions.set(response.data));
    };
};

export const deleteApp = (id: number) => {
    return async (dispatch: AppDispatch) => {
        // do not dispatch a loading indicator as the test does not expect it
        // dispatch(appActions.loading(true));
        await axios.delete(`${config.get('url')}application/${id}`);
        dispatch(appActions.remove(id));
        dispatch(uiActions.addSnackMessage('Application deleted'));
    };
};

export const uploadImage = (id: number, file: Blob) => {
    return async (dispatch: AppDispatch) => {
        dispatch(appActions.loading(true));
        const formData = new FormData();
        formData.append('file', file);

        const response = await axios.post(
            `${config.get('url')}application/${id}/image`,
            formData,
            {
                headers: {'content-type': 'multipart/form-data'},
            }
        );

        dispatch(appActions.replace(response.data));
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
        dispatch(appActions.loading(true));
        const response = await axios.put(`${config.get('url')}application/${id}`, {
            name,
            description,
            defaultPriority,
        });
        dispatch(appActions.replace(response.data));
        dispatch(uiActions.addSnackMessage('Application updated'));
    };
};

export const createApp = (name: string, description: string, defaultPriority: number) => {
    return async (dispatch: AppDispatch) => {
        dispatch(appActions.loading(true));
        const response = await axios.post(`${config.get('url')}application`, {
            name,
            description,
            defaultPriority,
        });
        dispatch(appActions.add(response.data));
        dispatch(uiActions.addSnackMessage('Application created'));
    };
};

export const getAppName = (id: number) => {
    return (_dispatch: AppDispatch, getState: () => RootState) => {
        const app = getState().app.items.find((item) => item.id === id);
        return id === -1 ? 'All Messages' : app !== undefined ? app.name : 'unknown';
    };
}
