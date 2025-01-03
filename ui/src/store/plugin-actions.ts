import axios from 'axios';
import * as config from '../config.ts';
import {IPlugin} from '../types.ts';
import {AppDispatch, RootState} from './index.ts';
import {pluginActions} from './plugin-slice.ts';
import {uiActions} from './ui-slice.ts';

export const requestPluginConfig = (id: number) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.get(`${config.get('url')}plugin/${id}/config`);
            return response.data;
        };

        const data = await sendRequest();
        dispatch(pluginActions.setDisplay(data));
    };
};

export const requestPluginDisplay = (id: number) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.get(`${config.get('url')}plugin/${id}/display`);
            return response.data;
        };
        const data = await sendRequest();
        dispatch(pluginActions.setCurrentConfig(data));
    };
};

export const fetchPlugins = () => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.get<IPlugin[]>(`${config.get('url')}plugin`);
            return response.data;
        };
        const data = await sendRequest();
        dispatch(pluginActions.set(data));
    };
};

export const deletePlugin = () => {
    return async (dispatch: AppDispatch) => {
        dispatch(uiActions.addSnackMessage('Cannot delete plugin'));
        throw new Error('Cannot delete plugin');
    };
};

export const getPluginName = (id: number) => {
    return (_dispatch: AppDispatch, getState: () => RootState) => {
        const plugin = getState().plugin.items.filter(plugin => plugin.id === id);

        return id === -1 ? 'All Plugins': plugin[0].name || 'unknown';
    };
};

export const updatePluginConfig = (id: number, newConfig: string) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.post(`${config.get('url')}plugin/${id}/config`, newConfig, {
                headers: {'content-type': 'application/x-yaml'},
            });
            return response.data;
        }
        await sendRequest();
        dispatch(uiActions.addSnackMessage('Plugin config updated'));
        await dispatch(fetchPlugins())
    }
}

export const changePluginEnableState = (id: number, enabled: boolean) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.post(`${config.get('url')}plugin/${id}/${enabled ? 'enable' : 'disable'}`);
            return response.data;
        }
        await sendRequest();
        dispatch(uiActions.addSnackMessage(`Plugin ${enabled ? 'enabled' : 'disabled'}`));
        await dispatch(fetchPlugins());
    }
}
