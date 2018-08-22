import axios, {AxiosResponse} from 'axios';
import * as config from '../config';
import dispatcher from '../stores/dispatcher';
import {snack} from './GlobalAction';

/** Fetches all applications. */
export function fetchApps() {
    axios.get(config.get('url') + 'application').then((resp: AxiosResponse<IMessage[]>) => {
        dispatcher.dispatch({type: 'UPDATE_APPS', payload: resp.data});
    });
}

/**
 * Delete an application by id.
 * @param {int} id the application id
 */
export function deleteApp(id: number) {
    axios
        .delete(config.get('url') + 'application/' + id)
        .then(() => {
            fetchApps();
            dispatcher.dispatch({type: 'DELETE_MESSAGES', payload: id});
        })
        .then(() => snack('Application deleted'));
}

/**
 * Create an application.
 * @param {string} name the application name
 * @param {string} description the description of the application.
 */
export function createApp(name: string, description: string) {
    axios
        .post(config.get('url') + 'application', {name, description})
        .then(fetchApps)
        .then(() => snack('Application created'));
}

/**
 * Upload an image for an application.
 * @param {int} id the application id
 * @param {Blob} file the description of the application.
 */
export function uploadImage(id: number, file: Blob) {
    const formData = new FormData();
    formData.append('file', file);
    axios
        .post(config.get('url') + 'application/' + id + '/image', formData, {
            headers: {'content-type': 'multipart/form-data'},
        })
        .then(fetchApps)
        .then(() => snack('Application image updated'));
}
