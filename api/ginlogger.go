package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func logWithZeroLog(param gin.LogFormatterParams) string {
	msg := fmt.Sprintf("%s - %s %s %d %s %s",
		param.ClientIP,
		param.Method,
		param.Path,
		param.StatusCode,
		param.Latency,
		param.ErrorMessage,
	)
	log.Debug().Msg(msg)
	return ""
}
