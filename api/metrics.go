package api

import (
	"github.com/penglongli/gin-metrics/ginmetrics"
)

const (
	slowTime = 10
)

func (api *Api) declareMetricsRoutes() {
	api.metricsMonitor = ginmetrics.GetMonitor()
	api.metricsMonitor.SetMetricPath("/metrics")
	api.metricsMonitor.SetSlowTime(slowTime)
	api.metricsMonitor.Use(api.router)
}

func (api *Api) addGaugeMetricForEndpoint(metricName, metricFieldName, metricDescription string) {
	gaugeMetric := &ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        metricName,
		Description: metricDescription,
		Labels:      []string{metricFieldName},
	}
	_ = api.metricsMonitor.AddMetric(gaugeMetric)
}

func (api *Api) incGaugeMetricForEndpoint(metricName, metricFieldName string) {
	_ = api.metricsMonitor.GetMetric(metricName).Inc([]string{metricFieldName})
}
