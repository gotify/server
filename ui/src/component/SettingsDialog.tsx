import Button from 'material-ui/Button';
import Dialog, {DialogActions, DialogContent, DialogTitle} from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import Tooltip from 'material-ui/Tooltip';
import React, {ChangeEvent, Component} from 'react';
import * as UserAction from '../actions/UserAction';

interface IState {
    pass: string
}

interface IProps {
    fClose: VoidFunction
}

export default class SettingsDialog extends Component<IProps, IState> {
    public state = {pass: ''};

    public render() {
        const {pass} = this.state;
        const {fClose} = this.props;
        const submitAndClose = () => {
            UserAction.changeCurrentUser(pass);
            fClose();
        };
        return (
            <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Change Password</DialogTitle>
                <DialogContent>
                    <TextField autoFocus margin="dense" type="password" label="New Pass *" value={pass}
                               onChange={this.handleChange.bind(this, 'pass')} fullWidth/>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip title={pass.length !== 0 ? '' : 'pass is required'}>
                        <div>
                            <Button disabled={pass.length === 0} onClick={submitAndClose} color="primary"
                                    variant="raised">
                                Change
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }

    private handleChange(propertyName: string, event: ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }
}
