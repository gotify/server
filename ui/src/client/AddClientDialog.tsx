import React, {useState} from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';

interface IProps {
    fClose: VoidFunction;
    fOnSubmit: (name: string) => Promise<void>;
}

const AddClientDialog = ({fClose, fOnSubmit}: IProps) => {
    const [name, setName] = useState('');

    const submitEnabled = name.length !== 0;
    const submitAndClose = async () => {
        await fOnSubmit(name);
        fClose();
    };

    return (
        <Dialog
            open={true}
            onClose={fClose}
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
            </DialogContent>
            <DialogActions>
                <Button onClick={fClose}>Cancel</Button>
                <Tooltip
                    placement={'bottom-start'}
                    title={submitEnabled ? '' : 'name is required'}>
                    <div>
                        <Button
                            className="create"
                            disabled={!submitEnabled}
                            onClick={submitAndClose}
                            color="primary"
                            variant="contained">
                            Create
                        </Button>
                    </div>
                </Tooltip>
            </DialogActions>
        </Dialog>
    );
}

export default AddClientDialog;
