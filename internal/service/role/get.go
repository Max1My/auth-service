package role

import (
	"auth-service/internal/model"
	"context"
)

func (s *serv) GetUserRole(ctx context.Context) (*model.RoleInfo, error) {
	role, err := s.roleRepository.FindByName(ctx, "user")
	if err != nil {
		return nil, err
	}
	return role, nil
}
