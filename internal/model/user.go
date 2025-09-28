package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
