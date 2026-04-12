import React, {useState} from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import {observer} from 'mobx-react-lite';
import {useStores} from '../stores';
import ElevationForm from './ElevationForm';

interface IProps {
    fClose: VoidFunction;
}

const SettingsDialog = observer(({fClose}: IProps) => {
    const [pass, setPass] = useState('');
    const {currentUser, elevateStore} = useStores();

    const handleClose = () => {
        elevateStore.cleanupOidcElevate();
        fClose();
    };

    const submitAndClose = () => {
        currentUser.changePassword(pass);
        fClose();
    };

    return (
        <Dialog
            open={true}
            onClose={handleClose}
            aria-labelledby="form-dialog-title"
            id="changepw-dialog">
            <DialogTitle id="form-dialog-title">Change Password</DialogTitle>
            <DialogContent>
                {elevateStore.elevated ? (
                    <TextField
                        className="newpass"
                        autoFocus
                        margin="dense"
                        type="password"
                        label="New Password *"
                        value={pass}
                        onChange={(e) => setPass(e.target.value)}
                        fullWidth
                    />
                ) : (
                    <ElevationForm />
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>Cancel</Button>
                {elevateStore.elevated && (
                    <Tooltip title={pass.length !== 0 ? '' : 'Password is required'}>
                        <div>
                            <Button
                                className="change"
                                disabled={pass.length === 0}
                                onClick={submitAndClose}
                                color="primary"
                                variant="contained">
                                Change
                            </Button>
                        </div>
                    </Tooltip>
                )}
            </DialogActions>
        </Dialog>
    );
});

export default SettingsDialog;
