package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var ZapLog *zap.Logger

func InitLogger(modeDev bool, logPaths string) error {
	var err error

	if logPaths == "" {
		return fmt.Errorf("the name of the logging file is not specified")
	}

	if modeDev {
		dev := zap.NewDevelopmentConfig()
		ZapLog, err = dev.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
	} else {
		prod := zap.NewDevelopmentConfig() // NewProductionConfig
		prod.OutputPaths = append(prod.OutputPaths, logPaths)
		ZapLog, err = prod.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
	}

	ZapLog.Info("logger is initialized")

	return err
}

func Info(message, text string) {
	ZapLog.Info(message, zap.String("text", text))
}

func Infof(message string, fields ...zap.Field) {
	ZapLog.Info(message, fields...)
}

func Debug(message string, val interface{}) {
	ZapLog.Debug(message, zap.Any("params", val))
}

func Error(message string, err error) {
	ZapLog.Error(message, zap.Error(err))
}

func Errorf(message string, fields ...zap.Field) {
	ZapLog.Error(message, fields...)
}

func ErrorAny(message string, value any) {
	ZapLog.Error(message, zap.Any("any", value))
}

func Fatal(message string, err error) {
	ZapLog.Fatal(message, zap.Error(err))
}
