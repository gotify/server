import React from 'react';
import {createHashRouter, RouterProvider} from 'react-router-dom';
import Applications from './application/Applications.tsx';
import Clients from './client/Clients.tsx';
import {checkAuthLoader} from './common/Auth.ts';
import Messages from './message/Messages.tsx';
import PluginsRootLayout from './pages/plugins.tsx';
import RootLayout from './pages/root';
import PluginDetailView from './plugin/PluginDetailView.tsx';
import Plugins from './plugin/Plugins.tsx';
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

const App = () => (<RouterProvider router={router} />);

export default App;
