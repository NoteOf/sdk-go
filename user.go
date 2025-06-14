package sdk

// UserMeta is a map of string keys to string values that can be used to store
// arbitrary metadata about a note that is not part of the note itself
// and is useful for implementation specific values.
//
// The server enforces a 255 character limit on both keys and values.
// The key is expected to be ascii only and the value is expected to be utf-8
type UserMeta map[string]string

type User struct {
	Username     string `json:"username"`
	EmailAddress string `json:"email"`

	Meta UserMeta `json:"meta"`
}
