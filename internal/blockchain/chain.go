package blockchain

import "time"

type Chain struct {
	Blocks []Block
	Forks  [][]Block
}

func NewChain() *Chain {
	genesis := Block{
		Index:     0,
		Timestamp: time.Now().Unix(),
		PrevHash:  "",
	}
	genesis.Hash = CalculateHash(genesis)

	return &Chain{
		Blocks: []Block{genesis},
	}
}

func (c *Chain) LastBlock() Block {
	return c.Blocks[len(c.Blocks)-1]
}

func (c *Chain) AddBlock(b Block) {
	lastBlock := c.LastBlock()

	if b.Index != lastBlock.Index+1 {
		return
	}

	if b.PrevHash == lastBlock.Hash {
		c.Blocks = append(c.Blocks, b)
	} else {
		c.Forks = append(c.Forks, []Block{b})
	}
}
