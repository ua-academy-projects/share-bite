package comment

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"fmt"
)

func (s *service) List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error) {
	postIDStr := strconv.FormatInt(in.PostID, 10)
	if _, err := s.postSvc.Get(ctx, postIDStr, ""); err != nil {
		return entity.ListCommentsOutput{}, fmt.Errorf("check post existence for comments list: %w", err)
	}

	if in.PageToken != "" {
		decoded, err := base64.URLEncoding.DecodeString(in.PageToken)
		if err == nil {
			in.PageToken = string(decoded)
		}
	}

	out, err := s.commentRepo.List(ctx, in)
	if err != nil {
		return entity.ListCommentsOutput{}, fmt.Errorf("get list of comments from repo: %w", err)
	}

	if out.NextPageToken != "" {
		out.NextPageToken = base64.URLEncoding.EncodeToString([]byte(out.NextPageToken))
	}

	return out, nil
}
