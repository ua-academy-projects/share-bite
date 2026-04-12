package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/client/business/api/client"
	"github.com/ua-academy-projects/share-bite/internal/guest/client/business/api/client/venues"
	"github.com/ua-academy-projects/share-bite/internal/guest/client/business/api/models"
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
	params := venues.NewPostBusinessOrgUnitsVenuesParamsWithContext(ctx).WithRequest(&models.BusinessGetVenuesByIDsRequest{
		Ids: []int64{venueID},
	})

	resp, err := c.apiClient.Venues.PostBusinessOrgUnitsVenues(params)
	if err != nil {
		var badRequest *venues.PostBusinessOrgUnitsVenuesBadRequest
		if errors.As(err, &badRequest) {
			return false, nil
		}

		type statusCoder interface {
			Code() int
		}

		var coder statusCoder
		if errors.As(err, &coder) && coder.Code() == 404 {
			return false, nil
		}

		return false, fmt.Errorf("post business org units venues: %w: %w", apperror.ErrUpstreamError, err)
	}

	if resp == nil || len(resp.Payload) == 0 {
		return false, nil
	}

	for _, venue := range resp.Payload {
		if venue != nil && venue.ID == venueID {
			return true, nil
		}
	}

	return false, nil
}
