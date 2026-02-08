package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	Index     int              `json:"index"`
	Timestamp int64            `json:"timestamp"`
	Txs       []DIDTransaction `json:"txs"`
	PrevHash  string           `json:"prev_hash"`
	Hash      string           `json:"hash"`
}

func CalculateHash(block Block) string {
	record := (fmt.Sprintf("%d%d%v%s", block.Index, block.Timestamp, block.Txs, block.PrevHash))
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
