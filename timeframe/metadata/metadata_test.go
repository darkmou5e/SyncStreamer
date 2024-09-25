package metadata

import (
	"reflect"
	"testing"
)

var metadataRecordSample = MetadataRecord{
	OffsetInData: 4,
	ChannelId:    "id",
	ChannelType:  "type",
}

var metadataRecordBinarySample = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // OffsetInData
	0x00, 0x00, 0x00, 0x02, // ChannelIdSize "id" -> 2 byte
	0x00, 0x00, 0x00, 0x04, // ChannelType "type" -> 4 byte
	105, 100, // Id string data "id"
	116, 121, 112, 101, // Type string data "type"
}

func TestMetadataRecordEncoding(t *testing.T) {
	binaryMetadataRecord := Encode(&metadataRecordSample)

	t.Run("Encoded metadata record binary size", func(t *testing.T) {
		if len(binaryMetadataRecord) != len(metadataRecordBinarySample) {
			t.Fail()
		}
	})

	t.Run("Encoded metadata record binary data", func(t *testing.T) {
		if !reflect.DeepEqual(binaryMetadataRecord, metadataRecordBinarySample) {
			t.Errorf("%v != %v\n", binaryMetadataRecord, metadataRecordBinarySample)
		}
	})
}

func TestMetadataRecordDecoding(t *testing.T) {
	t.Run("Detect error on decoding buffer less than size declared", func(t *testing.T) {
		t.Run("Size less than known minimum", func(t *testing.T) {
			metadataRecord, err := Decode([]byte{})

			if metadataRecord != nil || err == nil {
				t.Fail()
			}
		})
		t.Run("Size less than declared size", func(t *testing.T) {
			metadataRecord, err := Decode(metadataRecordBinarySample[:MinMetadataRecordSize])

			if metadataRecord != nil || err == nil {
				t.Fail()
			}
		})
	})

	t.Run("Decoding MetadataRecord buffer", func(t *testing.T) {
		mdRecord, err := Decode(metadataRecordBinarySample)

		if mdRecord == nil || err != nil {
			t.Fail()
		}

		if mdRecord.OffsetInData != metadataRecordSample.OffsetInData {
			t.Errorf("OffsetInData %v != %v", mdRecord.OffsetInData, metadataRecordSample.OffsetInData)
		}
		if mdRecord.ChannelId != metadataRecordSample.ChannelId {
			t.Errorf("ChannelId %v := %v", mdRecord.ChannelId, metadataRecordSample.ChannelId)
		}
		if mdRecord.ChannelType != metadataRecordSample.ChannelType {
			t.Errorf("ChannelType %v := %v", mdRecord.ChannelType, metadataRecordSample.ChannelType)
		}
	})
}
