package role

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"auth-service/internal/client/db"
	"auth-service/internal/model"
	"auth-service/internal/repository"
)

type repo struct {
	db db.Client
}

const (
	tableName  = "roles"
	idColumn   = "id"
	nameColumn = "role_name"
)

func NewRepository(db db.Client) repository.RoleRepository {
	return &repo{db: db}
}

func (r *repo) FindByName(ctx context.Context, name string) (*model.RoleInfo, error) {
	builder := sq.Select(idColumn, nameColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{nameColumn: name}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "role_repository.Get",
		QueryRaw: query,
	}

	var role model.RoleInfo
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}

	return &role, err
}

func (r *repo) ExistsByName(ctx context.Context, name string) (bool, error) {
	// Строим SQL-запрос
	builder := sq.Select("1").
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{nameColumn: name}).
		Limit(1)

	// Преобразуем запрос в строку и аргументы
	query, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	// Создаем структуру для запроса
	q := db.Query{
		Name:     "role_repository.ExistsByName",
		QueryRaw: query,
	}

	// Выполняем запрос
	var exists int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&exists)

	if err != nil {
		// Если запись не найдена, возвращаем false и nil для ошибки
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		// Любая другая ошибка обрабатывается как обычно
		return false, err
	}

	// Если нашли запись, возвращаем true
	return exists == 1, nil
}

func (r *repo) Create(ctx context.Context, name string) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn).
		Values(name).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "role_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
