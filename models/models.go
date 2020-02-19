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
	ID       int    `json:"id"`
	UserID   string `json:"user_id"`
	UpdateID int    `json:"update_id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	News     string `json:"news"`
	Date     string `json:"date"`
}

type FirebaseTokens struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type Timing struct {
	ID          int    `json:"id"`
	MobID       int    `json:"mob_id"`
	AccID       int    `json:"acc_id"`
	UserID      string `json:"user_id"`
	Date        string `json:"date"`
	Status      string `json:"status"`
	IsTurnstile bool   `json:"is_turnstile"`
	StartedAt   string `json:"started_at"`
	EndedAt     string `json:"ended_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}
