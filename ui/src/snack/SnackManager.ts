import {enqueueSnackbar, SnackbarKey} from 'notistack';
import {ReactNode} from 'react';

export interface SnackReporter {
    (message: string, options?: SnackOptions): SnackbarKey;
}

export interface SnackOptions {
    action?: (key: number | string) => ReactNode;
    autoHideDuration?: number;
    variant?: 'default' | 'error' | 'success' | 'warning' | 'info';
}

export class SnackManager {
    public snack: SnackReporter = (message: string, options?: SnackOptions): SnackbarKey => {
        return enqueueSnackbar({
            message,
            variant: options?.variant ?? 'info',
            action: options?.action,
            autoHideDuration: options?.autoHideDuration,
        });
    };
}
