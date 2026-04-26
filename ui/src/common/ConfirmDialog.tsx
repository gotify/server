import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import React from 'react';
import {observer} from 'mobx-react-lite';
import {useStores} from '../stores';
import ElevationForm from './ElevationForm';

interface IProps {
    title: string;
    text: string;
    fClose: VoidFunction;
    fOnSubmit: VoidFunction;
    requireElevated?: boolean;
}

const ConfirmDialog = observer(({title, text, fClose, fOnSubmit, requireElevated}: IProps) => {
    const {elevateStore} = useStores();

    const needsElevation = requireElevated && !elevateStore.elevated;

    const submitAndClose = () => {
        fOnSubmit();
        fClose();
    };

    const handleClose = () => {
        elevateStore.cleanupOidcElevate();
        fClose();
    };

    return (
        <Dialog
            open={true}
            onClose={handleClose}
            aria-labelledby="form-dialog-title"
            className="confirm-dialog">
            <DialogTitle id="form-dialog-title">{title}</DialogTitle>
            <DialogContent>
                {needsElevation ? <ElevationForm /> : <DialogContentText>{text}</DialogContentText>}
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} className="cancel">
                    Cancel
                </Button>
                {!needsElevation && (
                    <Button
                        onClick={submitAndClose}
                        autoFocus
                        color="primary"
                        variant="contained"
                        className="confirm">
                        Yes
                    </Button>
                )}
            </DialogActions>
        </Dialog>
    );
});

export default ConfirmDialog;
