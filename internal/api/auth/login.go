package auth

import (
	"auth-service/internal/model"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"

	"auth-service/internal/utils"
	desc "auth-service/pkg/auth"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	var user *model.UserInfo
	var err error

	// Поиск пользователя по имени или email
	if req.GetUsername() != "" {
		user, err = i.userRepo.FindByName(ctx, req.GetUsername())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user data by username: %w", err)
		}
	} else if req.GetEmail() != "" {
		user, err = i.userRepo.FindByEmail(ctx, req.GetEmail())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user data by email: %w", err)
		}
	} else {
		return nil, errors.New("either username or email must be provided")
	}

	// Проверка наличия пользователя и валидации пароля
	if user == nil || !utils.VerifyPassword(user.PasswordHash, req.GetPassword()) {
		return nil, errors.New("invalid username or password")
	}

	// Проверка на наличие роли у пользователя
	if user.RoleId == 0 {
		return nil, errors.New("user role not found")
	}

	// Генерация refresh-токена
	refreshToken, err := utils.GenerateToken(model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	},
		[]byte(i.config.RefreshTokenSecretKey),
		i.config.RefreshTokenExpiration,
	)
	fmt.Println(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Возвращение токена в ответе
	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}

func (i *Implementation) AuthenticateTelegramUser(ctx context.Context, req *desc.AuthenticateTelegramUserRequest) (*desc.AuthenticateTelegramUserResponse, error) {
	// Поиск пользователя по имени
	user, err := i.userRepo.FindByTelegramId(ctx, req.GetTelegramId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}

	userRole, err := i.roleRepo.FindByName(ctx, "telegram_user")
	if err != nil {
		log.Printf("Error finding user role: %v", err)
		return nil, err
	}

	if userRole == nil {
		log.Println("User role not found")
		return nil, errors.New("Doesn't have user role")
	}

	if user == nil {
		_, err = i.userRepo.Create(ctx, model.UserInfo{
			Role:       userRole.Name,
			RoleId:     userRole.ID,
			TelegramID: req.GetTelegramId(),
		})
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return nil, err
		}
	}

	// Генерация refresh-токена
	refreshToken, err := utils.GenerateToken(model.UserInfo{
		ID:         user.ID,
		TelegramID: req.GetTelegramId(),
		Role:       userRole.Name,
	},
		[]byte(i.config.RefreshTokenSecretKey),
		i.config.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Возвращение токена в ответе
	return &desc.AuthenticateTelegramUserResponse{RefreshToken: refreshToken}, nil
}

func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {

	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(i.config.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		ID:       claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	},
		[]byte(i.config.RefreshTokenSecretKey),
		i.config.RefreshTokenExpiration,
	)

	if err != nil {
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {

	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(i.config.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.UserInfo{
		ID:       claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	},
		[]byte(i.config.AccessTokenSecretKey),
		i.config.AccessTokenExpiration,
	)

	fmt.Println(claims.UserID)

	if err != nil {
		return nil, err
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}

func (i *Implementation) getUserIdFromContext(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("no metadata in context")
	}

	fmt.Println(md)

	token, err := utils.GetTokenFromContext(ctx)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "failed to retrieve user ID: %v", err)
	}

	claims, err := utils.ParseToken(token, i.config.AccessTokenSecretKey)
	if err != nil {
		return 0, err
	}

	userId := claims.UserID
	if userId == 0 {
		return 0, errors.New("user_id not found in token")
	}

	return userId, nil
}
