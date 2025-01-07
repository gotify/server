import {configureStore} from '@reduxjs/toolkit';
import {useDispatch, useSelector} from 'react-redux';

import uiReducer from './ui-slice';
import authReducer from '../user/auth-slice';
import appReducer from '../application/app-slice';
import userReducer from '../user/user-slice';
import clientReducer from '../client/client-slice';
import pluginsReducer from '../plugin/plugin-slice';
import messageReducer from '../message/message-slice';

const store = configureStore({
    reducer: {
        ui: uiReducer,
        auth: authReducer,
        app: appReducer,
        user: userReducer,
        client: clientReducer,
        plugin: pluginsReducer,
        message: messageReducer,
    }
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch

export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();

export default store;
