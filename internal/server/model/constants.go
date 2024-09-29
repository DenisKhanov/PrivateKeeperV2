package model

// CTXKey is the type used as a context key for storing user ID.
type CTXKey string

const (
	UserIDKey CTXKey = "userID"  // UserIDKey is the specific key used in the context to store user ID.
	UserKey   CTXKey = "userKey" // UserKey is the specific key used in the context to store user key.
)
