import axios, {AxiosResponse} from 'axios';
import * as config from '../config.ts';
import {IApplication, IMessage, IPagedMessages} from '../types.ts';
import {AppDispatch, RootState} from './index.ts';
import {messageActions} from './message-slice.ts';
import {uiActions} from './ui-slice.ts';

const AllMessages = -1;

const refreshByApp = (appId: number) => {
    return async (dispatch: AppDispatch) => {

    }
}

export const fetchMessages = (appId: number = AllMessages, since: number = 0) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async (url: string) => {
            const response = await axios.get(url);
            return response.data;
        };
        dispatch(messageActions.loading(true));
        let url;
        if (appId === AllMessages) {
            url = config.get('url') + 'message?since=' + since;
        } else {
            url = config.get('url') + 'application/' + appId + '/message?since=' + since;
        }
        const data = await sendRequest(url);
        dispatch(messageActions.set(data));
    }
};

export const removeSingleMessage = (message: IMessage) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async () => {
            const response = await axios.delete(config.get('url') + 'message/' + message.id);
            return response.data;
        }
        await sendRequest();
        dispatch(messageActions.remove(message.id));
        dispatch(uiActions.addSnackMessage('Message deleted'));
    }
}

export const removeMessagesByApp = (app: IApplication | undefined) => {
    return async (dispatch: AppDispatch) => {
        const sendRequest = async (url: string) => {
            const response = await axios.delete(url);
            return response.data;
        }
        let url;
        if (app === undefined) {
            url = config.get('url') + 'message';
            await sendRequest(url);
            dispatch(messageActions.clear());
            dispatch(uiActions.addSnackMessage('Deleted all messages'));
        } else {
            url = config.get('url') + 'application/' + app.id + '/message';
            await sendRequest(url);
            dispatch(messageActions.removeByAppId(app.id));
            dispatch(uiActions.addSnackMessage(`Deleted all messages from ${app.name}`));
        }
    }
}
