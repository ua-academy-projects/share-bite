package post

import (
	"context"
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"

	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

// list returns paginated published posts.
//
//	@Summary		List posts
//	@Description	Returns paginated list of published posts with authors information.
//	@Tags			guest-posts
//	@Produce		json
//	@Param			limit	query		int						false	"Max items per page (1..100)"	default(20)
//	@Param			offset	query		int						false	"Offset (0..1000)"				default(0)
//	@Success		200		{object}	listResponse			"Successfully retrieved posts"
//	@Failure		400		{object}	response.ErrorResponse	"Invalid query parameters"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/posts/ [get]
func (h *handler) list(c *gin.Context) {
	var req listRequest
	if err := request.BindQuery(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()

	var customerID string
	optionalCustomerID, err := httpctx.GetOptionalCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}
	if optionalCustomerID != nil {
		customerID = *optionalCustomerID
	}

	in := dto.ListPostsInput{
		Limit:      req.Limit,
		Offset:     req.Offset,
		CustomerID: customerID,
		AuthorID:   req.AuthorID,
	}
	out, err := h.service.List(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp, err := listPostsOutToResponse(
		ctx,
		out,
		h.storage,
		h.customerService,
		h.service,
	)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

type listRequest struct {
	Limit    int    `form:"limit,default=20" binding:"gte=1,lte=100"`
	Offset   int    `form:"offset,default=0" binding:"gte=0,lte=1000"`
	AuthorID string `form:"customer_id" binding:"omitempty"`
}

type listResponse struct {
	Posts []postResponse `json:"posts"`
	Total int            `json:"total"`
}

func listPostsOutToResponse(ctx context.Context, out dto.ListPostsOutput, storage storage.ObjectStorage, customerService customerService, postService postService) (listResponse, error) {
	customerIDSet := make(map[string]struct{})

	postAuthors := make(map[string][]string)

	for _, p := range out.Posts {
		customerIDSet[p.CustomerID] = struct{}{}

		authorIDs, err := postService.GetPostAuthors(ctx, p.ID)
		if err != nil {
			return listResponse{}, err
		}

		postAuthors[p.ID] = authorIDs

		for _, authorID := range authorIDs {
			customerIDSet[authorID] = struct{}{}
		}
	}

	customerIDs := make([]string, 0, len(customerIDSet))

	for id := range customerIDSet {
		customerIDs = append(customerIDs, id)
	}

	customers, err := customerService.GetByIDs(ctx, customerIDs)
	if err != nil {
		return listResponse{}, err
	}

	customerMap := make(map[string]entity.Customer, len(customers))

	for _, c := range customers {
		customerMap[c.ID] = c
	}

	list := make([]postResponse, 0, len(out.Posts))

	for _, p := range out.Posts {
		customer, ok := customerMap[p.CustomerID]
		if !ok {
			customer = entity.Customer{ID: p.CustomerID}
		}

		authorResponses := make([]authorResponse, 0, len(postAuthors[p.ID]))

		for _, authorID := range postAuthors[p.ID] {
			author, ok := customerMap[authorID]
			if !ok {
				continue
			}

			var avatarURL string

			if author.AvatarObjectKey != nil && storage != nil {
				avatarURL = storage.BuildURL(
					*author.AvatarObjectKey,
				)
			}

			authorResponses = append(authorResponses, authorResponse{
				ID:        author.ID,
				UserName:  author.UserName,
				AvatarURL: avatarURL,
			})
		}

		list = append(list, postToResponse(
			p,
			storage,
			customer,
			authorResponses,
		))
	}

	return listResponse{
		Posts: list,
		Total: out.Total,
	}, nil
}
