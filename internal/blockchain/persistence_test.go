package blockchain

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/Max-65/blockchain-go/internal/storage"
)

func TestBlockchainSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.json")

	bc := NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())
	bc.AddBlockAt([]Transaction{
		NewTransaction("tx-1", "alice", "bob", 10, time.Unix(1700000100, 0).UTC()),
	}, time.Unix(1700000100, 0).UTC())

	if err := bc.SaveFile(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadBlockchainFile(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(loaded.Blocks()) != len(bc.Blocks()) {
		t.Fatalf("unexpected block count: want %d got %d", len(bc.Blocks()), len(loaded.Blocks()))
	}

	if err := loaded.Validate(); err != nil {
		t.Fatalf("loaded chain is invalid: %v", err)
	}

	for i, block := range bc.Blocks() {
		if loaded.Blocks()[i].Hash != block.Hash {
			t.Fatalf("block %d hash mismatch", i)
		}
	}
}

func TestLoadBlockchainFileMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")

	_, err := LoadBlockchainFile(path)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected storage.ErrNotFound, got: %v", err)
	}
}
