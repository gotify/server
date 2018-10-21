import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import React, {Component} from 'react';
import {currentUser} from '../stores/CurrentUser';
import {observable} from 'mobx';
import {observer} from 'mobx-react';

interface IProps {
    fClose: VoidFunction;
}

@observer
export default class SettingsDialog extends Component<IProps> {
    @observable
    private pass = '';

    public render() {
        const {pass} = this;
        const {fClose} = this.props;
        const submitAndClose = () => {
            currentUser.changePassword(pass);
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
                        label="New Pass *"
                        value={pass}
                        onChange={(e) => (this.pass = e.target.value)}
                        fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip title={pass.length !== 0 ? '' : 'pass is required'}>
                        <div>
                            <Button
                                className="change"
                                disabled={pass.length === 0}
                                onClick={submitAndClose}
                                color="primary"
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
