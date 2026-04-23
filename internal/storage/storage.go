package storage

import "errors"

var ErrNotFound = errors.New("storage file not found")

type Store interface {
	Save(data []byte) error
	Load() ([]byte, error)
}
