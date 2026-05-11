package business

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type searchVenuesRequest struct {
	Query string `form:"q"`
	Tags  string `form:"tags"`
	Skip  int    `form:"skip"`
	Limit int    `form:"limit"`
}

// searchVenues searches venues by keyword and/or tags.
//
//	@Summary		Search venues
//	@Description	Search venues by keyword (`q`) in name/description and optional tags (`tags=tag1,tag2`).
//	@Tags			venues
//	@Produce		json
//	@Param			q		query		string	false	"Keyword for name/description"
//	@Param			tags	query		string	false	"Comma-separated location tag slugs"
//	@Param			skip	query		int		false	"Number of items to skip (default: 0)"
//	@Param			limit	query		int		false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	listResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/venues/search [get]
func (h *handler) searchVenues(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var req searchVenuesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query params"})
		return
	}

	if req.Skip < 0 {
		req.Skip = 0
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit < 0 {
		req.Limit = 1
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	var tags []string
	if req.Tags != "" {
		seen := make(map[string]struct{})
		for _, tag := range strings.Split(req.Tags, ",") {
			tag = strings.ToLower(strings.TrimSpace(tag))
			if tag == "" {
				continue
			}
			if _, ok := seen[tag]; ok {
				continue
			}
			seen[tag] = struct{}{}
			tags = append(tags, tag)
		}
	}

	log.Info("search venues", "q", req.Query, "tags", tags, "skip", req.Skip, "limit", req.Limit)

	query := strings.TrimSpace(req.Query)
	if query == "" && len(tags) == 0 {
		c.Error(apperror.BadRequest("at least one search filter is required: q or tags"))
		return
	}

	result, err := h.service.SearchVenues(ctx, query, req.Skip, req.Limit, tags)
	if err != nil {
		log.Error("failed to search venues", "error", err)
		c.Error(err)
		return
	}

	items := make([]listItem, 0, len(result.Items))
	for _, u := range result.Items {
		items = append(items, listItem{
			ID:          u.Id,
			Name:        u.Name,
			Avatar:      u.Avatar,
			Description: u.Description,
			Latitude:    u.Latitude,
			Longitude:   u.Longitude,
			Tags:        normalizeTags(u.Tags),
		})
	}

	c.JSON(http.StatusOK, listResponse{
		Items: items,
		Total: result.Total,
	})
}
