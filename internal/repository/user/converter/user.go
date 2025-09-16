package converter

import (
	"github.com/dratum/auth/internal/repository/user/model"
	"github.com/dratum/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromRepo(user *model.User, role int32) *auth_v1.GetResponse {
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
