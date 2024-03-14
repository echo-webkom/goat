package domain

import "time"

type Session struct {
	SessionToken string    `json:"session_token"`
	UserID       string    `json:"user_id"`
	Expires      time.Time `json:"expires"`
}
