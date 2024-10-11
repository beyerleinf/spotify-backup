package spotify

import (
	"beyerleinf/spotify-backup/internal/config"
	"beyerleinf/spotify-backup/internal/global"
	http_utils "beyerleinf/spotify-backup/pkg/http"
	"beyerleinf/spotify-backup/pkg/models"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var tokenMutex sync.RWMutex

const tokenFile = "token.bin"

var authToken *models.AuthToken

func (s *SpotifyService) GetAuthUrl() string {
	scope := url.QueryEscape("playlist-read-private user-read-private")

	return fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		config.AppConfig.Spotify.ClientId, scope, url.QueryEscape(s.redirectUri), s.state,
	)
}

func (s *SpotifyService) HandleAuthCallback(code string, state string) error {
	if state != s.state {
		return fmt.Errorf("state mismatch")
	}

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", s.redirectUri)

	clientIdAndSecret := fmt.Sprintf("%s:%s", config.AppConfig.Spotify.ClientId, config.AppConfig.Spotify.ClientSecret)
	authHeaderValue := base64.StdEncoding.EncodeToString([]byte(clientIdAndSecret))

	headers := map[string][]string{
		"Authorization": {fmt.Sprintf("Basic %s", authHeaderValue)},
	}

	data, status, err := http_utils.PostForm("https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()), headers)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("token request failed: %d - %s", status, string(data))
	}

	var tokenResponse models.AuthTokenResponse
	err = json.Unmarshal(data, &tokenResponse)
	if err != nil {
		s.slogger.Error("Failed to unmarshal response", "err", err)
		return err
	}

	tokenMutex.Lock()
	authToken = &models.AuthToken{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Second * time.Duration(tokenResponse.ExpiresIn)),
	}
	tokenMutex.Unlock()

	s.saveToken()

	s.slogger.Verbose("Successfully authenticated with Spotify!")

	return nil
}

func (s *SpotifyService) GetAccessToken() (string, error) {
	tokenMutex.RLock()
	token := authToken
	tokenMutex.RUnlock()

	if token == nil {
		s.loadToken()

		tokenMutex.RLock()
		token = authToken
		tokenMutex.RUnlock()
	}

	if token != nil && time.Now().Before(token.ExpiresAt) {
		return token.AccessToken, nil
	}

	if token != nil && time.Now().After(token.ExpiresAt) {
		err := s.RefreshAccessToken(token.RefreshToken)
		if err != nil {
			return "", err
		}

		tokenMutex.RLock()
		token = authToken
		tokenMutex.RUnlock()

		return token.AccessToken, nil
	}

	if token == nil {
		return "", &SpotifyUnauthenticatedError{}
	}

	return "", fmt.Errorf("something went wrong")
}

func (s *SpotifyService) RefreshAccessToken(refreshToken string) error {
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("client_id", config.AppConfig.Spotify.ClientId)

	clientIdAndSecret := fmt.Sprintf("%s:%s", config.AppConfig.Spotify.ClientId, config.AppConfig.Spotify.ClientSecret)
	authHeaderValue := base64.StdEncoding.EncodeToString([]byte(clientIdAndSecret))

	headers := map[string][]string{
		"Authorization": {fmt.Sprintf("Basic %s", authHeaderValue)},
	}

	data, status, err := http_utils.PostForm("https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()), headers)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("token request failed: %d - %s", status, string(data))
	}

	var tokenResponse models.AuthTokenResponse
	err = json.Unmarshal(data, &tokenResponse)
	if err != nil {
		return err
	}

	tokenMutex.Lock()
	authToken = &models.AuthToken{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Second * time.Duration(tokenResponse.ExpiresIn)),
	}
	tokenMutex.Unlock()

	s.saveToken()

	return nil
}

func (s *SpotifyService) saveToken() {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()

	jsonData, err := json.Marshal(authToken)
	if err != nil {
		s.slogger.Error("Error marshaling auth token", "err", err)
		return
	}

	encryptedData, err := s.encryptToken(jsonData)
	if err != nil {
		s.slogger.Error("Error encrypting auth token", "err", err)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		s.slogger.Error("Error getting home directory", "err", err)
	}

	tokenPath := filepath.Join(homeDir, global.StorageDir, tokenFile)

	err = os.WriteFile(tokenPath, encryptedData, 0600)
	if err != nil {
		s.slogger.Error("Error writing encrypted auth token", "err", err)
		return
	}
}

func (s *SpotifyService) loadToken() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		s.slogger.Error("Error getting home directory", "err", err)
	}

	tokenPath := filepath.Join(homeDir, global.StorageDir, tokenFile)

	encryptedData, err := os.ReadFile(tokenPath)
	if err != nil {
		if !os.IsNotExist(err) {
			s.slogger.Error("Error reading encrypted auth token", "err", err)
		}
		return
	}

	decryptedData, err := s.decryptToken(encryptedData)
	if err != nil {
		s.slogger.Error("Error decrypting auth token", "err", err)
		return
	}

	var token models.AuthToken
	err = json.Unmarshal(decryptedData, &token)
	if err != nil {
		s.slogger.Error("Error unmarshaling auth token", "err", err)
		return
	}

	tokenMutex.Lock()
	authToken = &token
	tokenMutex.Unlock()
}

func (s *SpotifyService) encryptToken(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(config.AppConfig.EncryptionKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (s *SpotifyService) decryptToken(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(config.AppConfig.EncryptionKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
