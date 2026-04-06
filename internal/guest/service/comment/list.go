package comment

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error) {
	postIDStr := strconv.FormatInt(in.PostID, 10)
	if _, err := s.postSvc.Get(ctx, postIDStr, ""); err != nil {
		return entity.ListCommentsOutput{}, errwrap.Wrap("check post existence for comments list", err)
	}

	if in.PageToken != "" {
		decoded, err := base64.URLEncoding.DecodeString(in.PageToken)
		if err == nil {
			in.PageToken = string(decoded)
		}
	}

	out, err := s.commentRepo.List(ctx, in)
	if err != nil {
		return entity.ListCommentsOutput{}, errwrap.Wrap("get list of comments from repo", err)
	}

	if out.NextPageToken != "" {
		out.NextPageToken = base64.URLEncoding.EncodeToString([]byte(out.NextPageToken))
	}

	return out, nil
}
