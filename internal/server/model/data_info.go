package model

import "time"

type DataInfo struct {
	ID        string    `json:"id"`
	DataType  string    `json:"data_type"`
	MetaData  string    `json:"meta_data"`
	CreatedAt time.Time `json:"created_at"`
}
