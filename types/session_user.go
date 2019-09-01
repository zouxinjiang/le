package types

import (
	"time"
)

type UserSession struct {
	UserName  string
	LoginTime time.Time
}
