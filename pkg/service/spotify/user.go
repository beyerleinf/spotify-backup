package spotify

import (
	"beyerleinf/spotify-backup/pkg/request"
	"context"
	"encoding/json"
)

// GetUserProfile returns a [UserProfile] from Spotify's API.
// [Get User Profile API]: https://developer.spotify.com/documentation/web-api/reference/get-current-users-profile
func (s *Service) GetUserProfile() (UserProfile, error) {
	ctx := context.Background()

	token, err := s.GetAccessToken()
	if err != nil {
		return UserProfile{}, err
	}

	headers := map[string][]string{
		"Authorization": {"Bearer " + token},
	}

	data, _, err := request.Get(ctx, "https://api.spotify.com/v1/me", headers)
	if err != nil {
		return UserProfile{}, err
	}

	var profile UserProfile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}
