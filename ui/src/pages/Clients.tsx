import Delete from 'material-ui-icons/Delete';
import Grid from 'material-ui/Grid';
import IconButton from 'material-ui/IconButton';
import Paper from 'material-ui/Paper';
import Table, {TableBody, TableCell, TableHead, TableRow} from 'material-ui/Table';
import React, {Component, SFC} from 'react';
import * as ClientAction from '../actions/ClientAction';
import ConfirmDialog from '../component/ConfirmDialog';
import DefaultPage from '../component/DefaultPage';
import ToggleVisibility from '../component/ToggleVisibility';
import ClientStore from '../stores/ClientStore';
import AddClientDialog from './dialog/AddClientDialog';

interface IState {
    clients: IClient[];
    showDialog: boolean;
    deleteId: number;
}

class Clients extends Component<{}, IState> {
    public state = {clients: [], showDialog: false, deleteId: -1};

    public componentWillMount() {
        ClientStore.on('change', this.updateClients);
        this.updateClients();
    }

    public componentWillUnmount() {
        ClientStore.removeListener('change', this.updateClients);
    }

    public render() {
        const {clients, deleteId, showDialog} = this.state;
        return (
            <DefaultPage
                title="Clients"
                buttonTitle="Create Client"
                fButton={this.showCreateDialog}>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table>
                            <TableHead>
                                <TableRow style={{textAlign: 'center'}}>
                                    <TableCell>Name</TableCell>
                                    <TableCell style={{width: 200}}>token</TableCell>
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {clients.map((client: IClient) => {
                                    return (
                                        <Row
                                            key={client.id}
                                            name={client.name}
                                            value={client.token}
                                            fDelete={() => this.showDeleteDialog(client.id)}
                                        />
                                    );
                                })}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
                {showDialog && (
                    <AddClientDialog
                        fClose={this.hideCreateDialog}
                        fOnSubmit={ClientAction.createClient}
                    />
                )}
                {deleteId !== -1 && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete ' + ClientStore.getById(this.state.deleteId).name + '?'}
                        fClose={this.hideDeleteDelete}
                        fOnSubmit={this.deleteClient}
                    />
                )}
            </DefaultPage>
        );
    }

    private deleteClient = () => ClientAction.deleteClient(this.state.deleteId);

    private updateClients = () => this.setState({...this.state, clients: ClientStore.get()});
    private showCreateDialog = () => this.setState({...this.state, showDialog: true});

    private hideCreateDialog = () => this.setState({...this.state, showDialog: false});
    private showDeleteDialog = (deleteId: number) => this.setState({...this.state, deleteId});

    private hideDeleteDelete = () => this.setState({...this.state, deleteId: -1});
}

interface IRowProps {
    name: string;
    value: string;
    fDelete: VoidFunction;
}

const Row: SFC<IRowProps> = ({name, value, fDelete}) => (
    <TableRow>
        <TableCell>{name}</TableCell>
        <TableCell>
            <ToggleVisibility
                value={value}
                style={{display: 'flex', alignItems: 'center', width: 200}}
            />
        </TableCell>
        <TableCell numeric padding="none">
            <IconButton onClick={fDelete}>
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

export default Clients;
