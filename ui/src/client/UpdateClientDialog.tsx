import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import React, {Component} from 'react';

interface IProps {
    fClose: VoidFunction;
    fOnSubmit: (name: string) => void;
    initialName: string;
}

interface IState {
    name: string;
}

export default class UpdateDialog extends Component<IProps, IState> {
    public state = {name: ''};

    constructor(props: IProps) {
        super(props);
        this.state = {
            name: props.initialName,
        };
    }

    public render() {
        const {fClose, fOnSubmit} = this.props;
        const {name} = this.state;
        const submitEnabled = this.state.name.length !== 0;
        const submitAndClose = () => {
            fOnSubmit(name);
            fClose();
        };
        return (
            <Dialog
                open={true}
                onClose={fClose}
                aria-labelledby="form-dialog-title"
                id="client-dialog">
                <DialogTitle id="form-dialog-title">Update a Client</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        A client manages messages, clients, applications and users (with admin
                        permissions).
                    </DialogContentText>
                    <TextField
                        autoFocus
                        margin="dense"
                        className="name"
                        label="Name *"
                        type="text"
                        value={name}
                        onChange={this.handleChange.bind(this, 'name')}
                        fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip title={submitEnabled ? '' : 'name is required'}>
                        <div>
                            <Button
                                className="update"
                                disabled={!submitEnabled}
                                onClick={submitAndClose}
                                color="primary"
                                variant="contained">
                                Update
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }

    private handleChange(propertyName: 'name', event: React.ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }
}
