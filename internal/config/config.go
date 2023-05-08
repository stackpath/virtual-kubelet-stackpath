package config

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/virtual-kubelet/virtual-kubelet/log"
	"gopkg.in/yaml.v2"
)

const (
	// default to the official StackPath API
	defaultAPIHost = "gateway.stackpath.com"
)

// Config is the provider's configuration
type Config struct {
	// A string specifies the unique identifier of your StackPath account.
	// This field is required for the provider to authenticate with
	// the StackPath API and perform operations on your behalf.
	AccountID string `yaml:"account_id"`

	// A string that specifies the unique identifier of your StackPath Edge Computing stack.
	// This field is required for the provider to create and manage virtual nodes within your stack.
	StackID string `yaml:"stack_id"`

	// A string that specifies the client ID of your StackPath API credentials.
	// This field is required for the provider to authenticate with the StackPath API and
	// perform operations on your behalf.
	ClientID string `yaml:"client_id"`

	// A string that specifies the client secret of your StackPath API credentials.
	// This field is required for the provider to authenticate with the StackPath API and
	// perform operations on your behalf.
	ClientSecret string `yaml:"client_secret"`

	// A string that specifies the API host of the StackPath Edge Computing API.
	// This field is optional and defaults to the production API URL.
	ApiHost string `yaml:"base_url"`

	// A string that specifies the StackPath Point of Presence (PoP) where pods
	// will be created, updated, or deleted. This setting allows you to take
	// advantage of StackPath's global network of edge servers to deploy workloads
	// closer to end-users and reduce latency. The value of this field should correspond
	// to the code name or identifier of a specific StackPath edge location, such as
	// "lax", "dfw", "ord", "iad", "atl", "mia", "ams", "fra", "cdg", "sin", "nrt", etc
	CityCode string `yaml:"city_code"`
}

// NewConfig creates and loads configuration from either a YAML file or environment variables
func NewConfig(ctx context.Context) (*Config, error) {
	c := &Config{}

	if configFilepath := os.Getenv("SP_CONFIG_LOCATION"); configFilepath != "" {
		log.G(ctx).Info("getting StackPath config from YAML file: %s", configFilepath)
		if err := c.newConfigFromFile(configFilepath); err != nil {
			return nil, err
		}
		return c, nil
	}

	c.StackID = os.Getenv("SP_STACK_ID")
	c.ClientID = os.Getenv("SP_CLIENT_ID")
	c.ClientSecret = os.Getenv("SP_CLIENT_SECRET")
	c.ApiHost = os.Getenv("SP_API_HOST")
	c.CityCode = strings.ToUpper(os.Getenv("SP_CITY_CODE"))

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) newConfigFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	if err = c.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates the configuration parameters
func (config *Config) Validate() error {
	if config.StackID == "" {
		return errors.New("must provide a stack ID")
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		// Require the user to provide client ID and client secret
		return errors.New("must provide a client ID and client secret")
	}

	if config.CityCode == "" {
		return errors.New("must provide a city code")
	} else if !isValidLocation(config.CityCode) {
		return errors.New("must provide a valid city code")
	}

	if config.ApiHost == "" {
		// if the API host is not set, use the default one
		config.ApiHost = defaultAPIHost
	} else {
		// Validate the API host
		if strings.HasSuffix(config.ApiHost, "/") {
			return errors.New("API host must not end in a trailing slash")
		}
	}

	return nil
}
