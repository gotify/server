import {reaction} from 'mobx';
import * as Notifications from './snack/browserNotification';
import {StoreMapping} from './stores';

const AUDIO_REPEAT_DELAY = 1000;

export const registerReactions = (stores: StoreMapping) => {
    const clearAll = () => {
        stores.messagesStore.clearAll();
        stores.appStore.clear();
        stores.clientStore.clear();
        stores.userStore.clear();
        stores.wsStore.close();
    };

    let audio: HTMLAudioElement | undefined;
    let lastAudio = 0;

    const loadAll = () => {
        stores.wsStore.listen((message) => {
            stores.messagesStore.publishSingleMessage(message);
            Notifications.notifyNewMessage(message);
            if (message.priority >= 4 && Date.now() > lastAudio + AUDIO_REPEAT_DELAY) {
                lastAudio = Date.now();

                audio ??= new Audio('static/notification.ogg');
                audio.currentTime = 0;
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
                stores.currentUser.refreshKey++;
            }
        }
    );
};
