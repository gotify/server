import React, {ChangeEvent, useEffect, useRef, useState} from 'react';
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
import CloudUpload from '@mui/icons-material/CloudUpload';
import DragIndicator from '@mui/icons-material/DragIndicator';
import Button from '@mui/material/Button';
import {
    DndContext,
    closestCenter,
    KeyboardSensor,
    PointerSensor,
    useSensor,
    useSensors,
    DragEndEvent,
} from '@dnd-kit/core';
import {SortableContext, useSortable, verticalListSortingStrategy} from '@dnd-kit/sortable';
import {CSS} from '@dnd-kit/utilities';

import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import CopyableSecret from '../common/CopyableSecret';
import {AddApplicationDialog} from './AddApplicationDialog';
import * as config from '../config';
import {UpdateApplicationDialog} from './UpdateApplicationDialog';
import {IApplication} from '../types';
import {LastUsedCell} from '../common/LastUsedCell';
import {useStores} from '../stores';
import {observer} from 'mobx-react-lite';
import {makeStyles} from 'tss-react/mui';
import {ButtonBase, Tooltip} from '@mui/material';

const useStyles = makeStyles()((theme) => ({
    imageContainer: {
        '&::after': {
            content: '"×"',
            position: 'absolute',
            top: 0,
            left: 0,
            width: 40,
            height: 40,
            background: theme.palette.error.main,
            color: theme.palette.getContrastText(theme.palette.error.main),
            fontSize: 40,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            opacity: 0,
        },
        '&:hover::after': {opacity: 1},
    },
}));

const Applications = observer(() => {
    const {appStore} = useStores();
    const apps = appStore.getItems();
    const [toDeleteApp, setToDeleteApp] = useState<IApplication>();
    const [toDeleteImage, setToDeleteImage] = useState<IApplication>();
    const [toUpdateApp, setToUpdateApp] = useState<IApplication>();
    const [createDialog, setCreateDialog] = useState<boolean>(false);

    const fileInputRef = useRef<HTMLInputElement>(null);
    const uploadId = useRef(-1);

    const sensors = useSensors(useSensor(PointerSensor), useSensor(KeyboardSensor));

    useEffect(() => void appStore.refresh(), []);

    const validExtensions = ['.gif', '.png', '.jpg', '.jpeg'];

    const handleImageUploadClick = (id: number) => {
        uploadId.current = id;
        if (fileInputRef.current) {
            fileInputRef.current.click();
        }
    };

    const onUploadImage = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) {
            return;
        }
        appStore.uploadImage(uploadId.current, file);
    };

    const handleDragEnd = (event: DragEndEvent) => {
        const {active, over} = event;

        if (over && active.id !== over.id) {
            appStore.reorder(active.id as number, over.id as number);
        }
    };

    return (
        <DefaultPage
            title="Applications"
            rightControl={
                <Button
                    id="create-app"
                    variant="contained"
                    color="primary"
                    onClick={() => setCreateDialog(true)}>
                    Create Application
                </Button>
            }
            maxWidth={1000}>
            <Grid size={12}>
                <Paper elevation={6} style={{overflowX: 'auto'}}>
                    <DndContext
                        sensors={sensors}
                        collisionDetection={closestCenter}
                        onDragEnd={handleDragEnd}>
                        <Table id="app-table">
                            <TableHead>
                                <TableRow>
                                    <TableCell padding="none" style={{width: 0}} />
                                    <TableCell padding="checkbox" style={{width: 80}} />
                                    <TableCell>Name</TableCell>
                                    <TableCell>Token</TableCell>
                                    <TableCell>Description</TableCell>
                                    <TableCell>Priority</TableCell>
                                    <TableCell>Last Used</TableCell>
                                    <TableCell />
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <SortableContext items={apps} strategy={verticalListSortingStrategy}>
                                <TableBody>
                                    {apps.map((app: IApplication) => (
                                        <Row
                                            key={app.id}
                                            app={app}
                                            fUpload={() => handleImageUploadClick(app.id)}
                                            fDeleteImage={() => setToDeleteImage(app)}
                                            fDelete={() => setToDeleteApp(app)}
                                            fEdit={() => setToUpdateApp(app)}
                                        />
                                    ))}
                                </TableBody>
                            </SortableContext>
                        </Table>
                    </DndContext>
                    <input
                        ref={fileInputRef}
                        type="file"
                        accept={validExtensions.join(',')}
                        style={{display: 'none'}}
                        onChange={onUploadImage}
                    />
                </Paper>
            </Grid>
            {createDialog && (
                <AddApplicationDialog
                    fClose={() => setCreateDialog(false)}
                    fOnSubmit={appStore.create}
                />
            )}
            {toUpdateApp != null && (
                <UpdateApplicationDialog
                    fClose={() => setToUpdateApp(undefined)}
                    fOnSubmit={(name, description, defaultPriority) =>
                        appStore.update({...toUpdateApp, name, description, defaultPriority})
                    }
                    initialDescription={toUpdateApp?.description}
                    initialName={toUpdateApp?.name}
                    initialDefaultPriority={toUpdateApp?.defaultPriority}
                />
            )}
            {toDeleteApp != null && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete ' + toDeleteApp.name + '?'}
                    fClose={() => setToDeleteApp(undefined)}
                    fOnSubmit={() => appStore.remove(toDeleteApp.id)}
                />
            )}
            {toDeleteImage != null && (
                <ConfirmDialog
                    title="Confirm Delete Image"
                    text={'Delete image for ' + toDeleteImage.name + '?'}
                    fClose={() => setToDeleteImage(undefined)}
                    fOnSubmit={() => appStore.deleteImage(toDeleteImage.id)}
                />
            )}
        </DefaultPage>
    );
});

interface IRowProps {
    app: IApplication;
    fUpload: VoidFunction;
    fDeleteImage: VoidFunction;
    fDelete: VoidFunction;
    fEdit: VoidFunction;
}

const Row = ({app, fDelete, fUpload, fDeleteImage, fEdit}: IRowProps) => {
    const {classes} = useStyles();
    const isDefaultImage = app.image === 'static/defaultapp.png';

    const {attributes, listeners, setNodeRef, transform, transition, isDragging} = useSortable({
        id: app.id,
    });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
        backgroundColor: isDragging ? '#f5f5f5' : 'transparent',
    };

    return (
        <TableRow ref={setNodeRef} style={style}>
            <TableCell padding="none" style={{paddingLeft: 5}}>
                <div
                    {...attributes}
                    {...listeners}
                    style={{
                        cursor: 'grab',
                        display: 'flex',
                        alignItems: 'center',
                        touchAction: 'none',
                    }}>
                    <DragIndicator style={{color: '#999'}} />
                </div>
            </TableCell>
            <TableCell padding="normal">
                <div style={{display: 'flex'}}>
                    <Tooltip title="Delete image" placement="top" arrow>
                        <ButtonBase
                            className={classes.imageContainer}
                            onClick={fDeleteImage}
                            disabled={isDefaultImage}>
                            <img
                                src={config.get('url') + app.image}
                                alt="app logo"
                                width="40"
                                height="40"
                            />
                        </ButtonBase>
                    </Tooltip>
                    <IconButton onClick={fUpload} style={{height: 40}}>
                        <CloudUpload />
                    </IconButton>
                </div>
            </TableCell>
            <TableCell>{app.name}</TableCell>
            <TableCell>
                <CopyableSecret value={app.token} style={{display: 'flex', alignItems: 'center'}} />
            </TableCell>
            <TableCell>{app.description}</TableCell>
            <TableCell>{app.defaultPriority}</TableCell>
            <TableCell>
                <LastUsedCell lastUsed={app.lastUsed} />
            </TableCell>
            <TableCell align="right" padding="none">
                <IconButton onClick={fEdit} className="edit">
                    <Edit />
                </IconButton>
            </TableCell>
            <TableCell align="right" padding="none">
                <IconButton onClick={fDelete} className="delete" disabled={app.internal}>
                    <Delete />
                </IconButton>
            </TableCell>
        </TableRow>
    );
};

export default Applications;
