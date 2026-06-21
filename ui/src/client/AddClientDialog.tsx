import React, {useState} from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import {NumberField} from '../common/NumberField';
import {DialogContentText, Typography} from '@mui/material';
import {copyToClipboard} from '../clipboard';

interface IProps {
    fClose: VoidFunction;
    fOnSubmit: (name: string, expiresAfterInactivitySeconds: number) => Promise<string>;
}

const AddClientDialog = ({fClose, fOnSubmit}: IProps) => {
    const [returnToken, setReturnToken] = useState('');
    const [name, setName] = useState('');
    const [expiresAfter, setExpiresAfter] = useState(0);

    const submitEnabled = name.length !== 0;
    const submitAndNext = async () => {
        const token = await fOnSubmit(name, Math.max(0, expiresAfter));
        setReturnToken(token);
    };

    return (
        <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title" id="client-dialog">
            <DialogTitle id="form-dialog-title">Create a client</DialogTitle>
            <DialogContent>
                {returnToken ? (
                    <>
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
                                {returnToken}
                            </Typography>
                        </span>
                    </>
                ) : (
                    <>
                        <TextField
                            autoFocus
                            margin="dense"
                            className="name"
                            label="Name *"
                            type="email"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            fullWidth
                        />
                        <NumberField
                            margin="dense"
                            className="expires-after"
                            label="Expires after inactivity (seconds, 0 = never)"
                            value={expiresAfter}
                            onChange={(value) => setExpiresAfter(value)}
                            fullWidth
                        />
                    </>
                )}
            </DialogContent>
            <DialogActions>
                {returnToken ? (
                    <Button onClick={() => copyToClipboard(returnToken)}>Copy to clipboard</Button>
                ) : (
                    <Button onClick={fClose}>Cancel</Button>
                )}
                {returnToken ? (
                    <Button className="finish" onClick={fClose} color="primary" variant="contained">
                        Finish
                    </Button>
                ) : (
                    <Tooltip
                        placement={'bottom-start'}
                        title={submitEnabled ? '' : 'name is required'}>
                        <div>
                            <Button
                                className="create"
                                disabled={!submitEnabled}
                                onClick={submitAndNext}
                                color="primary"
                                variant="contained">
                                Create
                            </Button>
                        </div>
                    </Tooltip>
                )}
            </DialogActions>
        </Dialog>
    );
};

export default AddClientDialog;
