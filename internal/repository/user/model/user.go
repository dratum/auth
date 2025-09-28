package modelRepo

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Role      string       `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
