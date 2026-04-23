package blockchain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrEmptyChain      = errors.New("blockchain is empty")
	ErrInvalidHash     = errors.New("invalid block hash")
	ErrBrokenLink      = errors.New("broken chain link")
	ErrInvalidBlockSeq = errors.New("invalid block index sequence")
)

type Blockchain struct {
	blocks []Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []Block{
			NewGenesisBlock(time.Unix(0, 0).UTC()),
		},
	}
}

func NewBlockchainWithGenesis(timestamp time.Time) *Blockchain {
	return &Blockchain{
		blocks: []Block{
			NewGenesisBlock(timestamp),
		},
	}
}

func (bc *Blockchain) Blocks() []Block {
	out := make([]Block, len(bc.blocks))
	copy(out, bc.blocks)
	return out
}

func (bc *Blockchain) LastBlock() (Block, error) {
	if len(bc.blocks) == 0 {
		return Block{}, ErrEmptyChain
	}
	return bc.blocks[len(bc.blocks)-1], nil
}

func (bc *Blockchain) AddBlock(transactions []Transaction) Block {
	return bc.AddBlockAt(transactions, time.Now().UTC())
}

func (bc *Blockchain) AddBlockAt(transactions []Transaction, timestamp time.Time) Block {
	prev, err := bc.LastBlock()
	if err != nil {
		genesis := NewGenesisBlock(timestamp)
		bc.blocks = append(bc.blocks, genesis)
		return genesis
	}

	block := NewBlock(prev.Index+1, prev.Hash, transactions, timestamp)
	bc.blocks = append(bc.blocks, block)
	return block
}

func (bc *Blockchain) Validate() error {
	if len(bc.blocks) == 0 {
		return ErrEmptyChain
	}

	for i, block := range bc.blocks {
		expectedHash := HashBlock(block)
		if block.Hash != expectedHash {
			return fmt.Errorf("%w at index %d", ErrInvalidHash, i)
		}

		if i == 0 {
			if block.Index != 0 {
				return fmt.Errorf("%w at index %d: genesis index must be 0", ErrInvalidBlockSeq, i)
			}
			continue
		}

		prev := bc.blocks[i-1]

		if block.Index != prev.Index+1 {
			return fmt.Errorf("%w at index %d: got %d, want %d", ErrInvalidBlockSeq, i, block.Index, prev.Index+1)
		}

		if block.PrevHash != prev.Hash {
			return fmt.Errorf("%w at index %d", ErrBrokenLink, i)
		}
	}

	return nil
}
