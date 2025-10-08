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
import {
    arrayMove,
    SortableContext,
    sortableKeyboardCoordinates,
    useSortable,
    verticalListSortingStrategy,
} from '@dnd-kit/sortable';
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

const Applications = observer(() => {
    const {appStore} = useStores();
    const apps = appStore.getItems();
    const [toDeleteApp, setToDeleteApp] = useState<IApplication>();
    const [toUpdateApp, setToUpdateApp] = useState<IApplication>();
    const [createDialog, setCreateDialog] = useState<boolean>(false);
    const [localApps, setLocalApps] = useState<IApplication[]>([]);

    const fileInputRef = useRef<HTMLInputElement>(null);
    const uploadId = useRef(-1);

    const sensors = useSensors(
        useSensor(PointerSensor, {
            activationConstraint: {
                distance: 8,
            },
        }),
        useSensor(KeyboardSensor, {
            coordinateGetter: sortableKeyboardCoordinates,
        })
    );

    useEffect(() => void appStore.refresh(), []);

    useEffect(() => {
        setLocalApps(apps);
    }, [apps]);

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
        if (['image/png', 'image/jpeg', 'image/gif'].indexOf(file.type) !== -1) {
            appStore.uploadImage(uploadId.current, file);
        } else {
            alert('Uploaded file must be of type png, jpeg or gif.');
        }
    };

    const handleDragEnd = (event: DragEndEvent) => {
        const {active, over} = event;

        if (over && active.id !== over.id) {
            setLocalApps((items) => {
                const oldIndex = items.findIndex((item) => item.id === active.id);
                const newIndex = items.findIndex((item) => item.id === over.id);

                const newItems = arrayMove(items, oldIndex, newIndex);
                
                const applicationIds = newItems.map((app) => app.id);
                appStore.reorder(applicationIds).catch((err) => {
                    console.error('Failed to reorder applications:', err);
                    setLocalApps(apps);
                });

                return newItems;
            });
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
                                    <TableCell padding="checkbox" style={{width: 60}} />
                                    <TableCell padding="checkbox" style={{width: 80}} />
                                    <TableCell>Name</TableCell>
                                    <TableCell>Token</TableCell>
                                    <TableCell>Description</TableCell>
                                    <TableCell>Priority</TableCell>
                                    <TableCell>Sort Order</TableCell>
                                    <TableCell>Last Used</TableCell>
                                    <TableCell />
                                    <TableCell />
                                </TableRow>
                            </TableHead>
                            <SortableContext
                                items={localApps.map((app) => app.id)}
                                strategy={verticalListSortingStrategy}>
                                <TableBody>
                                    {localApps.map((app: IApplication) => (
                                        <SortableRow
                                            key={app.id}
                                            app={app}
                                            fUpload={() => handleImageUploadClick(app.id)}
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
                    fOnSubmit={(name, description, defaultPriority, sortOrder) =>
                        appStore.update(toUpdateApp.id, name, description, defaultPriority, sortOrder)
                    }
                    initialDescription={toUpdateApp?.description}
                    initialName={toUpdateApp?.name}
                    initialDefaultPriority={toUpdateApp?.defaultPriority}
                    initialSortOrder={appStore.getByID(toUpdateApp.id).sortOrder}
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
        </DefaultPage>
    );
});

interface ISortableRowProps {
    app: IApplication;
    fUpload: VoidFunction;
    fDelete: VoidFunction;
    fEdit: VoidFunction;
}

const SortableRow = ({app, fUpload, fDelete, fEdit}: ISortableRowProps) => {
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
            <TableCell padding="checkbox">
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
                    <img
                        src={config.get('url') + app.image}
                        alt="app logo"
                        width="40"
                        height="40"
                    />
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
            <TableCell>{app.sortOrder}</TableCell>
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