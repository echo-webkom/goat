package providers

import (
	"net/http"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

const (
	authURL      string = "https://discord.com/api/oauth2/authorize"
	tokenURL     string = "https://discord.com/api/oauth2/token"
	userEndpoint string = "https://discord.com/api/users/@me"
)

const (
	ScopeEmail   = "email"
	ScopeOpenID  = "openid"
	ScopeProfile = "profile"
	ScopeGroups  = "groups"
)

type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
	permissions  string
}

func (p *Provider) Name() string {
	return p.providerName
}

func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) SetPermissions(permissions string) {
	p.permissions = permissions
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}
