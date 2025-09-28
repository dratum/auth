package user

import (
	"github.com/dratum/auth/internal/repository"
	"github.com/dratum/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func newService(userRepository repository.UserRepository) service.AuthService {
	return &serv{
		userRepository: userRepository,
	}
}
