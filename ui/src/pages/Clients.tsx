import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Delete from '@material-ui/icons/Delete';
import React, {Component, SFC} from 'react';
import ConfirmDialog from '../component/ConfirmDialog';
import DefaultPage from '../component/DefaultPage';
import ToggleVisibility from '../component/ToggleVisibility';
import AddClientDialog from './dialog/AddClientDialog';
import ClientStore from '../stores/ClientStore';
import {observer} from 'mobx-react';
import {observable} from 'mobx';

@observer
class Clients extends Component {
    @observable
    private showDialog = false;
    @observable
    private deleteId: false | number = false;

    public componentDidMount = ClientStore.refresh;

    public render() {
        const {deleteId, showDialog} = this;
        const clients = ClientStore.getItems();

        return (
            <DefaultPage
                title="Clients"
                buttonTitle="Create Client"
                buttonId="create-client"
                fButton={() => (this.showDialog = true)}>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table id="client-table">
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
                                            fDelete={() => (this.deleteId = client.id)}
                                        />
                                    );
                                })}
                            </TableBody>
                        </Table>
                    </Paper>
                </Grid>
                {showDialog && (
                    <AddClientDialog
                        fClose={() => (this.showDialog = false)}
                        fOnSubmit={ClientStore.create}
                    />
                )}
                {deleteId !== false && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete ' + ClientStore.getByID(deleteId).name + '?'}
                        fClose={() => (this.deleteId = false)}
                        fOnSubmit={() => ClientStore.remove(deleteId)}
                    />
                )}
            </DefaultPage>
        );
    }
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
            <IconButton onClick={fDelete} className="delete">
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

export default Clients;
