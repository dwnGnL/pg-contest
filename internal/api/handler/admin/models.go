package admin

import (
	"errors"
	"time"
)

type AdminAccessDetails struct {
	ID   int64  `json:"id"`
	User string `json:"user"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
}

func (a AdminAccessDetails) Valid() error {
	now := time.Now().Unix()
	if a.Iat > now {
		return errors.New("token is not valid yet")
	}
	if a.Exp < now {
		return errors.New("token is expired")
	}
	return nil
}
