package business

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

const maxErrorBodySize = 2048

type BusinessAPIClient struct {
	client  *http.Client
	baseURL string
}

func NewBusinessAPIClient(baseURL string, client *http.Client) *BusinessAPIClient {
	return &BusinessAPIClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *BusinessAPIClient) CheckExists(ctx context.Context, venueID string) (bool, error) {
	reqURL, err := url.JoinPath(c.baseURL, "api", "internal", "venues", venueID)
	if err != nil {
		return false, fmt.Errorf("join url path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("execute request: %w: %w", apperror.ErrUpstreamError, err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	statusCode := resp.StatusCode

	if statusCode == http.StatusOK {
		var dto businessVenueResponseDTO
		if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
			return false, fmt.Errorf("decode response: %w: %w", apperror.ErrUpstreamError, err)
		}
		return dto.Status == "active", nil
	}

	if statusCode == http.StatusNotFound {
		return false, nil
	}

	errBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))

	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		return false, fmt.Errorf("client error %d: %s: %w", statusCode, string(errBody), apperror.ErrUpstreamError)
	}

	return false, fmt.Errorf("server error %d: %s: %w", statusCode, string(errBody), apperror.ErrUpstreamError)
}
