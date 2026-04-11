package provider

import (
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/google"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
)

type Factory struct {
	provider map[string]authsvc.OAuthProvider
}

func NewFactory(googleCfg google.Config) *Factory {
	return &Factory{
		provider: map[string]authsvc.OAuthProvider{
			"google": google.New(googleCfg),
		},
	}
}

func (f *Factory) Get(name string) (authsvc.OAuthProvider, error) {
	p, ok := f.provider[name]
	if !ok {
		return nil, apperr.ErrUnsupportedProvider
	}
	return p, nil
}
