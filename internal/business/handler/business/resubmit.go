package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type resubmitRequest struct {
	ID int `uri:"id" binding:"required"`
}

type resubmitResponse struct {
	Message string `json:"message" example:"Business verification request resubmitted successfully."`
}

// resubmitVerification updates organization unit status back to pending.
//
//	@Summary      Resubmit business verification
//	@Description  Allows a business user to resubmit a previously rejected organization unit back to pending status.
//	@Tags         org-units
//	@Produce      json
//	@Param        id    path      int    true   "Business Org Unit ID"
//	@Success      200   {object}  resubmitResponse
//	@Failure      400   {object}  errorResponse
//	@Failure      401   {object}  errorResponse
//	@Failure      500   {object}  errorResponse
//	@Router       /business/org-units/{id}/resubmit [post]
func (h *handler) resubmitVerification(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	req := new(resubmitRequest)
	if err := c.ShouldBindUri(req); err != nil {
		_ = c.Error(err)
		return
	}

	userUUID, err := h.extractUserUUID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	log.Info("resubmitting business verification", "id", req.ID, "userId", userUUID.String())

	err = h.service.ResubmitVerification(ctx, req.ID, userUUID.String())
	if err != nil {
		log.Error("failed to resubmit business verification", "id", req.ID, "userId", userUUID.String(), "error", err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resubmitResponse{
		Message: "Business verification request resubmitted successfully.",
	})
}
