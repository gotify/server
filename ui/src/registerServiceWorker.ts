export function unregister() {
    if ('serviceWorker' in navigator) {
        navigator.serviceWorker.ready.then((registration) => {
            registration.unregister();
        });
    }
}

export function registerNotificationWorker(key: string) {
    if ('serviceWorker' in navigator) {
        // provide the gotify-login-key as query parameter, as the service worker cannot access
        // localStorage. There is no need to implement a mechanism to update the key as the
        // worker will be unregistered on logout.
        navigator.serviceWorker.register("static/notification-worker.js?key=" + key, {
            scope: "/static/notification-worker"
        })
            .catch(console.error)
    } else {
        console.error("Service workers are not supported in your browser!")
    }
}

export function unregisterNotificationWorker() {
    if ('serviceWorker' in navigator) {
        // get service worker by scope
        navigator.serviceWorker.getRegistration("/static/notification-worker").then((reg) => {
            if (reg) {
                reg.unregister().then(ok => { // ok: bool === true, if unregister was successfull
                    if (!ok) {
                        console.error("Error unregistering service worker")
                    }
                })
            } else {
                console.error("Error finding service worker by scope")
            }
        })
    } else {
        console.error("Service workers are not supported in your browser!")
    }
}