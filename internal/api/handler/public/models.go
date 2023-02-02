package public

import (
	"errors"
	"time"
)

type PublicAccessDetails struct {
	ID   int64  `json:"id"`
	User string `json:"user"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
}

func (p PublicAccessDetails) Valid() error {
	now := time.Now().Unix()
	if p.Iat > now {
		return errors.New("token is not valid yet")
	}
	if p.Exp < now {
		return errors.New("token is expired")
	}
	return nil
}
