package models

import "time"

type RefreshToken struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	ClientIP         string
	Used             bool
	CreatedAt        time.Time
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
