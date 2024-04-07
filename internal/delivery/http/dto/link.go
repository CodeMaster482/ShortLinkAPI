package dto

import "time"

type CreateLinkRequest struct {
	Link string `json:"link"`
}

type CreateLinkResponse struct {
	ShortLink string    `json:"short_link"`
	ExpiresAt time.Time `json:"expires_at"`
}
