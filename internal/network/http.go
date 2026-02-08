package network

import (
	"blockchaingo/internal/blockchain"
	"encoding/json"
	"net/http"
)

func (n *Node) Start() error {
	http.HandleFunc("/block", n.HandleBlock)
	http.HandleFunc("/chain", n.HandleChain)
	return http.ListenAndServe(n.Addr, nil)
}

func (n *Node) HandleBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var block blockchain.Block
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.Chain.AddBlock(block)
	w.WriteHeader(http.StatusOK)
}

func (n *Node) HandleChain(w http.ResponseWriter, r *http.Request) {
	n.mu.Lock()
	defer n.mu.Unlock()

	json.NewEncoder(w).Encode(n.Chain.Blocks)
}
