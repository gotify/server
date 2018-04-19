import Notify from 'notifyjs';
import dispatcher, {IEvent} from './dispatcher';

export function requestPermission() {
    if (Notify.needsPermission && Notify.isSupported()) {
        Notify.requestPermission(() => console.log('granted notification permissions'),
            () => console.log('notification permission denied'));
    }
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

dispatcher.register((data: IEvent): void => {
    if (data.type === 'ONE_MESSAGE') {
        const msg = data.payload;

        const notify = new Notify(msg.title, {
            body: msg.message,
            icon: msg.image,
            notifyClick: closeAndFocus,
            notifyShow: closeAfterTimeout,
        });
        notify.show();
    }
});
