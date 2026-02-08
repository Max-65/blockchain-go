package network

import (
	"blockchaingo/internal/blockchain"
	"sync"
)

type Node struct {
	ID    string
	Addr  string
	Peers []string

	Chain *blockchain.Chain
	mu    sync.Mutex
}
