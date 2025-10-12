package converter

import (
	"github.com/dratum/auth/internal/model"
	"github.com/dratum/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User, role int32) *auth_v1.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &auth_v1.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      auth_v1.Role(role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToCreateFromAuth(user *auth_v1.CreateRequest) *model.User {
	return &model.User{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            user.Role.Enum().String(),
	}
}

func ToUpdateFromAuth(user *auth_v1.UpdateRequest) *model.UserUpdate {
	result := &model.UserUpdate{
		Id: user.Id,
	}

	// Проверяем, задано ли поле name
	if user.Name != nil {
		result.Name = &user.Name.Value
	}

	// Проверяем, задано ли поле email
	if user.Email != nil {
		result.Email = &user.Email.Value
	}

	return result
}
