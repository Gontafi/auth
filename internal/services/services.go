package services

import (
	"auth/internal/models"
	"auth/internal/notify"
	"auth/internal/repo"
	"auth/internal/utils"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Services struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Services {
	return &Services{repo: repo}
}

func (s Services) Authenticate(ctx context.Context, userID string, userIP string) (*models.LoginResponse, error) {
	accessToken, err := utils.GenerateToken(
		userID,
		userIP, time.Now().Add(time.Minute*30),
		os.Getenv("SECRET_KEY"),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRandomString(50)
	if err != nil {
		return nil, err
	}
	hashedToken, err := utils.HashToken(refreshToken)
	if err != nil {
		return nil, err
	}

	err = s.repo.AddToken(ctx,
		models.RefreshToken{
			UserID:           userID,
			RefreshTokenHash: hashedToken,
			ClientIP:         userIP,
		})
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s Services) RefreshToken(ctx context.Context, userID string, userIP string, refreshToken string) (*models.LoginResponse, error) {
	token, err := s.repo.GetToken(ctx, userID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if token == nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	err = bcrypt.CompareHashAndPassword([]byte(token.RefreshTokenHash), []byte(refreshToken))
	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	if token.UserID != userID {
		return nil, fmt.Errorf("Unauthorized")
	}

	if token.ClientIP != userIP {
		notify.SendEmailMessage(fmt.Sprintf("Wrong IP address before:%s after:%s", token.ClientIP, userIP))
	}

	err = s.repo.UpdateTokenUsedInfo(ctx, token.ID)
	if err != nil {
		return nil, err
	}

	return s.Authenticate(ctx, userID, userIP)
}
