package business

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client/venues"
	business_dto "github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/dto"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

const maxErrorBodySize = 2048

type businessAPIClient struct {
	api    *business_client.ShareBiteBusinessAPI
	scheme string

	// TODO: remove to use swagger autogen files only
	client  *http.Client
	baseURL string
}

func NewBusinessAPIClient(baseURL string, basePath string, httpClient *http.Client) (*businessAPIClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse business baseURL: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid business baseURL %q: scheme and host are required", baseURL)
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
		baseURL: baseURL,
	}, nil
}

func (c *businessAPIClient) CheckExists(ctx context.Context, venueID string) (bool, error) {
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
