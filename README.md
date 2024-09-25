# SyncStreamer

SyncStreamer is an open-source project written in Go that provides efficient streaming of synchronized data in a time-based format. It offers both JSON and binary data streaming over HTTP, allowing data to be transmitted in a flexible timeframe format.

## Table of Contents
- [Overview](#overview)
- [TimeFrame Format Specification](#timeframe-format-specification)
- [API Endpoints](#api-endpoints)
  - [Index Endpoint](#index-endpoint)
  - [GetFrame Endpoint](#getframe-endpoint)
  - [Channel POST Endpoint](#channel-post-endpoint)
- [Building and Running](#building-and-running)
- [Contributing](#contributing)
- [License](#license)

## Overview

SyncStreamer is designed to efficiently transmit data that is organized based on timeframes. It uses a custom binary format to handle metadata and data sections, enabling precise synchronization of channels and data streams.

The project exposes an HTTP API for fetching and submitting data frames. It supports both JSON and binary data formats, making it flexible for various use cases, such as media synchronization, sensor data streaming, and time-series data.

## TimeFrame Format Specification

The SyncStreamer system uses a custom **TimeFrame** format to organize and transmit data. Below is a detailed explanation of this format, including the structure of the header, metadata, and data sections.

### Section Header

The header contains basic information about the data and metadata size, and the timeframe during which the data was recorded.

- **Version (2 octets)**: The version number of the TimeFrame format.
- **Metadata Size (4 octets)**: The size of the metadata section in bytes.
- **Data Size (8 octets)**: The size of the data section in bytes.
- **Start Timestamp (8 octets)**: The UNIX timestamp in milliseconds representing when the timeframe started.
- **End Timestamp (8 octets)**: The UNIX timestamp in milliseconds representing when the timeframe ended.

### Section Metadata

Each metadata record describes a data channel within the timeframe.

- **Offset in Data (8 octets)**: The offset position of this record's data in the data section.
- **Channel Type Size (4 octets)**: The length of the channel type string (in bytes).
- **Channel ID Size (4 octets)**: The length of the channel ID string (in bytes).
- **Channel Type (chan_type_size octets)**: A string representing the type of the channel (e.g., "audio", "video").
- **Channel ID (chan_id_size octets)**: A string representing the unique ID of the channel.

### Section Data

The data section contains multiple `data_item` records, where each record represents a piece of data with a timestamp.

- **Timestamp (8 octets)**: The UNIX timestamp in milliseconds for the data item.
- **Record Data Size (8 octets)**: The size of the data record in bytes.
- **Record Data (record_data_size octets)**: The actual data for the record.

Each `data_item` is an atomic unit of data within the timeframe, and multiple items are used to represent the entire stream of data.

## API Endpoints

SyncStreamer provides a simple HTTP-based API for accessing and transmitting data.

### Index Endpoint

- **URL**: `GET [host]/frame`
- **Method**: `GET`
- **Response Format**: `application/json`
- **Response**: Returns an array of available `TimeFrameItem` objects in JSON format. Each `TimeFrameItem` contains:
  - `startAt`: The start time (in milliseconds) of the frame.
  - `endAt`: The end time (in milliseconds) of the frame.
  - `id`: The unique identifier for the frame.

#### TimeFrameItem JSON Structure:

```json
{
  "startAt": float64,
  "endAt": float64,
  "id": "string"
}
```

**Note**: Caching is not allowed on this endpoint as the available frames may change frequently.

### GetFrame Endpoint

- **URL**: `GET [host]/{frameId}`
- **Method**: `GET`
- **Response Format**: `application/octet-stream`
- **Response**: Retrieves a specific frame identified by the `frameId` in binary format (following the **TimeFrame Format Specification**).

This endpoint should be cached by the client as the frame content is static once generated.

### Channel POST Endpoint

- **URL**: `POST [host]/channel/{id}`
- **Method**: `POST`
- **Request Format**: Either `application/json` or `application/octet-stream`
- **Request Body**: The body can either be a JSON or binary data stream that corresponds to the specified channel ID.

This endpoint allows clients to submit data to a specific channel within a timeframe.

## Building and Running

To build and run the SyncStreamer server locally:

1. Clone the repository:
   ```sh
   git clone https://github.com/darkmou5e/syncstreamer.git
   cd syncstreamer
   ```

2. Build the Go binary:
   ```sh
   go build -o syncstreamer
   ```

3. Run the server:
   ```sh
   ./syncstreamer
   ```

4. The server will start and listen for HTTP requests on the configured port (default: `8080`).

### Environment Variables (not yet)

- `PORT`: Define the port on which the server should listen (default is `8080`).

## Contributing

Contributions are welcome! Please submit a pull request or open an issue if you encounter any bugs or have suggestions for improvements.

1. Fork the repository.
2. Create a new feature branch.
3. Commit your changes.
4. Open a pull request to the main repository.

Before submitting, make sure your code follows the existing style and passes all tests.

## License

SyncStreamer is licensed under the GNU GPL 2 License. See the [LICENSE](LICENSE) file for more information.
