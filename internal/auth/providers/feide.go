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

// Name gets the name used to retrieve this provider.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

// SetPermissions is to update the bot permissions (used for when ScopeBot is set)
func (p *Provider) SetPermissions(permissions string) {
	p.permissions = permissions
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

// BeginAuth asks Discord for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {

	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("prompt", "none"),
	}

	if p.permissions != "" {
		opts = append(opts, oauth2.SetAuthURLParam("permissions", p.permissions))
	}

	url := p.config.AuthCodeURL(state, opts...)

	s := &Session{
		AuthURL: url,
	}
	return s, nil
}

// FetchUser will go to Discord and access basic info about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	s := session.(*Session)

	// Introduced by : Yyewolf

	if u.AvatarID != "" {
		avatarExtension := ".jpg"
		prefix := "a_"
		if len(u.AvatarID) >= len(prefix) && u.AvatarID[0:len(prefix)] == prefix {
			avatarExtension = ".gif"
		}
		user.AvatarURL = "https://media.discordapp.net/avatars/" + u.ID + "/" + u.AvatarID + avatarExtension
	}

	user.Name = u.Name
	user.Email = u.Email
	user.UserID = u.ID

	return nil
}

func newConfig(p *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     p.ClientKey,
		ClientSecret: p.Secret,
		RedirectURL:  p.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	if len(scopes) > 0 {
		for _, scope := range scopes {
			c.Scopes = append(c.Scopes, scope)
		}
	} else {
		c.Scopes = []string{ScopeIdentify}
	}

	return c
}

// RefreshTokenAvailable refresh token is provided by auth provider or not
func (p *Provider) RefreshTokenAvailable() bool {
	return true
}

// RefreshToken get new access token based on the refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.config.TokenSource(oauth2.NoContext, token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}
