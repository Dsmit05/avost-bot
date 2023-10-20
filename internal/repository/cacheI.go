package repository

import (
	"github.com/Dsmit05/avost-bot/internal/models"
)

type CacheInterface interface {
	AddUser(user models.User) bool
	GetUser(userID int64) (models.User, bool)
	DelUser(userID int64) bool
	UpdateUserRole(userID int64, newRole uint8) bool
	UpdateUserSubManageType(userID int64, newSubManageType uint8) bool
	GetUserSubManageType(userID int64) uint8
	AddAnime(userID int64, animeID int) bool
	DelAnime(userID int64, animeID int) bool
	GetAllUserAnime(userID int64) ([]int, bool)
	GetAllMap() map[int64]models.User
	CheckUsersAnime(userID int64, animeID int) bool
	GetUsersID() ([]int64, bool)
	GetUsers() ([]models.User, bool)
	GetUsersAdmins() []models.User
	GetCountUsers() int
	CheckUserRole(id int64, role uint8) bool
}
