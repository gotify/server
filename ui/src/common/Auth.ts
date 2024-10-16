import {redirect} from 'react-router-dom';
import {tokenKey} from '../store/auth-actions.ts';

export function getAuthToken() {
    const token = localStorage.getItem(tokenKey);
    return token;
}

export function tokenLoader() {
    return getAuthToken();
}

export function checkAuthLoader() {
    const token = getAuthToken();
    if (!token) {
        return redirect('/login');
    }

    return null;
}
