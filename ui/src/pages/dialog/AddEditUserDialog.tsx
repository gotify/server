import Button from 'material-ui/Button';
import Dialog, {DialogActions, DialogContent, DialogTitle} from 'material-ui/Dialog';
import {FormControlLabel} from 'material-ui/Form';
import Switch from 'material-ui/Switch';
import TextField from 'material-ui/TextField';
import Tooltip from 'material-ui/Tooltip';
import React, {ChangeEvent, Component} from 'react';

interface IProps {
    name?: string
    admin?: boolean
    fClose: VoidFunction
    fOnSubmit: (name: string, pass: string, admin: boolean) => void
    isEdit?: boolean
}

interface IState {
    name: string
    pass: string
    admin: boolean
}

export default class AddEditDialog extends Component<IProps, IState> {
    public state = {
        name: '',
        pass: '',
        admin: false,
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
            <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">{isEdit ? 'Edit ' + this.props.name : 'Add a user'}</DialogTitle>
                <DialogContent>
                    <TextField autoFocus margin="dense" id="name" label="Name *" type="email" value={name}
                               onChange={this.handleChange.bind(this, 'name')} fullWidth/>
                    <TextField margin="dense" id="description" type="password" value={pass} fullWidth
                               label={isEdit ? 'Pass (empty if no change)' : 'Pass *'}
                               onChange={this.handleChange.bind(this, 'pass')}/>
                    <FormControlLabel
                        control={<Switch checked={admin} onChange={this.handleChecked.bind(this, 'admin')}
                                         value="admin"/>} label="has administrator rights"/>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip placement={'bottom-start'}
                             title={namePresent ? (passPresent ? '' : 'password is required') : 'name is required'}>
                        <div>
                            <Button disabled={!passPresent || !namePresent} onClick={submitAndClose}
                                    color="primary" variant="raised">
                                {isEdit ? 'Save' : 'Create'}
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }

    public componentWillMount() {
        const {name, admin} = this.props;
        this.setState({...this.state, name: name || '', admin: admin || false});
    }

    private handleChange(propertyName: string, event: ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    private handleChecked(propertyName: string, event: ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.checked;
        this.setState(state);
    }
}