import React, {useEffect, useState} from 'react';
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
import Security from '@mui/icons-material/Security';
import Button from '@mui/material/Button';
import Tooltip from '@mui/material/Tooltip';
import TimeAgo from 'react-timeago';
import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import AddClientDialog from './AddClientDialog';
import UpdateClientDialog from './UpdateClientDialog';
import ElevateClientDialog from './ElevateClientDialog';
import {IClient} from '../types';
import CopyableSecret from '../common/CopyableSecret';
import {LastUsedCell} from '../common/LastUsedCell';
import {TimeAgoFormatter} from '../common/TimeAgoFormatter';
import {observer} from 'mobx-react-lite';
import {useStores} from '../stores';

const Clients = observer(() => {
    const {clientStore} = useStores();
    const [toDeleteClient, setToDeleteClient] = useState<IClient>();
    const [toUpdateClient, setToUpdateClient] = useState<IClient>();
    const [toElevateClient, setToElevateClient] = useState<IClient>();
    const [createDialog, setCreateDialog] = useState<boolean>(false);
    const clients = clientStore.getItems();

    useEffect(() => void clientStore.refresh(), []);

    return (
        <DefaultPage
            title="Clients"
            maxWidth={1000}
            rightControl={
                <Button
                    id="create-client"
                    variant="contained"
                    color="primary"
                    onClick={() => setCreateDialog(true)}>
                    Create Client
                </Button>
            }>
            <Grid size={12}>
                <Paper elevation={6} style={{overflowX: 'auto'}}>
                    <Table id="client-table">
                        <TableHead>
                            <TableRow style={{textAlign: 'center'}}>
                                <TableCell>Name</TableCell>
                                <TableCell style={{width: 200}}>Token</TableCell>
                                <TableCell>Last Used</TableCell>
                                <TableCell>Elevation ends</TableCell>
                                <TableCell />
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
                                    elevatedUntil={client.elevatedUntil}
                                    fEdit={() => setToUpdateClient(client)}
                                    fDelete={() => setToDeleteClient(client)}
                                    fElevate={() => setToElevateClient(client)}
                                />
                            ))}
                        </TableBody>
                    </Table>
                </Paper>
            </Grid>
            {createDialog && (
                <AddClientDialog
                    fClose={() => setCreateDialog(false)}
                    fOnSubmit={clientStore.create}
                />
            )}
            {toUpdateClient != null && (
                <UpdateClientDialog
                    fClose={() => setToUpdateClient(undefined)}
                    fOnSubmit={(name) => clientStore.update(toUpdateClient.id, name)}
                    initialName={toUpdateClient.name}
                />
            )}
            {toDeleteClient != null && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete ' + toDeleteClient.name + '?'}
                    fClose={() => setToDeleteClient(undefined)}
                    fOnSubmit={() => clientStore.remove(toDeleteClient.id)}
                    requireElevated
                />
            )}
            {toElevateClient != null && (
                <ElevateClientDialog
                    clientName={toElevateClient.name}
                    clientId={toElevateClient.id}
                    fClose={() => setToElevateClient(undefined)}
                />
            )}
        </DefaultPage>
    );
});

interface IRowProps {
    name: string;
    value: string;
    lastUsed: string | null;
    elevatedUntil?: string;
    fEdit: VoidFunction;
    fDelete: VoidFunction;
    fElevate: VoidFunction;
}

const Row = ({name, value, lastUsed, elevatedUntil, fEdit, fDelete, fElevate}: IRowProps) => (
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
        <TableCell>
            {elevatedUntil && Date.parse(elevatedUntil) > Date.now() ? (
                <TimeAgo date={elevatedUntil} formatter={TimeAgoFormatter.longMinutes} />
            ) : (
                '-'
            )}
        </TableCell>
        <TableCell align="right" padding="none">
            <Tooltip title="Elevate">
                <IconButton onClick={fElevate} className="elevate">
                    <Security />
                </IconButton>
            </Tooltip>
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
