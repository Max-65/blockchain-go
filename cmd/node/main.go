package main

import (
	"fmt"
	"time"

	"github.com/Max-65/blockchain-go/internal/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	now := time.Now().UTC()

	bc.AddBlockAt([]blockchain.Transaction{
		blockchain.NewTransaction("tx-1", "alice", "bob", 10, now),
	}, now)

	bc.AddBlockAt([]blockchain.Transaction{
		blockchain.NewTransaction("tx-2", "bob", "carol", 3.5, now.Add(time.Minute)),
	}, now.Add(time.Minute))

	if err := bc.Validate(); err != nil {
		panic(err)
	}

	for _, block := range bc.Blocks() {
		fmt.Printf(
			"index=%d time=%s prev=%s hash=%s txs=%d\n",
			block.Index,
			block.Timestamp.Format(time.RFC3339Nano),
			block.PrevHash,
			block.Hash,
			len(block.Transactions),
		)
	}
}
