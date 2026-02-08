package main

import (
	"blockchaingo/internal/blockchain"
	"blockchaingo/internal/network"
	"log"
	"os"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	node := &network.Node{
		ID:    "node-1",
		Addr:  addr,
		Chain: blockchain.NewChain(),
	}

	log.Fatal(node.Start())
}
