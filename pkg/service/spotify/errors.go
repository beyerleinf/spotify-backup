package spotify

// A UnauthenticatedError is returned when Authentication
// with the Spotify API failed.
type UnauthenticatedError struct{}

func (e *UnauthenticatedError) Error() string {
	return "Authentication with Spotify failed! Try signing into your Account again."
}
