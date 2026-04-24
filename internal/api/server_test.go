package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/Max-65/blockchain-go/internal/blockchain"
)

func TestCreateBlockAndReadChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.json")

	chain := blockchain.NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())
	srv := NewServer(":0", chain, path, []string{})

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	body := createBlockRequest{
		Transactions: []blockchain.Transaction{
			blockchain.NewTransaction("tx-1", "alice", "bob", 10, time.Unix(1700000100, 0).UTC()),
		},
		Timestamp: time.Unix(1700000100, 0).UTC().Format(time.RFC3339Nano),
	}

	buf, _ := json.Marshal(body)
	resp, err := http.Post(ts.URL+"/blocks", "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	resp2, err := http.Get(ts.URL + "/chain")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp2.StatusCode)
	}

	if chain.Len() != 2 {
		t.Fatalf("expected 2 blocks, got %d", chain.Len())
	}

	if err := chain.Validate(); err != nil {
		t.Fatalf("chain invalid: %v", err)
	}
}
