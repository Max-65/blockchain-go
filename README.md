# Minimal Blockchain Prototype in Go

A small blockchain prototype written in Go.

The goal of this repository is to keep the core logic compact and testable. The project starts with a single-node chain, then adds persistence, network sync, HTTP control endpoints, peer discovery, and Docker-based verification.

## What this project does

- stores blocks in memory;
- validates chain integrity;
- persists the chain to a JSON file;
- syncs chains between nodes over TCP;
- exposes a small HTTP API for inspection and control;
- exchanges peer lists so nodes can discover each other from a seed set;
- runs in Docker and Docker Compose;
- includes unit tests and a simple end-to-end script.

## What is intentionally left out for now

This prototype does not yet include:

- proof-of-work mining;
- wallets or addresses;
- digital signatures;
- mempool / transaction relay;
- Merkle trees;
- fork choice rules beyond a simple deterministic tie-breaker;
- a full P2P protocol;
- a database engine.

The project is kept small on purpose so each layer can be added and tested separately.

## Current architecture

The system is split into a few small parts:

- `internal/blockchain` for blocks, chain validation, persistence hooks, and chain selection;
- `internal/storage` for file-based JSON storage;
- `internal/network` for TCP chain sync and peer exchange;
- `internal/peers` for peer registry and normalization;
- `internal/api` for the HTTP control plane;
- `cmd/node/main.go` for wiring everything together.

## HTTP endpoints

The node exposes a small HTTP API:

- `GET /health`
- `GET /chain`
- `POST /blocks`
- `POST /sync`
- `GET /peers`
- `POST /peers`

## TCP behavior

The TCP side is used for chain transfer between nodes. A node can request a full chain, push a full chain, or push a single block with fallback to full-chain sync when needed.

## Persistence

The current storage layer writes the chain to a JSON file on disk. On startup the node tries to load the chain from that file. If the file does not exist, a new chain is created.

## Peer discovery

The node starts from a small seed list and exchanges peer lists over HTTP. This keeps the topology from being fully static while staying simple enough to reason about.