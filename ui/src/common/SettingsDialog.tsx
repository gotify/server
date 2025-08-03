import React, {useState} from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import {observer} from 'mobx-react';
import {useStores} from '../stores';

interface IProps {
    fClose: VoidFunction;
}

const SettingsDialog = observer(({fClose}: IProps) => {
    const [pass, setPass] = useState('');
    const {currentUser} = useStores();

    const submitAndClose = async () => {
        currentUser.changePassword(pass);
        fClose();
    };

    return (
        <Dialog
            open={true}
            onClose={fClose}
            aria-labelledby="form-dialog-title"
            id="changepw-dialog">
            <DialogTitle id="form-dialog-title">Change Password</DialogTitle>
            <DialogContent>
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
            </DialogContent>
            <DialogActions>
                <Button onClick={fClose}>Cancel</Button>
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
            </DialogActions>
        </Dialog>
    );
});

export default SettingsDialog;
