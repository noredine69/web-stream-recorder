package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
}

func NewGormLogger() logger.Interface {
	return &gormLogger{}
}

func (logger *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Info().Msg(msg)
}

func (logger gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warn().Msg(msg)
}

func (logger gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Error().Msg(msg)
}

func (logger gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	if rows == -1 {
		log.Trace().Msgf("[-], %v", sql)
	} else {
		log.Trace().Msgf("[%v], %v", rows, sql)
	}
}
