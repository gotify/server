import axios from 'axios';
import * as config from '../config.ts';
import {IUser} from '../types.ts';
import {AppDispatch} from './index.ts';
import {uiActions} from './ui-slice.ts';
import {userActions} from './user-slice.ts';

export const fetchUsers = () => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.get<IUser[]>(`${config.get('url')}user`);

            return response.data;
        };

        const data = await sendRequest();
        dispatch(userActions.set(data));
    };
};

export const deleteUser = (id: number) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.delete(`${config.get('url')}user/${id}`);

            return response.data;
        };

        await sendRequest();
        dispatch(userActions.remove(id));
        dispatch(uiActions.addSnackMessage('User deleted'));
    };
};

export const createUser = (name: string, pass: string | null, admin: boolean) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.post(`${config.get('url')}user`, {name, pass, admin});

            return response.data;
        }

        const data = await sendRequest();
        dispatch(userActions.add(data));
        dispatch(uiActions.addSnackMessage('User created'));
    }
}

export const updateUser = (id: number, name: string, pass: string | null, admin: boolean) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.post(config.get('url') + 'user/' + id, {name, pass, admin});
            return response.data;
        }
        const data = await sendRequest();
        dispatch(userActions.replace(data));
        dispatch(uiActions.addSnackMessage('User updated'));
    }
}
