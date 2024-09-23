package model

import "time"

type Data struct {
	ID        string    `db:"id"`
	OwnerID   string    `db:"owner_id"`
	Type      string    `db:"type"`
	Data      []byte    `db:"data"`
	MetaData  string    `db:"meta_data"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
