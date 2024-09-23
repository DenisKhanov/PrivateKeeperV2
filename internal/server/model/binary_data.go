package model

import "time"

type BinaryDataPostRequest struct {
	Name      string `validate:"required"`
	Extension string `validate:"required"`
	Data      []byte `validate:"required"`
	MetaData  string
}

type BinaryData struct {
	ID        string
	OwnerID   string
	Name      string
	Extension string
	Data      []byte
	MetaData  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BinaryCryptData struct {
	Name      string
	Extension string
	Data      []byte
}

type BinaryDataDB struct {
	ID        string    `db:"id"`
	CryptData []byte    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	MetaData  string    `db:"meta_data"`
}
