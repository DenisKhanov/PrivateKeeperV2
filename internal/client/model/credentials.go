package model

type CredentialsPostRequest struct {
	Login    string
	Password string
	MetaData string
}

type Credentials struct {
	Login    string
	Password string
	MetaData string
}

type CredentialsLoadRequest struct {
	Login    string
	Password string
	MetaData string
}
