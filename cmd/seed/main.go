package main

import (
	"auth-service/internal/client/db/pg"
	"auth-service/internal/config"
	"auth-service/internal/config/env"
	"auth-service/internal/repository"
	roleRepository "auth-service/internal/repository/role"
	"context"
	"fmt"
	"log"
)

func main() {
	// Инициализируем контекст
	ctx := context.Background()

	// Загружаем конфигурацию
	err := config.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	// Подключаемся к базе данных
	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("Failed to load PG config: %s", err.Error())
	}

	dbClient, err := pg.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}

	// Создаем репозиторий ролей
	roleRepo := roleRepository.NewRepository(dbClient)

	// Наполняем базу данных начальными данными
	err = seedRoles(ctx, roleRepo)
	if err != nil {
		log.Fatalf("Failed to seed roles: %s", err.Error())
	}

	fmt.Println("Seed completed successfully!")
}

func seedRoles(ctx context.Context, repo repository.RoleRepository) error {
	roles := []string{"admin", "user", "telegram_user"}

	for _, roleName := range roles {
		// Проверяем, существует ли роль
		exists, err := repo.ExistsByName(ctx, roleName)
		if err != nil {
			return fmt.Errorf("failed to check if role exists: %w", err)
		}

		// Если роль не существует, добавляем её
		if !exists {
			_, err = repo.Create(ctx, roleName)
			if err != nil {
				return fmt.Errorf("failed to create role '%s': %w", roleName, err)
			}
			fmt.Printf("Role '%s' created\n", roleName)
		} else {
			fmt.Printf("Role '%s' already exists\n", roleName)
		}
	}

	return nil
}
