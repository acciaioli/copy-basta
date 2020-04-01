package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	var err error
	err = os.Setenv("TESTING_SCHEME", "http")
	require.Nil(t, err)
	err = os.Setenv("TESTING_HOST", "0.0.0.0")
	require.Nil(t, err)
	err = os.Setenv("TESTING_PORT", "3000")
	require.Nil(t, err)
	err = os.Setenv("TESTING_ENVNAMESPACE", "dev")
	require.Nil(t, err)

	config, err := LoadConfig("testing")
	require.Nil(t, err)

	require.Equal(t, "http", config.Scheme)
	require.Equal(t, "0.0.0.0", config.Host)
	require.Equal(t, 3000, *config.Port)
	require.Equal(t, "dev", *config.EnvNamespace)
}

func TestConfig_BuildURL_error(t *testing.T) {
	config := Config{
		Scheme: "http",
		Host:   "0.0.0.0",
	}
	tests := []struct {
		name     string
		endpoint string
	}{
		{name: "empty", endpoint: ""},
		{name: "no starting slash", endpoint: "api/tests"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := config.BuildURL(tt.endpoint)
			require.NotNil(t, err)
		})
	}
}

func TestConfig_BuildURL(t *testing.T) {
	endpoint := "/api/tests"

	tests := []struct {
		name        string
		config      Config
		expectedURL string
	}{
		{
			name: "local",
			config: Config{
				Scheme:       "http",
				Host:         "0.0.0.0",
				Port:         func() *int { v := 3000; return &v }(),
				EnvNamespace: nil,
			},
			expectedURL: "http://0.0.0.0:3000/api/tests",
		},
		{
			name: "dev",
			config: Config{
				Scheme:       "https",
				Host:         "101address.domain",
				Port:         nil,
				EnvNamespace: func() *string { v := "dev"; return &v }(),
			},
			expectedURL: "https://101address.domain/dev/api/tests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			url, err := tt.config.BuildURL(endpoint)
			require.Nil(t, err)
			require.Equal(t, tt.expectedURL, *url)
		})
	}
}
