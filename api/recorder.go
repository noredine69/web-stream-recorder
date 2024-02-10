package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	METRICS_ADD_ENTRY_END_POINT = "add_entry_endpoint"
	METRICS_ADD_ENTRY_NB_CALL   = "add_entry_endpoint_nb_call"
	METRICS_ADD_ENTRY_DESC      = "Metrics for Add Entry endpoint"
)

func (api *Api) declareRecorderRoutes() {
	privateRoutes := api.router.Group("/recorder/")
	{
		privateRoutes.GET("entry", api.scheduleRecording)
	}
	api.addGaugeMetricForEndpoint(METRICS_ADD_ENTRY_END_POINT, METRICS_ADD_ENTRY_NB_CALL, METRICS_ADD_ENTRY_DESC)
}

func (api *Api) scheduleRecording(ginContext *gin.Context) {
	api.incGaugeMetricForEndpoint(METRICS_ADD_ENTRY_END_POINT, METRICS_ADD_ENTRY_NB_CALL)
	channel, errChannel := strconv.ParseInt(ginContext.Query("channel"), 10, 0)
	if errChannel != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error parsing channel": errChannel.Error()})
		return
	}
	duration, errDuration := strconv.ParseInt(ginContext.Query("duration"), 10, 0)
	if errDuration != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error parsing duration": errDuration.Error()})
		return
	}
	delay, errDelay := strconv.ParseInt(ginContext.Query("delay"), 10, 0)
	if errDelay != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error parsing delay": errDelay.Error()})
		return
	}
	errRecord := api.recorderApi.ScheduleRecording(ginContext, delay, channel, duration)
	if errRecord != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errRecord.Error()})
	} else {
		ginContext.JSON(http.StatusOK, gin.H{"message": "Recording scheduled"})
	}
}
