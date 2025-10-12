package auth

import "context"

func (s *serv) Update(ctx context.Context, id int64, opts ...map[string]string) error {
	err := s.userRepository.Update(ctx, id, opts...)
	if err != nil {
		return err
	}

	return nil
}
