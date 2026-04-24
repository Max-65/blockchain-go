package peers

import "testing"

func TestRegistryNormalizesAndDeduplicates(t *testing.T) {
	r := NewRegistry("node1:8080", "http://node1:8080/", "https://node2:8080")

	got := r.List()
	if len(got) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(got))
	}

	if got[0] == "" || got[1] == "" {
		t.Fatalf("empty peer in registry")
	}
}
