import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';

interface IProps {
    title: string;
    text: string;
    fClose: VoidFunction;
    fOnSubmit: VoidFunction;
}

export default function ConfirmDialog({title, text, fClose, fOnSubmit}: IProps) {
    const submitAndClose = () => {
        fOnSubmit();
        fClose();
    };
    return (
        <Dialog
            open={true}
            onClose={fClose}
            aria-labelledby="form-dialog-title"
            className="confirm-dialog">
            <DialogTitle id="form-dialog-title">{title}</DialogTitle>
            <DialogContent>
                <DialogContentText>{text}</DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={fClose} className="cancel">
                    No
                </Button>
                <Button
                    onClick={submitAndClose}
                    autoFocus
                    color="primary"
                    variant="contained"
                    className="confirm">
                    Yes
                </Button>
            </DialogActions>
        </Dialog>
    );
}
