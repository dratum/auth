package auth

import (
	"context"

	"github.com/dratum/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := s.authService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
