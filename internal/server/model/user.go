package model

import "time"

type UserRegisterRequest struct {
	Login    string `validate:"email"`
	Password string `validate:"required"`
}

type UserLoginRequest struct {
	Login    string `validate:"email"`
	Password string `validate:"required"`
}

type User struct {
	ID        string    `db:"id"`
	Login     string    `db:"login"`
	Password  []byte    `db:"password"`
	CryptKey  []byte    `db:"crypt_key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
