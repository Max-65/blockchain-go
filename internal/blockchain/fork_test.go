package blockchain

import (
	"testing"
	"time"
)

func TestReplaceIfBetterUsesTieBreakerAtEqualLength(t *testing.T) {
	baseTime := time.Unix(1700000000, 0).UTC()

	local := NewBlockchainWithGenesis(baseTime)
	local.AddBlockAt([]Transaction{
		NewTransaction("tx-local-1", "alice", "bob", 10, baseTime.Add(time.Minute)),
	}, baseTime.Add(time.Minute))

	candidate := NewBlockchainWithGenesis(baseTime)
	candidate.AddBlockAt([]Transaction{
		NewTransaction("tx-candidate-1", "alice", "bob", 11, baseTime.Add(time.Minute)),
	}, baseTime.Add(time.Minute))

	if len(local.Blocks()) != len(candidate.Blocks()) {
		t.Fatalf("setup error: chains must have equal length")
	}

	localDigestBefore := local.ChainDigest()
	candidateDigest := candidate.ChainDigest()

	if err := local.ReplaceIfBetter(candidate.Blocks()); err != nil {
		t.Fatalf("replace failed: %v", err)
	}

	localDigestAfter := local.ChainDigest()

	switch {
	case candidateDigest < localDigestBefore:
		if localDigestAfter != candidateDigest {
			t.Fatalf("expected candidate chain to replace local chain")
		}
	case candidateDigest > localDigestBefore:
		if localDigestAfter != localDigestBefore {
			t.Fatalf("expected local chain to remain unchanged")
		}
	default:
		if localDigestAfter != localDigestBefore {
			t.Fatalf("expected stable result on identical digest")
		}
	}
}
