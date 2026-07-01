import React, {useState} from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import {NumberField} from '../common/NumberField';

interface IProps {
    fClose: (token: string | null) => void;
    fOnSubmit: (name: string, expiresAfterInactivitySeconds: number) => Promise<string>;
}

const AddClientDialog = ({fClose, fOnSubmit}: IProps) => {
    const [name, setName] = useState('');
    const [expiresAfter, setExpiresAfter] = useState(0);
    const submitEnabled = name.length !== 0;
    const submitAndNext = async () => {
        const token = await fOnSubmit(name, Math.max(0, expiresAfter));
        fClose(token);
    };

    return (
        <Dialog
            open={true}
            onClose={() => fClose(null)}
            aria-labelledby="form-dialog-title"
            id="client-dialog">
            <DialogTitle id="form-dialog-title">Create a client</DialogTitle>
            <DialogContent>
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
            </DialogContent>
            <DialogActions>
                <Button onClick={() => fClose(null)}>Cancel</Button>
                <Tooltip placement={'bottom-start'} title={submitEnabled ? '' : 'name is required'}>
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
            </DialogActions>
        </Dialog>
    );
};

export default AddClientDialog;
