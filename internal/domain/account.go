package domain

type Account struct {
	UserID            string `json:"user_id"`
	Type              string `json:"type"`
	Provider          string `json:"provider"`
	ProviderAccountID string `json:"provider_account_id"`
	RefreshToken      string `json:"refresh_token"`
	AccessToken       string `json:"access_token"`
	ExpiresAt         int64  `json:"expires_at"`
	TokenType         string `json:"token_type"`
	Scope             string `json:"scope"`
	IDToken           string `json:"id_token"`
	SessionState      string `json:"session_state"`
}
