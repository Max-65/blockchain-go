package blockchain

import (
	"testing"
	"time"
)

func TestHashBlockChangesWhenTransactionChanges(t *testing.T) {
	ts := time.Unix(1700000000, 0).UTC()

	b1 := NewBlock(1, "prev", []Transaction{
		NewTransaction("tx-1", "alice", "bob", 10, ts),
	}, ts)

	b2 := NewBlock(1, "prev", []Transaction{
		NewTransaction("tx-1", "alice", "bob", 11, ts),
	}, ts)

	if b1.Hash == b2.Hash {
		t.Fatalf("expected hashes to differ when transaction changes")
	}
}

func TestHashBlockStableForSameData(t *testing.T) {
	ts := time.Unix(1700000000, 0).UTC()

	txs := []Transaction{
		NewTransaction("tx-1", "alice", "bob", 10, ts),
		NewTransaction("tx-2", "bob", "carol", 5, ts),
	}

	b1 := NewBlock(2, "abc", txs, ts)
	b2 := NewBlock(2, "abc", txs, ts)

	if b1.Hash != b2.Hash {
		t.Fatalf("expected identical hashes for identical block data")
	}
}
