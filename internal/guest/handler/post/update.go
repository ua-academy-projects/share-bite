package post

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

type updateRequest struct {
	VenueID *int64                  `form:"venue_id" binding:"omitempty"`
	Text    *string                 `form:"text" binding:"omitempty,max=2000"`
	Rating  *int16                  `form:"rating" binding:"omitempty,min=1,max=5"`
	Status  *entity.PostStatus      `form:"status" binding:"omitempty,oneof=draft published archived"`
	Images  []*multipart.FileHeader `form:"images" binding:"omitempty"`
}

type updateURIRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

type updateResponse struct {
	Post postResponse `json:"post"`
}

// update updates a guest post owned by the authenticated customer.
//
// @Summary      Update post
// @Description  Updates post fields and optionally rewrites images.
// @Tags         guest-posts
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      int     true   "Post ID"
// @Param        venue_id  formData  int     false  "Venue ID"
// @Param        text      formData  string  false  "Post text"
// @Param        rating    formData  int     false  "Rating (1..5)"
// @Param        status    formData  string  false  "Allowed: draft,published,archived"
// @Param        images    formData  file    false  "Images field presence triggers rewrite"
// @Success      200       {object}  updateResponse
// @Failure      400       {object}  errorResponse
// @Failure      401       {object}  errorResponse
// @Failure      404       {object}  errorResponse
// @Failure      500       {object}  errorResponse
// @Router       /posts/{id} [patch]
func (h *handler) update(c *gin.Context) {
	if !strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
		c.Error(apperror.BadRequest("content type must be multipart/form-data"))
		return
	}

	var uriReq updateURIRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	if strings.EqualFold(strings.TrimSpace(c.PostForm("status")), string(entity.PostStatusDeleted)) {
		c.Error(apperror.BadRequest("status deleted is not allowed in patch"))
		return
	}

	var req updateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(apperror.BadRequest(err.Error()))
		return
	}

	if req.Status != nil && *req.Status == entity.PostStatusDeleted {
		c.Error(apperror.BadRequest("status deleted is not allowed in patch"))
		return
	}

	images, err := buildUploadImages(req.Images)
	if err != nil {
		c.Error(err)
		return
	}

	rewriteImages := multipartFieldProvided(c, "images")

	ctx := c.Request.Context()
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.customerService.GetByUserID(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	in := entity.UpdatePostInput{
		ID:            uriReq.ID,
		CustomerID:    customer.ID,
		VenueID:       req.VenueID,
		Text:          req.Text,
		Rating:        req.Rating,
		Status:        req.Status,
		Images:        images,
		RewriteImages: rewriteImages,
	}

	post, err := h.service.Update(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := updateResponse{Post: postToResponse(post, h.storage)}
	c.JSON(http.StatusOK, resp)
}
