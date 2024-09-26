package email

import (
	"auth-service/internal/service"
)

type serv struct {
}

func NewService() service.EmailService {
	return &serv{}
}
