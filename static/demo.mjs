import { LedDisplay } from './LedDisplay.mjs'
import { TimeframePlayer } from './TimeframePlayer.mjs'

let player = null
let led = null

let streamtimeEl = null
let localtimeEl = null
let deltaEl = null
let frameProgressEl = null
let frameProgressBarEl = null


function setDisplay(streamTime) {
    streamtimeEl.innerText = streamTime
    localtimeEl.innerText = Date.now()
    deltaEl.innerText = Date.now() - streamTime

    const progress = player.getFrameProgress()
    frameProgressEl.innerText = `${Math.trunc(progress)}%`
    frameProgressBarEl.style.width = `${90 * (progress / 100)}%`
}

function playEvent(channelId, event) {
    switch (channelId) {
        case 'led':
            led.update(event.data)
            break
        case 'streamtime':
            setDisplay(event.data.time)
            break
    }
}


function main() {
    streamtimeEl = document.getElementById("streamtime")
    localtimeEl = document.getElementById("localtime")
    deltaEl = document.getElementById("delta")
    frameProgressEl = document.getElementById("frameProgress")
    frameProgressBarEl = document.getElementById("frameProgressBar")

    const rootEl = document.getElementById("led")
    led = new LedDisplay(47, 6)
    led.mount(rootEl)

    player = new TimeframePlayer(playEvent)
    player.start()
}

main()
