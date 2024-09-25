package header

import (
	"testing"
	"reflect"
)

var headerSample = Header{
	Version:        1,
	MetadataSize:   2,
	DataSize:       3,
	StartTimestamp: 4,
	EndTimestamp:   5,
}

var headerBinaryDataSample = []byte{
	0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x02, // metadata_size
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // data_size
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // start_timestamp
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // end_timestamp
}

func TestHeaderEncoding(t *testing.T) {
	binaryHeader := Encode(&headerSample)

	t.Run("Encoded header binary size", func(t *testing.T) {
		if len(binaryHeader) != HeaderSize {
			t.Fail()
		}
	})

	t.Run("Encoded header binary data", func(t *testing.T) {
		if !reflect.DeepEqual(binaryHeader, headerBinaryDataSample) {
			t.Errorf("%v != %v\n", binaryHeader, headerBinaryDataSample)
		}
	})
}

func TestHeaderDecoding(t *testing.T) {
	t.Run("Detect error on decoding buffer less than HeaderSize", func(t *testing.T) {
		header, err := Decode([]byte{})

		if header != nil || err == nil {
			t.Fail()
		}
	})

	t.Run("Decoding Header buffer", func(t *testing.T) {
		header, err := Decode(headerBinaryDataSample)

		if header == nil || err != nil {
			t.Fail()
		}

		if header.Version != headerSample.Version {
			t.Errorf("Header decoding error: Version should be equal %v, but got %v", headerSample.Version, headerSample.Version)
		}
		if header.MetadataSize != headerSample.MetadataSize {
			t.Errorf("Header decoding error: MetadataSize should be equal %v, but got %v", headerSample.MetadataSize, headerSample.MetadataSize)
		}
		if header.DataSize != headerSample.DataSize {
			t.Errorf("Header decoding error: DataSize should be equal %v, but got %v", headerSample.DataSize, headerSample.DataSize)
		}
		if header.StartTimestamp != headerSample.StartTimestamp {
			t.Errorf("Header decoding error: StartTimestamp should be equal %v, but got %v", headerSample.StartTimestamp, headerSample.StartTimestamp)
		}
		if header.EndTimestamp != headerSample.EndTimestamp {
			t.Errorf("Header decoding error: EndTimestamp should be equal %v, but got %v", headerSample.EndTimestamp, headerSample.EndTimestamp)
		}
	})
}
