package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Dsmit05/avost-bot/internal/models"
	"github.com/gin-gonic/gin"
)

type botI interface {
	Replay(msgID int, tgID int64, text string) error
	Send(tgID int64, text string) error
	SendAll(text string) error
	GetUserInfo(id int64) (*models.UserFullInfo, error)
}

type statI interface {
	Chart(start string, stop string, delim string, w io.Writer) error
}

type Bot struct {
	bot  botI
	stat statI
}

func NewBot(bot botI, stat statI) *Bot {
	return &Bot{bot: bot, stat: stat}
}

type SendMsgUserInput struct {
	TgID int64  `json:"tgid" binding:"required"`
	Text string `json:"text" binding:"required"`
}

// SendMsgUser
//
// @Summary      SendMsgUser
// @Tags         bot
// @Description  Send message user
// @ID           protected-send-msg-user
// @Accept       json
// @Produce      json
// @Param        input  body      SendMsgUserInput  true  "credentials"
// @Success      200    {string}  string            "data"
// @Failure      400    {object}  string            "error"
// @Failure      417    {object}  string            "error"
// @Security     ApiKeyAuth
// @Router       /bot/send [POST]
func (b *Bot) SendMsgUser(c *gin.Context) {
	var inputData SendMsgUserInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := b.bot.Send(inputData.TgID, inputData.Text); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, "Сообщение отправлено")
}

type ReplayMsgUserInput struct {
	MsgID int    `json:"msgid" binding:"required"`
	TgID  int64  `json:"tgid" binding:"required"`
	Text  string `json:"text" binding:"required"`
}

// ReplayMsgUser
//
// @Summary      ReplayMsgUser
// @Tags         bot
// @Description  Replay msg user
// @ID           protected-replay-msg-user
// @Accept       json
// @Produce      json
// @Param        input  body      ReplayMsgUserInput  true  "credentials"
// @Success      200    {string}  string              "data"
// @Failure      400    {string}  string              "error"
// @Failure      417    {string}  string              "error"
// @Security     ApiKeyAuth
// @Router       /bot/replay [POST]
func (b *Bot) ReplayMsgUser(c *gin.Context) {
	var inputData ReplayMsgUserInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := b.bot.Replay(inputData.MsgID, inputData.TgID, inputData.Text); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, "Сообщение отправлено")
}

type SendAllMsgUserInput struct {
	Text string `json:"text" binding:"required"`
}

// SendAllMsgUser
//
// @Summary      SendAllMsgUser
// @Tags         bot
// @Description  Send all msg User
// @ID           protected-send-all-msg-user
// @Accept       json
// @Produce      json
// @Param        input  body      SendAllMsgUserInput  true  "credentials"
// @Success      200    {string}  string               "data"
// @Failure      400    {string}  string               "error"
// @Failure      417    {string}  string               "error"
// @Security     ApiKeyAuth
// @Router       /bot/sends [POST]
func (b *Bot) SendAllMsgUser(c *gin.Context) {
	var inputData SendAllMsgUserInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := b.bot.SendAll(inputData.Text); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, "Сообщения отправлены")
}

// GetUserInfo
//
// @Summary      GetUserInfo
// @Tags         bot
// @Description  get user info
// @ID           get-user-info
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  []models.UserFullInfo
// @Failure      400  {string}  string  "error"
// @Failure      417  {string}  string  "error"
// @Security     ApiKeyAuth
// @Router       /bot/user/{id} [GET]
func (b *Bot) GetUserInfo(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.String(http.StatusBadRequest, "not have params")
		return
	}

	id, err := strconv.ParseInt(userID, 10, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid params")
		return
	}

	users, err := b.bot.GetUserInfo(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, &users)
}

// GetStat
//
// @Summary      GetStat
// @Tags         bot
// @Description  get stat
// @ID           get-stat-info
// @Accept       html
// @Produce      html
// @Param        start   query     string  false  "start at"       Format(date)
// @Param        stop    query     string  false  "end at"         Format(date)
// @Param        format  query     string  false  "one of h or m"  Format(string)
// @Success      200     {html}    html    "data"
// @Failure      400     {string}  string  "error"
// @Failure      417     {string}  string  "error"
// @Security     ApiKeyAuth
// @Router       /bot/stat [GET]
func (b *Bot) GetStat(c *gin.Context) {
	start := c.Query("start")
	stop := c.Query("stop")
	format := c.Query("format")

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "data.html"))
	err := b.stat.Chart(start, stop, format, c.Writer)
	if err != nil {
		c.String(400, "error %v", err)
	}
}
