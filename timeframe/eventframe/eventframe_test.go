package eventframe

import (
	"testing"

	"github.com/syncstreamer/server/timestamp"
)

func TestEventFrame(t *testing.T) {
	prepareEventFrame := func() *EventFrame {
		eventFrame := MakeEventFrame()
		now := timestamp.Now()
		eventFrame.StartAt = now
		eventFrame.EndAt = now.Add(1000)
		return eventFrame
	}

	t.Run("Add event to frame", func(t *testing.T) {
		eventFrame := prepareEventFrame()
		if len(eventFrame.Channels) != 0 {
			t.Fail()
		}

		err := eventFrame.AddEventNow(&Event{"position", "number", []byte{10}})
		if err != nil {
			t.Fail()
		}

		err = eventFrame.AddEventNow(&Event{"position", "number", []byte{20}})
		if err != nil {
			t.Fail()
		}

		if len(eventFrame.Channels) != 1 {
			t.Fail()
		}

		channel, ok := eventFrame.Channels["position"]
		if !ok {
			t.Fail()
		}

		if channel.Id != "position" || channel.Type != "number" || len(channel.Events) != 2 || channel.Events[0].Data[0] != 10 {
			t.Errorf("%v\n", channel)
		}
	})

	t.Run("Error on channel type mismatch", func(t *testing.T) {
		eventFrame := prepareEventFrame()
		err := eventFrame.AddEventNow(&Event{"position", "number", []byte{10}})
		if err != nil {
			t.Fail()
		}

		err = eventFrame.AddEventNow(&Event{"position", "number", []byte{20}})
		if err != nil {
			t.Fail()
		}

		err = eventFrame.AddEventNow(&Event{"position", "string", []byte{30}})
		if err != TypeMismatchError {
			t.Fail()
		}
	})

	t.Run("Error on event out of timeframe type mismatch", func(t *testing.T) {
		eventFrame := prepareEventFrame()
		err := eventFrame.AddEvent(eventFrame.StartAt.Add(-100), &Event{"position", "number", []byte{10}})
		if err != OutOfTimeframeError {
			t.Fail()
		}

		err = eventFrame.AddEvent(eventFrame.EndAt.Add(100), &Event{"position", "number", []byte{10}})
		if err != OutOfTimeframeError {
			t.Fail()
		}
	})
}
