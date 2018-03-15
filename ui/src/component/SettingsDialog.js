import React, {Component} from 'react';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import Tooltip from 'material-ui/Tooltip';
import Dialog, {DialogActions, DialogContent, DialogTitle} from 'material-ui/Dialog';
import PropTypes from 'prop-types';
import * as UserAction from '../actions/UserAction';

export default class SettingsDialog extends Component {
    static propTypes = {
        fClose: PropTypes.func.isRequired,
    };

    constructor() {
        super();
        this.state = {pass: ''};
    }

    handleChange(propertyName, event) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    render() {
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
                    <Tooltip title={pass.length === 0 ? '' : 'pass is required'}>
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
}
