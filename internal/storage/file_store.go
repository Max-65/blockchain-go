package storage

import (
	"errors"
	"os"
	"path/filepath"
)

type FileStore struct {
	Path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{Path: path}
}

func (fs *FileStore) Save(data []byte) error {
	if fs.Path == "" {
		return errors.New("storage path is empty")
	}

	dir := filepath.Dir(fs.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, filepath.Base(fs.Path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}

	if _, err := tmp.Write(data); err != nil {
		cleanup()
		return err
	}

	if err := tmp.Sync(); err != nil {
		cleanup()
		return err
	}

	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}

	if err := os.Rename(tmpName, fs.Path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}

	return nil
}

func (fs *FileStore) Load() ([]byte, error) {
	if fs.Path == "" {
		return nil, errors.New("storage path is empty")
	}

	data, err := os.ReadFile(fs.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return data, nil
}
