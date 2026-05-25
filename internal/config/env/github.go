package env

type githubConfig struct {
	ClientID           string `env:"GITHUB_CLIENT_ID,required"`
	ClientSecret       string `env:"GITHUB_CLIENT_SECRET,required"`
	RedirectURL        string `env:"GITHUB_REDIRECT_URL,required"`
	SuccessRedirectURL string `env:"GITHUB_SUCCESS_REDIRECT_URL"`
}

func NewGitHubConfig(opts ...Options) (*githubConfig, error) {
	config := new(githubConfig)
	if err := Parse(config, opts...); err != nil {
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
