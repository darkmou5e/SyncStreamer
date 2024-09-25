package dataitem

import (
	"reflect"
	"testing"
)

var dataItemSample = DataItem{
	Timestamp: 1,
	Data:      []byte{0x01, 0x02, 0x03, 0x04},
}

var dataItemSampleBinary = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Timestamp
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // record_data_size
	0x01, 0x02, 0x03, 0x04, // data {1, 2, 3, 4} -> 4 byte
}

func TestDataItemEncoding(t *testing.T) {
	dataItemBinary := Encode(&dataItemSample)

	t.Run("Encoded data item binary size", func(t *testing.T) {
		if len(dataItemBinary) != len(dataItemSampleBinary) {
			t.Fail()
		}
	})
	t.Run("Encoded data item binary data", func(t *testing.T) {
		if !reflect.DeepEqual(dataItemBinary, dataItemSampleBinary) {
			t.Errorf("%v != %v\n", dataItemBinary, dataItemSampleBinary)
		}
	})
}

func TestDataItemDecoding(t *testing.T) {
	t.Run("Detect error on decoding buffer less than size declared", func(t *testing.T) {
		t.Run("Size less than known minimum", func(t *testing.T) {
			dataItem, err := Decode([]byte{})

			if dataItem != nil || err == nil {
				t.Fail()
			}
		})
		t.Run("Size less than declared size", func(t *testing.T) {
			dataItem, err := Decode(dataItemSampleBinary[:MinDataItemSize])

			if dataItem != nil || err == nil {
				t.Fail()
			}
		})
	})

	t.Run("Decoding DataItem buffer", func(t *testing.T) {
		dataItem, err := Decode(dataItemSampleBinary)

		if dataItem == nil || err != nil {
			t.Fail()
		}

		if dataItem.Timestamp != dataItemSample.Timestamp {
			t.Errorf("Timestamp %v != %v", dataItem.Timestamp, dataItemSample.Timestamp)
		}
		if !reflect.DeepEqual(dataItem.Data, dataItemSample.Data) {
			t.Errorf("Data %v := %v", dataItem.Data, dataItemSample.Data)
		}
	})
}
