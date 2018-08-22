import Delete from 'material-ui-icons/Delete';
import Edit from 'material-ui-icons/Edit';
import Avatar from 'material-ui/Avatar';
import Grid from 'material-ui/Grid';
import IconButton from 'material-ui/IconButton';
import Paper from 'material-ui/Paper';
import Table, {TableBody, TableCell, TableHead, TableRow} from 'material-ui/Table';
import React, {ChangeEvent, Component, SFC} from 'react';
import * as AppAction from '../actions/AppAction';
import ConfirmDialog from '../component/ConfirmDialog';
import DefaultPage from '../component/DefaultPage';
import ToggleVisibility from '../component/ToggleVisibility';
import AppStore from '../stores/AppStore';
import AddApplicationDialog from './dialog/AddApplicationDialog';

interface IState {
    apps: IApplication[];
    createDialog: boolean;
    deleteId: number;
}

class Applications extends Component<{}, IState> {
    public state = {apps: [], createDialog: false, deleteId: -1};
    private uploadId: number = -1;
    private upload: HTMLInputElement | null = null;

    public componentWillMount() {
        AppStore.on('change', this.updateApps);
        this.updateApps();
    }

    public componentWillUnmount() {
        AppStore.removeListener('change', this.updateApps);
    }

    public render() {
        const {apps, createDialog, deleteId} = this.state;
        return (
            <DefaultPage
                title="Applications"
                buttonTitle="Create Application"
                maxWidth={1000}
                fButton={this.showCreateDialog}>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell padding="checkbox" style={{width: 80}} />
                                    <TableCell>Name</TableCell>
                                    <TableCell>Token</TableCell>
                                    <TableCell>Description</TableCell>
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {apps.map((app: IApplication) => {
                                    return (
                                        <Row
                                            key={app.id}
                                            description={app.description}
                                            image={app.image}
                                            name={app.name}
                                            value={app.token}
                                            fUpload={() => this.uploadImage(app.id)}
                                            fDelete={() => this.showCloseDialog(app.id)}
                                        />
                                    );
                                })}
                            </TableBody>
                        </Table>
                        <input
                            ref={(upload) => (this.upload = upload)}
                            type="file"
                            style={{display: 'none'}}
                            onChange={this.onUploadImage}
                        />
                    </Paper>
                </Grid>
                {createDialog && (
                    <AddApplicationDialog
                        fClose={this.hideCreateDialog}
                        fOnSubmit={AppAction.createApp}
                    />
                )}
                {deleteId !== -1 && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete ' + AppStore.getById(deleteId).name + '?'}
                        fClose={this.hideCloseDialog}
                        fOnSubmit={() => AppAction.deleteApp(deleteId)}
                    />
                )}
            </DefaultPage>
        );
    }

    private updateApps = () => this.setState({...this.state, apps: AppStore.get()});
    private uploadImage = (id: number) => {
        this.uploadId = id;
        if (this.upload) {
            this.upload.click();
        }
    };

    private onUploadImage = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files && e.target.files[0];
        if (!file) {
            return;
        }
        if (['image/png', 'image/jpeg', 'image/gif'].indexOf(file.type) !== -1) {
            AppAction.uploadImage(this.uploadId, file);
        } else {
            alert('Uploaded file must be of type png, jpeg or gif.');
        }
    };
    private showCreateDialog = () => this.setState({...this.state, createDialog: true});

    private hideCreateDialog = () => this.setState({...this.state, createDialog: false});
    private showCloseDialog = (deleteId: number) => this.setState({...this.state, deleteId});

    private hideCloseDialog = () => this.setState({...this.state, deleteId: -1});
}

interface IRowProps {
    name: string;
    value: string;
    description: string;
    fUpload: VoidFunction;
    image: string;
    fDelete: VoidFunction;
}

const Row: SFC<IRowProps> = ({name, value, description, fDelete, fUpload, image}) => (
    <TableRow>
        <TableCell padding="checkbox">
            <div style={{display: 'flex'}}>
                <Avatar src={image} />
                <IconButton onClick={fUpload} style={{height: 40}}>
                    <Edit />
                </IconButton>
            </div>
        </TableCell>
        <TableCell>{name}</TableCell>
        <TableCell>
            <ToggleVisibility value={value} style={{display: 'flex', alignItems: 'center'}} />
        </TableCell>
        <TableCell>{description}</TableCell>
        <TableCell numeric padding="none">
            <IconButton onClick={fDelete}>
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
);

export default Applications;
