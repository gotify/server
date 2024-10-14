import IconButton from '@material-ui/core/IconButton';
import Snackbar from '@material-ui/core/Snackbar';
import Close from '@material-ui/icons/Close';
import React, {Component} from 'react';
import { observable, reaction, toJS, action } from 'mobx';
import { observer } from 'mobx-react';
import { inject, Stores } from '../inject';

@observer
class SnackBarHandler extends Component<Stores<'snackManager'>> {
    private static MAX_VISIBLE_SNACK_TIME_IN_MS = 6000;
    private static MIN_VISIBLE_SNACK_TIME_IN_MS = 1000;

    @observable
    private open = false;
    @observable
    private openWhen = 0;
    @observable
    private snackManager: any = null;

    private dispose: () => void = () => {};

    @action
    public componentDidMount = () => {
        this.snackManager = this.props.snackManager;

        this.dispose = reaction(
            () => toJS(this.snackManager.counter),
            this.onNewSnack
        );
    }

    public componentWillUnmount = () => this.dispose();

    public render() {
        if (!this.snackManager) return null;

        const {message: current, hasNext} = this.snackManager;
        const duration = hasNext()
            ? SnackBarHandler.MIN_VISIBLE_SNACK_TIME_IN_MS
            : SnackBarHandler.MAX_VISIBLE_SNACK_TIME_IN_MS;

        return (
            <Snackbar
                anchorOrigin={{vertical: 'bottom', horizontal: 'left'}}
                open={this.open}
                autoHideDuration={duration}
                onClose={this.closeCurrentSnack}
                TransitionProps={{ onExited: this.openNextSnack }}
                message={<span id="message-id">{current}</span>}
                action={
                    <IconButton
                        key="close"
                        aria-label="Close"
                        color="inherit"
                        onClick={this.closeCurrentSnack}>
                        <Close />
                    </IconButton>
                }
            />
        );
    }

    @action
    private onNewSnack = () => {
        const {open, openWhen} = this;

        if (!open) {
            this.openNextSnack();
            return;
        }

        const snackOpenSince = Date.now() - openWhen;
        if (snackOpenSince > SnackBarHandler.MIN_VISIBLE_SNACK_TIME_IN_MS) {
            this.closeCurrentSnack();
        } else {
            setTimeout(
                this.closeCurrentSnack,
                SnackBarHandler.MIN_VISIBLE_SNACK_TIME_IN_MS - snackOpenSince
            );
        }
    };

    @action
    private openNextSnack = () => {
        if (this.snackManager?.hasNext()) {
            this.open = true;
            this.openWhen = Date.now();
            this.snackManager.next();
        }
    };

    @action
    private closeCurrentSnack = () => (this.open = false);
}

export default inject('snackManager')(SnackBarHandler);
