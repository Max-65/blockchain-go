package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Max-65/blockchain-go/internal/blockchain"
)

func FetchChain(addr string, timeout time.Duration) ([]blockchain.Block, error) {
	resp, err := request(addr, timeout, Message{Type: MsgGetChain})
	if err != nil {
		return nil, err
	}

	if resp.Type == MsgError {
		return nil, fmt.Errorf("%s", resp.Error)
	}

	if resp.Type != MsgChain {
		return nil, fmt.Errorf("unexpected response type: %s", resp.Type)
	}

	return resp.Blocks, nil
}

func PushChain(addr string, blocks []blockchain.Block, timeout time.Duration) error {
	resp, err := request(addr, timeout, Message{
		Type:   MsgPushChain,
		Blocks: blocks,
	})
	if err != nil {
		return err
	}

	if resp.Type == MsgError {
		return fmt.Errorf("%s", resp.Error)
	}

	if resp.Type != MsgChain {
		return fmt.Errorf("unexpected response type: %s", resp.Type)
	}

	return nil
}

func PushBlock(addr string, block blockchain.Block, timeout time.Duration) error {
	resp, err := request(addr, timeout, Message{
		Type:  MsgPushBlock,
		Block: &block,
	})
	if err != nil {
		return err
	}

	if resp.Type == MsgError {
		return fmt.Errorf("%s", resp.Error)
	}

	if resp.Type != MsgBlock {
		return fmt.Errorf("unexpected response type: %s", resp.Type)
	}

	return nil
}

func SyncChain(chain *blockchain.Blockchain, addr string, timeout time.Duration) error {
	blocks, err := FetchChain(addr, timeout)
	if err != nil {
		return err
	}

	return chain.ReplaceIfBetter(blocks)
}

func request(addr string, timeout time.Duration, req Message) (Message, error) {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return Message{}, err
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return Message{}, err
	}

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return Message{}, err
	}

	var resp Message
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return Message{}, err
	}

	return resp, nil
}
