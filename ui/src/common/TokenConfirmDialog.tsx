import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import React from 'react';
import {Typography} from '@mui/material';
import {useStores} from '../stores';
import {copyToClipboard} from '../clipboard';

interface IProps {
    token: string;
    fClose: VoidFunction;
}

export const TokenConfirmDialog = ({token, fClose}: IProps) => {
    const {snackManager} = useStores();

    return (
        <Dialog open={true} aria-labelledby="form-dialog-title" id="token-dialog">
            <DialogTitle id="form-dialog-title">Token Confirmation</DialogTitle>
            <DialogContent>
                <DialogContentText>Your token will only be shown once.</DialogContentText>
                <span
                    style={{padding: 16}}
                    onClick={(e) => {
                        window.getSelection()?.selectAllChildren(e.currentTarget);
                    }}>
                    <Typography
                        className="token"
                        variant="body1"
                        style={{fontFamily: 'monospace', fontSize: 16}}>
                        {token}
                    </Typography>
                </span>
            </DialogContent>
            <DialogActions>
                <Button
                    onClick={() =>
                        copyToClipboard(token).then(
                            () => snackManager.snack('Copied to clipboard'),
                            () => snackManager.snack('Cannot access clipboard.')
                        )
                    }>
                    Copy to clipboard
                </Button>
                <Button className="finish" onClick={fClose} color="primary" variant="contained">
                    Finish
                </Button>
            </DialogActions>
        </Dialog>
    );
};
