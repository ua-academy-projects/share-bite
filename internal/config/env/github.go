package env

import (
	"github.com/caarlos0/env/v11"
)

type githubConfig struct {
	ClientID           string `env:"GITHUB_CLIENT_ID,required"`
	ClientSecret       string `env:"GITHUB_CLIENT_SECRET,required"`
	RedirectURL        string `env:"GITHUB_REDIRECT_URL,required"`
	SuccessRedirectURL string `env:"GITHUB_SUCCESS_REDIRECT_URL"`
}

func NewGitHubConfig() (*githubConfig, error) {
	config := new(githubConfig)
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *githubConfig) GetClientID() string {
	return c.ClientID
}

func (c *githubConfig) GetClientSecret() string {
	return c.ClientSecret
}

func (c *githubConfig) GetRedirectURL() string {
	return c.RedirectURL
}

func (c *githubConfig) GetSuccessRedirectURL() string {
	return c.SuccessRedirectURL
}