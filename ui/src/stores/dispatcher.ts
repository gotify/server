import {Dispatcher} from 'flux';

export interface IEvent {
    type: string;
    // tslint:disable-next-line
    payload?: any;
}

export default new Dispatcher<IEvent>();
