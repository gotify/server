import React, {Component} from 'react';
import Grid from 'material-ui/Grid';
import Table, {TableBody, TableCell, TableHead, TableRow} from 'material-ui/Table';
import Paper from 'material-ui/Paper';
import Button from 'material-ui/Button';
import IconButton from 'material-ui/IconButton';
import Dialog, {DialogActions, DialogContent, DialogContentText, DialogTitle} from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import Tooltip from 'material-ui/Tooltip';
import AppStore from '../stores/AppStore';
import ToggleVisibility from '../component/ToggleVisibility';
import ConfirmDialog from '../component/ConfirmDialog';
import * as AppAction from '../actions/AppAction';
import DefaultPage from '../component/DefaultPage';
import PropTypes from 'prop-types';
import Delete from 'material-ui-icons/Delete';
import Avatar from 'material-ui/Avatar';
import Edit from 'material-ui-icons/Edit';

class Applications extends Component {
    constructor() {
        super();
        this.uploadId = -1;
        this.state = {apps: [], createDialog: false, deleteId: -1};
    }

    componentWillMount() {
        AppStore.on('change', this.updateApps);
        this.updateApps();
    }

    componentWillUnmount() {
        AppStore.removeListener('change', this.updateApps);
    }

    updateApps = () => this.setState({...this.state, apps: AppStore.get()});

    uploadImage = (id) => {
        this.uploadId = id;
        this.upload.click();
    };
    onUploadImage = (e) => {
        const file = e.target.files[0];
        if (!file) {
            return;
        }
        if (['image/png', 'image/jpeg', 'image/gif'].indexOf(file['type']) !== -1) {
            AppAction.uploadImage(this.uploadId, file);
        } else {
            alert('Uploaded file must be of type png, jpeg or gif.');
        }
    };

    showCreateDialog = () => this.setState({...this.state, createDialog: true});
    hideCreateDialog = () => this.setState({...this.state, createDialog: false});

    showCloseDialog = (deleteId) => this.setState({...this.state, deleteId});
    hideCloseDialog = () => this.setState({...this.state, deleteId: -1});

    render() {
        const {apps, createDialog, deleteId} = this.state;
        return (
            <DefaultPage title="Applications" buttonTitle="Create Application" maxWidth={1000}
                         fButton={this.showCreateDialog}>
                <Grid item xs={12}>
                    <Paper elevation={6}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell padding="checkbox" style={{width: 80}}/>
                                    <TableCell>Name</TableCell>
                                    <TableCell>Token</TableCell>
                                    <TableCell>Description</TableCell>
                                    <TableCell/>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {apps.map((app) => {
                                    return (
                                        <Row key={app.id} appId={app.id} description={app.description} image={app.image}
                                             name={app.name} value={app.token} fUpload={() => this.uploadImage(app.id)}
                                             fDelete={() => this.showCloseDialog(app.id)}/>
                                    );
                                })}
                            </TableBody>
                        </Table>
                        <input ref={(upload) => this.upload = upload} type="file" style={{display: 'none'}}
                               onChange={this.onUploadImage}/>
                    </Paper>
                </Grid>
                {createDialog && <AddDialog fClose={this.hideCreateDialog} fOnSubmit={AppAction.createApp}/>}
                {deleteId !== -1 && <ConfirmDialog title="Confirm Delete"
                                                   text={'Delete ' + AppStore.getById(deleteId).name + '?'}
                                                   fClose={this.hideCloseDialog}
                                                   fOnSubmit={() => AppAction.deleteApp(deleteId)}
                />}
            </DefaultPage>
        );
    }
}

class Row extends Component {
    static propTypes = {
        name: PropTypes.string.isRequired,
        value: PropTypes.string.isRequired,
        description: PropTypes.string.isRequired,
        fDelete: PropTypes.func.isRequired,
        fUpload: PropTypes.func.isRequired,
        image: PropTypes.string.isRequired,
    };

    render() {
        const {name, value, description, fDelete, fUpload, image} = this.props;
        return (
            <TableRow>
                <TableCell padding="checkbox">
                    <div style={{display: 'flex'}}>
                        <Avatar src={image}/><IconButton onClick={fUpload} style={{height: 40}}><Edit/></IconButton>
                    </div>
                </TableCell>
                <TableCell>{name}</TableCell>
                <TableCell>
                    <ToggleVisibility value={value} style={{display: 'flex', alignItems: 'center'}}/>
                </TableCell>
                <TableCell>{description}</TableCell>
                <TableCell numeric padding="none">
                    <IconButton onClick={fDelete}><Delete/></IconButton>
                </TableCell>
            </TableRow>
        );
    }
}

class AddDialog extends Component {
    static propTypes = {
        fClose: PropTypes.func.isRequired,
        fOnSubmit: PropTypes.func.isRequired,
    };

    constructor() {
        super();
        this.state = {name: '', description: ''};
    }

    handleChange(propertyName, event) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    render() {
        const {fClose, fOnSubmit} = this.props;
        const {name, description} = this.state;
        const submitEnabled = this.state.name.length !== 0;
        const submitAndClose = () => {
            fOnSubmit(name, description);
            fClose();
        };
        return (
            <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Create an application</DialogTitle>
                <DialogContent>
                    <DialogContentText>An application is allowed to send messages.</DialogContentText>
                    <TextField autoFocus margin="dense" id="name" label="Name *" type="email" value={name}
                               onChange={this.handleChange.bind(this, 'name')} fullWidth/>
                    <TextField margin="dense" id="description" label="Short Description" value={description}
                               onChange={this.handleChange.bind(this, 'description')} fullWidth multiline/>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fClose}>Cancel</Button>
                    <Tooltip title={submitEnabled ? '' : 'name is required'}>
                        <div>
                            <Button disabled={!submitEnabled} onClick={submitAndClose} color="primary" variant="raised">
                                Create
                            </Button>
                        </div>
                    </Tooltip>
                </DialogActions>
            </Dialog>
        );
    }
}

export default Applications;
