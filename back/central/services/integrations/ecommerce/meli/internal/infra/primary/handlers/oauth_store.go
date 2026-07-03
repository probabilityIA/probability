package handlers

import (
	"sync"
	"time"
)

type OAuthStateData struct {
	IntegrationName string
	UserID          uint
	BusinessID      uint
	CodeVerifier    string
	IsTesting       bool
	Expiry          time.Time
}

var (
	oauthStateStore = make(map[string]*OAuthStateData)
	oauthStateMutex sync.Mutex
)

func storeOAuthState(state string, data *OAuthStateData) {
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

func consumeOAuthState(state string) (*OAuthStateData, bool) {
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

type TokenExchangeData struct {
	AccessToken  string
	RefreshToken string
	ClientID     string
	ClientSecret string
	SellerID     int64
	IsTesting    bool
	ExpiresAt    time.Time
	Expiry       time.Time
}

var (
	tokenStore = make(map[string]TokenExchangeData)
	tokenMutex sync.RWMutex
)

func storeExchangeToken(token string, data TokenExchangeData) {
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

func retrieveExchangeToken(token string) (TokenExchangeData, bool) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	data, exists := tokenStore[token]
	if !exists || time.Now().After(data.Expiry) {
		delete(tokenStore, token)
		return TokenExchangeData{}, false
	}
	delete(tokenStore, token)
	return data, true
}
