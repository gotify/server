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
import LoadingSpinner from '../common/LoadingSpinner.tsx';
import {useAppDispatch, useAppSelector} from '../store';
import {uiActions} from '../store/ui-slice.ts';
import {createClient, deleteClient, fetchClients, updateClient} from './client-actions.ts';
import AddClientDialog from './AddClientDialog';
import UpdateClientDialog from './UpdateClientDialog';
import {IClient} from '../types';
import CopyableSecret from '../common/CopyableSecret';
import {LastUsedCell} from '../common/LastUsedCell';

const Clients = () => {
    const dispatch = useAppDispatch();
    const clients = useAppSelector((state) => state.client.items);
    const isLoading = useAppSelector((state) => state.client.isLoading);
    const reloadRequired = useAppSelector((state) => state.ui.reloadRequired);
    const [toDeleteClient, setToDeleteClient] = useState<IClient | null>();
    const [toUpdateClient, setToUpdateClient] = useState<IClient | null>();
    const [createDialog, setCreateDialog] = useState<boolean>(false);

    // handle a requested reload
    useEffect(() => {
        if (reloadRequired) {
            dispatch(uiActions.setReloadRequired(false));
            dispatch(fetchClients());
        }
    }, [dispatch, reloadRequired]);

    useEffect(() => {
        dispatch(fetchClients());
    }, [dispatch]);

    const handleCreateClient = async (name: string) => {
        await dispatch(createClient(name));
    }

    const handleUpdateClient = async (name: string) => {
        await dispatch(updateClient(toUpdateClient!.id, name));
    }

    const handleDeleteClient = async () => {
        await dispatch(deleteClient(toDeleteClient!.id));
    }

    return (
        <DefaultPage
            title="Clients"
            rightControl={
                <Button
                    id="create-client"
                    variant="contained"
                    color="primary"
                    onClick={() => setCreateDialog(true)}>
                    Create Client
                </Button>
            }>
            {isLoading ? (
                <LoadingSpinner />
            ): (
                <Grid size={12}>
                    <Paper elevation={6} style={{overflowX: 'auto'}}>
                        <Table id="client-table">
                            <TableHead>
                                <TableRow style={{textAlign: 'center'}}>
                                    <TableCell>Name</TableCell>
                                    <TableCell style={{width: 200}}>Token</TableCell>
                                    <TableCell>Last Used</TableCell>
                                    <TableCell />
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {clients.map((client: IClient) => (
                                    <Row
                                        key={client.id}
                                        name={client.name}
                                        value={client.token}
                                        lastUsed={client.lastUsed}
                                        fEdit={() => setToUpdateClient(client)}
                                        fDelete={() => setToDeleteClient(client)}
                                    />
                                ))}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
            )}
            {createDialog && (
                <AddClientDialog
                    fClose={() => setCreateDialog(false)}
                    fOnSubmit={handleCreateClient}
                />
            )}
            {toUpdateClient != null && (
                <UpdateClientDialog
                    fClose={() => setToUpdateClient(null)}
                    fOnSubmit={handleUpdateClient}
                    initialName={toUpdateClient.name}
                />
            )}
            {toDeleteClient != null && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete ' + toDeleteClient.name + '?'}
                    fClose={() => setToDeleteClient(null)}
                    fOnSubmit={handleDeleteClient}
                />
            )}
        </DefaultPage>
    );
};

interface IRowProps {
    name: string;
    value: string;
    lastUsed: string | null;
    fEdit: VoidFunction;
    fDelete: VoidFunction;
}

const Row = ({name, value, lastUsed, fEdit, fDelete}: IRowProps) => (
    <TableRow>
        <TableCell>{name}</TableCell>
        <TableCell>
            <CopyableSecret
                value={value}
                style={{display: 'flex', alignItems: 'center', width: 250}}
            />
        </TableCell>
        <TableCell>
            <LastUsedCell lastUsed={lastUsed} />
        </TableCell>
        <TableCell align="right" padding="none">
            <IconButton onClick={fEdit} className="edit">
                <Edit />
            </IconButton>
        </TableCell>
        <TableCell align="right" padding="none">
            <IconButton onClick={fDelete} className="delete">
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

export default Clients;
