import axios from 'axios';
import * as config from '../config.ts';
import {IClient} from '../types.ts';
import {clientActions} from './client-slice.ts';
import {AppDispatch} from './index.ts';
import {uiActions} from './ui-slice.ts';

export const fetchClients = () => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.get<IClient[]>(`${config.get('url')}client`);
            return response.data;
        };

        const data = await sendRequest();
        dispatch(clientActions.set(data));
    };
};

export const deleteClient = (id: number) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.delete<IClient>(`${config.get('url')}client/${id}`);
            return response.data;
        };
        await sendRequest();
        dispatch(clientActions.remove(id));
        dispatch(uiActions.addSnackMessage('Client deleted'));
    };
};

export const updateClient = (id: number, name: string) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.put<IClient>(`${config.get('url')}client/${id}`, {name});
            return response.data;
        }
        const data = await sendRequest();
        dispatch(clientActions.replace(data));
        dispatch(uiActions.addSnackMessage('Client deleted'));
    }
}

export const createClientNoNotification = (name: string) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.post<IClient>(`${config.get('url')}client`, {name});
            return response.data;
        }
        const data = await sendRequest();
        dispatch(clientActions.add(data));
    }
}

export const createClient = (name: string) => {
    return async (dispatch: AppDispatch) => {
        await dispatch(createClientNoNotification(name));
        dispatch(uiActions.addSnackMessage('Client added'));
    }
}
