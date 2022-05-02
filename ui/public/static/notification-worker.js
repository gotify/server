var gotifyKey = undefined

self.addEventListener("install", event => {
    event.waitUntil(new Promise((resolve, reject) => {
        try {
            // resolves install promise only if search param 'key' was provided
            gotifyKey = new URL(location).searchParams.get('key');
        } catch (e) {
            reject(e)
        }
        if (!gotifyKey) {
            reject("gotify-login-key not provided")
        }
        console.log("Worker recieved gotify-login-key, successfully installed")
        resolve()
    }))
})

self.addEventListener("activate", () => {

    console.log("Notification worker activated")

    const host = location.host
    const wsProto = location.protocol === "https:" ? "wss:" : "ws:"
    const ws = new WebSocket(`${wsProto}//${host}/stream?token=${gotifyKey}`)
    console.log("Notification worker connected to websocket, waiting for messages")

    ws.onmessage = (event) => {

        // check if any client is currently visible
        // if so, skip sending notification from worker
        self.clients.matchAll({
            type: "window",
            includeUncontrolled: true
        })
            .then((windowClients) => {
                var clientVisible = false
                for (var i = 0; i < windowClients.length; i++) {
                    const windowClient = windowClients[i]
                    // check if a client is visible, then break
                    if (windowClient.visibilityState === "visible") {
                        clientVisible = true
                        break
                    }
                }
                return clientVisible
            }) // use the bool to evaluate whether to send a notification 
            .then((clientVisible) => {
                if (!clientVisible) {
                    var msgObj = JSON.parse(event.data) // parse event data, only if not visible
                    self.registration.showNotification("WORKER: " + msgObj.title, {
                        body: msgObj.message
                    })
                } else {
                    console.log("not sending worker notification, as gotify window is visible")
                }
            })

    }

})