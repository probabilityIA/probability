package handlers

import (
	"sync"
	"time"
)

type oAuthStateData struct {
	IntegrationName string
	UserID          uint
	BusinessID      uint
	IsTesting       bool
	Expiry          time.Time
}

var (
	oauthStateStore = make(map[string]*oAuthStateData)
	oauthStateMutex sync.Mutex
)

func storeOAuthState(state string, data *oAuthStateData) {
	oauthStateMutex.Lock()
	defer oauthStateMutex.Unlock()
	now := time.Now()
	for k, v := range oauthStateStore {
		if now.After(v.Expiry) {
			delete(oauthStateStore, k)
		}
	}
	oauthStateStore[state] = data
}

func consumeOAuthState(state string) (*oAuthStateData, bool) {
	oauthStateMutex.Lock()
	defer oauthStateMutex.Unlock()
	data, exists := oauthStateStore[state]
	if !exists || time.Now().After(data.Expiry) {
		delete(oauthStateStore, state)
		return nil, false
	}
	delete(oauthStateStore, state)
	return data, true
}

type tokenExchangeData struct {
	AccessToken  string
	RefreshToken string
	IsTesting    bool
	ExpiresAt    time.Time
	Expiry       time.Time
}

var (
	tokenStore = make(map[string]tokenExchangeData)
	tokenMutex sync.Mutex
)

func storeExchangeToken(token string, data tokenExchangeData) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	now := time.Now()
	for k, v := range tokenStore {
		if now.After(v.Expiry) {
			delete(tokenStore, k)
		}
	}
	tokenStore[token] = data
}

func retrieveExchangeToken(token string) (tokenExchangeData, bool) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	data, exists := tokenStore[token]
	if !exists || time.Now().After(data.Expiry) {
		delete(tokenStore, token)
		return tokenExchangeData{}, false
	}
	delete(tokenStore, token)
	return data, true
}
