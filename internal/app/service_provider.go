package app

import (
	"auth-service/internal/api/access"
	"auth-service/internal/api/auth"
	"auth-service/internal/client/db"
	"auth-service/internal/client/db/pg"
	"auth-service/internal/client/db/transaction"
	"auth-service/internal/closer"
	"auth-service/internal/config"
	"auth-service/internal/config/env"
	"auth-service/internal/repository"
	roleRepository "auth-service/internal/repository/role"
	userRepository "auth-service/internal/repository/user"
	"auth-service/internal/service"
	emailService "auth-service/internal/service/email"
	roleService "auth-service/internal/service/role"
	userService "auth-service/internal/service/user"
	"context"
	"log"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	tokenConfig *env.TokenConfigData
	smtpConfig  *env.SmtpConfigData
	httpConfig  config.HTTPConfig

	dbClient  db.Client
	txManager db.TxManager

	authService service.AuthService

	authImpl   *auth.Implementation
	accessImpl *access.Implementation

	userService    service.UserService
	userRepository repository.UserRepository

	roleService    service.RoleService
	roleRepository repository.RoleRepository
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("Failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("Failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) TokenConfig() *env.TokenConfigData {
	if s.tokenConfig == nil {
		cfg, err := env.NewTokenConfig()
		if err != nil {
			log.Fatalf("Failed to get token config: %s", err.Error())
		}

		s.tokenConfig = cfg
	}

	return s.tokenConfig
}

func (s *serviceProvider) SmtpConfig() *env.SmtpConfigData {
	if s.smtpConfig == nil {
		cfg, err := env.NewSmtpConfig()
		if err != nil {
			log.Fatalf("Failed  to get smtp config: %s", err.Error())
		}
		s.smtpConfig = cfg
	}
	return s.smtpConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to create db client: %s", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) GetAuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		tokenConfig := s.TokenConfig()
		smtpConfig := s.SmtpConfig()
		userRepo := userRepository.NewRepository(s.DBClient(ctx))
		roleRepo := roleRepository.NewRepository(s.DBClient(ctx))
		emailServ := emailService.NewService()
		s.authImpl = auth.NewImplementation(tokenConfig, smtpConfig, userRepo, roleRepo, emailServ)
		if s.authImpl == nil {
			log.Fatalf("Failed to initialize authImpl")
		}
	}
	return s.authImpl
}

func (s *serviceProvider) GetAccessImpl() *access.Implementation {
	if s.accessImpl == nil {
		tokenConfig := s.TokenConfig()
		s.accessImpl = access.NewImplementation(tokenConfig)
	}

	return s.accessImpl
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) RoleRepository(ctx context.Context) repository.RoleRepository {
	if s.roleRepository == nil {
		s.roleRepository = roleRepository.NewRepository(s.DBClient(ctx))
	}

	return s.roleRepository
}

func (s *serviceProvider) RoleService(ctx context.Context) service.RoleService {
	if s.roleService == nil {
		s.roleService = roleService.NewService(
			s.RoleRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.roleService
}
