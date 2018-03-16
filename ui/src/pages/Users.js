import React, {Component} from 'react';
import Grid from 'material-ui/Grid';
import {withStyles} from 'material-ui/styles';
import Table, {TableBody, TableCell, TableHead, TableRow} from 'material-ui/Table';
import Paper from 'material-ui/Paper';
import Switch from 'material-ui/Switch';
import Button from 'material-ui/Button';
import IconButton from 'material-ui/IconButton';
import Dialog, {DialogActions, DialogContent, DialogTitle} from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import Tooltip from 'material-ui/Tooltip';
import UserStore from '../stores/UserStore';
import * as UserAction from '../actions/UserAction';
import {FormControlLabel} from 'material-ui/Form';
import ConfirmDialog from '../component/ConfirmDialog';
import DefaultPage from '../component/DefaultPage';
import Delete from 'material-ui-icons/Delete';
import Edit from 'material-ui-icons/Edit';
import PropTypes from 'prop-types';

const styles = () => ({
    wrapper: {
        margin: '0 auto',
        maxWidth: 700,
    },
});

class UserRow extends Component {
    static propTypes = {
        name: PropTypes.string.isRequired,
        admin: PropTypes.bool.isRequired,
        fDelete: PropTypes.func.isRequired,
        fEdit: PropTypes.func.isRequired,
    };

    render() {
        const {name, admin, fDelete, fEdit} = this.props;
        return (
            <TableRow>
                <TableCell>{name}</TableCell>
                <TableCell>{admin ? 'Yes' : 'No'}</TableCell>
                <TableCell numeric padding="none">
                    <IconButton onClick={fEdit}><Edit/></IconButton>
                    <IconButton onClick={fDelete}><Delete/></IconButton>
                </TableCell>
            </TableRow>
        );
    }
}

class Users extends Component {
    constructor() {
        super();
        this.state = {users: [], createDialog: false, deleteId: -1, editId: -1};
    }

    componentWillMount() {
        UserStore.on('change', this.updateUsers);
        this.updateUsers();
    }

    componentWillUnmount() {
        UserStore.removeListener('change', this.updateUsers);
    }

    updateUsers = () => this.setState({...this.state, users: UserStore.get()});

    showCreateDialog = () => this.setState({...this.state, createDialog: true});
    hideCreateDialog = () => this.setState({...this.state, createDialog: false});

    showEditDialog = (editId) => this.setState({...this.state, editId});
    hideEditDialog = () => this.setState({...this.state, editId: -1});

    showDeleteDialog = (deleteId) => this.setState({...this.state, deleteId});
    hideDeleteDialog = () => this.setState({...this.state, deleteId: -1});

    render() {
        const {users, deleteId, editId} = this.state;
        return (
            <DefaultPage title="Users" buttonTitle="Create User" fButton={this.showCreateDialog}>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table>
                            <TableHead>
                                <TableRow style={{textAlign: 'center'}}>
                                    <TableCell>Name</TableCell>
                                    <TableCell>Admin</TableCell>
                                    <TableCell/>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {users.map((user) => {
                                    return (
                                        <UserRow key={user.id} name={user.name} admin={user.admin}
                                                 fDelete={() => this.showDeleteDialog(user.id)}
                                                 fEdit={() => this.showEditDialog(user.id)}/>
                                    );
                                })}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
                {this.state.createDialog && <AddEditDialog fClose={this.hideCreateDialog}
                                                           fOnSubmit={UserAction.createUser}/>}
                {editId !== -1 && <AddEditDialog fClose={this.hideEditDialog}
                                                 fOnSubmit={UserAction.updateUser.bind(this, editId)}
                                                 name={UserStore.getById(this.state.editId).name}
                                                 admin={UserStore.getById(this.state.editId).admin}
                                                 isEdit={true}/>}
                {deleteId !== -1 && <ConfirmDialog title="Confirm Delete"
                                                   text={'Delete ' + UserStore.getById(this.state.deleteId).name + '?'}
                                                   fClose={this.hideDeleteDialog}
                                                   fOnSubmit={() => UserAction.deleteUser(this.state.deleteId)}
                />}
            </DefaultPage>
        );
    }
}

class AddEditDialog extends Component {
    static defaultProps = {
        name: '',
        pass: '',
        admin: false,
        isEdit: false,
    };

    static propTypes = {
        name: PropTypes.string.isRequired,
        admin: PropTypes.bool.isRequired,
        fClose: PropTypes.func.isRequired,
        fOnSubmit: PropTypes.func.isRequired,
        isEdit: PropTypes.bool.isRequired,
    };

    constructor() {
        super();
        this.state = {
            name: '',
            pass: '',
            admin: false,
        };
    }

    componentWillMount() {
        const {name, admin} = this.props;
        this.setState({...this.state, name, admin: admin});
    }


    handleChange(propertyName, event) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    handleChecked(propertyName, event) {
        const state = this.state;
        state[propertyName] = event.target.checked;
        this.setState(state);
    }

    render() {
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
}

export default withStyles(styles)(Users);
