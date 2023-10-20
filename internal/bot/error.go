package bot

import (
	"github.com/Dsmit05/avost-bot/internal/logger"
	tele "gopkg.in/telebot.v3"
)

func LogError(err error, c tele.Context) {
	if c != nil {
		var id int64
		if c.Sender() != nil {
			id = c.Sender().ID
		}

		var callback string
		if c.Callback() != nil {
			callback = c.Callback().Unique
		}

		logger.ErrorUser("bot error", id, c.Text(), c.Data(), callback, err)
	} else {
		logger.Error("bot error without ctx", err)
	}
}
