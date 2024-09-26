package model

import (
	modelRepo "auth-service/internal/repository/role/model"
)

type User struct {
	ID           int64           `db:"id"`
	Username     string          `db:"username"`
	PasswordHash string          `db:"password_hash"`
	Role         *modelRepo.Role `db:"role"` // Связь с моделью Role
}
