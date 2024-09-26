package auth

import (
	"auth-service/internal/model"
	"auth-service/internal/sys/codes"
	"auth-service/internal/utils"
	desc "auth-service/pkg/auth"
	"context"
	"errors"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) Register(ctx context.Context, req *desc.RegisterRequest) (*desc.RegisterResponse, error) {

	// Поиск пользователя по имени
	user, err := i.userRepo.FindByName(ctx, req.GetUsername())
	if user != nil {
		return nil, status.Errorf(grpcCodes.Code(codes.AlreadyExists), "%s already registered", user.Username)
	}

	// Проверка существования пользователя по email (если email указан)
	if req.GetEmail() != "" {
		user, err = i.userRepo.FindByEmail(ctx, req.GetEmail())
		if err != nil {
			log.Printf("Error finding user by email: %v", err)
			return nil, status.Errorf(grpcCodes.Code(codes.Internal), "error checking email: %v", err)
		}
		if user != nil {
			return nil, status.Errorf(grpcCodes.Code(codes.AlreadyExists), "email %s already registered", req.GetEmail())
		}
	}

	userRole, err := i.roleRepo.FindByName(ctx, "user")
	if err != nil {
		log.Printf("Error finding user role: %v", err)
		return nil, err
	}

	if userRole == nil {
		log.Println("User role not found")
		return nil, errors.New("Doesn't have user role")
	}

	password, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}

	// Создание пользователя
	newUserId, err := i.userRepo.Create(ctx, model.UserInfo{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: password,
		Role:         userRole.Name,
		RoleId:       userRole.ID,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	newUser, err := i.userRepo.Get(ctx, newUserId)
	if err != nil {
		log.Printf("Error getting newly created user: %v", err)
		return nil, err
	}

	// Генерация refresh-токена
	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: newUser.Username,
		Role:     newUser.Role,
	},
		[]byte(i.config.RefreshTokenSecretKey),
		i.config.RefreshTokenExpiration,
	)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, errors.New("failed to generate token")
	}

	// Возвращение токена в ответе
	return &desc.RegisterResponse{RefreshToken: refreshToken}, nil
}
