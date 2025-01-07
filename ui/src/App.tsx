import React, {useEffect} from 'react';
import {createHashRouter, RouterProvider} from 'react-router-dom';
import Applications from './application/Applications.tsx';
import Clients from './client/Clients.tsx';
import {checkAuthLoader} from './common/Auth.ts';
import Messages from './message/Messages.tsx';
import {WebSocketStore} from './message/WebSocketStore.ts';
import PluginsRootLayout from './pages/plugins.tsx';
import RootLayout from './pages/root';
import PluginDetailView from './plugin/PluginDetailView.tsx';
import Plugins from './plugin/Plugins.tsx';
import * as Notifications from './snack/browserNotification.ts';
import {useAppDispatch, useAppSelector} from './store';
import {messageActions} from './message/message-slice.ts';
import Login from './user/Login.tsx';
import Users from './user/Users.tsx';

const router = createHashRouter([
    {
        path: '/',
        element: <RootLayout />,
        children: [
            {
                index: true,
                element: <Messages />,
                loader: checkAuthLoader,
            },
            {
                path: 'messages',
                element: <Messages />,
                loader: checkAuthLoader,
                children: [
                    {
                        path: ':id',
                        element: <Messages />,
                    }
                ]
            },
            {
                path: 'login',
                element: <Login />,
            },
            {
                path: 'applications',
                element: <Applications />,
                loader: checkAuthLoader,
            },
            {
                path: 'users',
                element: <Users />,
                loader: checkAuthLoader,
            },
            {
                path: 'clients',
                element: <Clients />,
                loader: checkAuthLoader,
            },
            {
                path: 'plugins',
                element: <PluginsRootLayout />,
                loader: checkAuthLoader,
                children: [
                    { index: true, element: <Plugins /> },
                    {
                        path: ':id',
                        element: <PluginDetailView />,
                    }
                ],
            },
        ],
    },
]);

const ws = new WebSocketStore();

const App = () => {
    const dispatch = useAppDispatch();
    const loggedIn = useAppSelector((state) => state.auth.loggedIn);

    useEffect(() => {
        if (loggedIn) {
            ws.listen((message) => {
                dispatch(messageActions.loading(true));
                dispatch(messageActions.add(message));
                Notifications.notifyNewMessage(message);
                if (message.priority >= 4) {
                    const src = 'static/notification.ogg';
                    const audio = new Audio(src);
                    audio.play();
                }
            });
            window.onbeforeunload = () => {
                ws.close();
            };
        } else {
            ws.close();
        }
    }, [dispatch, loggedIn]);


    return (
        <RouterProvider router={router} />
    );
};

export default App;

/*
<Routes>
    {authenticating ? (<Route path="/" element={<LoadingSpinner />} />) : null}
    <Route path="/" element={<Messages />} />
    <Route path="messages/:id" element={<Messages />} />

</Routes>
*/
