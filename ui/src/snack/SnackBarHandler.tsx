import IconButton from '@material-ui/core/IconButton';
import Snackbar from '@material-ui/core/Snackbar';
import Close from '@material-ui/icons/Close';
import React, {Component} from 'react';
import {observable, reaction} from 'mobx';
import {observer} from 'mobx-react';
import {inject, Stores} from '../inject';

@observer
class SnackBarHandler extends Component<Stores<'snackManager'>> {
    private static MAX_VISIBLE_SNACK_TIME_IN_MS = 6000;
    private static MIN_VISIBLE_SNACK_TIME_IN_MS = 1000;

    @observable
    private open = false;
    @observable
    private openWhen = 0;

    private dispose: () => void = () => {};

    public componentDidMount = () =>
        (this.dispose = reaction(() => this.props.snackManager.counter, this.onNewSnack));

    public componentWillUnmount = () => this.dispose();

    public render() {
        const {message: current, hasNext} = this.props.snackManager;
        const duration = hasNext()
            ? SnackBarHandler.MIN_VISIBLE_SNACK_TIME_IN_MS
            : SnackBarHandler.MAX_VISIBLE_SNACK_TIME_IN_MS;

        return (
            <Snackbar
                anchorOrigin={{vertical: 'bottom', horizontal: 'left'}}
                open={this.open}
                autoHideDuration={duration}
                onClose={this.closeCurrentSnack}
                onExited={this.openNextSnack}
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

    private openNextSnack = () => {
        if (this.props.snackManager.hasNext()) {
            this.open = true;
            this.openWhen = Date.now();
            this.props.snackManager.next();
        }
    };

    private closeCurrentSnack = () => (this.open = false);
}

export default inject('snackManager')(SnackBarHandler);
