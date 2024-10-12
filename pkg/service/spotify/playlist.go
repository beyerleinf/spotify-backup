package spotify

import (
	"beyerleinf/spotify-backup/pkg/request"
	"context"
	"encoding/json"
)

func (s *Service) GetCurrentUserPlaylists() ([]Playlist, error) {
	ctx := context.Background()

	token, err := s.GetAccessToken()
	if err != nil {
		return []Playlist{}, err
	}

	headers := map[string][]string{
		"Authorization": {"Bearer " + token},
	}

	data, _, err := request.Get(ctx, "https://api.spotify.com/v1/me/playlists ", headers)
	if err != nil {
		return []Playlist{}, err
	}

	s.slogger.Verbose("playlists", "lists", data)

	var response GetPlaylistsResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return []Playlist{}, err
	}

	return response.Items, nil
}
