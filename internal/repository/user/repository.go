package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dratum/auth/internal/model"
	"github.com/dratum/auth/internal/repository"
	"github.com/dratum/auth/internal/repository/user/converter"
	modelRepo "github.com/dratum/auth/internal/repository/user/model"
	"github.com/dratum/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	query := `
			  select id
			       , name
			       , email
			       , role
			       , created_at
			       , updated_at
			  from users u
			  where u.id = $1
	`

	var user modelRepo.User
	var roleStr string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&roleStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	log.Print(roleStr)
	//TODO: Проблемы с определением роли при запроса GET
	roleValue, exists := auth_v1.Role_value[strings.ToUpper(roleStr)]
	if !exists {
		return nil, fmt.Errorf("invalid role value from database: %s", roleStr)
	}

	return converter.ToUserFromRepo(&user, roleValue), nil
}

func (r *repo) Create(ctx context.Context, fields *model.User) (int64, error) {
	if fields.Password != fields.PasswordConfirm {
		return 0, errors.New("passwords do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fields.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
			  INSERT INTO users
			  (
				name
			  , email
			  , password_hash
			  , role
			  , created_at
			  , updated_at
			  ) 
        	  VALUES (
			  	$1
			  , $2
			  , $3
			  , $4
			  , NOW()
			  , NOW()
			  )
        	  RETURNING id
		`

	var id int64

	err = r.db.QueryRow(
		ctx,
		query,
		fields.Name,
		fields.Email,
		hashedPassword,
		fields.Role,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create new user: %w", err)
	}

	return id, nil
}

// func (r *repo) Update(ctx context.Context, id int64) error {}

// func (r *repo) Delete(ctx context.Context, id int64) error {}
