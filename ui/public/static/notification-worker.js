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

    ws.onmessage = function(event) {
        msgObj = JSON.parse(event.data)
        self.registration.showNotification("WORKER: " + msgObj.message)
    }

})