package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SpotifyTokenManager struct {
	ClientID     string
	ClientSecret string
	Token        string
	ExpiresAt    time.Time
	mu           sync.Mutex
}

func (m *SpotifyTokenManager) GetToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if token is still valid (with 30 seconds buffer)
	if time.Now().Before(m.ExpiresAt.Add(-30 * time.Second)) {
		return m.Token, nil
	}

	// Otherwise, fetch new token
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", m.ClientID)
	data.Set("client_secret", m.ClientSecret)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %s", resp.Status)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	// Store the token and expiration
	m.Token = result.AccessToken
	m.ExpiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return m.Token, nil
}

func (m *SpotifyTokenManager) DoSpotifyAPIRequest(url string) (*http.Response, error) {
	token, err := m.GetToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	return client.Do(req)
}
