import axios from 'axios';
import * as config from '../config.ts';
import {IClient} from '../types.ts';
import {clientActions} from './client-slice.ts';
import {AppDispatch} from '../store/index.ts';
import {uiActions} from '../store/ui-slice.ts';

export const fetchClients = () => {
    return async (dispatch: AppDispatch) => {
        const response = await axios.get<IClient[]>(`${config.get('url')}client`);
        dispatch(clientActions.set(response.data));
    };
};

export const deleteClient = (id: number) => {
    return async (dispatch: AppDispatch) => {
        await axios.delete<IClient>(`${config.get('url')}client/${id}`);
        dispatch(clientActions.remove(id));
        dispatch(uiActions.addSnackMessage('Client deleted'));
    };
};

export const updateClient = (id: number, name: string) => {
    return async (dispatch: AppDispatch) => {
        const response = await axios.put<IClient>(`${config.get('url')}client/${id}`, {name});
        dispatch(clientActions.replace(response.data));
        dispatch(uiActions.addSnackMessage('Client deleted'));
    }
}

export const createClientNoNotification = (name: string) => {
    return async (dispatch: AppDispatch) => {
        const response = await axios.post<IClient>(`${config.get('url')}client`, {name});
        dispatch(clientActions.add(response.data));
    }
}

export const createClient = (name: string) => {
    return async (dispatch: AppDispatch) => {
        await dispatch(createClientNoNotification(name));
        dispatch(uiActions.addSnackMessage('Client added'));
    }
}
