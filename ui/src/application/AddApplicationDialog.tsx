import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import {NumberField} from '../common/NumberField';
import React, {useState} from 'react';
import {Typography} from '@mui/material';
import {copyToClipboard} from '../clipboard';

interface IProps {
    fClose: VoidFunction;
    fOnSubmit: (name: string, description: string, defaultPriority: number) => Promise<string>;
}

export const AddApplicationDialog = ({fClose, fOnSubmit}: IProps) => {
    const [returnToken, setReturnToken] = useState('');
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [defaultPriority, setDefaultPriority] = useState(0);

    const submitEnabled = name.length !== 0;
    const submitAndNext = async () => {
        const token = await fOnSubmit(name, description, defaultPriority);
        setReturnToken(token);
    };

    return (
        <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title" id="app-dialog">
            <DialogTitle id="form-dialog-title">Create an application</DialogTitle>
            <DialogContent>
                {returnToken ? (
                    <DialogContentText>
                        Your token is:
                        <Typography variant="body1" style={{fontFamily: 'monospace', fontSize: 16}}>
                            {returnToken}
                        </Typography>
                    </DialogContentText>
                ) : (
                    <>
                        <DialogContentText>
                            An application is allowed to send messages.
                        </DialogContentText>
                        <TextField
                            autoFocus
                            margin="dense"
                            className="name"
                            label="Name *"
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            fullWidth
                        />
                        <TextField
                            margin="dense"
                            className="description"
                            label="Short Description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            fullWidth
                            multiline
                        />
                        <NumberField
                            margin="dense"
                            className="priority"
                            label="Default Priority"
                            value={defaultPriority}
                            onChange={(value) => setDefaultPriority(value)}
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
                    <Button onClick={fClose} color="primary" variant="contained">
                        Close
                    </Button>
                ) : (
                    <Tooltip title={submitEnabled ? '' : 'name is required'}>
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
