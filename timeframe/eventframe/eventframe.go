package eventframe

import (
	"errors"

	"github.com/darkmou5e/syncstreamer/timeframe/dataitem"
	"github.com/darkmou5e/syncstreamer/timestamp"
	"github.com/darkmou5e/syncstreamer/types"
)

type Channel struct {
	Id     types.Id
	Type   types.ChannelType
	Events []*dataitem.DataItem
}

type EventFrame struct {
	StartAt      timestamp.Timestamp
	EndAt        timestamp.Timestamp
	Channels     map[types.Id]*Channel
	ChannelOrder []types.Id
}

func MakeChannel(id types.Id, channelType types.ChannelType) *Channel {
	return &Channel{
		Id:     id,
		Type:   channelType,
		Events: []*dataitem.DataItem{},
	}
}

func MakeEventFrame() *EventFrame {
	return &EventFrame{
		Channels:     make(map[types.Id]*Channel),
		ChannelOrder: []types.Id{},
	}
}

func StartEventFrame(duration timestamp.Duration) *EventFrame {
	eventFrame := MakeEventFrame()
	eventFrame.StartAt = timestamp.Now()
	eventFrame.EndAt = eventFrame.StartAt.Add(duration)
	return eventFrame
}

var TypeMismatchError = errors.New("Event and Channel types mismatch")
var OutOfTimeframeError = errors.New("The Event is out of timeframe")

type Event struct {
	ChannelId types.Id
	EventType types.ChannelType
	EventData []byte
}

func (r *EventFrame) AddEvent(timestamp timestamp.Timestamp, ev *Event) error {
	if timestamp < r.StartAt || timestamp > r.EndAt {
		return OutOfTimeframeError
	}
	channel, ok := r.Channels[ev.ChannelId]
	if !ok {
		channel = MakeChannel(ev.ChannelId, ev.EventType)
		r.Channels[ev.ChannelId] = channel
		r.ChannelOrder = append(r.ChannelOrder, types.Id(ev.ChannelId))
	}

	if channel.Type != ev.EventType {
		return TypeMismatchError
	}

	dataItem := dataitem.DataItem{
		Timestamp: timestamp,
		Data:      ev.EventData,
	}

	channel.Events = append(channel.Events, &dataItem)

	return nil
}

func (r *EventFrame) AddEventNow(ev *Event) error {
	return r.AddEvent(timestamp.Now(), ev)
}

func (r *EventFrame) IsActive() bool {
	now := timestamp.Now()
	return (now > r.StartAt) && (now < r.EndAt)
}
