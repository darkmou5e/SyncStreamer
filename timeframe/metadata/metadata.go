package metadata

import (
	"encoding/binary"
	"fmt"

	"github.com/darkmou5e/syncstreamer/types"
)

type MetadataRecord struct {
	OffsetInData int
	ChannelId    types.Id
	ChannelType  types.ChannelType
}

/*
   metadata record binary spec:
		offset_in_data/8              uint64 // 8 (offset 0)
		chan_id_size/4				  uint32 // 4 (offset 8)
		chan_type_size/4			  uint32 // 4 (offset 12)
		chan_id/chan_id_size		  [chan_id_size]byte // chan_id_size (offset 16)
		chan_type/chan_type_size	  [chan_type_size]byte // chan_type_size (offset 16 + chan_id_size)
*/

const MinMetadataRecordSize = 8 + 4 + 4

func CalculateBinarySize(mdRecord *MetadataRecord) int {
	return MinMetadataRecordSize + len(mdRecord.ChannelId) + len(mdRecord.ChannelType)
}

func Encode(mdRecord *MetadataRecord) []byte {
	channelId := []byte(mdRecord.ChannelId)
	channelType := []byte(mdRecord.ChannelType)
	recordSize := CalculateBinarySize(mdRecord)
	data := make([]byte, recordSize)

	binary.BigEndian.PutUint64(data[0:], uint64(mdRecord.OffsetInData))
	binary.BigEndian.PutUint32(data[8:], uint32(len(channelId)))
	binary.BigEndian.PutUint32(data[12:], uint32(len(channelType)))
	copy(data[16:], mdRecord.ChannelId)
	copy(data[16+len(channelId):], mdRecord.ChannelType)

	return data
}

func Decode(buffer []byte) (*MetadataRecord, error) {
	if len(buffer) < MinMetadataRecordSize {
		return nil, fmt.Errorf("MetadataRecord binary buffer has length less than min %d bytes", MinMetadataRecordSize)
	}

	offsetInData := int(binary.BigEndian.Uint64(buffer[0:]))
	channelIdSize := binary.BigEndian.Uint32(buffer[8:])
	channelTypeSize := binary.BigEndian.Uint32(buffer[12:])

	declaredSize := MinMetadataRecordSize + channelIdSize + channelTypeSize

	if len(buffer) < int(declaredSize) {
		return nil, fmt.Errorf("MetadataRecord binary buffer has length less than declared %d bytes", declaredSize)
	}

	channelIdOffset := MinMetadataRecordSize
	channelTypeOffset := MinMetadataRecordSize + channelIdSize
	metadataRecord := MetadataRecord{
		OffsetInData: offsetInData,
		ChannelId:    types.Id(buffer[channelIdOffset:channelTypeOffset]),
		ChannelType:  types.ChannelType(buffer[channelTypeOffset : channelTypeOffset+channelTypeSize]),
	}

	return &metadataRecord, nil
}
