package client

import (
	"sync"
	"time"
)

type tokenEntry struct {
	accessToken string
	expiresAt   time.Time
}

type TokenStore struct {
	entries map[string]tokenEntry
	mu      sync.RWMutex
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		entries: make(map[string]tokenEntry),
	}
}

func tokenKey(username, accountID, partnerID, baseURL string) string {
	return username + "|" + accountID + "|" + partnerID + "|" + baseURL
}

func (ts *TokenStore) Get(username, accountID, partnerID, baseURL string) (string, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	key := tokenKey(username, accountID, partnerID, baseURL)
	entry, ok := ts.entries[key]
	if !ok || entry.accessToken == "" || time.Now().After(entry.expiresAt) {
		return "", false
	}
	return entry.accessToken, true
}

func (ts *TokenStore) Set(username, accountID, partnerID, baseURL, accessToken string, expiresIn int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	buffer := 1800
	ttl := expiresIn - buffer
	if ttl <= 0 {
		ttl = expiresIn / 2
	}

	key := tokenKey(username, accountID, partnerID, baseURL)
	ts.entries[key] = tokenEntry{
		accessToken: accessToken,
		expiresAt:   time.Now().Add(time.Duration(ttl) * time.Second),
	}
}

func (ts *TokenStore) Clear(username, accountID, partnerID, baseURL string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	key := tokenKey(username, accountID, partnerID, baseURL)
	delete(ts.entries, key)
}
