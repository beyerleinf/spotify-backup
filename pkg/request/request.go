package request

import (
	"context"
	"io"
	"net/http"
)

// Post sends a POST request.
func Post(ctx context.Context, url string, body io.Reader, headers map[string][]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header = headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return data, res.StatusCode, nil
}

// PostForm sends a POST request with a application/x-www-form-urlencoded body.
func PostForm(ctx context.Context, url string, body io.Reader, headers map[string][]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header = headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return data, res.StatusCode, nil
}

// Get sends a GET request.
func Get(ctx context.Context, url string, headers map[string][]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header = headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return data, res.StatusCode, nil
}
