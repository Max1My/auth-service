package role

import (
	"auth-service/internal/client/db"
	"auth-service/internal/repository"
	"auth-service/internal/service"
)

type serv struct {
	roleRepository repository.RoleRepository
	txManger       db.TxManager
}

func NewService(
	roleRepository repository.RoleRepository,
	txManager db.TxManager,
) service.RoleService {
	return &serv{
		roleRepository: roleRepository,
		txManger:       txManager,
	}
}
