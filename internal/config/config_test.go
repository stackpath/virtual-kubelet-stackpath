package config

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestNewConfigFromFile(t *testing.T) {
	testCases := []struct {
		description    string
		configFilename string
		init           func(string)
		expectedConfig *Config
		expectedError  error
	}{
		{
			description:    "fails to read config data from a file that doesn't exist",
			configFilename: "config_doesnt_exist.yaml",
			init:           func(configFilename string) { os.Setenv("SP_CONFIG_LOCATION", configFilename) },
			expectedConfig: nil,
			expectedError:  fmt.Errorf("open config_doesnt_exist.yaml: no such file or directory"),
		},
		{
			description:    "fails to read config data due to bad YAML format of the data in the file.",
			configFilename: "bad_content_config.yaml",
			init: func(configFilename string) {
				os.Setenv("SP_CONFIG_LOCATION", configFilename)
				err := os.WriteFile(configFilename, []byte("not YAML"), 0777)
				if err != nil {
					t.Fatalf("couldn't write data to the file %s", configFilename)
				}
			},
			expectedConfig: nil,
			expectedError:  fmt.Errorf("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `not YAML` into config.Config"),
		},
		{
			description:    "fails to read config data because some parameters are missing.",
			configFilename: "invalid_config.yaml",
			init: func(configFilename string) {
				os.Setenv("SP_CONFIG_LOCATION", configFilename)
				invalidConfig := Config{
					StackID:      "",
					ClientID:     "",
					ApiHost:      "",
					ClientSecret: "",
					CityCode:     "",
				}
				data, err := yaml.Marshal(&invalidConfig)
				if err != nil {
					t.Fatal("couldn't serialize the config")
				}
				err = os.WriteFile(configFilename, data, 0777)
				if err != nil {
					t.Fatalf("couldn't write data to the file %s", configFilename)
				}
			},
			expectedConfig: nil,
			expectedError:  fmt.Errorf("must provide a stack ID"),
		},
		{
			description:    "successfully loads config from a file.",
			configFilename: "valid_config.yaml",
			init: func(configFilename string) {
				os.Setenv("SP_CONFIG_LOCATION", configFilename)
				config := Config{
					AccountID:    "a7188caa-e29b-11ed-b5ea-0242ac120002",
					StackID:      "a7188caa-e29b-11ed-b5ea-0242ac120003",
					ClientID:     "123",
					ApiHost:      "",
					ClientSecret: "123",
					CityCode:     "DFW",
				}
				data, err := yaml.Marshal(&config)
				if err != nil {
					t.Fatal("couldn't serialize the config")
				}
				err = os.WriteFile(configFilename, data, 0777)
				if err != nil {
					t.Fatalf("couldn't write the config to the file %s", configFilename)
				}
			},
			expectedConfig: &Config{
				AccountID:    "a7188caa-e29b-11ed-b5ea-0242ac120002",
				StackID:      "a7188caa-e29b-11ed-b5ea-0242ac120003",
				ClientID:     "123",
				ApiHost:      "gateway.stackpath.com",
				ClientSecret: "123",
				CityCode:     "DFW",
			},
			expectedError: nil,
		},
	}

	ctx := context.TODO()

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			c.init(c.configFilename)
			actualConfig, err := NewConfig(ctx)
			defer os.Remove(c.configFilename)
			defer os.Unsetenv("SP_CONFIG_LOCATION")

			if err != nil {
				assert.Equal(t, c.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, c.expectedError, nil)
			}
			assert.Equal(t, c.expectedConfig, actualConfig)
		})
	}
}

func TestNewConfigFromEnvVars(t *testing.T) {

	testCases := []struct {
		stackID       string
		clientID      string
		apiHost       string
		clientSecret  string
		cityCode      string
		expectedError error
	}{
		{
			stackID:       "",
			clientID:      "",
			apiHost:       "",
			clientSecret:  "",
			cityCode:      "",
			expectedError: fmt.Errorf("must provide a stack ID"),
		},
		{
			stackID:       "a7188caa-e29b-11ed-b5ea-0242ac120003",
			clientID:      "",
			apiHost:       "",
			clientSecret:  "",
			cityCode:      "",
			expectedError: fmt.Errorf("must provide a client ID and client secret"),
		},
		{
			stackID:       "a7188caa-e29b-11ed-b5ea-0242ac120003",
			clientID:      "123",
			apiHost:       "",
			clientSecret:  "",
			cityCode:      "",
			expectedError: fmt.Errorf("must provide a client ID and client secret"),
		},
		{
			stackID:       "a7188caa-e29b-11ed-b5ea-0242ac120003",
			clientID:      "",
			apiHost:       "",
			clientSecret:  "123",
			cityCode:      "",
			expectedError: fmt.Errorf("must provide a client ID and client secret"),
		},
		{
			stackID:       "a7188caa-e29b-11ed-b5ea-0242ac120003",
			clientID:      "123",
			apiHost:       "",
			clientSecret:  "1234",
			cityCode:      "",
			expectedError: fmt.Errorf("must provide a city code"),
		},
		{
			stackID:       "a7188caa-e29b-11ed-b5ea-0242ac120003",
			clientID:      "123",
			apiHost:       "",
			clientSecret:  "1234",
			cityCode:      "not valid",
			expectedError: fmt.Errorf("must provide a valid city code"),
		},
		{
			stackID:       "a7188caa-e29b-11ed-b5ea-0242ac120003",
			clientID:      "123",
			apiHost:       "asd/",
			clientSecret:  "1234",
			cityCode:      "dfw",
			expectedError: fmt.Errorf("API host must not end in a trailing slash"),
		},
		{
			stackID:       "a7188caa-e29b-11ed-b5ea-0242ac120003",
			clientID:      "123",
			apiHost:       "",
			clientSecret:  "1234",
			cityCode:      "dfw",
			expectedError: nil,
		},
	}

	ctx := context.TODO()

	for _, c := range testCases {
		os.Setenv("SP_STACK_ID", c.stackID)
		os.Setenv("SP_CLIENT_ID", c.clientID)
		os.Setenv("SP_CLIENT_SECRET", c.clientSecret)
		os.Setenv("SP_API_HOST", c.apiHost)
		os.Setenv("SP_CITY_CODE", c.cityCode)

		_, err := NewConfig(ctx)
		if c.expectedError != nil || err != nil {
			assert.Equal(t, c.expectedError.Error(), err.Error())
		}
	}
}
