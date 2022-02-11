package model

type RowInURLStorage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	UserID string `json:"userId"`
}

type URLInDB struct {
	URL    string
	UserID string
}
