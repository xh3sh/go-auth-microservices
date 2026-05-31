package models

// User представляет публичный профиль пользователя
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
