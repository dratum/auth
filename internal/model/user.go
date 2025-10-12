package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id              int64
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
	CreatedAt       time.Time
	UpdatedAt       sql.NullTime
}

type UserUpdate struct {
	Id        int64
	Name      *string
	Email     *string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
