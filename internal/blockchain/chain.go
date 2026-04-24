package blockchain

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrEmptyChain      = errors.New("blockchain is empty")
	ErrInvalidHash     = errors.New("invalid block hash")
	ErrBrokenLink      = errors.New("broken chain link")
	ErrInvalidBlockSeq = errors.New("invalid block index sequence")
)

type Blockchain struct {
	mu     sync.RWMutex
	blocks []Block
}

func NewBlockchain() *Blockchain {
	return NewBlockchainWithGenesis(time.Unix(0, 0).UTC())
}

func NewBlockchainWithGenesis(timestamp time.Time) *Blockchain {
	return &Blockchain{
		blocks: []Block{
			NewGenesisBlock(timestamp),
		},
	}
}

func NewBlockchainFromBlocks(blocks []Block) (*Blockchain, error) {
	if len(blocks) == 0 {
		return nil, ErrEmptyChain
	}

	bc := &Blockchain{
		blocks: cloneBlockSlice(blocks),
	}

	if err := bc.Validate(); err != nil {
		return nil, err
	}

	return bc, nil
}

func (bc *Blockchain) Blocks() []Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return cloneBlockSlice(bc.blocks)
}

func (bc *Blockchain) Len() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return len(bc.blocks)
}

func (bc *Blockchain) LastBlock() (Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.blocks) == 0 {
		return Block{}, ErrEmptyChain
	}

	return bc.blocks[len(bc.blocks)-1], nil
}

func (bc *Blockchain) AddBlock(transactions []Transaction) Block {
	return bc.AddBlockAt(transactions, time.Now().UTC())
}

func (bc *Blockchain) AddBlockAt(transactions []Transaction, timestamp time.Time) Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.blocks) == 0 {
		genesis := NewGenesisBlock(timestamp)
		bc.blocks = append(bc.blocks, genesis)
		return genesis
	}

	prev := bc.blocks[len(bc.blocks)-1]
	block := NewBlock(prev.Index+1, prev.Hash, transactions, timestamp)
	bc.blocks = append(bc.blocks, block)

	return block
}

func (bc *Blockchain) TryAppendBlock(block Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.blocks) == 0 {
		if block.Index != 0 {
			return fmt.Errorf("%w: expected genesis block", ErrInvalidBlockSeq)
		}
		if HashBlock(block) != block.Hash {
			return ErrInvalidHash
		}
		bc.blocks = append(bc.blocks, cloneBlock(block))
		return nil
	}

	last := bc.blocks[len(bc.blocks)-1]

	if block.Index != last.Index+1 {
		return fmt.Errorf("%w: got %d, want %d", ErrInvalidBlockSeq, block.Index, last.Index+1)
	}

	if block.PrevHash != last.Hash {
		return ErrBrokenLink
	}

	if HashBlock(block) != block.Hash {
		return ErrInvalidHash
	}

	bc.blocks = append(bc.blocks, cloneBlock(block))
	return nil
}

func (bc *Blockchain) ReplaceIfBetter(blocks []Block) error {
	if len(blocks) == 0 {
		return ErrEmptyChain
	}

	candidate, err := NewBlockchainFromBlocks(blocks)
	if err != nil {
		return err
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(candidate.blocks) <= len(bc.blocks) {
		return nil
	}

	bc.blocks = cloneBlockSlice(candidate.blocks)
	return nil
}

func (bc *Blockchain) Validate() error {
	bc.mu.RLock()
	blocks := cloneBlockSlice(bc.blocks)
	bc.mu.RUnlock()

	if len(blocks) == 0 {
		return ErrEmptyChain
	}

	for i, block := range blocks {
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

		prev := blocks[i-1]

		if block.Index != prev.Index+1 {
			return fmt.Errorf("%w at index %d: got %d, want %d", ErrInvalidBlockSeq, i, block.Index, prev.Index+1)
		}

		if block.PrevHash != prev.Hash {
			return fmt.Errorf("%w at index %d", ErrBrokenLink, i)
		}
	}

	return nil
}

func cloneBlock(b Block) Block {
	b.Transactions = cloneTransactions(b.Transactions)
	return b
}

func cloneBlockSlice(src []Block) []Block {
	if len(src) == 0 {
		return nil
	}

	dst := make([]Block, len(src))
	copy(dst, src)

	for i := range dst {
		dst[i].Transactions = cloneTransactions(dst[i].Transactions)
	}

	return dst
}
