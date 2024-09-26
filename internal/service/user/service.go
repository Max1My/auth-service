package user

import (
	"auth-service/internal/client/db"
	"auth-service/internal/repository"
	"auth-service/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManger       db.TxManager
}

func NewService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManger:       txManager,
	}
}
