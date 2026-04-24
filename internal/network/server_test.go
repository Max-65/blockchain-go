package network

import (
	"net"
	"testing"
	"time"

	"github.com/Max-65/blockchain-go/internal/blockchain"
)

func startTestServer(t *testing.T, chain *blockchain.Blockchain) (addr string, cleanup func()) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}

	srv := NewServer("", chain)
	done := make(chan struct{})

	go func() {
		defer close(done)
		_ = srv.Serve(ln)
	}()

	cleanup = func() {
		_ = srv.Close()
		<-done
	}

	return ln.Addr().String(), cleanup
}

func TestPushBlock(t *testing.T) {
	chain := blockchain.NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())
	addr, cleanup := startTestServer(t, chain)
	defer cleanup()

	block := blockchain.NewBlock(
		1,
		chain.Blocks()[0].Hash,
		[]blockchain.Transaction{
			blockchain.NewTransaction("tx-1", "alice", "bob", 10, time.Unix(1700000100, 0).UTC()),
		},
		time.Unix(1700000100, 0).UTC(),
	)

	if err := PushBlock(addr, block, 250*time.Millisecond); err != nil {
		t.Fatalf("push block failed: %v", err)
	}

	if chain.Len() != 2 {
		t.Fatalf("expected chain length 2, got %d", chain.Len())
	}

	if err := chain.Validate(); err != nil {
		t.Fatalf("chain invalid after push: %v", err)
	}
}

func TestPushBlockRejectsBadBlock(t *testing.T) {
	chain := blockchain.NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())
	addr, cleanup := startTestServer(t, chain)
	defer cleanup()

	bad := blockchain.NewBlock(
		5,
		"wrong-prev-hash",
		nil,
		time.Unix(1700000100, 0).UTC(),
	)

	if err := PushBlock(addr, bad, 250*time.Millisecond); err == nil {
		t.Fatalf("expected push block to fail")
	}

	if chain.Len() != 1 {
		t.Fatalf("chain should not change after bad block")
	}
}
