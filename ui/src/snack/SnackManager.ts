import {enqueueSnackbar} from 'notistack';

export interface SnackReporter {
    (message: string): void;
}

export class SnackManager {
    public snack: SnackReporter = (message: string): void => {
        enqueueSnackbar({message, variant: 'info'});
    };
}
