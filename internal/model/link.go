package model

import "time"

type Link struct {
	OriginalLink string `db:"original_link"`
	ShortLink    string
	Token        string    `db:"token"`
	ExpiresAt    time.Time `db:"expires_at"`
}

// func (l *Link) Expired(now string) bool {
// 	return now > l.ExpiresAt
// }
