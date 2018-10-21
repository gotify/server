import Avatar from '@material-ui/core/Avatar';
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
import React, {ChangeEvent, Component, SFC} from 'react';
import ConfirmDialog from '../component/ConfirmDialog';
import DefaultPage from '../component/DefaultPage';
import ToggleVisibility from '../component/ToggleVisibility';
import AddApplicationDialog from './AddApplicationDialog';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';

@observer
class Applications extends Component<Stores<'appStore'>> {
    @observable
    private deleteId: number | false = false;
    @observable
    private createDialog = false;

    private uploadId = -1;
    private upload: HTMLInputElement | null = null;

    public componentDidMount = () => this.props.appStore.refresh();

    public render() {
        const {
            createDialog,
            deleteId,
            props: {appStore},
        } = this;
        const apps = appStore.getItems();
        return (
            <DefaultPage
                title="Applications"
                buttonTitle="Create Application"
                buttonId="create-app"
                maxWidth={1000}
                fButton={() => (this.createDialog = true)}>
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
                                            fDelete={() => (this.deleteId = app.id)}
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
                        fClose={() => (this.createDialog = false)}
                        fOnSubmit={appStore.create}
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
        const file = e.target.files && e.target.files[0];
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
    description: string;
    fUpload: VoidFunction;
    image: string;
    fDelete: VoidFunction;
}

const Row: SFC<IRowProps> = observer(({name, value, description, fDelete, fUpload, image}) => (
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
            <IconButton onClick={fDelete} className="delete">
                <Delete />
            </IconButton>
        </TableCell>
    </TableRow>
));

export default inject('appStore')(Applications);
