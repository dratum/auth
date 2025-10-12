package auth

import (
	"context"

	"github.com/dratum/auth/internal/converter"
	"github.com/dratum/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	err := s.authService.Update(ctx, converter.ToUpdateFromAuth(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
