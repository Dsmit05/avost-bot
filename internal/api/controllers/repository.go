package controllers

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"

	"github.com/Dsmit05/avost-bot/internal/models"
	"github.com/gin-gonic/gin"
)

type repositoryI interface {
	UpdateUserRole(userID int64, newRole uint8) bool
	UpdateUserSubManageType(userID int64, newSubManageType uint8) bool
	GetUsers() ([]models.User, bool)
	GetUsersAdmins() []models.User
}

type Repository struct {
	cache repositoryI
}

func NewRepository(cache repositoryI) *Repository {
	return &Repository{cache: cache}
}

// GetUsers
//
// @Summary      GetUsers
// @Tags         repository
// @Description  get all users
// @ID           get-users
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.User
// @Failure      400  {string}  string  "error"
// @Failure      417  {string}  string  "error"
// @Security     ApiKeyAuth
// @Router       /repository/users [GET]
func (r *Repository) GetUsers(c *gin.Context) {
	users, ok := r.cache.GetUsers()
	if !ok {
		c.String(400, "not have users")
		return
	}

	c.JSON(http.StatusOK, &users)
}

// GetUsersAdmins
//
// @Summary      GetUsersAdmins
// @Tags         repository
// @Description  get admins
// @ID           get-users-admins
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.User
// @Failure      400  {string}  string  "error"
// @Failure      417  {string}  string  "error"
// @Security     ApiKeyAuth
// @Router       /repository/admins [GET]
func (r *Repository) GetUsersAdmins(c *gin.Context) {
	users := r.cache.GetUsersAdmins()

	c.JSON(http.StatusOK, &users)
}

type UpdateUserRoleInput struct {
	ID   int64 `json:"id" binding:"required"`
	Role uint8 `json:"role"`
}

// UpdateUserRole
//
// @Summary      UpdateUserRole
// @Tags         repository
// @Description  update user role
// @ID           update-user-role
// @Accept       json
// @Produce      json
// @Param        input  body      UpdateUserRoleInput  true  "credentials"
// @Success      200    {string}  string               "data"
// @Failure      400    {string}  string               "error"
// @Failure      417    {string}  string               "error"
// @Security     ApiKeyAuth
// @Router       /repository/role [PATCH]
func (r *Repository) UpdateUserRole(c *gin.Context) {
	var inputData UpdateUserRoleInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ok := r.cache.UpdateUserRole(inputData.ID, inputData.Role)
	if !ok {
		c.String(http.StatusBadRequest, "Нет такого пользователя")
		return
	}

	c.String(http.StatusOK, "Роль успешно обновилась")
}

type UpdateUserSubInput struct {
	ID  int64 `json:"id" binding:"required"`
	Sub uint8 `json:"sub"`
}

// UpdateUserSub
//
// @Summary      UpdateUserSub
// @Tags         repository
// @Description  update user sub
// @ID           update-user-sub
// @Accept       json
// @Produce      json
// @Param        input  body      UpdateUserSubInput  true  "credentials"
// @Success      200    {string}  string              "data"
// @Failure      400    {string}  string              "error"
// @Failure      417    {string}  string              "error"
// @Security     ApiKeyAuth
// @Router       /repository/sub [PATCH]
func (r *Repository) UpdateUserSub(c *gin.Context) {
	var inputData UpdateUserSubInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ok := r.cache.UpdateUserSubManageType(inputData.ID, inputData.Sub)
	if !ok {
		c.String(http.StatusBadRequest, "Нет такого пользователя")
		return
	}

	c.String(http.StatusOK, "Роль успешно обновилась")
}

// GetUsersFile
//
// @Summary      GetUsersFile
// @Tags         repository
// @Description  get all users
// @ID           get-users-file
// @Accept       xml
// @Produce      xml
// @Success      200  {file}    string  "data"
// @Failure      400  {string}  string  "error"
// @Failure      417  {string}  string  "error"
// @Security     ApiKeyAuth
// @Router       /repository/users-file [GET]
func (r *Repository) GetUsersFile(c *gin.Context) {
	users, ok := r.cache.GetUsers()
	if !ok {
		c.String(400, "not have users")
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	f.NewSheet("users")
	f.SetCellValue("users", fmt.Sprintf("A%v", 1), "TelegramID")
	f.SetCellValue("users", fmt.Sprintf("B%v", 1), "Role")
	f.SetCellValue("users", fmt.Sprintf("C%v", 1), "SubManageType")

	for n, v := range users {
		f.SetCellValue("users", fmt.Sprintf("A%v", n+2), v.TelegramID)
		f.SetCellValue("users", fmt.Sprintf("B%v", n+2), v.Role)
		f.SetCellValue("users", fmt.Sprintf("C%v", n+2), v.SubManageType)
	}

	b, err := f.WriteToBuffer()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error %v", err))
		return
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "users.xls"))
	c.Data(http.StatusOK, "application/vnd.ms-excel", b.Bytes())
}
