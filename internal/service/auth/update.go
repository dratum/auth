package auth

import (
	"context"

	"github.com/dratum/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, fields *model.UserUpdate) error {
	err := s.userRepository.Update(ctx, fields)
	if err != nil {
		return err
	}

	return nil
}
