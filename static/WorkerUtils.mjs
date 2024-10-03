let nextMessageId = 0
const workerWaiters = {}

function onMessageHandler(msg) {
    const [msgId, _methodName, resp] = msg.data
    const handler = workerWaiters[msgId]
    if (!handler) {
        console.error(`No one callback is waiting for the msgId: '${msgId}'`)
    } else {
        handler(resp)
    }
}


// PUBLIC INTERFACE
export function createMessageId() {
    return nextMessageId++
}

export function createWorker(url) {
    const worker = new Worker(url, { type: "module" })
    worker.onmessage = onMessageHandler
    return worker
}

export function forIt(msgId) {
    return new Promise((resolve, _reject) => {
        // TODO: handle exceptions
        workerWaiters[msgId] = (data) => {
            resolve(data)
        }
    })
}
