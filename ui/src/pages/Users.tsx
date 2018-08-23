import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import {withStyles, WithStyles} from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Delete from '@material-ui/icons/Delete';
import Edit from '@material-ui/icons/Edit';
import React, {Component, SFC} from 'react';
import * as UserAction from '../actions/UserAction';
import ConfirmDialog from '../component/ConfirmDialog';
import DefaultPage from '../component/DefaultPage';
import UserStore from '../stores/UserStore';
import AddEditDialog from './dialog/AddEditUserDialog';

const styles = () => ({
    wrapper: {
        margin: '0 auto',
        maxWidth: 700,
    },
});

interface IRowProps {
    name: string;
    admin: boolean;
    fDelete: VoidFunction;
    fEdit: VoidFunction;
}

const UserRow: SFC<IRowProps> = ({name, admin, fDelete, fEdit}) => (
    <TableRow>
        <TableCell>{name}</TableCell>
        <TableCell>{admin ? 'Yes' : 'No'}</TableCell>
        <TableCell numeric padding="none">
            <IconButton onClick={fEdit}>
                <Edit />
            </IconButton>
            <IconButton onClick={fDelete}>
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

interface IState {
    users: IUser[];
    createDialog: boolean;
    deleteId: number;
    editId: number;
}

class Users extends Component<WithStyles<'wrapper'>, IState> {
    public state = {users: [], createDialog: false, deleteId: -1, editId: -1};

    public componentWillMount() {
        UserStore.on('change', this.updateUsers);
        this.updateUsers();
    }

    public componentWillUnmount() {
        UserStore.removeListener('change', this.updateUsers);
    }

    public render() {
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
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {users.map((user: IUser) => {
                                    return (
                                        <UserRow
                                            key={user.id}
                                            name={user.name}
                                            admin={user.admin}
                                            fDelete={() => this.showDeleteDialog(user.id)}
                                            fEdit={() => this.showEditDialog(user.id)}
                                        />
                                    );
                                })}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
                {this.state.createDialog && (
                    <AddEditDialog
                        fClose={this.hideCreateDialog}
                        fOnSubmit={UserAction.createUser}
                    />
                )}
                {editId !== -1 && (
                    <AddEditDialog
                        fClose={this.hideEditDialog}
                        fOnSubmit={UserAction.updateUser.bind(this, editId)}
                        name={UserStore.getById(this.state.editId).name}
                        admin={UserStore.getById(this.state.editId).admin}
                        isEdit={true}
                    />
                )}
                {deleteId !== -1 && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete ' + UserStore.getById(this.state.deleteId).name + '?'}
                        fClose={this.hideDeleteDialog}
                        fOnSubmit={() => UserAction.deleteUser(this.state.deleteId)}
                    />
                )}
            </DefaultPage>
        );
    }

    private updateUsers = () => this.setState({...this.state, users: UserStore.get()});
    private showCreateDialog = () => this.setState({...this.state, createDialog: true});

    private hideCreateDialog = () => this.setState({...this.state, createDialog: false});
    private showEditDialog = (editId: number) => this.setState({...this.state, editId});

    private hideEditDialog = () => this.setState({...this.state, editId: -1});
    private showDeleteDialog = (deleteId: number) => this.setState({...this.state, deleteId});

    private hideDeleteDialog = () => this.setState({...this.state, deleteId: -1});
}

export default withStyles(styles)(Users);
