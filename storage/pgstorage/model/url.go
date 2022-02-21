package model

import "github.com/google/uuid"

type URLInDB struct {
	ID        uint      `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	URL       string    `db:"url"`
	IsDeleted bool      `db:"is_deleted"`
	ShortURL  string    `db:"short_url"`
}
