package user

import (
	"auth-service/internal/model"
	"context"
)

func (s *serv) FindByName(ctx context.Context, username string) (*model.UserInfo, error) {
	user, err := s.userRepository.FindByName(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *serv) FindByTelegramId(ctx context.Context, telegramId int64) (*model.UserInfo, error) {
	user, err := s.userRepository.FindByTelegramId(ctx, telegramId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *serv) Get(ctx context.Context, id int64) (*model.UserInfo, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
