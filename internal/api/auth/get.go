package auth

import (
	"context"
	"fmt"

	"github.com/dratum/auth/internal/converter"
	"github.com/dratum/auth/pkg/auth_v1"
)

func (s *Server) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	user, err := s.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	roleValue, exists := auth_v1.Role_value[user.Role]
	if !exists {
		return nil, fmt.Errorf("invalid role value from database: %s", user.Role)
	}
	return converter.ToUserFromService(user, roleValue), nil
}
