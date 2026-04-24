package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Max-65/blockchain-go/internal/api"
	"github.com/Max-65/blockchain-go/internal/blockchain"
	"github.com/Max-65/blockchain-go/internal/network"
	"github.com/Max-65/blockchain-go/internal/peers"
	"github.com/Max-65/blockchain-go/internal/storage"
)

func main() {
	tcpAddr := env("NODE_ADDR", ":3000")
	httpAddr := env("HTTP_ADDR", ":8080")
	storePath := env("CHAIN_PATH", "./data/chain.json")
	tcpPort := env("NODE_TCP_PORT", "3000")

	seedPeers := parsePeers(os.Getenv("PEER_SEEDS"))
	peerRegistry := peers.NewRegistry(seedPeers...)

	chain, err := blockchain.LoadBlockchainFile(storePath)
	switch {
	case err == nil:
	case errors.Is(err, storage.ErrNotFound):
		chain = blockchain.NewBlockchain()
	default:
		log.Fatal(err)
	}

	if err := chain.SaveFile(storePath); err != nil {
		log.Fatal(err)
	}

	tcpServer := network.NewServer(tcpAddr, chain)
	httpServer := api.NewServer(httpAddr, chain, storePath, peerRegistry)

	errCh := make(chan error, 2)

	go func() {
		if err := tcpServer.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("tcp server: %w", err)
		}
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()

	go peerGossipLoop(chain, peerRegistry, tcpPort)

	log.Printf("tcp listening on %s", tcpAddr)
	log.Printf("http listening on %s", httpAddr)
	log.Printf("storage: %s", storePath)
	if len(seedPeers) > 0 {
		log.Printf("seed peers: %s", strings.Join(seedPeers, ", "))
	}

	printChain(chain)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-stop:
		log.Printf("shutdown signal: %s", sig)
	case err := <-errCh:
		log.Printf("server stopped: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("http shutdown failed: %v", err)
	}
	if err := tcpServer.Close(); err != nil {
		log.Printf("tcp close failed: %v", err)
	}

	if err := chain.SaveFile(storePath); err != nil {
		log.Printf("save on shutdown failed: %v", err)
	}
}

func peerGossipLoop(chain *blockchain.Blockchain, registry *peers.Registry, tcpPort string) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	// Initial gossip immediately.
	runPeerGossip(chain, registry, tcpPort)

	for range ticker.C {
		runPeerGossip(chain, registry, tcpPort)
	}
}

func runPeerGossip(chain *blockchain.Blockchain, registry *peers.Registry, tcpPort string) {
	snapshot := registry.List()

	for _, peer := range snapshot {
		remotePeers, err := network.ExchangePeers(peer, snapshot, 3*time.Second)
		if err == nil {
			registry.Merge(remotePeers)
		}

		tcpAddr, err := network.TCPAddrFromPeerURL(peer, tcpPort)
		if err == nil {
			_ = network.SyncChain(chain, tcpAddr, 3*time.Second)
		}
	}
}

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

func parsePeers(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func appendUnique(items []string, v string) []string {
	for _, existing := range items {
		if existing == v {
			return items
		}
	}
	return append(items, v)
}

func printChain(chain *blockchain.Blockchain) {
	for _, block := range chain.Blocks() {
		fmt.Printf(
			"index=%d time=%s prev=%s hash=%s txs=%d\n",
			block.Index,
			block.Timestamp.Format(time.RFC3339Nano),
			block.PrevHash,
			block.Hash,
			len(block.Transactions),
		)
	}
}
