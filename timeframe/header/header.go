package header

import (
	"encoding/binary"
	"fmt"

	"github.com/syncstreamer/server/timestamp"
)

const (
	HeaderSize = 2 + 4 + 8 + 8 + 8
)

// Size in bytes
type Header struct {
	Version        int                 // 2 (offset 0)
	MetadataSize   int                 // 4 (offset 2)
	DataSize       int                 // 8 (offset 6)
	StartTimestamp timestamp.Timestamp // 8 (offset 14)
	EndTimestamp   timestamp.Timestamp // 8 (offset 22)
}

func Encode(header *Header) []byte {
	data := make([]byte, HeaderSize)

	binary.BigEndian.PutUint16(data[0:], uint16(header.Version))
	binary.BigEndian.PutUint32(data[2:], uint32(header.MetadataSize))
	binary.BigEndian.PutUint64(data[6:], uint64(header.DataSize))
	binary.BigEndian.PutUint64(data[14:], uint64(header.StartTimestamp))
	binary.BigEndian.PutUint64(data[22:], uint64(header.EndTimestamp))

	return data
}

func Decode(buffer []byte) (*Header, error) {
	if len(buffer) < HeaderSize {
		return nil, fmt.Errorf("Header binary buffer has length less than %d bytes", HeaderSize)
	}

	header := Header{
		Version:        int(binary.BigEndian.Uint16(buffer[0:])),
		MetadataSize:   int(binary.BigEndian.Uint32(buffer[2:])),
		DataSize:       int(binary.BigEndian.Uint64(buffer[6:])),
		StartTimestamp: timestamp.Timestamp(binary.BigEndian.Uint64(buffer[14:])),
		EndTimestamp:   timestamp.Timestamp(binary.BigEndian.Uint64(buffer[22:])),
	}

	return &header, nil
}
