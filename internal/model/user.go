package model

type UserInfo struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	RoleId       int64  `json:"role_id"`
	PasswordHash string `json:"password_hash"`
	TelegramID   int64  `json:"telegram_id"`
	Email        string `json:"email"`
}
