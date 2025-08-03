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
import AddClientDialog from './AddClientDialog';
import UpdateDialog from './UpdateClientDialog';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';
import {IClient} from '../types';
import CopyableSecret from '../common/CopyableSecret';
import {LastUsedCell} from '../common/LastUsedCell';

@observer
class Clients extends Component<Stores<'clientStore'>> {
    @observable
    private showDialog = false;
    @observable
    private deleteId: false | number = false;
    @observable
    private updateId: false | number = false;

    public componentDidMount = () => this.props.clientStore.refresh();

    public render() {
        const {
            deleteId,
            updateId,
            showDialog,
            props: {clientStore},
        } = this;
        const clients = clientStore.getItems();

        return (
            <DefaultPage
                title="Clients"
                rightControl={
                    <Button
                        id="create-client"
                        variant="contained"
                        color="primary"
                        onClick={() => (this.showDialog = true)}>
                        Create Client
                    </Button>
                }>
                <Grid size={{xs: 12}}>
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
                                        fEdit={() => (this.updateId = client.id)}
                                        fDelete={() => (this.deleteId = client.id)}
                                    />
                                ))}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
                {showDialog && (
                    <AddClientDialog
                        fClose={() => (this.showDialog = false)}
                        fOnSubmit={clientStore.create}
                    />
                )}
                {updateId !== false && (
                    <UpdateDialog
                        fClose={() => (this.updateId = false)}
                        fOnSubmit={(name) => clientStore.update(updateId, name)}
                        initialName={clientStore.getByID(updateId).name}
                    />
                )}
                {deleteId !== false && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete ' + clientStore.getByID(deleteId).name + '?'}
                        fClose={() => (this.deleteId = false)}
                        fOnSubmit={() => clientStore.remove(deleteId)}
                    />
                )}
            </DefaultPage>
        );
    }
}

interface IRowProps {
    name: string;
    value: string;
    lastUsed: string | null;
    fEdit: VoidFunction;
    fDelete: VoidFunction;
}

const Row: SFC<IRowProps> = ({name, value, lastUsed, fEdit, fDelete}) => (
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
            <IconButton onClick={fEdit} className="edit" size="large">
                <Edit />
            </IconButton>
        </TableCell>
        <TableCell align="right" padding="none">
            <IconButton onClick={fDelete} className="delete" size="large">
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

export default inject('clientStore')(Clients);
