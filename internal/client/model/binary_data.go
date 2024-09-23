package model

type BinaryDataPostRequest struct {
	Name      string
	Extension string
	Data      []byte
	MetaData  string
}

type BinaryData struct {
	Name      string
	Extension string
	Data      []byte
	MetaData  string
}

type BinaryDataLoadRequest struct {
	Name     string
	MetaData string
}
