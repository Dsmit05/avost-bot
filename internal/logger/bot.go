package logger

import "go.uber.org/zap"

func ErrorUser(message string, userID int64, text string, data string, callback string, err error) {
	ZapLog.Error(
		message,
		zap.Int64("userID", userID),
		zap.String("data", data),
		zap.String("text", text),
		zap.String("callback", callback),
		zap.String("err", err.Error()))
}
