package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PeersMessage struct {
	Peers []string `json:"peers"`
}

func ExchangePeers(peer string, localPeers []string, timeout time.Duration) ([]string, error) {
	endpoint, err := peerEndpoint(peer)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(PeersMessage{Peers: localPeers})
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: timeout}
	if timeout <= 0 {
		client.Timeout = 3 * time.Second
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var msg PeersMessage
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		if len(msg.Peers) > 0 {
			return msg.Peers, fmt.Errorf("peer exchange failed with status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("peer exchange failed with status %d", resp.StatusCode)
	}

	var out PeersMessage
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out.Peers, nil
}

func TCPAddrFromPeerURL(peer string, tcpPort string) (string, error) {
	if tcpPort == "" {
		tcpPort = "3000"
	}

	u, err := normalizePeerURL(peer)
	if err != nil {
		return "", err
	}

	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("invalid peer host")
	}

	return net.JoinHostPort(host, tcpPort), nil
}

func peerEndpoint(raw string) (string, error) {
	u, err := normalizePeerURL(raw)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(u.String(), "/") + "/peers", nil
}

func normalizePeerURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("peer is empty")
	}

	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Host == "" {
		return nil, fmt.Errorf("peer host is empty")
	}

	u.Path = strings.TrimRight(u.Path, "/")
	return u, nil
}
