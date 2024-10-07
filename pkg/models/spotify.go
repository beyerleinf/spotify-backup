package models

type SpotifyImage struct {
	Url    string `json:"url"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

type SpotifyUserProfile struct {
	Id          string         `json:"string"`
	DisplayName string         `json:"display_name"`
	Images      []SpotifyImage `json:"images"`
}

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token+type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
