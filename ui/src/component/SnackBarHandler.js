import React, {Component} from 'react';
import SnackBarStore from '../stores/SnackBarStore';
import Snackbar from 'material-ui/Snackbar';
import IconButton from 'material-ui/IconButton';
import Close from 'material-ui-icons/Close';

class SnackBarHandler extends Component {
    static MAX_VISIBLE_SNACK_TIME_IN_MS = 6000;
    static MIN_VISIBLE_SNACK_TIME_IN_MS = 1000;

    state = {
        current: '',
        open: false,
        openWhen: 0,
    };

    componentWillMount = () => SnackBarStore.on('change', this.onNewSnack);
    componentWillUnmount = () => SnackBarStore.removeListener('change', this.onNewSnack);

    onNewSnack = () => {
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

    openNextSnack = () => {
        if (SnackBarStore.hasNext()) {
            this.setState({
                ...this.state,
                open: true,
                openWhen: Date.now(),
                current: SnackBarStore.next(),
            });
        }
    };

    closeCurrentSnack = () => this.setState({...this.state, open: false});

    render() {
        const {open, current} = this.state;

        return (
            <Snackbar
                anchorOrigin={{vertical: 'bottom', horizontal: 'left'}}
                open={open} autoHideDuration={SnackBarHandler.MAX_VISIBLE_SNACK_TIME_IN_MS}
                onClose={this.closeCurrentSnack} onExited={this.openNextSnack}
                message={<span id="message-id">{current}</span>}
                action={[
                    <IconButton key="close" aria-label="Close" color="inherit" onClick={this.closeCurrentSnack}>
                        <Close/>
                    </IconButton>,
                ]}
            />
        );
    }
}

export default SnackBarHandler;
