import axios from 'axios';
import {action, observable, runInAction} from 'mobx';
import * as config from './config';
import {SnackReporter} from './snack/SnackManager';
import {CurrentUser} from './CurrentUser';

export class ElevateStore {
    @observable accessor elevated = false;
    @observable accessor oidcElevatePending = false;
    private oidcPollIntervalId: number | undefined = undefined;
    private oidcPopup: Window | null = null;

    public constructor(
        private readonly snack: SnackReporter,
        private readonly currentUser: CurrentUser
    ) {}

    @action
    public refreshElevated = (): number => {
        const elevatedUntil = this.currentUser.user.elevatedUntil;
        if (!elevatedUntil) {
            this.elevated = false;
            return 0;
        }
        const ms = new Date(elevatedUntil).getTime() - 30_000 - Date.now();
        if (ms <= 0) {
            this.elevated = false;
            return 0;
        }
        this.elevated = true;
        return ms;
    };

    public localElevate = async (password: string, durationSeconds: number): Promise<void> => {
        await axios.create().request({
            url: `${config.get('url')}client/${this.currentUser.user.clientId}/elevate`,
            method: 'POST',
            data: {durationSeconds},
            headers: {
                Authorization: 'Basic ' + btoa(this.currentUser.user.name + ':' + password),
            },
        });
        await this.currentUser.tryAuthenticate();
        this.cleanupOidcElevate();
    };

    public oidcElevate = (durationSeconds: number): void => {
        // prevent double execution
        if (this.oidcElevatePending) return;

        const url =
            config.get('url') +
            'auth/oidc/elevate?id=' +
            this.currentUser.user.clientId +
            '&durationSeconds=' +
            durationSeconds;

        this.oidcPopup = window.open(url, 'gotify-oidc-elevate', 'width=600,height=700');
        if (!this.oidcPopup) {
            this.snack('Popup was blocked. Please allow popups for this site and try again.');
            return;
        }

        runInAction(() => (this.oidcElevatePending = true));

        this.oidcPollIntervalId = window.setInterval(this.checkOidcPopup, 500);
    };

    private checkOidcPopup = async () => {
        if (this.oidcPopup && !this.oidcPopup.closed) {
            // waiting for the popup to close.
            return;
        }

        window.clearInterval(this.oidcPollIntervalId);
        this.oidcPollIntervalId = undefined;

        try {
            await this.currentUser.tryAuthenticate();
        } catch {
            // errors handled in tryAuthenticate
        }

        if (!this.elevated) {
            this.snack('OIDC elevation was not completed.');
        }
        this.cleanupOidcElevate();
    };

    public cleanupOidcElevate = () => {
        window.clearInterval(this.oidcPollIntervalId);
        this.oidcPollIntervalId = undefined;

        if (this.oidcPopup && !this.oidcPopup.closed) {
            this.oidcPopup.close();
        }
        this.oidcPopup = null;
        runInAction(() => {
            this.oidcElevatePending = false;
        });
    };
}
