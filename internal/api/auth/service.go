package auth

import (
	"auth-service/internal/config/env"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	desc "auth-service/pkg/auth"
)

type Implementation struct {
	desc.UnimplementedAuthServer
	config     *env.TokenConfigData
	smtpConfig *env.SmtpConfigData
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	emailServ  service.EmailService
}

func NewImplementation(
	config *env.TokenConfigData,
	smtpConfig *env.SmtpConfigData,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	emailService service.EmailService,
) *Implementation {
	return &Implementation{
		config:     config,
		smtpConfig: smtpConfig,
		userRepo:   userRepository,
		roleRepo:   roleRepository,
		emailServ:  emailService,
	}
}
