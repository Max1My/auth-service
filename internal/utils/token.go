package utils

import (
	"auth-service/internal/model"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

func GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		UserID:     info.ID,
		Username:   info.Username,
		Role:       info.Role,
		TelegramID: info.TelegramID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("Invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("Invalid token claims")
	}
	return claims, nil
}

// Функция для парсинга и проверки JWT токена
func ParseToken(tokenString string, secretKey string) (*model.UserClaims, error) {
	// Парсим токен с помощью библиотеки jwt-go
	token, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи, должен быть HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Возвращаем секретный ключ для проверки подписи
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена и извлекаем claims
	if claims, ok := token.Claims.(*model.UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// Пример обработки заголовка аутентификации
func GetTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata found")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return "", errors.New("authorization token is not provided")
	}

	token, err := ExtractToken(authHeaders[0])
	if err != nil {
		return "", err
	}

	return token, nil
}

func ExtractToken(authHeader string) (string, error) {
	// Проверяем, что заголовок начинается с "Bearer "
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		// Удаляем префикс "Bearer " и возвращаем токен
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	return "", errors.New("tokenstring should not contain 'bearer '")
}
