import test from 'node:test'
import assert from 'node:assert'

import { headerSize, decodeHeader, decodeMetadataRecord, decodeDataItem, decodeTimeframe } from './timeframe.mjs'

test("Header decoder", (t) => {
  const binaryData = new Uint8Array([
    0x00, 0x01, // version
    0x00, 0x00, 0x00, 0x02, // metadata_size
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // data_size
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // start_timestamp
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05  // end_timestamp
  ])

  const sample = {
    version: 1,
    metadataSize: 2,
    dataSize: 3,
    startTimestamp: 4,
    endTimestamp: 5
  }
  assert.deepEqual(sample, decodeHeader(binaryData.buffer))
})


test("Metadata record decoder", (t) => {
  const binaryData = new Uint8Array([
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // OffsetInData
    0x00, 0x00, 0x00, 0x02, // ChannelIdSize "id" -> 2 byte
    0x00, 0x00, 0x00, 0x04, // ChannelType "type" -> 4 byte
    105, 100, // Id string data "id"
    116, 121, 112, 101, // Type string data "type"
  ])

  const sample = {
    offsetInData: 4,
    channelId: "id",
    channelType: "type",
  }
  assert.deepEqual(sample, decodeMetadataRecord(binaryData.buffer))
})

test("Data Item With Binary Payload Decoder", (t) => {
  const binaryData = new Uint8Array([
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Timestamp
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // record_data_size
	0x01, 0x02, 0x03, 0x04, // data {1, 2, 3, 4} -> 4 byte
  ])

  const sample = {
	timestamp: 1,
	data: new Uint8Array([1, 2, 3, 4]),
  }
  assert.deepEqual(sample, decodeDataItem(binaryData.buffer))
})

test("Data Item With JSON Payload Decoder", (t) => {
  const binaryData = new Uint8Array([
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Timestamp
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 14, // record_data_size
    123, 34, 65, 110, 115, 119, 101, 114, 34, 58, 32, 52, 50, 125,
  ])

  const sample = {
	timestamp: 1,
	data: {Answer: 42},
  }
  assert.deepEqual(sample, decodeDataItem(binaryData.buffer, true))
})

test("Timeframe Decoder", (t) => {
  const binaryData = new Uint8Array([
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

  ])


  const sample = {
	startAt: 1,
	endAt: 10,
    channels: {
      position: {
        type: "number",
        events: [
          {timestamp: 1, data: new Uint8Array([100])},
          {timestamp: 2, data: new Uint8Array([200])},
        ]
      },
      color: {
        type: "number",
        events: [
          {timestamp: 1, data: new Uint8Array([1])},
          {timestamp: 3, data: new Uint8Array([2])},
        ]
      }
    }
  }

  assert.deepEqual(sample, decodeTimeframe(binaryData.buffer))
})
