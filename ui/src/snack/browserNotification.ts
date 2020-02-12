import Notify from 'notifyjs';
import removeMarkdown from 'remove-markdown';
import {IMessage} from '../types';

export function mayAllowPermission(): boolean {
    return Notify.needsPermission && Notify.isSupported() && Notification.permission !== 'denied';
}

export function requestPermission() {
    if (Notify.needsPermission && Notify.isSupported()) {
        Notify.requestPermission(
            () => console.log('granted notification permissions'),
            () => console.log('notification permission denied')
        );
    }
}

export function notifyNewMessage(msg: IMessage) {
    const notify = new Notify(msg.title, {
        body: removeMarkdown(msg.message),
        icon: msg.image,
        silent: true,
        notifyClick: closeAndFocus,
        notifyShow: closeAfterTimeout,
    });
    notify.show();
}

function closeAndFocus(event: Event) {
    if (window.parent) {
        window.parent.focus();
    }
    window.focus();
    window.location.href = '/';
    const target = event.target as Notification;
    target.close();
}

function closeAfterTimeout(event: Event) {
    setTimeout(() => {
        const target = event.target as Notification;
        target.close();
    }, 5000);
}
