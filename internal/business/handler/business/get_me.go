package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type onboardingContextResponse struct {
	BrandID *int `json:"brandId"`
	VenueID *int `json:"venueId"`
}

// getMe returns the authenticated business user's brand and first venue ids.
// @Summary Get business onboarding context
// @Description Returns brand and venue ids for the authenticated business owner, if they exist.
// @Tags business
// @Produce json
// @Success 200 {object} onboardingContextResponse
// @Failure 401 {object} errorResponse
// @Router /business/me [get]
// @Security BearerAuth
func (h *handler) getMe(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	brandID, venueID, err := h.service.GetOnboardingContext(ctx, userID)
	if err != nil {
		log.Error("failed to get onboarding context", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load business profile"})
		return
	}

	resp := onboardingContextResponse{}
	if brandID > 0 {
		resp.BrandID = &brandID
	}
	if venueID > 0 {
		resp.VenueID = &venueID
	}

	c.JSON(http.StatusOK, resp)
}
