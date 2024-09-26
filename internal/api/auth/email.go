package auth

import (
	"auth-service/internal/model"
	"auth-service/internal/utils"
	desc "auth-service/pkg/auth"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (i *Implementation) SendConfirmationEmail(ctx context.Context, req *desc.SendConfirmationEmailRequest) (*desc.SendConfirmationEmailResponse, error) {
	// Извлекаем userId из контекста авторизованного пользователя
	userId, err := i.getUserIdFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to retrieve user ID: %v", err)
	}

	user, err := i.userRepo.Get(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed  to get user: %v", err)
	}

	// Генерируем токен для подтверждения email
	token, err := utils.GenerateToken(model.UserInfo{
		Username: user.Username,
		Role:     user.Role,
	}, []byte(i.config.EmailTokenSecretKey), i.config.EmailTokenExpiration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate email confirmation token: %v", err)
	}

	// Сохраняем токен в базе данных
	err = i.userRepo.CreateToken(ctx, model.MailTokenInfo{
		UserID:    userId,
		Token:     token,
		ExpiresAt: time.Now().Add(i.config.EmailTokenExpiration),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save email token: %v", err)
	}

	// Отправляем email с токеном пользователю (псевдокод)
	err = i.emailServ.SendConfirmationEmail(model.MailMessageInfo{
		From:         "noreply@yourservice.com",
		To:           req.GetEmail(),
		SmtpHost:     i.smtpConfig.SmtpHost,
		SmtpPort:     i.smtpConfig.SmtpPort,
		SmtpUser:     i.smtpConfig.SmtpUser,
		SmtpPassword: i.smtpConfig.SmtpPassword,
		Token:        token,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send email: %v", err)
	}

	return &desc.SendConfirmationEmailResponse{Send: true}, nil
}

func (i *Implementation) ConfirmEmail(ctx context.Context, req *desc.ConfirmEmailRequest) (*desc.ConfirmEmailResponse, error) {
	// Валидация токена
	user, err := i.userRepo.FindByEmailToken(ctx, req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid token")
	}

	// Обновление статуса email_verified
	err = i.userRepo.UpdateEmailVerified(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to verify email")
	}

	return &desc.ConfirmEmailResponse{Success: true}, nil
}
