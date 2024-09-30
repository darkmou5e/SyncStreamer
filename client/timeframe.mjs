/*
Timeframe format v1
Port of Go code

Size in bytes

HeaderSize = 2 + 4 + 8 + 8 + 8

type Header struct {
  Version        int                 2 (offset 0)
  MetadataSize   int                 4 (offset 2)
  DataSize       int                 8 (offset 6)
  StartTimestamp timestamp.Timestamp 8 (offset 14)
  EndTimestamp   timestamp.Timestamp 8 (offset 22)
}

*/

// HEADER
export const headerSize = 2 + 4 + 8 + 8 + 8

// buffer: ArrayBuffer
export function decodeHeader(buffer) {
  const dataView = new DataView(buffer)

  return {
    version: Number(dataView.getInt16(0)),
    metadataSize: Number(dataView.getInt32(2)),
    dataSize: Number(dataView.getBigInt64(6)),
    startTimestamp: Number(dataView.getBigInt64(14)),
    endTimestamp: Number(dataView.getBigInt64(22))
  }
}


// METADATA
/*
   metadata record binary spec:
    offset_in_data/8              uint64 // 8 (offset 0)
    chan_id_size/4				  uint32 // 4 (offset 8)
    chan_type_size/4			  uint32 // 4 (offset 12)
    chan_id/chan_id_size		  [chan_id_size]byte // chan_id_size (offset 16)
    chan_type/chan_type_size	  [chan_type_size]byte // chan_type_size (offset 16 + chan_id_size)

type MetadataRecord struct {
  OffsetInData int
  ChannelId    types.Id
  ChannelType  types.ChannelType
}
*/

export function decodeMetadataRecord(buffer) {
  const dataView = new DataView(buffer)
  const textDecoder = new TextDecoder()

  const channelIdSize = dataView.getInt32(8)
  const channelTypeSize = dataView.getInt32(12)
  const channelTypeOffset = 16 + channelIdSize

  return {
    offsetInData: Number(dataView.getBigInt64(0)),
    channelId: textDecoder.decode(buffer.slice(16, channelTypeOffset)),
    channelType: textDecoder.decode(buffer.slice(channelTypeOffset, channelTypeOffset + channelTypeSize))
  }
}



// DATA


/*
   data item record binary spec:
    timestamp/8              	int64 // 8 (offset 0)
    record_data_size/8			int64 // 8 (offset 8)
    data/record_data_size		[record_data_size]byte // record_data_size (offset 16)


type DataItem struct {
  Timestamp timestamp.Timestamp
  Data      []byte
}
*/

export function decodeDataItem(buffer, parseJSON = false) {
  const dataView = new DataView(buffer)
  const dataSize = Number(dataView.getBigInt64(8))
  const rawData = buffer.slice(16, 16 + dataSize)

  let data

  if (parseJSON) {
    const textDecoder = new TextDecoder()
    data = JSON.parse(textDecoder.decode(rawData))
  } else {
    data = new Uint8Array(rawData)
  }

  const itemSize = minDataItemSize + rawData.byteLength

  return [
    {
      timestamp: dataView.getBigInt64(0),
      data,
    },
    itemSize
  ]
}


// TIMEFRAME
//

const minMetadataRecordSize = 8 + 4 + 4

function calculateMetadataRecordBinarySize(mdRecord) {
  const textEncoder = new TextEncoder()
  return minMetadataRecordSize + textEncoder.encode(mdRecord.channelId).length + textEncoder.encode(mdRecord.channelType).length
}

const minDataItemSize = 8 + 8


export function decodeTimeframe(buffer) {
  let currentPosition = 0
  const header = decodeHeader(buffer)
  currentPosition = headerSize

  const metadataRecords = []

  while (currentPosition < headerSize + header.metadataSize) {
    const metadataRecord = decodeMetadataRecord(buffer.slice(currentPosition))
    currentPosition = currentPosition + calculateMetadataRecordBinarySize(metadataRecord)
    metadataRecords.push(metadataRecord)
  }

  const channels = {}

  metadataRecords.forEach(x => channels[x.channelId] = { type: x.channelType, events: [] })

  metadataRecords.forEach((mdRecord, i) => {
    const isLastDataSection = i == (metadataRecords.length - 1)
    currentPosition = headerSize + header.metadataSize + mdRecord.offsetInData

    let dataSectionEndPosition = 0
    if (isLastDataSection) {
      dataSectionEndPosition = buffer.byteLength
    } else {
      dataSectionEndPosition = headerSize + header.metadataSize + metadataRecords[i + 1].offsetInData
    }

    while (currentPosition < dataSectionEndPosition) {
      const isJSON = mdRecord.channelType === "application/json"
      const [dataItem, dataItemBinarySize] = decodeDataItem(buffer.slice(currentPosition), isJSON)
      currentPosition = currentPosition + dataItemBinarySize
      channels[mdRecord.channelId].events.push(dataItem)
    }
  })


  return {
    startAt: header.startTimestamp,
    endAt: header.endTimestamp,
    channels
  }
}
