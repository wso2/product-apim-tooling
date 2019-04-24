package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/Jeffail/gabs"

	"gopkg.in/yaml.v2"
)

// Match for $VAR and capture VAR inside a group
var re = regexp.MustCompile(`\$(\w+)`)

// ErrRequiredEnvKeyMissing represents error used for indicate environment key missing
type ErrRequiredEnvKeyMissing struct {
	// Key is the missing entity
	Key string
}

func (e ErrRequiredEnvKeyMissing) Error() string {
	return fmt.Sprintf("%s is required, please set the environment variable", e.Key)
}

// Configuration represents endpoint config
type Configuration struct {
	// RetryTimeOut for endpoint
	RetryTimeOut *int `yaml:"retryTimeOut" json:"retryTimeOut,string"`
	// RetryDelay for endpoint
	RetryDelay *int `yaml:"retryDelay" json:"retryDelay,string"`
	// Factor used for config
	Factor *int `yaml:"factor" json:"factor,string"`
}

// Endpoint details
type Endpoint struct {
	// Url of the endpoint
	Url *string `yaml:"url" json:"url"`
	// Config of endpoint
	Config *Configuration `yaml:"config" json:"config"`
}

// EndpointData contains details about endpoints
type EndpointData struct {
	// Production endpoint
	Production *Endpoint `yaml:"production" json:"production_endpoints"`
	// Sandbox endpoint
	Sandbox *Endpoint `yaml:"sandbox" json:"sandbox_endpoints"`
}

// Environment represents an api environment
type Environment struct {
	// Name of the environment
	Name string `yaml:"name"`
	// Status of the API
	Status string `yaml:"status"`
	// Endpoints contain details about endpoints in a configuration
	Endpoints *EndpointData `yaml:"endpoints"`
}

// APIConfig represents environments defined in configuration file
type APIConfig struct {
	// Environments contains all environments in a configuration
	Environments []Environment `yaml:"environments"`
}

// APIEndpointConfig contains details about endpoints in an API
type APIEndpointConfig struct {
	// EPConfig is representing endpoint configuration
	EPConfig string `json:"endpointConfig"`
}

// injectEnv injects variables from environment to the content. It uses regex to match variables and look up them in the
// environment before processing.
// returns an error if anything happen
func injectEnv(content string) (string, error) {
	matches := re.FindAllStringSubmatch(content, -1) // matches is [][]string

	for _, match := range matches {
		Logln("Looking for: ", match[0])
		if os.Getenv(match[1]) == "" {
			return "", &ErrRequiredEnvKeyMissing{Key: match[0]}
		}
	}

	expanded := os.ExpandEnv(content)
	return expanded, nil
}

// LoadConfig loads an configuration from a reader. It returns an error or a valid APIConfig
func LoadConfig(r io.Reader) (*APIConfig, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	str, err := injectEnv(string(data))
	if err != nil {
		return nil, err
	}

	config := &APIConfig{}
	err = yaml.Unmarshal([]byte(str), &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// LoadConfigFromFile loads a configuration YAML file located in path. It returns an error or a valid APIConfig
func LoadConfigFromFile(path string) (*APIConfig, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	apiConfig, err := LoadConfig(r)
	_ = r.Close()

	return apiConfig, err
}

// LoadAPIFromFile loads API file from the path and returns a slice of bytes or an error
func LoadAPIFromFile(path string) ([]byte, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	_ = r.Close()

	return data, err
}

// ExtractAPIEndpointConfig extracts API endpoint information from a slice of byte b
func ExtractAPIEndpointConfig(b []byte) (string, error) {
	apiConfig := &APIEndpointConfig{}
	err := json.Unmarshal(b, &apiConfig)
	if err != nil {
		return "", err
	}

	return apiConfig.EPConfig, err
}

// MergeJSON secondSource with firstSource and returns merged JSON string
// Note: Fields in firstSource are merged with secondSource.
// If a field is not presented in secondSource, the one in firstSource will be preserved.
// If not a field from secondSource will replace it.
func MergeJSON(firstSource, secondSource []byte) ([]byte, error) {
	secondSourceJSON, err := gabs.ParseJSON(secondSource)
	if err != nil {
		return nil, err
	}

	firstSourceJSON, err := gabs.ParseJSON(firstSource)
	if err != nil {
		return nil, err
	}

	err = firstSourceJSON.MergeFn(secondSourceJSON, func(destination, source interface{}) interface{} {
		if source == nil {
			return destination
		}
		return source
	})

	return firstSourceJSON.Bytes(), nil
}

// GetEnv returns the EndpointData associated for key in the APIConfig, if not found returns nil
func (config APIConfig) GetEnv(key string) *Environment {
	for index, env := range config.Environments {
		if env.Name == key {
			return &config.Environments[index]
		}
	}
	return nil
}
