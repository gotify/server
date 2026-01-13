import {CloseReason, enqueueSnackbar, SnackbarKey} from 'notistack';
import type {ReactNode, SyntheticEvent} from 'react';

export interface SnackReporter {
    (message: string, options?: SnackOptions): SnackbarKey;
}

export interface SnackOptions {
    action?: (key: number | string) => ReactNode;
    autoHideDuration?: number;
    variant?: 'default' | 'error' | 'success' | 'warning' | 'info';
    onClose?: (event: SyntheticEvent | null, reason: CloseReason, key?: SnackbarKey) => void;
    onExited?: (node: HTMLElement, key: SnackbarKey) => void;
}

export class SnackManager {
    public snack: SnackReporter = (message: string, options?: SnackOptions): SnackbarKey => {
        return enqueueSnackbar(message, {
            variant: options?.variant ?? 'info',
            action: options?.action,
            autoHideDuration: options?.autoHideDuration,
            onClose: options?.onClose,
            onExited: options?.onExited,
        });
    };
}
