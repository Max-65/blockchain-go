package network

import (
	"encoding/json"
	"errors"
	"net"
	"sync"

	"github.com/Max-65/blockchain-go/internal/blockchain"
)

type Server struct {
	Addr  string
	Chain *blockchain.Blockchain

	mu sync.Mutex
	ln net.Listener
}

func NewServer(addr string, chain *blockchain.Blockchain) *Server {
	return &Server{
		Addr:  addr,
		Chain: chain,
	}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

func (s *Server) Serve(ln net.Listener) error {
	if ln == nil {
		return errors.New("listener is nil")
	}
	if s.Chain == nil {
		return errors.New("chain is nil")
	}

	s.mu.Lock()
	s.ln = ln
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		if s.ln == ln {
			s.ln = nil
		}
		s.mu.Unlock()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		go s.handleConn(conn)
	}
}

func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ln == nil {
		return nil
	}

	err := s.ln.Close()
	s.ln = nil
	return err
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	var req Message
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		return
	}

	enc := json.NewEncoder(conn)

	switch req.Type {
	case MsgGetChain:
		_ = enc.Encode(Message{
			Type:   MsgChain,
			Blocks: s.Chain.Blocks(),
		})

	case MsgPushChain:
		if err := s.Chain.ReplaceIfBetter(req.Blocks); err != nil {
			_ = enc.Encode(Message{
				Type:  MsgError,
				Error: err.Error(),
			})
			return
		}

		_ = enc.Encode(Message{
			Type:   MsgChain,
			Blocks: s.Chain.Blocks(),
		})

	default:
		_ = enc.Encode(Message{
			Type:  MsgError,
			Error: "unknown message type",
		})
	}
}
