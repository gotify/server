import {action, observable} from 'mobx';

export interface SnackReporter {
    (message: string): void;
}

export class SnackManager {
    @observable
    private messages: string[] = [];
    @observable
    public message: string | null = null;
    @observable
    public counter = 0;

    @action
    public next = (): void => {
        if (!this.hasNext()) {
            throw new Error('There is nothing here :(');
        }
        this.message = this.messages.shift() as string;
    };

    public hasNext = () => this.messages.length > 0;

    @action
    public snack: SnackReporter = (message: string): void => {
        this.messages.push(message);
        this.counter++;
    };
}
