package env

type GoogleConfig struct {
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,required"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URI,required"`
}

func NewGoogleConfig(opts ...Options) (*GoogleConfig, error) {
	config := new(GoogleConfig)
	if err := Parse(config, opts...); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *GoogleConfig) ClientID() string     { return c.GoogleClientID }
func (c *GoogleConfig) ClientSecret() string { return c.GoogleClientSecret }
func (c *GoogleConfig) RedirectURL() string  { return c.GoogleRedirectURL }
