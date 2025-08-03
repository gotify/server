import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Delete from '@mui/icons-material/Delete';
import Edit from '@mui/icons-material/Edit';
import React, {Component, SFC} from 'react';
import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import Button from '@mui/material/Button';
import AddEditDialog from './AddEditUserDialog';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';
import {IUser} from '../types';

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
            <IconButton onClick={fEdit} className="edit" size="large">
                <Edit />
            </IconButton>
            <IconButton onClick={fDelete} className="delete" size="large">
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

@observer
class Users extends Component<Stores<'userStore'>> {
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
                <Grid size={{xs: 12}}>
                    <Paper elevation={6} style={{overflowX: 'auto'}}>
                        <Table id="user-table">
                            <TableHead>
                                <TableRow style={{textAlign: 'center'}}>
                                    <TableCell>Username</TableCell>
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

export default inject('userStore')(Users);
