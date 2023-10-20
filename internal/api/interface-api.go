package api

import (
	"time"

	"github.com/Dsmit05/avost-bot/internal/models"
)

type RepositoryI interface {
	UpdateUserRole(userID int64, newRole uint8) bool
	UpdateUserSubManageType(userID int64, newSubManageType uint8) bool
	GetUsers() ([]models.User, bool)
	GetUsersAdmins() []models.User
}

type BotI interface {
	Replay(msgID int, tgID int64, text string) error
	Send(tgID int64, text string) error
	SendAll(text string) error
	GetUserInfo(id int64) (*models.UserFullInfo, error)
}

type cryptographyI interface {
	CreateToken(email string, role string, ttl time.Duration) (string, error)
	ParseToken(inputToken string) (email string, role string, err error)
}
