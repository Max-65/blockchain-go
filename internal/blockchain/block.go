package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Transaction struct {
	ID        string  `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
}

type Block struct {
	Index        int           `json:"index"`
	Timestamp    time.Time     `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prev_hash"`
	Hash         string        `json:"hash"`
}

type blockHashPayload struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prev_hash"`
}

func NewTransaction(id, from, to string, amount float64, timestamp time.Time) Transaction {
	return Transaction{
		ID:        id,
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: timestamp.UTC().UnixNano(),
	}
}

func NewBlock(index int, prevHash string, transactions []Transaction, timestamp time.Time) Block {
	b := Block{
		Index:        index,
		Timestamp:    timestamp.UTC(),
		Transactions: cloneTransactions(transactions),
		PrevHash:     prevHash,
	}
	b.Hash = HashBlock(b)
	return b
}

func NewGenesisBlock(timestamp time.Time) Block {
	return NewBlock(0, "", nil, timestamp)
}

func HashBlock(b Block) string {
	payload := blockHashPayload{
		Index:        b.Index,
		Timestamp:    b.Timestamp.UTC().UnixNano(),
		Transactions: cloneTransactions(b.Transactions),
		PrevHash:     b.PrevHash,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func cloneTransactions(src []Transaction) []Transaction {
	if len(src) == 0 {
		return nil
	}
	dst := make([]Transaction, len(src))
	copy(dst, src)
	return dst
}
