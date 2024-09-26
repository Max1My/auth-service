package user

import (
	"auth-service/internal/client/db"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	roleModelRepo "auth-service/internal/repository/role/model"
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

const (
	// Users table constants
	usersTableName           = "users"
	usersIdColumn            = "id"
	usersUsernameColumn      = "username"
	usersRoleIdColumn        = "role_id"
	usersPasswordHashColumn  = "password_hash"
	usersTelegramIdColumn    = "telegram_id"
	usersEmailVerifiedColumn = "email_verified"
	usersEmailColumn         = "email"
)

const (
	// Password reset tokens table constants
	passwordResetTokensTableName = "password_reset_tokens"
	resetTokenIdColumn           = "id"
	resetTokenUserIdColumn       = "user_id"
	resetTokenColumn             = "token"
	resetTokenExpiresAtColumn    = "expires_at"
)

const (
	// Role table constants
	rolesTableName  = "roles"
	rolesIdColumn   = "id"
	rolesNameColumn = "role_name"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) FindByName(ctx context.Context, username string) (*model.UserInfo, error) {
	var user model.UserInfo
	var role roleModelRepo.Role

	// Строим запрос с использованием JOIN для получения роли
	builder := sq.Select("u."+usersIdColumn, "u."+usersUsernameColumn, "u."+usersRoleIdColumn, "u."+usersPasswordHashColumn, "r."+rolesNameColumn).
		PlaceholderFormat(sq.Dollar).
		From(usersTableName + " u").
		Join("roles r ON u." + usersRoleIdColumn + " = r." + rolesIdColumn).
		Where(sq.Eq{"u." + usersUsernameColumn: username}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.FindByName",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&user.ID, &user.Username, &user.RoleId, &user.PasswordHash, &role.Name,
	)
	if err != nil {
		return nil, err
	}

	user.Role = role.Name
	return &user, nil
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*model.UserInfo, error) {
	var user model.UserInfo
	var role roleModelRepo.Role

	// Строим запрос с использованием JOIN для получения роли
	builder := sq.Select("u."+usersUsernameColumn, "u."+usersRoleIdColumn, "u."+usersPasswordHashColumn, "r."+rolesNameColumn).
		PlaceholderFormat(sq.Dollar).
		From(usersTableName + " u").
		Join("roles r ON u." + usersRoleIdColumn + " = r." + rolesIdColumn).
		Where(sq.Eq{"u." + usersEmailColumn: email}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.FindByName",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&user.Username, &user.RoleId, &user.PasswordHash, &role.Name,
	)
	if err != nil {
		return nil, err
	}

	user.Role = role.Name
	return &user, nil
}

func (r *repo) FindByTelegramId(ctx context.Context, telegramId int64) (*model.UserInfo, error) {
	var user model.UserInfo
	var role roleModelRepo.Role

	// Строим запрос с использованием JOIN для получения роли
	builder := sq.Select("u."+usersTelegramIdColumn, "u."+usersRoleIdColumn, "r."+rolesNameColumn).
		PlaceholderFormat(sq.Dollar).
		From(usersTableName + " u").
		Join("roles r ON u." + usersRoleIdColumn + " = r." + rolesIdColumn).
		Where(sq.Eq{"u." + usersTelegramIdColumn: telegramId}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.FindByTelegramId",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&user.TelegramID, &user.RoleId, &role.Name,
	)
	if err != nil {
		// Если запись не найдена, возвращаем false и nil для ошибки
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		// Любая другая ошибка обрабатывается как обычно
		return nil, err
	}

	user.Role = role.Name
	return &user, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.UserInfo, error) {
	var user model.UserInfo
	var role roleModelRepo.Role

	// Строим запрос с использованием JOIN для получения роли
	builder := sq.Select("u."+usersUsernameColumn, "u."+usersRoleIdColumn, "u."+usersPasswordHashColumn, "r."+rolesNameColumn).
		PlaceholderFormat(sq.Dollar).
		From(usersTableName + " u").
		Join("roles r ON u." + usersRoleIdColumn + " = r." + rolesIdColumn).
		Where(sq.Eq{"u." + usersIdColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&user.Username, &user.RoleId, &user.PasswordHash, &role.Name,
	)
	if err != nil {
		return nil, err
	}

	user.Role = role.Name
	return &user, nil
}

func (r *repo) Create(ctx context.Context, info model.UserInfo) (int64, error) {
	// Создаём SQL-конструктор
	builder := sq.Insert(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usersUsernameColumn, usersRoleIdColumn) // Начальные колонки
	values := []interface{}{info.Username, info.RoleId} // Начальные значения

	// Проверяем наличие пароля и добавляем его в запрос
	if info.PasswordHash != "" {
		builder = builder.Columns(usersPasswordHashColumn) // Добавляем колонку
		values = append(values, info.PasswordHash)         // Добавляем значение
	}

	// Проверяем наличие telegram_id и добавляем его в запрос
	if info.TelegramID != 0 {
		builder = builder.Columns(usersTelegramIdColumn) // Добавляем колонку
		values = append(values, info.TelegramID)         // Добавляем значение
	}

	// Проверяем наличие email и добавляем его в запрос
	if info.Email != "" {
		builder = builder.Columns(usersEmailColumn) // Добавляем колонку
		values = append(values, info.Email)         // Добавляем значение
	}

	// Если не переданы ни пароль, ни telegram_id, возвращаем ошибку
	if info.PasswordHash == "" && info.TelegramID == 0 {
		return 0, errors.New("either password_hash or telegram_id must be provided")
	}

	// Применяем значения к запросу
	builder = builder.Values(values...)

	// Добавляем RETURNING для получения id
	builder = builder.Suffix("RETURNING id")

	// Преобразуем builder в SQL-запрос
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	// Выполняем запрос и сканируем результат
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) CreateToken(ctx context.Context, info model.MailTokenInfo) error {
	// Создаём SQL-конструктор

	builder := sq.Insert(passwordResetTokensTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(resetTokenUserIdColumn, resetTokenColumn, resetTokenExpiresAtColumn).
		Values(info.UserID, info.Token, info.ExpiresAt)

	// Добавляем RETURNING для получения id
	builder = builder.Suffix("RETURNING id")

	// Преобразуем builder в SQL-запрос
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	// Выполняем запрос и сканируем результат
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindByEmailToken(ctx context.Context, token string) (*model.UserInfo, error) {
	var user model.UserInfo

	// Строим запрос с использованием JOIN для получения роли
	builder := sq.Select("t." + passwordResetTokensTableName).
		PlaceholderFormat(sq.Dollar).
		From(passwordResetTokensTableName + " t").
		Join("users r ON u." + usersIdColumn + " = t." + resetTokenIdColumn).
		Where(sq.Eq{"t." + resetTokenColumn: token}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.FindByEmailToken",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&user.ID,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) UpdateEmailVerified(ctx context.Context, userID int64) error {
	// Построение SQL запроса на обновление статуса подтверждения почты
	query, args, err := sq.Update(usersTableName).
		Set(usersEmailVerifiedColumn, true).
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// Выполнение запроса
	q := db.Query{
		Name:     "user_repository.UpdateEmailVerified",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}
