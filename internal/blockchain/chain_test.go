package blockchain

import (
	"testing"
	"time"
)

func TestNewBlockchainHasValidGenesis(t *testing.T) {
	bc := NewBlockchain()

	if len(bc.Blocks()) != 1 {
		t.Fatalf("expected 1 block, got %d", len(bc.Blocks()))
	}

	if err := bc.Validate(); err != nil {
		t.Fatalf("expected valid chain, got error: %v", err)
	}
}

func TestAddBlockAndValidate(t *testing.T) {
	bc := NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())

	ts := time.Unix(1700000100, 0).UTC()
	bc.AddBlockAt([]Transaction{
		NewTransaction("tx-1", "alice", "bob", 12.5, ts),
	}, ts)

	if len(bc.Blocks()) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(bc.Blocks()))
	}

	if err := bc.Validate(); err != nil {
		t.Fatalf("expected valid chain, got error: %v", err)
	}
}

func TestValidateDetectsTampering(t *testing.T) {
	bc := NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())

	ts := time.Unix(1700000100, 0).UTC()
	bc.AddBlockAt([]Transaction{
		NewTransaction("tx-1", "alice", "bob", 12.5, ts),
	}, ts)

	blocks := bc.blocks
	blocks[1].Transactions[0].Amount = 999
	bc.blocks = blocks

	if err := bc.Validate(); err == nil {
		t.Fatalf("expected validation error after tampering")
	}
}

func TestValidateDetectsBrokenPrevHash(t *testing.T) {
	bc := NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())

	ts := time.Unix(1700000100, 0).UTC()
	bc.AddBlockAt(nil, ts)

	blocks := bc.blocks
	blocks[1].PrevHash = "broken"
	bc.blocks = blocks

	if err := bc.Validate(); err == nil {
		t.Fatalf("expected validation error for broken previous hash")
	}
}
