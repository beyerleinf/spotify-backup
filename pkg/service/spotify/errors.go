package spotify

type SpotifyUnauthenticatedError struct{}

func (e *SpotifyUnauthenticatedError) Error() string {
	return "Authentication with Spotify failed! Try signing into your Account again."
}
