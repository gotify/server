import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';
import {NumberField} from '../common/NumberField';
import React, {Component} from 'react';

interface IProps {
    fClose: VoidFunction;
    fOnSubmit: (name: string, description: string, defaultPriority: number) => void;
}

interface IState {
    name: string;
    description: string;
    defaultPriority: number;
}

export default class AddDialog extends Component<IProps, IState> {
    public state = {name: '', description: '', defaultPriority: 0};

    public render() {
        const {fClose, fOnSubmit} = this.props;
        const {name, description, defaultPriority} = this.state;
        const submitEnabled = this.state.name.length !== 0;
        const submitAndClose = () => {
            fOnSubmit(name, description, defaultPriority);
            fClose();
        };
        return (
            <Dialog
                open={true}
                onClose={fClose}
                aria-labelledby="form-dialog-title"
                id="app-dialog">
                <DialogTitle id="form-dialog-title">Create an application</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        An application is allowed to send messages.
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
                    <TextField
                        margin="dense"
                        className="description"
                        label="Short Description"
                        value={description}
                        onChange={this.handleChange.bind(this, 'description')}
                        fullWidth
                        multiline
                    />
                    <NumberField
                        margin="dense"
                        className="priority"
                        label="Default Priority"
                        value={defaultPriority}
                        onChange={(value) => this.setState({defaultPriority: value})}
                        fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip title={submitEnabled ? '' : 'name is required'}>
                        <div>
                            <Button
                                className="create"
                                disabled={!submitEnabled}
                                onClick={submitAndClose}
                                color="primary"
                                variant="contained">
                                Create
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }

    private handleChange(
        propertyName: 'description' | 'name',
        event: React.ChangeEvent<HTMLInputElement>
    ) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }
}
