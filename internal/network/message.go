package network

import "github.com/Max-65/blockchain-go/internal/blockchain"

const (
	MsgGetChain  = "get_chain"
	MsgPushChain = "push_chain"
	MsgPushBlock = "push_block"

	MsgChain = "chain"
	MsgBlock = "block"
	MsgError = "error"
)

type Message struct {
	Type   string             `json:"type"`
	Blocks []blockchain.Block `json:"blocks,omitempty"`
	Block  *blockchain.Block  `json:"block,omitempty"`
	Error  string             `json:"error,omitempty"`
}
