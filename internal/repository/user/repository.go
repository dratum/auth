package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dratum/auth/internal/repository"
	"github.com/dratum/auth/internal/repository/user/model"
	"github.com/dratum/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Get(ctx context.Context, id int64) (*auth_v1.GetResponse, error) {
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

	var user model.User
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

	roleValue, exists := auth_v1.Role_value[strings.ToUpper(roleStr)]
	if !exists {
		return nil, fmt.Errorf("invalid role value from database: %s", roleStr)
	}

	return &auth_v1.GetResponse{
		Id:        user.Id,
		Email:     user.Email,
		Name:      user.Name,
		Role:      auth_v1.Role(roleValue),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (r *repo) Create(ctx context.Context, fields *auth_v1.CreateRequest) (int64, error) {
	if fields.Password != fields.PasswordConfirm {
		return 0, errors.New("passwords do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fields.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
			  insert name          = $1
			  	   , email         = $2
			  	   , password_hash = $3
			  	   , role          = $4
			  into users	
	`

	var id int64

	err = r.db.QueryRow(
		ctx,
		query,
		fields.Email,
		fields.Name,
		hashedPassword,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create new user: %w", err)
	}

	return id, nil
}

func (r *repo) Update(ctx context.Context, id int64) error {}

func (r *repo) Delete(ctx context.Context, id int64) error {}
