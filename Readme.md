# Consistent Hashing Lab

A simple implementation of consistent hashing in Go, demonstrating hash rings, virtual nodes and key redistribution.

## Overview

Consistent hashing is a technique for distributing keys across multiple servers while minimizing the number of keys that need to move when servers are added or removed.

Traditional modulo-based hashing commonly uses:

`hash(key) % numberOfServers`

Changing the number of servers can remap many keys. Consistent hashing places servers and keys on a logical hash ring, assigning each key to the next server position clockwise.

## Features

- Hash ring implementation using Go
- Deterministic key-to-server lookup
- Configurable virtual nodes
- Server addition and removal
- Empty-ring and duplicate-server handling
- Unit tests for core behavior
- Comparison of modulo hashing and consistent hashing

## Requirements

- Go 

## Getting Started

Clone the repository:

```bash
git clone https://github.com/praveshdev3/consistent-hashing-lab.git
cd consistent-hashing-lab
```

Run the tests:

```bash
go test ./...
```

Run the demonstration:

```bash
go run .
```

## How It Works

### 1. Hash Ring

Each server is assigned one or more positions on a logical ring using a deterministic hash function.

### 2. Virtual Nodes

Each physical server can own multiple virtual positions. This generally improves key distribution compared with assigning only one position per server.

### 3. Key Lookup

The key is hashed to a ring position. The algorithm selects the first server position at or clockwise from that position, wrapping to the beginning of the ring if necessary.

### 4. Server Changes

When a server joins or leaves, only keys in the affected ring intervals need to be reassigned. The actual movement depends on the hash positions and the number of virtual nodes.

## Experiments

The project compares key movement when a server is added using:

- Modulo hashing
- Consistent hashing with virtual nodes

It also tests the removal of a server and checks that keys are no longer assigned to that server.

| Metric | Result |
|---|---|
| Number of keys tested | 100,000 |
| Initial servers | 3 |
| Final servers | 4 |
| Virtual nodes per server | 100 |
| Modulo hashing movement | 75137 (75.14%) |
| Consistent hashing movement | 28658 (28.66%) |

## Project Structure

```text
consistent-hashing-lab/
├── main.go
├── ring.go
├── ring_test.go
├── go.mod
└── README.md
```

## Testing

Run all tests:

```bash
go test -v ./...
```

Run Go's static analysis checks:

```bash
go vet ./...
```

## Limitations

This is an educational implementation, not a production-ready distributed storage system.

- Hash collisions require explicit handling.
- Concurrent access requires synchronization.
- Server membership changes are currently managed in memory.
- Virtual-node counts affect distribution and performance.
- No network communication or data replication is implemented.

## Future Improvements

- Benchmark different virtual-node counts
- Visualize the hash ring
- Add a command-line interface for experiments
- Compare key distribution across multiple configurations
- Add concurrency support
- Improve atomicity of server membership changes