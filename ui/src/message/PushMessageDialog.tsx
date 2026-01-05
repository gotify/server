import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import React, {useState} from 'react';
import {NumberField} from '../common/NumberField';

interface IProps {
    appName: string;
    defaultPriority: number;
    fClose: VoidFunction;
    fOnSubmit: (message: string, title: string, priority: number) => Promise<void>;
}

export const PushMessageDialog = ({appName, defaultPriority, fClose, fOnSubmit}: IProps) => {
    const [title, setTitle] = useState('');
    const [message, setMessage] = useState('');
    const [priority, setPriority] = useState(defaultPriority);

    const submitEnabled = message.trim().length !== 0;
    const submitAndClose = async () => {
        await fOnSubmit(message, title, priority);
        fClose();
    };

    return (
        <Dialog
            open={true}
            onClose={fClose}
            aria-labelledby="push-message-title"
            id="push-message-dialog">
            <DialogTitle id="push-message-title">Push message</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    Send a push message via {appName}. Leave the title empty to use the
                    application name.
                </DialogContentText>
                <TextField
                    margin="dense"
                    className="title"
                    label="Title"
                    type="text"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    fullWidth
                />
                <TextField
                    autoFocus
                    margin="dense"
                    className="message"
                    label="Message *"
                    type="text"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    fullWidth
                    multiline
                    minRows={4}
                />
                <NumberField
                    margin="dense"
                    className="priority"
                    label="Priority"
                    value={priority}
                    onChange={(value) => setPriority(value)}
                    fullWidth
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={fClose}>Cancel</Button>
                <Tooltip title={submitEnabled ? '' : 'message is required'}>
                    <div>
                        <Button
                            className="send"
                            disabled={!submitEnabled}
                            onClick={submitAndClose}
                            color="primary"
                            variant="contained">
                            Send
                        </Button>
                    </div>
                </Tooltip>
            </DialogActions>
        </Dialog>
    );
};
