package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	METRICS_HEALTH_DESC      = "Metrics for health (liveness, readyness) endpoint"
	METRICS_HEALTH_END_POINT = "healthz_endpoint"
	METRICS_HEALTH_NB_CALL   = "healthz_nb_call"
)

func (api *Api) declareHealthRoutes() {
	_ = api.router.GET("/healthz", api.healthz)
	api.addGaugeMetricForEndpoint(METRICS_HEALTH_END_POINT, METRICS_HEALTH_NB_CALL, METRICS_HEALTH_DESC)
}

func (api *Api) healthz(c *gin.Context) {
	api.incGaugeMetricForEndpoint(METRICS_HEALTH_END_POINT, METRICS_HEALTH_NB_CALL)
	c.Status(http.StatusOK)
}
