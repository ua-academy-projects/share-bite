package business

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

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
	address := normalizeDialAddress(c.baseURL)
	url := fmt.Sprintf("http://%s/api/internal/venues/%s", address, venueID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, errwrap.Wrap("create request", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, errwrap.Wrap("execute request", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dto businessVenueResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return false, errwrap.Wrap("decode response", err)
	}

	isActive := dto.Status == "active"

	return isActive, nil
}

func normalizeDialAddress(address string) string {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return address
	}

	if host == "0.0.0.0" || host == "::" || host == "" {
		host = "localhost"
	}

	return net.JoinHostPort(host, port)
}
