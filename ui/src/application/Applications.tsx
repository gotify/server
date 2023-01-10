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
import CloudUpload from '@material-ui/icons/CloudUpload';
import React, {ChangeEvent, Component, SFC} from 'react';
import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import Button from '@material-ui/core/Button';
import ToggleVisibility from '../common/ToggleVisibility';
import AddApplicationDialog from './AddApplicationDialog';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';
import * as config from '../config';
import UpdateDialog from './UpdateApplicationDialog';
import {IApplication} from '../types';

@observer
class Applications extends Component<Stores<'appStore'>> {
    @observable
    private deleteId: number | false = false;
    @observable
    private updateId: number | false = false;
    @observable
    private createDialog = false;

    private uploadId = -1;
    private upload: HTMLInputElement | null = null;

    public componentDidMount = () => this.props.appStore.refresh();

    public render() {
        const {
            createDialog,
            deleteId,
            updateId,
            props: {appStore},
        } = this;
        const apps = appStore.getItems();
        return (
            <DefaultPage
                title="Applications"
                rightControl={
                    <Button
                        id="create-app"
                        variant="contained"
                        color="primary"
                        onClick={() => (this.createDialog = true)}>
                        Create Application
                    </Button>
                }
                maxWidth={1000}>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table id="app-table">
                            <TableHead>
                                <TableRow>
                                    <TableCell padding="checkbox" style={{width: 80}} />
                                    <TableCell>Name</TableCell>
                                    <TableCell>Token</TableCell>
                                    <TableCell>Description</TableCell>
                                    <TableCell />
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {apps.map((app: IApplication) => (
                                    <Row
                                        key={app.id}
                                        description={app.description}
                                        image={app.image}
                                        name={app.name}
                                        value={app.token}
                                        fUpload={() => this.uploadImage(app.id)}
                                        fDelete={() => (this.deleteId = app.id)}
                                        fEdit={() => (this.updateId = app.id)}
                                        noDelete={app.internal}
                                    />
                                ))}
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
                        fClose={() => (this.createDialog = false)}
                        fOnSubmit={appStore.create}
                    />
                )}
                {updateId !== false && (
                    <UpdateDialog
                        fClose={() => (this.updateId = false)}
                        fOnSubmit={(name, description) =>
                            appStore.update(updateId, name, description)
                        }
                        initialDescription={appStore.getByID(updateId).description}
                        initialName={appStore.getByID(updateId).name}
                    />
                )}
                {deleteId !== false && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete ' + appStore.getByID(deleteId).name + '?'}
                        fClose={() => (this.deleteId = false)}
                        fOnSubmit={() => appStore.remove(deleteId)}
                    />
                )}
            </DefaultPage>
        );
    }

    private uploadImage = (id: number) => {
        this.uploadId = id;
        if (this.upload) {
            this.upload.click();
        }
    };

    private onUploadImage = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) {
            return;
        }
        if (['image/png', 'image/jpeg', 'image/gif'].indexOf(file.type) !== -1) {
            this.props.appStore.uploadImage(this.uploadId, file);
        } else {
            alert('Uploaded file must be of type png, jpeg or gif.');
        }
    };
}

interface IRowProps {
    name: string;
    value: string;
    noDelete: boolean;
    description: string;
    fUpload: VoidFunction;
    image: string;
    fDelete: VoidFunction;
    fEdit: VoidFunction;
}

const Row: SFC<IRowProps> = observer(
    ({name, value, noDelete, description, fDelete, fUpload, image, fEdit}) => (
        <TableRow>
            <TableCell padding="default">
                <div style={{display: 'flex'}}>
                    <img src={config.get('url') + image} alt="app logo" width="40" height="40" />
                    <IconButton onClick={fUpload} style={{height: 40}}>
                        <CloudUpload />
                    </IconButton>
                </div>
            </TableCell>
            <TableCell>{name}</TableCell>
            <TableCell>
                <ToggleVisibility value={value} style={{display: 'flex', alignItems: 'center'}} />
            </TableCell>
            <TableCell>{description}</TableCell>
            <TableCell align="right" padding="none">
                <IconButton onClick={fEdit} className="edit">
                    <Edit />
                </IconButton>
            </TableCell>
            <TableCell align="right" padding="none">
                <IconButton onClick={fDelete} className="delete" disabled={noDelete}>
                    <Delete />
                </IconButton>
            </TableCell>
        </TableRow>
    )
);

export default inject('appStore')(Applications);
