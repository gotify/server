import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import FormControlLabel from '@mui/material/FormControlLabel';
import Switch from '@mui/material/Switch';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import {NumberField} from '../common/NumberField';
import React, {useState} from 'react';

interface IProps {
    fClose: VoidFunction;
    fOnSubmit: (
        name: string,
        description: string,
        defaultPriority: number,
        allowMessageDelete: boolean
    ) => Promise<void>;
    initialName: string;
    initialDescription: string;
    initialDefaultPriority: number;
    initialAllowMessageDelete: boolean;
}

export const UpdateApplicationDialog = ({
    initialName,
    initialDescription,
    initialDefaultPriority,
    initialAllowMessageDelete,
    fClose,
    fOnSubmit,
}: IProps) => {
    const [name, setName] = useState(initialName);
    const [description, setDescription] = useState(initialDescription);
    const [defaultPriority, setDefaultPriority] = useState(initialDefaultPriority);
    const [allowMessageDelete, setAllowMessageDelete] = useState(initialAllowMessageDelete);

    const submitEnabled = name.length !== 0;
    const submitAndClose = async () => {
        await fOnSubmit(name, description, defaultPriority, allowMessageDelete);
        fClose();
    };

    return (
        <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title" id="app-dialog">
            <DialogTitle id="form-dialog-title">Update an application</DialogTitle>
            <DialogContent>
                <DialogContentText>An application is allowed to send messages.</DialogContentText>
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
                    onChange={(e) => setDefaultPriority(e)}
                    fullWidth
                />
                <FormControlLabel
                    control={
                        <Switch
                            checked={allowMessageDelete}
                            className="allow-message-delete"
                            onChange={(e) => setAllowMessageDelete(e.target.checked)}
                        />
                    }
                    label="Allow the application to delete its own messages"
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={fClose}>Cancel</Button>
                <Tooltip title={submitEnabled ? '' : 'name is required'}>
                    <div>
                        <Button
                            className="update"
                            disabled={!submitEnabled}
                            onClick={submitAndClose}
                            color="primary"
                            variant="contained">
                            Update
                        </Button>
                    </div>
                </Tooltip>
            </DialogActions>
        </Dialog>
    );
};
