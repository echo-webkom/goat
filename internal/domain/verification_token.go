package domain

import "time"

// VerificationToken represents a token used to verify a user's email address.
// Not really needed/used in the current implementation, but it's here for future use.

type VerificationToken struct {
	Identifier string    `json:"identifier"`
	UserID     string    `json:"user_id"`
	Expires    time.Time `json:"expires"`
}
