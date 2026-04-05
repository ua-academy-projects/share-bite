package gateway

import (
	"context"
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type businessAPIClient struct {
	client  *http.Client
	baseURL string
}

func NewBusinessAPIClient(baseURL string, client *http.Client) *businessAPIClient {
	return &businessAPIClient{
		client:  client,
		baseURL: baseURL,
	}
}

// TODO: replace this stub with an actual HTTP call to the Business API
// once the endpoint GET /venues is implemented
func (c *businessAPIClient) ListVenues(ctx context.Context, venueIDs []string) (map[string]entity.Venue, error) {
	return nil, nil
}
