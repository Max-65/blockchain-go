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
	"github.com/Max-65/blockchain-go/internal/storage"
)

func main() {
	tcpAddr := env("NODE_ADDR", ":3000")
	httpAddr := env("HTTP_ADDR", ":8080")
	storePath := env("CHAIN_PATH", "./data/chain.json")

	peers := parsePeers(os.Getenv("PEERS"))
	if singlePeer := strings.TrimSpace(os.Getenv("PEER_ADDR")); singlePeer != "" {
		peers = appendUnique(peers, singlePeer)
	}

	chain, err := blockchain.LoadBlockchainFile(storePath)
	switch {
	case err == nil:
	case errors.Is(err, storage.ErrNotFound):
		chain = blockchain.NewBlockchain()
	default:
		log.Fatal(err)
	}

	for _, peer := range peers {
		if err := network.SyncChain(chain, peer, 3*time.Second); err != nil {
			log.Printf("startup sync from %s failed: %v", peer, err)
		}
	}

	if err := chain.SaveFile(storePath); err != nil {
		log.Fatal(err)
	}

	tcpServer := network.NewServer(tcpAddr, chain)
	httpServer := api.NewServer(httpAddr, chain, storePath, peers)

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

	log.Printf("tcp listening on %s", tcpAddr)
	log.Printf("http listening on %s", httpAddr)
	log.Printf("storage: %s", storePath)
	if len(peers) > 0 {
		log.Printf("peers: %s", strings.Join(peers, ", "))
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
			out = appendUnique(out, p)
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
