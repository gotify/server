import axios from 'axios';
import * as config from '../config.ts';
import {IUser} from '../types.ts';
import {AppDispatch} from './index.ts';
import {uiActions} from './ui-slice.ts';
import {userActions} from './user-slice.ts';

export const fetchUsers = () => {
    return async (dispatch: AppDispatch) => {
        const response = await axios.get<IUser[]>(`${config.get('url')}user`);
        dispatch(userActions.set(response.data));
    };
};

export const deleteUser = (id: number) => {
    return async (dispatch: AppDispatch) => {
        await axios.delete(`${config.get('url')}user/${id}`);
        dispatch(userActions.remove(id));
        dispatch(uiActions.addSnackMessage('User deleted'));
    };
};

export const createUser = (name: string, pass: string | null, admin: boolean) => {
    return async (dispatch: AppDispatch) => {
        const response = await axios.post(`${config.get('url')}user`, {name, pass, admin});
        dispatch(userActions.add(response.data));
        dispatch(uiActions.addSnackMessage('User created'));
    }
}

export const updateUser = (id: number, name: string, pass: string | null, admin: boolean) => {
    return async (dispatch: AppDispatch) => {
        const response = await axios.post(config.get('url') + 'user/' + id, {name, pass, admin});
        dispatch(userActions.replace(response.data));
        dispatch(uiActions.addSnackMessage('User updated'));
    }
}
