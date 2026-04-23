package storage

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestFileStoreSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.json")

	store := NewFileStore(path)

	want := []byte(`{"hello":"world"}`)
	if err := store.Save(want); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if string(got) != string(want) {
		t.Fatalf("unexpected data:\nwant: %s\ngot:  %s", want, got)
	}
}

func TestFileStoreLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")

	store := NewFileStore(path)

	_, err := store.Load()
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}
