import { createMessageId, createWorker, forIt } from './WorkerUtils.mjs'


export class TimeframePlayer {
    _playEventHandler = null
    _worker = null

    _debugLocalTime = null
    _debugLocalTimeLastValue = null
    _debugLocalTimeLastValueTime = 0

    _currentTimeframe = null
    _nextTimeframe = null
    _playbackDelay = null // ms

    // Data loading
    _isTimeframeLoading = false
    _isGetIndexLoading = false

    _channelsPlaybackIndex = {} // {chanId: int}

    _nextTimeframeIsLoaded = false

    _timeframeDuration = 0

    constructor(playEventHandler) {
        this._playEventHandler = playEventHandler
        this._worker = createWorker("DemoWorker.mjs")
    }


    async start() {
        await this._loadInitialTimeframe()
        requestAnimationFrame(runLoopFunction)
    }

    getFrameProgress() {
        return (this._playbackTime() - currentTimeframe.startAt) / timeframeDuration * 100
    }


    async _getIndex() {
        this._isGetIndexLoading = true
        const msgId = createMessageId()
        this._worker.postMessage([msgId, "getIndex"])
        const result = await forIt(msgId)
        this._isGetIndexLoading = false
        return result
    }



    async _getTimeframeworker() {
        this._isTimeframeLoading = true
        const msgId = createMessageId()
        this._worker.postMessage([msgId, "getTimeframe", timeframeId])
        const result = await forIt(msgId)
        this._isTimeframeLoading = false
        return result
    }

    _playbackTime() {
        return Date.now() - this._playbackDelay
    }

    _isCurrentTimeframeEnded() {
        return this._playbackTime() >= this._currentTimeframe.endAt
    }

    _runLoopFunction = () => {
        const time = this._playbackTime()
        const eventsToPlay = []

        if ((this._debugLocalTime && this._debugLocalTimeLastValue) &&
            (this._debugLocalTime === this._debugLocalTimeLastValue)) {
            if ((Date.now() - this._debugLocalTimeLastValue) > 2000) {
                debugger
            }
        } else {
            this._debugLocalTimeLastValue = this._debugLocalTime
        }
        this._debugLocalTime = Date.now()

        Object.keys(this._currentTimeframe.data.channels).forEach(chanId => {
            if (!this._channelsPlaybackIndex.hasOwnProperty(chanId)) {
                this._channelsPlaybackIndex[chanId] = 0
            }
            const chan = this._currentTimeframe.data.channels[chanId]

            while ((this._channelsPlaybackIndex[chanId] < chan.events.length) &&
                (chan.events[this._channelsPlaybackIndex[chanId]].timestamp <= time)) {
                eventsToPlay.push({ chanId, ev: chan.events[this._channelsPlaybackIndex[chanId]] })
                this._channelsPlaybackIndex[chanId]++
            }
        })

        eventsToPlay.forEach(it => this._playEventHandler(it.chanId, it.ev))

        if (!this._nextTimeframeIsLoaded &&
            this._isItTimeToLoadNextTimeframe(time, this._currentTimeframe)) {
            console.log("Loading next")
            this._nextTimeframeIsLoaded = true
            this._loadNext()
        }

        if (this._isCurrentTimeframeEnded()) {
            this._switchToNextTimeframe()
        }

        requestAnimationFrame(this._runLoopFunction)
    }


    // pure
    _getNextTimeframeInfo(time, timeframeIndex) {
        const next = timeframeIndex.find(tf => (time >= tf.startAt) && (time <= tf.endAt))
        if (next) {
            return next
        } else {
            console.error("All timeframes in index are stale")
            return null
        }
    }

    // pure
    _isItTimeToLoadNextTimeframe(currentTime, currentTimeframe) {
        return currentTime > (currentTimeframe.startAt + timeframeDuration / 3)
    }

    async _loadNext() {
        const freshIndex = await this._getIndex()
        const shouldBeValidAtTime = this._playbackTime() + this._timeframeDuration
        const nextTimeframeInfo = this._getNextTimeframeInfo(shouldBeValidAtTime, freshIndex)
        const nextTimeframeData = await this._getTimeframe(nextTimeframeInfo.id)
        // console.log("find next", playbackTime(), shouldBeValidAtTime, freshIndex, currentTimeframe, nextTimeframeInfo, nextTimeframeData, nextTimeframe)
        this._nextTimeframe = nextTimeframeInfo
        this._nextTimeframe.data = nextTimeframeData
    }


    async loadInitialTimeframe() {
        const freshIndex = await this._getIndex()
        const nextTimeframeInfo = freshIndex[2]
        const nextTimeframeData = await this._getTimeframe(nextTimeframeInfo.id)
        this._currentTimeframe = nextTimeframeInfo
        this._currentTimeframe.data = nextTimeframeData
        this._timeframeDuration = (nextTimeframeInfo.endAt - nextTimeframeInfo.startAt)
        this._playbackDelay = this._timeframeDuration * 3 // 3 frames delay
    }


    _switchToNextTimeframe() {
        // console.log("switchh!", playbackTime(), currentTimeframe, nextTimeframe, channelsPlaybackIndex)
        this._currentTimeframe = this._nextTimeframe
        this._channelsPlaybackIndex = {}
        this._nextTimeframeIsLoaded = false
    }
}
