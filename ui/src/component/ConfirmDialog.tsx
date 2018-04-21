import Button from 'material-ui/Button';
import Dialog, {DialogActions, DialogContent, DialogContentText, DialogTitle} from 'material-ui/Dialog';
import React from 'react';

interface IProps {
    title: string
    text: string
    fClose: VoidFunction
    fOnSubmit: VoidFunction
}

export default function ConfirmDialog({title, text, fClose, fOnSubmit}: IProps) {
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