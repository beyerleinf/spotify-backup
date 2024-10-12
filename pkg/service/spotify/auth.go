package spotify

import (
	"beyerleinf/spotify-backup/pkg/assert"
	"beyerleinf/spotify-backup/pkg/request"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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

var authToken *AuthToken

// GetAuthURL returns a URL to redirect a user to sign in with Spotify.
func (s *Service) GetAuthURL() string {
	scope := url.QueryEscape("playlist-read-private user-read-private")

	return fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		s.config.Spotify.ClientID, scope, url.QueryEscape(s.redirectURI), s.state,
	)
}

// HandleAuthCallback handles a callback request from Spotify's Auth API.
// It takes a code and the state used to initiate the authentication flow
// and follows Spotify's requirements to request an Access Token.
// [Spotify Authorization Code Flow]: https://developer.spotify.com/documentation/web-api/tutorials/code-flow
func (s *Service) HandleAuthCallback(code string, state string) error {
	ctx := context.Background()

	if state != s.state {
		return errors.New("state mismatch")
	}

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", s.redirectURI)

	clientIDAndSecret := fmt.Sprintf("%s:%s", s.config.Spotify.ClientID, s.config.Spotify.ClientSecret)
	authHeaderValue := base64.StdEncoding.EncodeToString([]byte(clientIDAndSecret))

	headers := map[string][]string{
		"Authorization": {"Basic " + authHeaderValue},
	}

	data, status, err := request.PostForm(ctx, "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()), headers)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("token request failed: %d - %s", status, string(data))
	}

	var tokenResponse AuthTokenResponse
	err = json.Unmarshal(data, &tokenResponse)
	if err != nil {
		s.slogger.Error("Failed to unmarshal response", "err", err)
		return err
	}

	tokenMutex.Lock()
	authToken = &AuthToken{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    s.calculateExpiresAt(tokenResponse.ExpiresIn),
	}
	tokenMutex.Unlock()

	s.saveToken()

	s.slogger.Verbose("Successfully authenticated with Spotify!")

	return nil
}

// GetAccessToken tries to read the current Access Token from an encrypted file
// on disk. If that fails or of the Access Token expired, it will request
// a new Access Token using [RefreshAccessToken].
func (s *Service) GetAccessToken() (string, error) {
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
		assert.NotEqual("", token.AccessToken, "existing access token should not be an empty string")
		assert.NotNil(token.AccessToken, "existing access token is nil")

		return token.AccessToken, nil
	}

	if token != nil && time.Now().After(token.ExpiresAt) {
		assert.NotEqual("", token.RefreshToken, "stored refresh token should not be an empty string")

		err := s.RefreshAccessToken(token.RefreshToken)
		if err != nil {
			return "", err
		}

		tokenMutex.RLock()
		token = authToken
		tokenMutex.RUnlock()

		assert.NotEqual("", token.AccessToken, "new access token should not be an empty string")
		assert.NotNil(token.AccessToken, "new access token is nil")

		return token.AccessToken, nil
	}

	if token == nil {
		return "", &UnauthenticatedError{}
	}

	assert.Assert(true, "We should not have reached this")
	return "", errors.New("something went wrong")
}

// RefreshAccessToken makes a call to Spotify's Authentication API using
// the Refresh Token obtained on the last authentication request.
// It will request a new Access Token using the Refresh Token.
// [Refreshing Tokens]: https://developer.spotify.com/documentation/web-api/tutorials/refreshing-tokens
func (s *Service) RefreshAccessToken(refreshToken string) error {
	assert.NotEqual("", refreshToken, "RefreshToken should not be an empty string")

	ctx := context.Background()

	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("client_id", s.config.Spotify.ClientID)

	clientIDAndSecret := fmt.Sprintf("%s:%s", s.config.Spotify.ClientID, s.config.Spotify.ClientSecret)
	authHeaderValue := base64.StdEncoding.EncodeToString([]byte(clientIDAndSecret))

	headers := map[string][]string{
		"Authorization": {"Basic " + authHeaderValue},
	}

	data, status, err := request.PostForm(ctx, "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()), headers)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("token request failed: %d - %s", status, string(data))
	}

	var tokenResponse AuthTokenResponse
	err = json.Unmarshal(data, &tokenResponse)
	if err != nil {
		return err
	}

	assert.NotEqual("", tokenResponse.AccessToken, "AccessToken should not be empty")
	assert.NotNil(tokenResponse.ExpiresIn, "ExpiresAt should not be nil")

	tokenMutex.Lock()

	authToken.AccessToken = tokenResponse.AccessToken
	authToken.ExpiresAt = s.calculateExpiresAt(tokenResponse.ExpiresIn)

	if tokenResponse.RefreshToken != "" {
		authToken.RefreshToken = tokenResponse.RefreshToken
	}

	tokenMutex.Unlock()

	s.saveToken()

	return nil
}

func (s *Service) saveToken() {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()

	assert.NotEqual("", authToken.AccessToken, "AccessToken should not be empty")
	assert.NotNil(authToken.ExpiresAt, "ExpiresAt should not be nil")
	assert.NotEqual("", authToken.RefreshToken, "RefreshToken should not be empty")

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

	tokenPath := filepath.Join(s.storageDir, tokenFile)

	err = os.WriteFile(tokenPath, encryptedData, 0600)
	if err != nil {
		s.slogger.Error("Error writing encrypted auth token", "err", err)
		return
	}
}

func (s *Service) loadToken() {
	tokenPath := filepath.Join(s.storageDir, tokenFile)

	_, err := os.Stat(tokenPath)
	if err != nil {
		return
	}

	encryptedData, err := os.ReadFile(tokenPath)
	if err != nil {
		if !os.IsNotExist(err) {
			s.slogger.Error("Error reading encrypted auth token", "err", err)
		}
		return
	}

	assert.NotEqual(0, len(encryptedData), "token file exists but is empty")

	decryptedData, err := s.decryptToken(encryptedData)
	if err != nil {
		s.slogger.Error("Error decrypting auth token", "err", err)
		return
	}

	var token AuthToken
	err = json.Unmarshal(decryptedData, &token)
	if err != nil {
		s.slogger.Error("Error unmarshaling auth token", "err", err)
		return
	}

	assert.NotEqual("", token.AccessToken, "AccessToken should not be empty")
	assert.NotNil(token.ExpiresAt, "ExpiresAt should not be nil")
	assert.NotEqual("", token.RefreshToken, "RefreshToken should not be empty")

	tokenMutex.Lock()
	authToken = &token
	tokenMutex.Unlock()
}

func (s *Service) encryptToken(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(s.config.EncryptionKey))
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

func (s *Service) decryptToken(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(s.config.EncryptionKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (s *Service) calculateExpiresAt(expiresIn int) time.Time {
	return time.Now().Add(time.Second * time.Duration(expiresIn))
}
