package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log      *zap.Logger
	initOnce sync.Once
)

func Init() error {
	var err error
	initOnce.Do(func() {
		var cfg zap.Config
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		log, err = cfg.Build()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка инициализации логгера: %v\n", err)
		}

		log.Info("Успешная инициализация логгера")
	})
	return err
}

func L() *zap.Logger {
	if log == nil {
		panic("Logger not initialized")
	}
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
