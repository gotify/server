import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import React, {Component} from 'react';

interface IProps {
    fClose: VoidFunction;
    fOnSubmit: (name: string, description: string, defaultPriority: number) => void;
    initialName: string;
    initialDescription: string;
    initialDefaultPriority: number;
}

interface IState {
    name: string;
    description: string;
    defaultPriority: number;
}

export default class UpdateDialog extends Component<IProps, IState> {
    public state = {name: '', description: 'some', defaultPriority: 0};

    constructor(props: IProps) {
        super(props);
        this.state = {
            name: props.initialName,
            description: props.initialDescription,
            defaultPriority: props.initialDefaultPriority,
        };
    }

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
                <DialogTitle id="form-dialog-title">Update an application</DialogTitle>
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
                    <TextField
                        margin="dense"
                        className="DefaultPriority"
                        label="Default Priority"
                        type="number"
                        value={defaultPriority}
                        onChange={this.handleChangeNumeric.bind(this, 'defaultPriority')}
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

    private handleChange(propertyName: string, event: React.ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    private handleChangeNumeric(propertyName: string, event: React.ChangeEvent<HTMLInputElement>) {
        const state = this.state;

        const regex = /^[0-9\b]+$/;
        if (event.target.value === "" || regex.test(event.target.value)) {
            state[propertyName] = event.target.value;
            this.setState(state);
        }
    }
}
