import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
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
import ToggleVisibility from '../common/ToggleVisibility';
import AddClientDialog from './AddClientDialog';
import UpdateDialog from './UpdateClientDialog';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';
import {IClient} from '../types';

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
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table id="client-table">
                            <TableHead>
                                <TableRow style={{textAlign: 'center'}}>
                                    <TableCell>Name</TableCell>
                                    <TableCell style={{width: 200}}>Token</TableCell>
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
    fEdit: VoidFunction;
    fDelete: VoidFunction;
}

const Row: SFC<IRowProps> = ({name, value, fEdit, fDelete}) => (
    <TableRow>
        <TableCell>{name}</TableCell>
        <TableCell>
            <ToggleVisibility
                value={value}
                style={{display: 'flex', alignItems: 'center', width: 200}}
            />
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

export default inject('clientStore')(Clients);
