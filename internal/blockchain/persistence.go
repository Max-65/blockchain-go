package blockchain

import (
	"encoding/json"
	"fmt"

	"github.com/Max-65/blockchain-go/internal/storage"
)

type chainSnapshot struct {
	Blocks []Block `json:"blocks"`
}

func (bc *Blockchain) SaveTo(store storage.Store) error {
	if store == nil {
		return fmt.Errorf("store is nil")
	}

	snapshot := chainSnapshot{
		Blocks: bc.Blocks(),
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	return store.Save(data)
}

func (bc *Blockchain) SaveFile(path string) error {
	return bc.SaveTo(storage.NewFileStore(path))
}

func LoadFrom(store storage.Store) (*Blockchain, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}

	data, err := store.Load()
	if err != nil {
		return nil, err
	}

	var snapshot chainSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, err
	}

	return NewBlockchainFromBlocks(snapshot.Blocks)
}

func LoadBlockchainFile(path string) (*Blockchain, error) {
	return LoadFrom(storage.NewFileStore(path))
}
