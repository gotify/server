import React, {Component} from 'react';
import Grid from 'material-ui/Grid';
import Table, {TableBody, TableCell, TableHead, TableRow} from 'material-ui/Table';
import Paper from 'material-ui/Paper';
import Button from 'material-ui/Button';
import IconButton from 'material-ui/IconButton';
import Dialog, {DialogActions, DialogContent, DialogTitle} from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import Tooltip from 'material-ui/Tooltip';
import ClientStore from '../stores/ClientStore';
import ToggleVisibility from '../component/ToggleVisibility';
import * as ClientAction from '../actions/ClientAction';
import DefaultPage from '../component/DefaultPage';
import ConfirmDialog from '../component/ConfirmDialog';
import PropTypes from 'prop-types';
import Delete from 'material-ui-icons/Delete';

class Clients extends Component {
    constructor() {
        super();
        this.state = {clients: [], showDialog: false, deleteId: -1};
    }

    componentWillMount() {
        ClientStore.on('change', this.updateClients);
        this.updateClients();
    }

    componentWillUnmount() {
        ClientStore.removeListener('change', this.updateClients);
    }

    updateClients = () => this.setState({...this.state, clients: ClientStore.get()});

    showCreateDialog = () => this.setState({...this.state, showDialog: true});
    hideCreateDialog = () => this.setState({...this.state, showDialog: false});

    showDeleteDialog = (deleteId) => this.setState({...this.state, deleteId});
    hideDeleteDelete = () => this.setState({...this.state, deleteId: -1});

    render() {
        const {clients, deleteId, showDialog} = this.state;
        return (
            <DefaultPage title="Clients" buttonTitle="Create Client" fButton={this.showCreateDialog}>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table>
                            <TableHead>
                                <TableRow style={{textAlign: 'center'}}>
                                    <TableCell>Name</TableCell>
                                    <TableCell style={{width: 200}}>token</TableCell>
                                    <TableCell/>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {clients.map((client) => {
                                    return (
                                        <Row key={client.id} name={client.name}
                                             value={client.token} fDelete={() => this.showDeleteDialog(client.id)}/>
                                    );
                                })}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
                {showDialog && <AddDialog fClose={this.hideCreateDialog} fOnSubmit={ClientAction.createClient}/>}
                {deleteId !== -1 && <ConfirmDialog title="Confirm Delete"
                                                   text={'Delete ' + ClientStore.getById(this.state.deleteId).name + '?'}
                                                   fClose={this.hideDeleteDelete}
                                                   fOnSubmit={() => ClientAction.deleteClient(deleteId)}/>}
            </DefaultPage>
        );
    }
}

class Row extends Component {
    static propTypes = {
        name: PropTypes.string.isRequired,
        value: PropTypes.string.isRequired,
        fDelete: PropTypes.func.isRequired,
    };

    render() {
        const {name, value, fDelete} = this.props;
        return (
            <TableRow>
                <TableCell>{name}</TableCell>
                <TableCell>
                    <ToggleVisibility value={value} style={{display: 'flex', alignItems: 'center', width: 200}}/>
                </TableCell>
                <TableCell numeric padding="none">
                    <IconButton onClick={fDelete}><Delete/></IconButton>
                </TableCell>
            </TableRow>
        );
    }
}

class AddDialog extends Component {
    static propTypes = {
        fClose: PropTypes.func.isRequired,
        fOnSubmit: PropTypes.func.isRequired,
    };

    constructor() {
        super();
        this.state = {name: ''};
    }

    handleChange(propertyName, event) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    render() {
        const {fClose, fOnSubmit} = this.props;
        const {name} = this.state;
        const submitEnabled = this.state.name.length !== 0;
        const submitAndClose = () => {
            fOnSubmit(name);
            fClose();
        };
        return (
            <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Create a client</DialogTitle>
                <DialogContent>
                    <TextField autoFocus margin="dense" id="name" label="Name *" type="email" value={name}
                               onChange={this.handleChange.bind(this, 'name')} fullWidth/>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip placement={'bottom-start'} title={submitEnabled ? '' : 'name is required'}>
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
}

export default Clients;
