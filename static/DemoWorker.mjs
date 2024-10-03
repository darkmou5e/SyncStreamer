import { decodeTimeframe } from './timeframe.mjs'

async function getIndex() {
    const response = await fetch("/frame")
    const index = await response.json()
    const preparedIndex = index.map(it => {
        return {
            startAt: it.StartAt,
            endAt: it.EndAt,
            id: it.Id
        }
    })
    return preparedIndex
}

async function getTimeframe(timeframeId) {
    const frame = (await fetch(`/frame/${timeframeId}`))
    const arrayBuffer = await frame.arrayBuffer()
    const frameContent = decodeTimeframe(arrayBuffer)
    return frameContent
}

onmessage = async (msg) => {
    const [msgId, methodName, args] = msg.data
    switch (methodName) {
        case "getIndex": {
            const resp = await getIndex()
            postMessage([msgId, methodName, resp])
            break
        }
        case "getTimeframe": {
            const
                timeframeId = args,
                resp = await getTimeframe(timeframeId)
            postMessage([msgId, methodName, resp])
            break
        }
    }
}
