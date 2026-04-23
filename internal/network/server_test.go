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

func TestFetchChain(t *testing.T) {
	chain := blockchain.NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())
	chain.AddBlockAt(nil, time.Unix(1700000100, 0).UTC())
	chain.AddBlockAt(nil, time.Unix(1700000200, 0).UTC())

	addr, cleanup := startTestServer(t, chain)
	defer cleanup()

	var got []blockchain.Block
	var err error

	for i := 0; i < 20; i++ {
		got, err = FetchChain(addr, 250*time.Millisecond)
		if err == nil {
			break
		}
		time.Sleep(25 * time.Millisecond)
	}

	if err != nil {
		t.Fatalf("fetch failed: %v", err)
	}

	if len(got) != chain.Len() {
		t.Fatalf("unexpected block count: want %d got %d", chain.Len(), len(got))
	}

	want := chain.Blocks()
	for i := range want {
		if got[i].Hash != want[i].Hash {
			t.Fatalf("hash mismatch at index %d", i)
		}
	}
}

func TestSyncChain(t *testing.T) {
	remote := blockchain.NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())
	remote.AddBlockAt(nil, time.Unix(1700000100, 0).UTC())
	remote.AddBlockAt(nil, time.Unix(1700000200, 0).UTC())

	addr, cleanup := startTestServer(t, remote)
	defer cleanup()

	local := blockchain.NewBlockchainWithGenesis(time.Unix(1700000000, 0).UTC())

	if err := SyncChain(local, addr, 250*time.Millisecond); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	if local.Len() != remote.Len() {
		t.Fatalf("unexpected length after sync: want %d got %d", remote.Len(), local.Len())
	}

	if err := local.Validate(); err != nil {
		t.Fatalf("synced chain is invalid: %v", err)
	}
}
