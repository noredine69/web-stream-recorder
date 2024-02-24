package api

import (
	"web-stream-recorder/services/config"

	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	conf := initConfigHelperForHealth()
	api := New(conf, nil)
	ts := httptest.NewServer(api.router)
	defer ts.Close()
	body := checkLogsRouteCallStatusOk(t, fmt.Sprintf("%s/healthz", ts.URL), "GET")
	assert.Equal(t, "", body)
}

func initConfigHelperForHealth() config.Config {
	return config.Config{}
}
