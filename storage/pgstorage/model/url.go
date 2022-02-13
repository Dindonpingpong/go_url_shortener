package model

import "github.com/google/uuid"

type URLInDB struct {
	Id       uint      `db:"id"`
	UserID   uuid.UUID `db:"user_id"`
	Url      string    `db:"url"`
	ShortURL string    `db:"short_url"`
}
