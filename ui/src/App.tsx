import React, {useEffect} from 'react';
import {createBrowserRouter, RouterProvider} from 'react-router-dom';
import Applications from './application/Applications.tsx';
import Clients from './client/Clients.tsx';
import {checkAuthLoader} from './common/Auth.ts';
import Messages from './message/Messages.tsx';
import {WebSocketStore} from './message/WebSocketStore.ts';
import RootLayout from './pages/root';
import Plugins from './plugin/Plugins.tsx';
import * as Notifications from './snack/browserNotification.ts';
import {useAppDispatch, useAppSelector} from './store';
import {messageActions} from './store/message-slice.ts';
import Login from './user/Login.tsx';
import Users from './user/Users.tsx';

const router = createBrowserRouter([
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
                element: <Plugins />,
                loader: checkAuthLoader,
                children: [
                    {
                        path: ':id',
                        // TODO: fix problem in PluginDetailView
                        //element: <PluginDetailView />,
                    },
                ],
            },
        ],
    },
]);

const App = () => {
    const dispatch = useAppDispatch();
    const loggedIn = useAppSelector((state) => state.auth.loggedIn);
    const ws = new WebSocketStore();

    useEffect(() => {
        if (loggedIn) {
            ws.listen((message) => {
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
