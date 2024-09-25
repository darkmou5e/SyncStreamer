package timeframe

import (
	"reflect"
	"testing"

	"github.com/darkmou5e/syncstreamer/timeframe/eventframe"
)

func TestTimeFrame(t *testing.T) {
	var eventFrameSample = eventframe.MakeEventFrame()
	eventFrameSample.StartAt = 1
	eventFrameSample.EndAt = 10
	eventFrameSample.AddEvent(1, &eventframe.Event{"position", "number", []byte{100}}) // a1
	eventFrameSample.AddEvent(1, &eventframe.Event{"color", "number", []byte{1}})      // b1
	eventFrameSample.AddEvent(2, &eventframe.Event{"position", "number", []byte{200}}) // a2
	eventFrameSample.AddEvent(3, &eventframe.Event{"color", "number", []byte{2}})      // b2

	var timeframeSampleBinary = []byte{
		// header
		0x00, 0x01, // version
		0x00, 0x00, 0x00, 57, // metadata_size
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 68, // data_size
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // start_timestamp
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A, // end_timestamp

		// metadata a -> size 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // OffsetInData
		0x00, 0x00, 0x00, 0x08, // ChannelIdSize "position"
		0x00, 0x00, 0x00, 0x06, // ChannelType "number"
		112, 111, 115, 105, 116, 105, 111, 110, // Id string data "position"
		110, 117, 109, 98, 101, 114, // Type string data "number"

		// metadata b -> size 27
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 34, // OffsetInData
		0x00, 0x00, 0x00, 0x05, // ChannelIdSize "color"
		0x00, 0x00, 0x00, 0x06, // ChannelType "number"
		99, 111, 108, 111, 114, // Id string data "color"
		110, 117, 109, 98, 101, 114, // Type string data "number"

		// data a1 -> size 17
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Timestamp
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // record_data_size
		100, // data

		// data a2
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // Timestamp
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // record_data_size
		200, // data

		// data b1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Timestamp
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // record_data_size
		0x01, // data

		// data b2
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // Timestamp
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // record_data_size
		0x02, // data
	}

	t.Run("Encoding", func(t *testing.T) {
		timeFrameBinary, err := Encode(eventFrameSample)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(timeFrameBinary, timeframeSampleBinary) {
			t.Errorf("\n%v\n != \n%v\n", timeFrameBinary, timeframeSampleBinary)
		}
	})

	t.Run("Decoding", func(t *testing.T) {
		eventFrame, err := Decode(timeframeSampleBinary)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(eventFrame, eventFrameSample) {
			t.Errorf("\n%v\n != \n%v\n", eventFrame, eventFrameSample)
		}
	})
}
