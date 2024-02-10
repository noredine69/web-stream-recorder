package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-stream-recorder/services/config"
	"web-stream-recorder/services/recorder"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/rs/zerolog/log"
)

type Api struct {
	config         config.Config
	router         *gin.Engine
	recorderApi    recorder.RecorderAPIInterface
	metricsMonitor *ginmetrics.Monitor
}

func New(config config.Config) *Api {
	api := Api{config: config}
	api.recorderApi = recorder.New(config.Recorder)
	api.DeclareRoutes()
	return &api
}

const (
	formatUintBase = 10
)

func (api *Api) DeclareRoutes() {
	api.initGinEngine()
	api.declareBackEndRoutes()
}

func (api *Api) declareBackEndRoutes() {
	// Don't change the order, metrics routes must be declare first in order to be called by other endpoints
	api.declareMetricsRoutes()
	api.declareRecorderRoutes()
	api.declareHealthRoutes()
}

func (api *Api) initGinEngine() {
	if !api.config.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	api.router = gin.New()
	api.router.Use(gin.Recovery())
	api.router.Use(gin.LoggerWithFormatter(logWithZeroLog))

	// For profiling
	if api.config.DebugMode {
		pprof.Register(api.router)
	}
}

func (api *Api) Run() {
	port := strconv.FormatUint(uint64(api.config.Api.Port), formatUintBase)
	log.Info().Msg("Server Started on Port " + port)
	err := endless.ListenAndServe(":"+port, api.router)
	/*,
	csrf.Protect([]byte(api.config.Api.SecretKey),
		csrf.Secure(false),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Path("/"),
		csrf.ErrorHandler(http.HandlerFunc(csrfErrorHandlerFunc)))(api.router)
	)*/

	if err != nil {
		log.Error().Err(err).Msgf("Error while starting the web server")
	}
}

func csrfErrorHandlerFunc(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusForbidden)
	msg := "CSRF Token is invalid"
	log.Error().Msg(msg)
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msgf("Error while handling csrf token (Marshal)")
	}
	_, err = response.Write(jsonStr)
	if err != nil {
		log.Error().Err(err).Msgf("Error while handling csrf token (Write)")
	}
}
