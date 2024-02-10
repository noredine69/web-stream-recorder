package api

import (
	"io/ioutil"
	"net/http"
	"testing"
	"web-stream-recorder/services/config"

	"github.com/stretchr/testify/assert"
)

var (
	ERROR_PAGE_NOT_FOUND   = "404 page not found"
	ETH_PLORER_URL         = "ethplorer/"
	ETH_PLORER_API_KEY     = "AZERTY1234"
	ETH_GET_LAST_BLOCK_FCT = "lastblock"
)

func TestEthApi(t *testing.T) {
	/*
		expectedLastBlockId := uint(123456789)
		ethResponse := eth.LastBlock{
			LastBlockId: expectedLastBlockId,
		}
		expectedJsonStream, _ := json.Marshal(ethResponse)

		mux := http.NewServeMux()
		mux.HandleFunc(fmt.Sprintf("/%s%s", ETH_PLORER_URL, ETH_GET_LAST_BLOCK_FCT), func(w http.ResponseWriter, r *http.Request) {
			w.Write(expectedJsonStream)
		})
		muxServer := httptest.NewServer(mux)
		defer muxServer.Close()

		var INFURIA_CLOUD_ADDRESS = fmt.Sprintf("%s/%s", muxServer.URL, ETH_PLORER_URL)

		conf := initConfigHelper(INFURIA_CLOUD_ADDRESS)
		controller := New(conf)
		ts := httptest.NewServer(controller.router)
		defer ts.Close()

		body := checkLogsRouteCallStatusError(t, fmt.Sprintf("%s/rest/bad_endpoint", ts.URL), "GET")
		assert.Equal(t, ERROR_PAGE_NOT_FOUND, body)

		body = checkLogsRouteCallStatusOk(t, fmt.Sprintf("%s/eth/%s", ts.URL, ETH_GET_LAST_BLOCK_FCT), "GET")
		var result eth.LastBlock
		errUnMarshall := json.Unmarshal([]byte(body), &result)
		assert.Nil(t, errUnMarshall)
		assert.Equal(t, expectedLastBlockId, result.LastBlockId)
	*/
}

func initConfigHelper(cloudServiceAddress string) config.Config {
	return config.Config{
		/*
			Eth: config.EthConfig{
				Url:      cloudServiceAddress,
				ApiKey:   ETH_PLORER_API_KEY,
				Function: ETH_GET_LAST_BLOCK_FCT,
			},
		*/
	}
}

func checkLogsRouteCallStatusError(t *testing.T, url string, verb string) string {
	// nolint: noctx
	req, _ := http.NewRequest(verb, url, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	result := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted
	assert.False(t, result)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	return string(bodyBytes)
}

func checkLogsRouteCallStatusOk(t *testing.T, url string, verb string) string {
	// nolint: noctx
	req, _ := http.NewRequest(verb, url, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	result := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted
	assert.True(t, result)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	return string(bodyBytes)
}
