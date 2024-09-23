package model

import "time"

type TextDataPostRequest struct {
	Text     string `validate:"required"`
	MetaData string
}

type TextData struct {
	ID        string
	OwnerID   string
	Text      string
	MetaData  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TextCryptData struct {
	Text string
}

type TextDataDB struct {
	ID        string    `db:"id"`
	CryptData []byte    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	MetaData  string    `db:"meta_data"`
}
