package dataitem

import (
	"encoding/binary"
	"fmt"

	"github.com/syncstreamer/server/timestamp"
)

type DataItem struct {
	Timestamp timestamp.Timestamp
	Data      []byte
}

/*
   data item record binary spec:
		timestamp/8              	int64 // 8 (offset 0)
		record_data_size/8			int64 // 8 (offset 8)
		data/record_data_size		[record_data_size]byte // record_data_size (offset 16)
*/

const MinDataItemSize = 8 + 8

func CalculateBinarySize(dataItem *DataItem) int {
	return MinDataItemSize + len(dataItem.Data)
}

func Encode(dataItem *DataItem) []byte {
	recordSize := MinDataItemSize + len(dataItem.Data)
	data := make([]byte, recordSize)

	binary.BigEndian.PutUint64(data[0:], uint64(dataItem.Timestamp))
	binary.BigEndian.PutUint64(data[8:], uint64(len(dataItem.Data)))
	copy(data[16:], dataItem.Data)

	return data
}

func Decode(buffer []byte) (*DataItem, error) {
	if len(buffer) < MinDataItemSize {
		return nil, fmt.Errorf("Binary buffer has length less than min %d bytes", MinDataItemSize)
	}

	timestamp := timestamp.Timestamp(binary.BigEndian.Uint64(buffer[0:]))
	dataSize := int(binary.BigEndian.Uint64(buffer[8:]))

	declaredSize := MinDataItemSize + dataSize

	if len(buffer) < declaredSize {
		return nil, fmt.Errorf("Binary buffer has length less than declared %d bytes", declaredSize)
	}

	dataItem := DataItem{
		Timestamp: timestamp,
		Data:      buffer[MinDataItemSize:declaredSize],
	}

	return &dataItem, nil
}
