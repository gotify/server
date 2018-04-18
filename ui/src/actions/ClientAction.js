import dispatcher from '../stores/dispatcher';
import * as config from '../config';
import axios from 'axios';
import {snack} from './GlobalAction';

/** Fetches all clients. */
export function fetchClients() {
    axios.get(config.get('url') + 'client').then((resp) => {
        dispatcher.dispatch({
            type: 'UPDATE_CLIENTS',
            payload: resp.data,
        });
    });
}

/**
 * Delete a client by id.
 * @param {int} id the client id
 */
export function deleteClient(id) {
    axios.delete(config.get('url') + 'client/' + id).then(fetchClients).then(() => snack('Client deleted'));
}

/**
 * Create a client.
 * @param {string} name the client name
 */
export function createClient(name) {
    axios.post(config.get('url') + 'client', {name}).then(fetchClients).then(() => snack('Client created'));
}
