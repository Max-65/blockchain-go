# Minimal Blockchain Prototype in Go

A small blockchain prototype written in Go.

This project is intentionally simple. It focuses on the core data structures and validation logic first, without adding networking, storage, mining, or other extra pieces too early.

## Purpose

The goal of this repository is to build a minimal blockchain core that is:

- easy to understand;
- easy to test;
- easy to extend later;
- small enough to evolve step by step.

## What is included

- blocks with an index, timestamp, transactions, previous hash, and current hash;
- a simple transaction model;
- an in-memory blockchain;
- chain validation;
- unit tests for hashing, block addition, and chain integrity;
- GitHub Actions CI;
- a Dockerfile for containerized execution.

## What is not included yet

This project does not currently implement:

- peer-to-peer networking;
- proof-of-work mining;
- disk persistence;
- wallets or addresses;
- digital signatures;
- mempool;
- forks or chain selection rules;
- Merkle trees.

These parts may be added later, but they are not part of the first working version.