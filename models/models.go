package models

import "time"

type Config struct {
	Auth AuthConfig
	DB   DBConfig
}

type AuthConfig struct {
	User string
	Password string
}

type DBConfig struct {
	Name string
	User string
	Password string
}

type Channel struct {
	ExtID int `json:"ext_id"`
	UserID string `json:"user_id"`
	UpdateID string `json:"update_id"`
	Type string `json:"type"`
	Title string `json:"title"`
	New string `json:"news"`
	DateTime time.Time `json:"date_time"`
}

type FirebaseTokens struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	Token string `json:"token"`
}