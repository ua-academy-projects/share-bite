package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/client/business/api/client"
	"github.com/ua-academy-projects/share-bite/internal/guest/client/business/api/client/locations"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

type BusinessAPIClient struct {
	apiClient *client.ShareBiteBusinessAPI
}

func NewBusinessAPIClient(apiClient *client.ShareBiteBusinessAPI) *BusinessAPIClient {
	return &BusinessAPIClient{
		apiClient: apiClient,
	}
}

func (c *BusinessAPIClient) CheckExists(ctx context.Context, venueID int64) (bool, error) {
	return true, nil // remove later
	params := locations.NewGetBusinessIDParamsWithContext(ctx).WithID(venueID)

	_, err := c.apiClient.Locations.GetBusinessID(params)
	if err != nil {
		var notFound *locations.GetBusinessIDNotFound
		if errors.As(err, &notFound) {
			return false, nil
		}

		return false, fmt.Errorf("get business by id: %w: %w", apperror.ErrUpstreamError, err)
	}

	return true, nil
}
