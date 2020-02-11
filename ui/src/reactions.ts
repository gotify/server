import {StoreMapping} from './inject';
import {reaction} from 'mobx';
import * as Notifications from './snack/browserNotification';

export const registerReactions = (stores: StoreMapping) => {
    const clearAll = () => {
        stores.messagesStore.clearAll();
        stores.appStore.clear();
        stores.clientStore.clear();
        stores.userStore.clear();
        stores.wsStore.close();
    };
    const loadAll = () => {
        stores.wsStore.listen((message) => {
            stores.messagesStore.publishSingleMessage(message);
            Notifications.notifyNewMessage(message);
            if (message.priority >= 4) {
                const src = 'static/notification.ogg';
                const audio = new Audio(src);
                audio.play();
            }
        });
        stores.appStore.refresh();
    };

    reaction(
        () => stores.currentUser.loggedIn,
        (loggedIn) => {
            if (loggedIn) {
                loadAll();
            } else {
                clearAll();
            }
        }
    );

    reaction(
        () => stores.currentUser.connectionErrorMessage,
        (connectionErrorMessage) => {
            if (!connectionErrorMessage) {
                clearAll();
                loadAll();
            }
        }
    );
};
