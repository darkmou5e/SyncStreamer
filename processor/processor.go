package processor

import (
	"container/list"
	"errors"
	"fmt"
	"time"

	"github.com/syncstreamer/server/timeframe"
	"github.com/syncstreamer/server/timeframe/eventframe"
	"github.com/syncstreamer/server/timestamp"
)

type TimeframeItem struct {
	StartAt timestamp.Timestamp
	EndAt   timestamp.Timestamp
	Data    []byte
}

func convertToTimeframeItem(ef *eventframe.EventFrame) (*TimeframeItem, error) {
	data, err := timeframe.Encode(ef)

	if err != nil {
		return nil, err
	}

	tf := TimeframeItem{
		StartAt: ef.StartAt,
		EndAt:   ef.EndAt,
		Data:    data,
	}

	return &tf, nil
}

type Processor struct {
	timeframes        *list.List
	eventIn           chan *eventframe.Event
	tfRequest         chan (chan []*TimeframeItem)
	currentEventframe *eventframe.EventFrame
}

func (r Processor) AddEvent(ev *eventframe.Event) {
	r.eventIn <- ev
}

func (r Processor) GetTimeframes() []*TimeframeItem {
	responseChan := make(chan []*TimeframeItem)
	r.tfRequest <- responseChan
	return <-responseChan
}

func startNewEventframe() *eventframe.EventFrame {
	return eventframe.StartEventFrame(10_000)
}

func StartNewProcessor() *Processor {
	proc := Processor{
		timeframes:        list.New(),
		eventIn:           make(chan *eventframe.Event),
		tfRequest:         make(chan (chan []*TimeframeItem)),
		currentEventframe: startNewEventframe(),
	}

	completeCurrentEventFrame := func() {
		// TODO: use length from config
		if proc.timeframes.Len() == 5 {
			proc.timeframes.Remove(proc.timeframes.Back())
		}

		tf, err := convertToTimeframeItem(proc.currentEventframe)
		if err != nil {
			panic("Unknown error")
		}

		proc.timeframes.PushFront(tf)
		proc.currentEventframe = startNewEventframe()
	}

	checker := time.NewTicker(time.Duration(16) * time.Millisecond)
	go func() {
		for {
			select {
			case ev := <-proc.eventIn:
				{
				retry:
					err := proc.currentEventframe.AddEventNow(ev)
					if err != nil {
						if errors.Is(err, eventframe.OutOfTimeframeError) {
							completeCurrentEventFrame()
							goto retry
						} else {
							panic(fmt.Errorf("Error %v, the event %v", err, ev))
						}
					}
				}
			case respChan := <-proc.tfRequest:
				{
					tfs := make([]*TimeframeItem, proc.timeframes.Len())
					i := 0
					for tf := proc.timeframes.Front(); tf != nil; tf = tf.Next() {
						if tf.Value == nil {
							tfs[i] = nil
						} else {
							tfs[i] = tf.Value.(*TimeframeItem)
						}
						i = i + 1
					}
					respChan <- tfs
				}
			case <-checker.C:
				{
					if !proc.currentEventframe.IsActive() {
						completeCurrentEventFrame()
					}
				}
				// TODO: add context done cancelation
			}
		}
	}()

	return &proc
}
