package auth

import (
	"context"

	"github.com/dratum/auth/internal/converter"
	"github.com/dratum/auth/pkg/auth_v1"
)

func (s *Server) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	id, err := s.authService.Create(ctx, converter.ToUserFromAuth(req))
	if err != nil {
		return nil, err
	}

	return &auth_v1.CreateResponse{
		Id: id,
	}, nil
}
