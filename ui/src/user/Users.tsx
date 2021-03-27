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
import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import Button from '@material-ui/core/Button';
import AddEditDialog from './AddEditUserDialog';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';
import {IUser} from '../types';

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
        <TableCell align="right" padding="none">
            <IconButton onClick={fEdit} className="edit">
                <Edit />
            </IconButton>
            <IconButton onClick={fDelete} className="delete">
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

@observer
class Users extends Component<WithStyles<'wrapper'> & Stores<'userStore'>> {
    @observable
    private createDialog = false;
    @observable
    private deleteId: number | false = false;
    @observable
    private editId: number | false = false;

    public componentDidMount = () => this.props.userStore.refresh();

    public render() {
        const {
            deleteId,
            editId,
            createDialog,
            props: {userStore},
        } = this;
        const users = userStore.getItems();
        return (
            <DefaultPage
                title="Users"
                rightControl={
                    <Button
                        id="create-user"
                        variant="contained"
                        color="primary"
                        onClick={() => (this.createDialog = true)}>
                        Create User
                    </Button>
                }>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table id="user-table">
                            <TableHead>
                                <TableRow style={{textAlign: 'center'}}>
                                    <TableCell>Name</TableCell>
                                    <TableCell>Admin</TableCell>
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {users.map((user: IUser) => (
                                    <UserRow
                                        key={user.id}
                                        name={user.name}
                                        admin={user.admin}
                                        fDelete={() => (this.deleteId = user.id)}
                                        fEdit={() => (this.editId = user.id)}
                                    />
                                ))}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
                {createDialog && (
                    <AddEditDialog
                        fClose={() => (this.createDialog = false)}
                        fOnSubmit={userStore.create}
                    />
                )}
                {editId !== false && (
                    <AddEditDialog
                        fClose={() => (this.editId = false)}
                        fOnSubmit={userStore.update.bind(this, editId)}
                        name={userStore.getByID(editId).name}
                        admin={userStore.getByID(editId).admin}
                        isEdit={true}
                    />
                )}
                {deleteId !== false && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete ' + userStore.getByID(deleteId).name + '?'}
                        fClose={() => (this.deleteId = false)}
                        fOnSubmit={() => userStore.remove(deleteId)}
                    />
                )}
            </DefaultPage>
        );
    }
}

export default withStyles(styles)(inject('userStore')(Users));
