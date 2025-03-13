import axios from 'axios';
import * as config from '../config.ts';
import {IPlugin} from '../types.ts';
import {AppDispatch, RootState} from '../store/index.ts';
import {pluginActions} from './plugin-slice.ts';
import {uiActions} from '../store/ui-slice.ts';

export const requestPluginConfig = (id: number) => {
    return async (dispatch: AppDispatch) => {
        dispatch(pluginActions.loading(true));
        const response = await axios.get(`${config.get('url')}plugin/${id}/config`);
        dispatch(pluginActions.setDisplay(response.data));
    };
};

export const requestPluginDisplay = (id: number) => {
    return async (dispatch: AppDispatch) => {
        dispatch(pluginActions.loading(true));
        const response = await axios.get(`${config.get('url')}plugin/${id}/display`);
        dispatch(pluginActions.setCurrentConfig(response.data));
    };
};

export const fetchPlugins = () => {
    return async (dispatch: AppDispatch) => {
        dispatch(pluginActions.loading(true));
        const response = await axios.get<IPlugin[]>(`${config.get('url')}plugin`);
        dispatch(pluginActions.set(response.data));
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
        dispatch(pluginActions.loading(true));
        await axios.post(`${config.get('url')}plugin/${id}/config`, newConfig, {
            headers: {'content-type': 'application/x-yaml'},
        });
        dispatch(uiActions.addSnackMessage('Plugin config updated'));
        await dispatch(fetchPlugins())
    }
}

export const changePluginEnableState = (id: number, enabled: boolean) => {
    return async (dispatch: AppDispatch) => {
        dispatch(pluginActions.loading(true));
        await axios.post(`${config.get('url')}plugin/${id}/${enabled ? 'enable' : 'disable'}`);
        dispatch(uiActions.addSnackMessage(`Plugin ${enabled ? 'enabled' : 'disabled'}`));
        await dispatch(fetchPlugins());
    }
}
