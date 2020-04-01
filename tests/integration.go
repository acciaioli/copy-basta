package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func HealthCheckTest(envPrefix string) func(*testing.T) {
	return func(t *testing.T) {
		config, err := LoadConfig(envPrefix)
		require.Nil(t, err)

		c := http.Client{Timeout: 10 * time.Second}

		endpoint := "/api/healthcheck"
		url, err := config.BuildURL(endpoint)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodGet, *url, nil)
		require.Nil(t, err)
		t.Logf("Test Request: %s %s", req.Method, req.URL)

		resp, err := c.Do(req)
		require.Nil(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		data := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		require.Nil(t, err)

		require.Equal(t, 2, len(data))

		status, ok := data["status"]
		require.True(t, ok)
		require.Equal(t, "ok", status)

		_, ok = data["availableRoutes"]
		require.True(t, ok)
	}
}

func NotFoundTest(envPrefix string) func(t *testing.T) {
	return func(t *testing.T) {
		config, err := LoadConfig(envPrefix)
		require.Nil(t, err)

		c := http.Client{Timeout: 10 * time.Second}

		endpoint := "/api/banana404"
		url, err := config.BuildURL(endpoint)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodGet, *url, nil)
		require.Nil(t, err)
		t.Logf("Test Request: %s %s", req.Method, req.URL)

		resp, err := c.Do(req)
		require.Nil(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		data := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		require.Nil(t, err)

		require.Equal(t, 2, len(data))

		status, ok := data["requestedPath"]
		require.True(t, ok)
		require.Equal(t, endpoint, status)

		_, ok = data["availableRoutes"]
		require.True(t, ok)
	}
}
