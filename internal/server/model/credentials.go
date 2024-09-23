package model

import "time"

type CredentialsPostRequest struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
	MetaData string
}

type Credentials struct {
	ID        string
	OwnerID   string
	Login     string
	Password  string
	MetaData  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CredentialsCryptData struct {
	Login    string
	Password string
}

type CredentialsDB struct {
	ID        string    `db:"id"`
	OwnerID   string    `db:"owner_id"`
	CryptData []byte    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	MetaData  string    `db:"meta_data"`
}
