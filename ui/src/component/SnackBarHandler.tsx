import Close from 'material-ui-icons/Close';
import IconButton from 'material-ui/IconButton';
import Snackbar from 'material-ui/Snackbar';
import React, {Component} from 'react';
import SnackBarStore from '../stores/SnackBarStore';


interface IState {
    current: string
    hasNext: boolean
    open: boolean
    openWhen: number
}

class SnackBarHandler extends Component<{}, IState> {
    private static MAX_VISIBLE_SNACK_TIME_IN_MS = 6000;
    private static MIN_VISIBLE_SNACK_TIME_IN_MS = 1000;

    public state = {
        current: '',
        hasNext: false,
        open: false,
        openWhen: 0,
    };

    public componentWillMount() {
        SnackBarStore.on('change', this.onNewSnack);
    }

    public componentWillUnmount() {
        SnackBarStore.removeListener('change', this.onNewSnack);
    }

    public render() {
        const {open, current, hasNext} = this.state;
        const duration = hasNext
            ? SnackBarHandler.MIN_VISIBLE_SNACK_TIME_IN_MS
            : SnackBarHandler.MAX_VISIBLE_SNACK_TIME_IN_MS;

        return (
            <Snackbar
                anchorOrigin={{vertical: 'bottom', horizontal: 'left'}}
                open={open} autoHideDuration={duration}
                onClose={this.closeCurrentSnack} onExited={this.openNextSnack}
                message={<span id="message-id">{current}</span>}
                action={
                    <IconButton key="close" aria-label="Close" color="inherit" onClick={this.closeCurrentSnack}>
                        <Close/>
                    </IconButton>
                }
            />
        );
    }

    private onNewSnack = () => {
        const {open, openWhen} = this.state;

        if (!open) {
            this.openNextSnack();
            return;
        }

        const snackOpenSince = Date.now() - openWhen;
        if (snackOpenSince > SnackBarHandler.MIN_VISIBLE_SNACK_TIME_IN_MS) {
            this.closeCurrentSnack();
        } else {
            setTimeout(this.closeCurrentSnack, SnackBarHandler.MIN_VISIBLE_SNACK_TIME_IN_MS - snackOpenSince);
        }
    };

    private openNextSnack = () => {
        if (SnackBarStore.hasNext()) {
            this.setState({
                ...this.state,
                open: true,
                openWhen: Date.now(),
                current: SnackBarStore.next(),
                hasNext: SnackBarStore.hasNext(),
            });
        }
    };

    private closeCurrentSnack = () => this.setState({...this.state, open: false});
}

export default SnackBarHandler;
