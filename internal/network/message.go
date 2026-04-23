package network

import "github.com/Max-65/blockchain-go/internal/blockchain"

const (
	MsgGetChain  = "get_chain"
	MsgPushChain = "push_chain"
	MsgChain     = "chain"
	MsgError     = "error"
)

type Message struct {
	Type   string             `json:"type"`
	Blocks []blockchain.Block `json:"blocks,omitempty"`
	Error  string             `json:"error,omitempty"`
}
