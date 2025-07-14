# Peer-to-Peer Distributed File Server

A distributed file storage system built in Go that allows for secure peer-to-peer file storage, retrieval, and distribution across a network.

## Overview

This project implements a peer-to-peer distributed file server where files are stored across multiple nodes in the network. Each node can act as both a client and a server, allowing for decentralized file storage and retrieval. The system includes content-addressable storage (CAS) for efficient file management and basic peer discovery mechanisms.

## Features

- **Peer-to-Peer Architecture**: Nodes can connect directly to each other without a central server
- **Content-Addressable Storage**: Files are stored based on their content hash for efficient retrieval
- **File Replication**: Files are automatically replicated across the network for redundancy
- **Bootstrap Mechanism**: New nodes can join the network by connecting to known bootstrap nodes
- **TCP Transport Layer**: Reliable communication between nodes

## Installation

Clone the repository and build the project:

```bash
git clone https://github.com/aniketh3014/gstore
cd gstore
make build
```

## Usage

### Starting Nodes

Start a node by running the built binary:

```bash
make run
```

By default, this will start two nodes on your local machine:
- One node listening on port 4000
- Another node listening on port 5000, connecting to the first node

### Custom Configuration

To customize the node configuration, you can modify the `main.go` file to specify different:
- Port numbers
- Storage directories
- Bootstrap node addresses

### Storing and Retrieving Data

The file server provides functions for storing and retrieving data:

```go
// Store data
server.StoreData(key, dataReader)

// Retrieve data
size, reader, err := server.GetData(key)
```

## Project Structure

- `main.go`: Entry point and demo setup
- `fileserver.go`: Core file server implementation
- `store.go`: Storage management and content-addressable storage
- `crypto.go`: Cryptographic utilities
- `p2p/`: Peer-to-peer networking components
  - `tcp_transport.go`: TCP-based transport layer
  - `transport.go`: Transport interface definitions
  - `message.go`: Message format definitions
  - `encoding.go`: Encoding/decoding utilities
  - `handshaker.go`: Connection handshake protocols

## Running Tests

```bash
make test
```

## Architecture

The system is built around these main components:

1. **File Server**: Manages file storage, retrieval, and network operations
2. **Store**: Handles the local file system operations and implements content-addressable storage
3. **Transport Layer**: Manages network communications between nodes
4. **Peer Management**: Handles peer discovery and connection management

### Message Types

The system uses several message types for communication:
- `MessageStoreFile`: Request to store a file
- `MessageGetFile`: Request to retrieve a file

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Future Improvements

- DHT (Distributed Hash Table) implementation for better node discovery
- File encryption and secure transfer
- Authentication and access control mechanisms
- Bandwidth optimization and file chunking
- Web interface for file management
