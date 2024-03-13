package domain

import "time"

type User struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	EmilVerified     time.Time `json:"email_verified"`
	Image            string    `json:"image"`
	AlternativeEmail string    `json:"alternative_email"`
	DegreeID         string    `json:"degree_id"`
	Year             int       `json:"year"`
	Type             string    `json:"type"`
	IsBanned         bool      `json:"is_banned"`
	BannedFromStrike int       `json:"banned_from_strike"`
}
