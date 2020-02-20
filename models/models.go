package models

import (
	"fmt"
)

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
	ID          int64              `json:"id"`
	MobID       int64              `json:"mob_id"`
	AccID       string             `json:"acc_id"`
	UserID      string             `json:"user_id"`
	Date        NullTime           `json:"date"`
	Status      string             `json:"status"`
	IsTurnstile ConvertibleBoolean `json:"is_turnstile"`
	StartedAt   NullTime           `json:"started_at"`
	EndedAt     NullTime           `json:"ended_at"`
	CreatedAt   NullTime           `json:"created_at"`
	UpdatedAt   NullTime           `json:"updated_at"`
	DeletedAt   NullTime           `json:"deleted_at"`
}
