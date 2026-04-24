package peers

import (
	"net/url"
	"sort"
	"strings"
	"sync"
)

type Registry struct {
	mu    sync.RWMutex
	peers map[string]struct{}
}

func NewRegistry(seeds ...string) *Registry {
	r := &Registry{
		peers: make(map[string]struct{}),
	}
	r.Merge(seeds)
	return r
}

func (r *Registry) Add(peer string) bool {
	peer = normalize(peer)
	if peer == "" {
		return false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.peers == nil {
		r.peers = make(map[string]struct{})
	}

	if _, exists := r.peers[peer]; exists {
		return false
	}

	r.peers[peer] = struct{}{}
	return true
}

func (r *Registry) Remove(peer string) bool {
	peer = normalize(peer)
	if peer == "" {
		return false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.peers == nil {
		return false
	}

	if _, exists := r.peers[peer]; !exists {
		return false
	}

	delete(r.peers, peer)
	return true
}

func (r *Registry) Merge(peers []string) int {
	added := 0
	for _, peer := range peers {
		if r.Add(peer) {
			added++
		}
	}
	return added
}

func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.peers))
	for peer := range r.peers {
		out = append(out, peer)
	}

	sort.Strings(out)
	return out
}

func normalize(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Host == "" {
		return ""
	}

	return u.Scheme + "://" + u.Host
}
