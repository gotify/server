import axios from 'axios';
import * as config from '../config.ts';
import {IApplication, IMessage, IPagedMessages} from '../types.ts';
import {AppDispatch} from '../store/index.ts';
import {messageActions} from './message-slice.ts';
import {uiActions} from '../store/ui-slice.ts';

export const AllMessages = -1;

export const fetchMessages = (appId: number = AllMessages, since: number = 0) => {
    return async (dispatch: AppDispatch) => {
        dispatch(messageActions.loading(true));
        let url;
        if (appId === AllMessages) {
            url = config.get('url') + 'message?since=' + since;
        } else {
            url = config.get('url') + 'application/' + appId + '/message?since=' + since;
        }

        try {
            const response = await axios.get<IPagedMessages>(url);
            dispatch(messageActions.set({appId, pagedMessages: response.data}));
        } catch (error) {
            dispatch(messageActions.loading(false));
        }
    };
};

export const removeSingleMessage = (message: IMessage) => {
    return async (dispatch: AppDispatch) => {
        await axios.delete(config.get('url') + 'message/' + message.id);
        dispatch(messageActions.remove(message.id));
        dispatch(uiActions.addSnackMessage('Message deleted'));
    };
};

export const removeMessagesByApp = (app: IApplication | undefined) => {
    return async (dispatch: AppDispatch) => {
        let url;
        if (app === undefined) {
            url = config.get('url') + 'message';
            await axios.delete(url);
            dispatch(messageActions.clear());
            dispatch(uiActions.addSnackMessage('Deleted all messages'));
        } else {
            url = config.get('url') + 'application/' + app.id + '/message';
            await axios.delete(url);
            dispatch(messageActions.removeByAppId(app.id));
            dispatch(uiActions.addSnackMessage(`Deleted all messages from ${app.name}`));
        }
    };
};
