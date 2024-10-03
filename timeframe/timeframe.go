package timeframe

import (
	"fmt"

	"github.com/syncstreamer/server/timeframe/dataitem"
	"github.com/syncstreamer/server/timeframe/eventframe"
	"github.com/syncstreamer/server/timeframe/header"
	"github.com/syncstreamer/server/timeframe/metadata"
)

func Encode(eventFrame *eventframe.EventFrame) ([]byte, error) {
	dataSize := 0
	metadataSize := 0

	metadataSections := [][]byte{}
	offsetInData := 0

	dataSections := [][]byte{}

	for _, eventChannelId := range eventFrame.ChannelOrder {
		eventChannel := eventFrame.Channels[eventChannelId]
		mdRecord := metadata.MetadataRecord{
			OffsetInData: offsetInData,
			ChannelId:    eventChannel.Id,
			ChannelType:  eventChannel.Type,
		}

		metadataRecordBinary := metadata.Encode(&mdRecord)

		dataItemsBinarySize := 0
		for _, item := range eventChannel.Events {
			dataItemsBinarySize = dataItemsBinarySize + dataitem.CalculateBinarySize(item)
		}
		offsetInData = offsetInData + dataItemsBinarySize // next event channel data section offset
		dataSectionBinary := make([]byte, dataItemsBinarySize)

		currentDataSectionOffset := 0
		for _, dataItem := range eventChannel.Events {
			dataItemBinary := dataitem.Encode(dataItem)
			copy(dataSectionBinary[currentDataSectionOffset:], dataItemBinary)
			currentDataSectionOffset = currentDataSectionOffset + len(dataItemBinary)
		}

		metadataSections = append(metadataSections, metadataRecordBinary)
		metadataSize = metadataSize + len(metadataRecordBinary)

		dataSections = append(dataSections, dataSectionBinary)
		dataSize = dataSize + len(dataSectionBinary)
	}

	headerBinary := header.Encode(&header.Header{
		Version:        1,
		MetadataSize:   metadataSize,
		DataSize:       dataSize,
		StartTimestamp: eventFrame.StartAt,
		EndTimestamp:   eventFrame.EndAt,
	})

	totalSize := header.HeaderSize + metadataSize + dataSize
	timeframeBinary := make([]byte, totalSize)

	copy(timeframeBinary[0:], headerBinary)

	timeframeOffset := header.HeaderSize
	for _, metadataItemBinary := range metadataSections {
		copy(timeframeBinary[timeframeOffset:], metadataItemBinary)
		timeframeOffset = timeframeOffset + len(metadataItemBinary)
	}
	for _, dataSectionBinary := range dataSections {
		copy(timeframeBinary[timeframeOffset:], dataSectionBinary)
		timeframeOffset = timeframeOffset + len(dataSectionBinary)
	}

	return timeframeBinary, nil
}

func Decode(timeframeBinary []byte) (*eventframe.EventFrame, error) {
	currentPosition := 0
	headerRecord, err := header.Decode(timeframeBinary)
	currentPosition = header.HeaderSize

	metadataRecords := []*metadata.MetadataRecord{}

	for currentPosition < header.HeaderSize+headerRecord.MetadataSize {
		metadataRecord, err := metadata.Decode(timeframeBinary[currentPosition:])
		// TODO: handle error?
		if err != nil {
			return nil, err
		}
		currentPosition = currentPosition + metadata.CalculateBinarySize(metadataRecord)
		metadataRecords = append(metadataRecords, metadataRecord)
	}

	positionShouldBe := header.HeaderSize + headerRecord.MetadataSize
	if currentPosition != positionShouldBe {
		return nil, fmt.Errorf("Error decoding metadata section. Current position %v != %v", currentPosition, positionShouldBe)
	}

	eventFrame := eventframe.MakeEventFrame()
	eventFrame.StartAt = headerRecord.StartTimestamp
	eventFrame.EndAt = headerRecord.EndTimestamp

	for _, mdRecord := range metadataRecords {
		eventFrame.Channels[mdRecord.ChannelId] = eventframe.MakeChannel(mdRecord.ChannelId, mdRecord.ChannelType)
		eventFrame.ChannelOrder = append(eventFrame.ChannelOrder, mdRecord.ChannelId)
	}

	for i, mdRecord := range metadataRecords {
		isLastDataSection := i == (len(metadataRecords) - 1)
		currentPosition = header.HeaderSize + headerRecord.MetadataSize + mdRecord.OffsetInData

		var dataSectionEndPosition int
		if isLastDataSection {
			dataSectionEndPosition = len(timeframeBinary)
		} else {
			dataSectionEndPosition = header.HeaderSize + headerRecord.MetadataSize + metadataRecords[i+1].OffsetInData
		}

		for currentPosition < dataSectionEndPosition {
			dataItem, err := dataitem.Decode(timeframeBinary[currentPosition:])
			// TODO: handle error?
			if err != nil {
				return nil, err
			}
			currentPosition = currentPosition + dataitem.CalculateBinarySize(dataItem)
			eventFrame.Channels[mdRecord.ChannelId].Events = append(eventFrame.Channels[mdRecord.ChannelId].Events, dataItem)
		}

		positionShouldBe := dataSectionEndPosition
		if currentPosition != positionShouldBe {
			return nil, fmt.Errorf("Error decoding data section. Current position %v != %v", currentPosition, positionShouldBe)
		}
	}

	return eventFrame, err
}
