package model

type CreditCardPostRequest struct {
	Number    string
	OwnerName string
	ExpiresAt string
	CVV       string
	PinCode   string
	MetaData  string
}

type CreditCardLoadRequest struct {
	Number        string
	Owner         string
	CvvCode       string
	PinCode       string
	Metadata      string
	ExpiresAfter  string
	ExpiresBefore string
}

type CreditCard struct {
	Number    string
	OwnerName string
	ExpiresAt string
	CVV       string
	PinCode   string
	MetaData  string
}
