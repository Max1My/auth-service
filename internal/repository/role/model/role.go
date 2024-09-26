package model

type Role struct {
	ID   int64  `db:"id"`
	Name string `db:"role_name"`
}
