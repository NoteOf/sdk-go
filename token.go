package sdk

import "time"

type Token struct {
	APIToken   string     `json:"token"`
	Created    time.Time  `json:"created"`
	Expiration *time.Time `json:"expiration,omitempty"`
}
