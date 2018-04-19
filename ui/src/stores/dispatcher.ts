import {Dispatcher} from 'flux';

export interface IEvent {
    type: string
    payload?: any
}

export default new Dispatcher<IEvent>();
