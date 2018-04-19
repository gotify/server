import Button from 'material-ui/Button';
import Dialog, {DialogActions, DialogContent, DialogContentText, DialogTitle} from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import Tooltip from 'material-ui/Tooltip';
import React, {Component} from 'react';

interface IProps {
    fClose: VoidFunction
    fOnSubmit: (name: string, description: string) => void
}

interface IState {
    name: string
    description: string
}

export default class AddDialog extends Component<IProps, IState> {
    public state = {name: '', description: ''};

    public render() {
        const {fClose, fOnSubmit} = this.props;
        const {name, description} = this.state;
        const submitEnabled = this.state.name.length !== 0;
        const submitAndClose = () => {
            fOnSubmit(name, description);
            fClose();
        };
        return (
            <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Create an application</DialogTitle>
                <DialogContent>
                    <DialogContentText>An application is allowed to send messages.</DialogContentText>
                    <TextField autoFocus margin="dense" id="name" label="Name *" type="email" value={name}
                               onChange={this.handleChange.bind(this, 'name')} fullWidth/>
                    <TextField margin="dense" id="description" label="Short Description" value={description}
                               onChange={this.handleChange.bind(this, 'description')} fullWidth multiline/>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip title={submitEnabled ? '' : 'name is required'}>
                        <div>
                            <Button disabled={!submitEnabled} onClick={submitAndClose} color="primary" variant="raised">
                                Create
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }

    private handleChange(propertyName: string, event: React.ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }
}