package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const (
	providerName = "google"
	tokenURL     = "https://oauth2.googleapis.com/token"
	userInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type Provider struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) *Provider {
	return &Provider{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *Provider) ExchangeCode(ctx context.Context, code string) (*dto.OAuthUserInfo, error) {
	accessToken, err := p.exchangeToken(ctx, code)
	if err != nil {
		return nil, apperr.ErrProviderExchangeFail
	}
	info, err := p.fetchUserInfo(ctx, accessToken)
	if err != nil {
		return nil, apperr.ErrProviderUserInfoFail
	}
	return &dto.OAuthUserInfo{
		Provider:   providerName,
		ProviderID: info.Sub,
		Email:      info.Email,
	}, nil
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

type googleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (p *Provider) exchangeToken(ctx context.Context, code string) (string, error) {
	body := url.Values{
		"code":          {code},
		"client_id":     {p.cfg.ClientID},
		"client_secret": {p.cfg.ClientSecret},
		"redirect_uri":  {p.cfg.RedirectURL},
		"grant_type":    {"authorization_code"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do token request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.WarnKV(ctx, "failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("google returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	var t tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if t.Error != "" {
		return "", fmt.Errorf("google token error: %s", t.Error)
	}

	return t.AccessToken, nil
}

func (p *Provider) fetchUserInfo(ctx context.Context, accessToken string) (*googleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do userinfo request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.WarnKV(ctx, "failed to close userinfo response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google userinfo returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode userinfo response: %w", err)
	}

	if !info.EmailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	return &info, nil
}
