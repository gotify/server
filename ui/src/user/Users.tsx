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
import React from 'react';
import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import Button from '@mui/material/Button';
import AddEditDialog from './AddEditUserDialog';
import {IUser} from '../types';
import {useStores} from '../stores';
import {observer} from 'mobx-react-lite';

interface IRowProps {
    name: string;
    admin: boolean;
    fDelete: VoidFunction;
    fEdit: VoidFunction;
}

const UserRow: React.FC<IRowProps> = ({name, admin, fDelete, fEdit}) => (
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

const Users = observer(() => {
    const [deleteUser, setDeleteUser] = React.useState<IUser>();
    const [editUser, setEditUser] = React.useState<IUser>();
    const [createDialog, setCreateDialog] = React.useState(false);
    const {userStore} = useStores();
    React.useEffect(() => void userStore.refresh(), []);
    const users = userStore.getItems();
    return (
        <DefaultPage
            title="Users"
            rightControl={
                <Button
                    id="create-user"
                    variant="contained"
                    color="primary"
                    onClick={() => setCreateDialog(true)}>
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
                                    fDelete={() => setDeleteUser(user)}
                                    fEdit={() => setEditUser(user)}
                                />
                            ))}
                        </TableBody>
                    </Table>
                </Paper>
            </Grid>
            {createDialog && (
                <AddEditDialog fClose={() => setCreateDialog(false)} fOnSubmit={userStore.create} />
            )}
            {editUser && (
                <AddEditDialog
                    fClose={() => setEditUser(undefined)}
                    fOnSubmit={userStore.update.bind(this, editUser.id)}
                    name={editUser.name}
                    admin={editUser.admin}
                    isEdit={true}
                />
            )}
            {deleteUser && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete ' + deleteUser.name + '?'}
                    fClose={() => setDeleteUser(undefined)}
                    fOnSubmit={() => userStore.remove(deleteUser.id)}
                />
            )}
        </DefaultPage>
    );
});

export default Users;
