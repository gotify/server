import React, {Component} from 'react';
import Button from 'material-ui/Button';
import Dialog, {DialogActions, DialogContent, DialogContentText, DialogTitle} from 'material-ui/Dialog';
import PropTypes from 'prop-types';

export default class ConfirmDialog extends Component {
    static propTypes = {
        title: PropTypes.string.isRequired,
        text: PropTypes.string.isRequired,
        fClose: PropTypes.func.isRequired,
        fOnSubmit: PropTypes.func.isRequired,
    };

    render() {
        const {title, text, fClose, fOnSubmit} = this.props;
        const submitAndClose = () => {
            fOnSubmit();
            fClose();
        };
        return (
            <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">{title}</DialogTitle>
                <DialogContent>
                    <DialogContentText>{text}</DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>No</Button>
                    <Button onClick={submitAndClose} color="primary" variant="raised">Yes</Button>
                </DialogActions>
            </Dialog>
        );
    }
}
