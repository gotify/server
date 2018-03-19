import dispatcher from './dispatcher';
import Notify from 'notifyjs';

export function requestPermission() {
    if (Notify.needsPermission && Notify.isSupported()) {
        Notify.requestPermission(() => console.log('granted notification permissions'),
            () => console.log('notification permission denied'));
    }
}

function closeAndFocus(event) {
    if (window.parent) {
        window.parent.focus();
    }
    window.focus();
    window.location.href = '/';
    event.target.close();
}

function closeAfterTimeout(event) {
    setTimeout(() => {
        event.target.close();
    }, 5000);
}

dispatcher.register((data) => {
    if (data.type === 'ONE_MESSAGE') {
        const msg = data.payload;

        const notify = new Notify(msg.title, {
            body: msg.message,
            icon: '/static/favicon.ico',
            notifyClick: closeAndFocus,
            notifyShow: closeAfterTimeout,
        });
        notify.show();
    }
});
