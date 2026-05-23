package admin

import (
	"context"
	"fmt"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/admin/client/admin_auth_client/user"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"
)

type Client struct {
	api    user.ClientService
	scheme string
	policy *resilience.Policy
}

func New(host, basePath, scheme string, policy *resilience.Policy) *Client {
	transport := httptransport.New(host, basePath, []string{scheme})
	api := user.New(transport, strfmt.Default)
	return &Client{api: api, scheme: scheme, policy: policy}
}

func (c *Client) execute(ctx context.Context, operation func() error) error {
	if c.policy == nil {
		return operation()
	}
	return c.policy.Execute(ctx, operation)
}

func (c *Client) getClientOption() user.ClientOption {
	return func(op *runtime.ClientOperation) {
		op.Schemes = []string{c.scheme}
	}
}

func (c *Client) GetUserEmail(ctx context.Context, userID, authToken string) (string, error) {
	params := user.NewGetUsersUserIDEmailParams().WithContext(ctx).WithUserID(userID)

	var auth runtime.ClientAuthInfoWriter
	if authToken != "" {
		auth = httptransport.BearerToken(authToken)
	}

	var email string
	err := c.execute(ctx, func() error {
		resp, err := c.api.GetUsersUserIDEmail(params, auth, c.getClientOption())
		if err != nil {
			return err
		}
		if resp == nil || resp.GetPayload() == nil {
			return fmt.Errorf("unexpected empty response")
		}
		email = resp.GetPayload().Email
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("get user email failed: %w", err)
	}
	return email, nil
}
