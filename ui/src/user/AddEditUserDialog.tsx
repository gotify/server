import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import FormControlLabel from '@mui/material/FormControlLabel';
import Switch from '@mui/material/Switch';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import React, {ChangeEvent, Component} from 'react';

interface IProps {
    name?: string;
    admin?: boolean;
    fClose: VoidFunction;
    fOnSubmit: (name: string, pass: string, admin: boolean) => void;
    isEdit?: boolean;
}

interface IState {
    name: string;
    pass: string;
    admin: boolean;
}

export default class AddEditDialog extends Component<IProps, IState> {
    public state = {
        name: this.props.name ?? '',
        pass: '',
        admin: this.props.admin ?? false,
    };

    public render() {
        const {fClose, fOnSubmit, isEdit} = this.props;
        const {name, pass, admin} = this.state;
        const namePresent = this.state.name.length !== 0;
        const passPresent = this.state.pass.length !== 0 || isEdit;
        const submitAndClose = () => {
            fOnSubmit(name, pass, admin);
            fClose();
        };
        return (
            <Dialog
                open={true}
                onClose={fClose}
                aria-labelledby="form-dialog-title"
                id="add-edit-user-dialog">
                <DialogTitle id="form-dialog-title">
                    {isEdit ? 'Edit ' + this.props.name : 'Add a user'}
                </DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        className="name"
                        label="Username *"
                        value={name}
                        name="username"
                        id="username"
                        onChange={this.handleChange.bind(this, 'name')}
                        fullWidth
                    />
                    <TextField
                        margin="dense"
                        className="password"
                        type="password"
                        value={pass}
                        fullWidth
                        label={isEdit ? 'Password (empty if no change)' : 'Password *'}
                        name="password"
                        id="password"
                        onChange={this.handleChange.bind(this, 'pass')}
                    />
                    <FormControlLabel
                        control={
                            <Switch
                                checked={admin}
                                className="admin-rights"
                                onChange={this.handleChecked.bind(this, 'admin')}
                                value="admin"
                            />
                        }
                        label="has administrator rights"
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
                                : 'username is required'
                        }>
                        <div>
                            <Button
                                className="save-create"
                                disabled={!passPresent || !namePresent}
                                onClick={submitAndClose}
                                color="primary"
                                variant="contained">
                                {isEdit ? 'Save' : 'Create'}
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }

    private handleChange(propertyName: 'name' | 'pass', event: ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    private handleChecked(propertyName: 'admin', event: ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.checked;
        this.setState(state);
    }
}
