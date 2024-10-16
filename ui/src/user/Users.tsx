import React, {useEffect, useState} from 'react';
import Grid from '@mui/material/Grid2';
import IconButton from '@mui/material/IconButton';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Delete from '@mui/icons-material/Delete';
import Edit from '@mui/icons-material/Edit';
import Button from '@mui/material/Button';
import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import {useAppDispatch, useAppSelector} from '../store';
import {createUser, deleteUser, fetchUsers, updateUser} from '../store/user-actions.ts';
import AddEditUserDialog from './AddEditUserDialog';
import {IUser} from '../types';

// const useStyles = makeStyles()(() => {
//     return {
//         wrapper: {
//             margin: '0 auto',
//             maxWidth: 700,
//         },
//     };
// });

interface IRowProps {
    name: string;
    admin: boolean;
    fDelete: VoidFunction;
    fEdit: VoidFunction;
}

const UserRow = ({name, admin, fDelete, fEdit}: IRowProps) => (
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

const Users = () => {
    const dispatch = useAppDispatch();

    const users = useAppSelector((state) => state.user.items);
    const [toDeleteUser, setToDeleteUser] = useState<IUser | null>();
    const [toUpdateUser, setToUpdateUser] = useState<IUser | null>();
    const [createDialog, setCreateDialog] = useState<boolean>(false);

    useEffect(() => {
        dispatch(fetchUsers());
    }, [dispatch]);

    const handleCreateUser = async (name: string, pass: string | null, admin: boolean) => {
        await dispatch(createUser(name, pass, admin))
    }

    const handleUpdateUser = async (name: string, pass: string | null, admin: boolean) => {
        await dispatch(updateUser(toUpdateUser!.id, name, pass, admin));
    };

    const handleDeleteUser = async() => {
        await dispatch(deleteUser(toDeleteUser!.id))
    }


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
            <Grid size={12}>
                    <Paper elevation={6} style={{overflowX: 'auto'}}>
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
                                    fDelete={() => setToDeleteUser(user)}
                                    fEdit={() => setToUpdateUser(user)}
                                />
                            ))}
                        </TableBody>
                    </Table>
                </Paper>
            </Grid>
            {createDialog && (
                <AddEditUserDialog
                    fClose={() => setCreateDialog(false)}
                    fOnSubmit={handleCreateUser}
                />
            )}
            {toUpdateUser != null && (
                <AddEditUserDialog
                    fClose={() => setToUpdateUser(null)}
                    fOnSubmit={handleUpdateUser}
                    name={toUpdateUser?.name}
                    admin={toUpdateUser?.admin}
                    isEdit={true}
                />
            )}
            {toDeleteUser != null && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete ' + toDeleteUser.name + '?'}
                    fClose={() => setToDeleteUser(null)}
                    fOnSubmit={handleDeleteUser}
                />
            )}
        </DefaultPage>
    );
}

export default Users;
