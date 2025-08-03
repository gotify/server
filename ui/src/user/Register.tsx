import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import React, {ChangeEvent, Component} from 'react';

interface IProps {
    name?: string;
    fClose: VoidFunction;
    fOnSubmit: (name: string, pass: string) => Promise<boolean>;
}

interface IState {
    name: string;
    pass: string;
}

export default class RegistrationDialog extends Component<IProps, IState> {
    public state = {
        name: '',
        pass: '',
    };

    public render() {
        const {fClose, fOnSubmit} = this.props;
        const {name, pass} = this.state;
        const namePresent = this.state.name.length !== 0;
        const passPresent = this.state.pass.length !== 0;
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
                        id="register-username"
                        margin="dense"
                        className="name"
                        label="Username *"
                        name="username"
                        value={name}
                        autoComplete="username"
                        onChange={this.handleChange.bind(this, 'name')}
                        fullWidth
                    />
                    <TextField
                        id="register-password"
                        margin="dense"
                        className="password"
                        type="password"
                        value={pass}
                        fullWidth
                        label="Password *"
                        name="password"
                        autoComplete="new-password"
                        onChange={this.handleChange.bind(this, 'pass')}
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
                                Register
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }

    private handleChange(propertyName: keyof IState, event: ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }
}
