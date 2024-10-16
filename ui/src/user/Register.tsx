import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import React, {useState} from 'react';

interface IProps {
    name?: string;
    fClose: VoidFunction;
    fOnSubmit: (name: string, pass: string) => Promise<boolean>;
}

const  RegistrationDialog = ({fClose, fOnSubmit, name: initialName = ''}: IProps) => {
    const [name, setName] = useState(initialName);
    const [pass, setPass] = useState('');
    const namePresent = name.length !== 0;
    const passPresent = pass.length !== 0;

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value);
    }

    const handlePassChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPass(e.target.value);
    }

    const submitAndClose = (): void => {
        fOnSubmit(name, pass).then((success) => {
            if (success) {
                fClose();
            }
        });
    };

    return (
        <Dialog
            open={true}
            onClose={fClose}
            aria-labelledby="form-dialog-title"
            id="add-edit-user-dialog">
            <DialogTitle id="form-dialog-title">Registration</DialogTitle>
            <DialogContent>
                <TextField
                    autoFocus
                    margin="dense"
                    className="name"
                    label="Name *"
                    type="email"
                    value={name}
                    onChange={handleNameChange}
                    fullWidth
                />
                <TextField
                    margin="dense"
                    className="password"
                    type="password"
                    value={pass}
                    fullWidth
                    label="Pass *"
                    onChange={handlePassChange}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={fClose}>Cancel</Button>
                <Tooltip
                    placement={'bottom-start'}
                    title={
                        namePresent
                            ? passPresent
                                ? ''
                                : 'password is required'
                            : 'name is required'
                    }>
                    <div>
                        <Button
                            className="save-create"
                            disabled={!passPresent || !namePresent}
                            onClick={submitAndClose}
                            color="primary"
                            variant="contained">
                            Register
                        </Button>
                    </div>
                </Tooltip>
            </DialogActions>
        </Dialog>
    );
}

export default RegistrationDialog;
