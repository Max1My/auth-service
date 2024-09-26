package service

import (
	"auth-service/internal/model"
	"context"
)

type AuthService interface {
	Login(ctx context.Context)
	GetAccessToken(ctx context.Context)
	GetRefreshToken(ctx context.Context)
	Check(ctx context.Context)
}

type UserService interface {
	Get(ctx context.Context, id int64) (*model.UserInfo, error)
	FindByName(ctx context.Context, username string) (*model.UserInfo, error)
	FindByTelegramId(ctx context.Context, telegramId int64) (*model.UserInfo, error)
}

type RoleService interface {
	GetUserRole(ctx context.Context) (*model.RoleInfo, error)
}

type EmailService interface {
	SendConfirmationEmail(message model.MailMessageInfo) error
}
