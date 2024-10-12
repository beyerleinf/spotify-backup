package spotify

import "time"

// Image is an image in a Spotify API response.
type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// UserProfile represents the logged in users' Spotify profile.
type UserProfile struct {
	ID          string  `json:"string"`
	DisplayName string  `json:"display_name"`
	Images      []Image `json:"images"`
}

// AuthTokenResponse is the response from Spotify's Authentication API.
type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// An AuthToken represents the current credentials.
type AuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GetPlaylistsResponse is the response of the get playlists endpoint
type GetPlaylistsResponse struct {
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Next   string     `json:"next"`
	Total  int        `json:"total"`
	Items  []Playlist `json:"items"`
}

// A Playlist represents a Spotify Playlist.
type Playlist struct {
	ID            string  `json:"id"`
	Collaborative bool    `json:"collaborative"`
	Public        bool    `json:"public"`
	Description   string  `json:"description"`
	Name          string  `json:"name"`
	TotalTracks   int     `json:"tracks.total"`
	Images        []Image `json:"images"`
}
