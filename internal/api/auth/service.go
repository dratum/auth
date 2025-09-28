package auth

import (
	"github.com/dratum/auth/internal/service"
	"github.com/dratum/auth/pkg/auth_v1"
)

type Server struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

func New(authService service.AuthService) *Server {
	return &Server{
		authService: authService,
	}
}
