import Notify from 'notifyjs';

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
        body: msg.message,
        icon: msg.image,
        notifyClick: (event: Event) => closeAndFocus(event, msg.pathonclick),
        notifyShow: closeAfterTimeout,
    });
    notify.show();
}

function closeAndFocus(event: Event, pathonclick: string) {
    if (window.parent) {
        window.parent.focus();
    }
    window.focus();
    if (pathonclick !== '') {
        window.location.href = pathonclick;
    } else {
        window.location.href = '/';
    }
    const target = event.target as Notification;
    target.close();
}

function closeAfterTimeout(event: Event) {
    setTimeout(() => {
        const target = event.target as Notification;
        target.close();
    }, 5000);
}
