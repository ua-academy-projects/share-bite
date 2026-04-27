package business

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client/locations"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/business_client/venues"
	business_dto "github.com/ua-academy-projects/share-bite/internal/guest/gateway/business/client/dto"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

const maxErrorBodySize = 2048

type businessAPIClient struct {
	api    *business_client.ShareBiteBusinessAPI
	scheme string
	policy *resilience.Policy

	// TODO: remove to use swagger autogen files only
	client  *http.Client
	baseURL string
}

type Option func(*businessAPIClient)

func WithResiliencePolicy(policy resilience.Policy) Option {
	return func(c *businessAPIClient) {
		c.policy = &policy
	}
}

func NewBusinessAPIClient(baseURL string, basePath string, httpClient *http.Client, opts ...Option) (*businessAPIClient, error) {
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

	out := &businessAPIClient{
		api:     api,
		scheme:  u.Scheme,
		client:  httpClient,
		baseURL: strings.TrimRight(normalizedBaseURL, "/"),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(out)
		}
	}

	return out, nil
}

func (c *businessAPIClient) CheckExists(ctx context.Context, venueID int64) (bool, error) {
	params := locations.NewGetBusinessOrgUnitsIDParamsWithContext(ctx).WithID(venueID)

	exists := false
	err := c.executeWithResilience(ctx, func() error {
		_, opErr := c.api.Locations.GetBusinessOrgUnitsID(params, schemeLocationClientOption(c.scheme))
		if opErr == nil {
			exists = true
			return nil
		}

		var notFoundErr *locations.GetBusinessOrgUnitsIDNotFound
		if errors.As(opErr, &notFoundErr) {
			exists = false
			return nil
		}

		mappedErr, retryable := mapCheckExistsError(opErr)
		if !retryable {
			return resilience.Permanent(mappedErr)
		}

		return mappedErr
	})
	if err != nil {
		return false, err
	}

	return exists, nil
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

	var resp *venues.PostBusinessOrgUnitsVenuesOK
	err := c.executeWithResilience(ctx, func() error {
		apiResp, opErr := c.api.Venues.PostBusinessOrgUnitsVenues(params, schemeClientOption(c.scheme))
		if opErr != nil {
			logger.ErrorKV(ctx, "get venues from business service",
				"error", opErr,
				"venue_ids_count", len(venueIDs),
			)

			wrappedErr := fmt.Errorf("get venues by IDs: %w", opErr)
			if !isRetryableError(opErr) {
				return resilience.Permanent(wrappedErr)
			}

			return wrappedErr
		}

		if !apiResp.IsSuccess() {
			wrappedErr := fmt.Errorf("get venues by IDs unexpected status code: %d", apiResp.Code())
			if !isRetryableStatus(apiResp.Code()) {
				return resilience.Permanent(wrappedErr)
			}

			return wrappedErr
		}

		resp = apiResp
		return nil
	})
	if err != nil {
		return nil, err
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

func (c *businessAPIClient) executeWithResilience(ctx context.Context, operation func() error) error {
	if c.policy == nil {
		return operation()
	}

	return c.policy.Execute(ctx, operation)
}

func mapCheckExistsError(err error) (error, bool) {
	statusCode := extractStatusCode(err)
	if statusCode == 0 {
		return fmt.Errorf("execute request: %w: %w", apperror.ErrUpstreamError, err), true
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
		return fmt.Errorf("client error %d: %s: %w", statusCode, errMsg, apperror.ErrUpstreamError), isRetryableStatus(statusCode)
	}

	return fmt.Errorf("server error %d: %s: %w", statusCode, errMsg, apperror.ErrUpstreamError), true
}

func extractStatusCode(err error) int {
	type withStatusCode interface {
		Code() int
	}

	var codeErr withStatusCode
	if errors.As(err, &codeErr) {
		return codeErr.Code()
	}

	return 0
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && (netErr.Timeout() || netErr.Temporary()) {
		return true
	}
	if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, io.EOF) {
		return true
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}
	statusCode := extractStatusCode(err)
	if statusCode == 0 {
		return true
	}

	return isRetryableStatus(statusCode)
}

func isRetryableStatus(statusCode int) bool {
	if statusCode == http.StatusRequestTimeout || statusCode == http.StatusTooManyRequests {
		return true
	}

	return statusCode >= http.StatusInternalServerError
}
