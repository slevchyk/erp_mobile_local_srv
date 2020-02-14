package models

type Config struct {
	Auth AuthConfig
	DB   DBConfig
}

type AuthConfig struct {
	User     string
	Password string
}

type DBConfig struct {
	Name     string
	User     string
	Password string
}

type Channel struct {
	ID       int       `json:"id"`
	UserID   string    `json:"user_id"`
	UpdateID string    `json:"update_id"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	News     string    `json:"news"`
	DateTime string `json:"date"`
}

type FirebaseTokens struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}
