package repository

import (
	"auth-service/internal/model"
	"context"
)

type UserRepository interface {
	FindByName(ctx context.Context, username string) (*model.UserInfo, error)
	FindByTelegramId(ctx context.Context, telegramId int64) (*model.UserInfo, error)
	FindByEmailToken(ctx context.Context, token string) (*model.UserInfo, error)
	FindByEmail(ctx context.Context, email string) (*model.UserInfo, error)
	Get(ctx context.Context, id int64) (*model.UserInfo, error)
	Create(ctx context.Context, info model.UserInfo) (int64, error)
	CreateToken(ctx context.Context, info model.MailTokenInfo) error
	UpdateEmailVerified(ctx context.Context, userID int64) error
}

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*model.RoleInfo, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, name string) (int64, error)
}
