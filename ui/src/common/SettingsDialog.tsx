import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import React, {Component} from 'react';
import { action, makeObservable, observable } from 'mobx';
import {observer} from 'mobx-react';
import {inject, Stores} from '../inject';

interface IProps {
    fClose: VoidFunction;
}

@observer
class SettingsDialog extends Component<IProps & Stores<'currentUser'>> {
    @observable
    private pass = '';

    constructor(props: any) {
        super(props);
        makeObservable(this);
    }

    @action
    private setPass = (value: string) => {
        this.pass = value;
    }

    public render() {
        this.setPass(this.pass);
        const {fClose, currentUser} = this.props;
        const submitAndClose = () => {
            currentUser.changePassword(this.pass);
            fClose();
        };
        return (
            <Dialog
                open={true}
                onClose={fClose}
                aria-labelledby="form-dialog-title"
                id="changepw-dialog">
                <DialogTitle id="form-dialog-title">Change Password</DialogTitle>
                <DialogContent>
                    <TextField
                        className="newpass"
                        autoFocus
                        margin="dense"
                        type="password"
                        label="New Password *"
                        value={this.pass}
                        onChange={(e) => (this.setPass(e.target.value))}
                        fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip title={this.pass.length !== 0 ? '' : 'Password is required'}>
                        <div>
                            <Button
                                className="change"
                                disabled={this.pass.length === 0}
                                onClick={submitAndClose}
                                color="primary"
                                variant="contained">
                                Change
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }
}

export default inject('currentUser')(SettingsDialog);
