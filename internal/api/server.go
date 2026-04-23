package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Max-65/blockchain-go/internal/blockchain"
	"github.com/Max-65/blockchain-go/internal/network"
)

type Server struct {
	Addr      string
	Chain     *blockchain.Blockchain
	StorePath string
	PeerAddr  string

	srv *http.Server
}

func NewServer(addr string, chain *blockchain.Blockchain, storePath, peerAddr string) *Server {
	return &Server{
		Addr:      addr,
		Chain:     chain,
		StorePath: storePath,
		PeerAddr:  peerAddr,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/chain", s.handleChain)
	mux.HandleFunc("/blocks", s.handleBlocks)
	mux.HandleFunc("/sync", s.handleSync)
	return mux
}

func (s *Server) ListenAndServe() error {
	s.srv = &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler(),
	}

	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	if s.Chain == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "chain is nil"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"length": len(s.Chain.Blocks()),
		"blocks": s.Chain.Blocks(),
	})
}

type createBlockRequest struct {
	Transactions []blockchain.Transaction `json:"transactions"`
	Timestamp    string                   `json:"timestamp,omitempty"`
}

func (s *Server) handleBlocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	if s.Chain == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "chain is nil"})
		return
	}

	var req createBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	ts := time.Now().UTC()
	if req.Timestamp != "" {
		parsed, err := time.Parse(time.RFC3339Nano, req.Timestamp)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid timestamp format"})
			return
		}
		ts = parsed.UTC()
	}

	block := s.Chain.AddBlockAt(req.Transactions, ts)

	if s.StorePath != "" {
		if err := s.Chain.SaveFile(s.StorePath); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"block":  block,
		"length": s.Chain.Len(),
	})
}

type syncRequest struct {
	PeerAddr string `json:"peer_addr,omitempty"`
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	if s.Chain == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "chain is nil"})
		return
	}

	var req syncRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	peer := req.PeerAddr
	if peer == "" {
		peer = s.PeerAddr
	}
	if peer == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "peer_addr is required"})
		return
	}

	if err := network.SyncChain(s.Chain, peer, 3*time.Second); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if s.StorePath != "" {
		if err := s.Chain.SaveFile(s.StorePath); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"length": s.Chain.Len(),
		"blocks": s.Chain.Blocks(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(payload)
}

func (s *Server) String() string {
	return fmt.Sprintf("api server on %s", s.Addr)
}
