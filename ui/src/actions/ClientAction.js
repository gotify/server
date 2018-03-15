import dispatcher from '../stores/dispatcher';
import config from 'react-global-configuration';
import axios from 'axios';

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
    axios.delete(config.get('url') + 'client/' + id).then(fetchClients);
}

/**
 * Create a client.
 * @param {string} name the client name
 */
export function createClient(name) {
    axios.post(config.get('url') + 'client', {name}).then(fetchClients);
}
