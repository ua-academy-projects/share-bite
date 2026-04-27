package business

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client/locations"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client/venues"
	business_dto "github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/dto"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

const maxErrorBodySize = 2048

type businessVenueResponseDTO struct {
	VenueID string `json:"venue_id"`
	Status  string `json:"status"`
}

type businessAPIClient struct {
	api    *business_client.ShareBiteBusinessAPI
	scheme string

	// TODO: remove to use swagger autogen files only
	client  *http.Client
	baseURL string
}

func NewBusinessAPIClient(baseURL string, basePath string, httpClient *http.Client) (*businessAPIClient, error) {
	normalizedBaseURL := strings.TrimSpace(baseURL)
	if !strings.Contains(normalizedBaseURL, "://") {
		normalizedBaseURL = "http://" + normalizedBaseURL
	}

	u, err := url.Parse(normalizedBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse business baseURL: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid business baseURL %q: scheme and host are required", normalizedBaseURL)
	}
	if httpClient == nil {
		return nil, fmt.Errorf("http client should be initialized")
	}

	transport := client.NewWithClient(u.Host, basePath, []string{u.Scheme}, httpClient)
	api := business_client.New(transport, strfmt.Default)

	return &businessAPIClient{
		api:     api,
		scheme:  u.Scheme,
		client:  httpClient,
		baseURL: strings.TrimRight(normalizedBaseURL, "/"),
	}, nil
}

func (c *businessAPIClient) CheckExists(ctx context.Context, venueID int64) (bool, error) {
	params := locations.NewGetBusinessOrgUnitsIDParamsWithContext(ctx).WithID(venueID)

	_, err := c.api.Locations.GetBusinessOrgUnitsID(params, schemeLocationClientOption(c.scheme))
	if err == nil {
		return true, nil
	}

	var notFoundErr *locations.GetBusinessOrgUnitsIDNotFound
	if errors.As(err, &notFoundErr) {
		return false, nil
	}

	statusCode := 0
	type withStatusCode interface {
		Code() int
	}
	var codeErr withStatusCode
	if errors.As(err, &codeErr) {
		statusCode = codeErr.Code()
	}

	if statusCode == 0 {
		return false, fmt.Errorf("execute request: %w: %w", apperror.ErrUpstreamError, err)
	}

	errMsg := ""
	var badRequestErr *locations.GetBusinessOrgUnitsIDBadRequest
	if errors.As(err, &badRequestErr) {
		if payload := badRequestErr.GetPayload(); payload != nil {
			errMsg = strings.TrimSpace(payload.Error)
		}
	}

	var internalErr *locations.GetBusinessOrgUnitsIDInternalServerError
	if errMsg == "" && errors.As(err, &internalErr) {
		if payload := internalErr.GetPayload(); payload != nil {
			errMsg = strings.TrimSpace(payload.Error)
		}
	}

	if errMsg == "" {
		errMsg = strings.TrimSpace(err.Error())
	}
	if errMsg == "" {
		errMsg = apperror.ErrUpstreamError.Error()
	}
	if len(errMsg) > maxErrorBodySize {
		errMsg = errMsg[:maxErrorBodySize]
	}

	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		return false, fmt.Errorf("client error %d: %s: %w", statusCode, errMsg, apperror.ErrUpstreamError)
	}

	return false, fmt.Errorf("server error %d: %s: %w", statusCode, errMsg, apperror.ErrUpstreamError)
}

func (c *businessAPIClient) ListVenuesByIDs(ctx context.Context, venueIDs []int64) (map[int64]entity.Venue, error) {
	if len(venueIDs) == 0 {
		return map[int64]entity.Venue{}, nil
	}
	venueIDs = removeDuplicates(venueIDs)

	payload := &business_dto.HandlerBusinessGetVenuesByIDsRequest{
		Ids: venueIDs,
	}
	params := venues.NewPostBusinessOrgUnitsVenuesParamsWithContext(ctx).WithRequest(payload)

	resp, err := c.api.Venues.PostBusinessOrgUnitsVenues(params, schemeClientOption(c.scheme))
	if err != nil {
		logger.ErrorKV(ctx, "get venues from business service",
			"error", err,
			"venue_ids_count", len(venueIDs),
		)

		return nil, fmt.Errorf("get venues by IDs: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("get venues by IDs unexpected status code: %d", resp.Code())
	}

	respPayload := resp.GetPayload()
	if respPayload == nil {
		return map[int64]entity.Venue{}, nil
	}

	out := make(map[int64]entity.Venue, len(respPayload))
	for _, v := range respPayload {
		if v == nil {
			continue
		}

		out[v.ID] = entity.Venue{
			ID:          v.ID,
			Name:        v.Name,
			Description: toNilStrPtr(v.Description),
			AvatarURL:   toNilStrPtr(v.Avatar),
			BannerURL:   toNilStrPtr(v.Banner),
		}
	}

	return out, nil
}

func toNilStrPtr(v string) *string {
	if v == "" {
		return nil
	}

	val := v
	return &val
}

func removeDuplicates(in []int64) []int64 {
	seen := make(map[int64]struct{}, len(in))

	out := make([]int64, 0, len(in))
	for _, id := range in {
		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}
		out = append(out, id)
	}

	return out
}

func schemeClientOption(scheme string) venues.ClientOption {
	return func(op *runtime.ClientOperation) {
		op.Schemes = []string{scheme}
	}
}

func schemeLocationClientOption(scheme string) locations.ClientOption {
	return func(op *runtime.ClientOperation) {
		op.Schemes = []string{scheme}
	}
}

func (c *businessAPIClient) GetNearbyVenues(ctx context.Context, lat, lon float64, limit int) ([]int64, error) {
	params := locations.NewGetBusinessLocationsNearbyParamsWithContext(ctx).
		WithLat(lat).
		WithLon(lon).
		WithLimit((*int64)(&[]int64{int64(limit)}[0]))

	resp, err := c.api.Locations.GetBusinessLocationsNearby(params, schemeLocationClientOption(c.scheme))
	if err != nil {
		return nil, fmt.Errorf("failed to call business api: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("get nearby venues unexpected status code: %d", resp.Code())
	}

	payload := resp.GetPayload()
	if payload == nil || len(payload.Items) == 0 {
		return []int64{}, nil
	}

	venueIDs := make([]int64, 0, len(payload.Items))
	for _, item := range payload.Items {
		venueIDs = append(venueIDs, int64(item.ID))
	}

	return venueIDs, nil
}
