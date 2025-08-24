package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/dratum/auth/internal/domain/models"
	"github.com/dratum/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthService struct {
	db     *pgx.Conn
	server *auth_v1.UnimplementedAuthV1Server
}

func (a *AuthService) New(db *pgx.Conn, server *auth_v1.UnimplementedAuthV1Server) *AuthService {
	return &AuthService{db: db, server: server}
}

func (a *AuthService) Get(
	ctx context.Context,
	req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
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

	var user models.User
	var roleStr string

	err := a.db.QueryRow(ctx, query, req.Id).Scan(
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

func (a *AuthService) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	return &auth_v1.CreateResponse{}, nil
}
