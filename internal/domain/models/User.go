package models

import "time"

type User struct {
	Id        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
