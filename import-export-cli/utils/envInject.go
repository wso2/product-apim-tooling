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
	Production *Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Sandbox endpoint
	Sandbox *Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
}

// Cert stores certificate details
type Cert struct {
	Host        string `yaml:"host" json:"hostName"`
	Alias       string `yaml:"alias" json:"alias"`
	Path        string `yaml:"path" json:"-"`
	Certificate string `json:"certificate"`
}

// Environment represents an api environment
type Environment struct {
	// Name of the environment
	Name string `yaml:"name"`
	// Endpoints contain details about endpoints in a configuration
	Endpoints *EndpointData `yaml:"endpoints"`
	// Certs for environment
	Certs []Cert `yaml:"certs"`
}

// ApiParams represents environments defined in configuration file
type ApiParams struct {
	// Environments contains all environments in a configuration
	Environments []Environment `yaml:"environments"`
}

// APIEndpointConfig contains details about endpoints in an API
type APIEndpointConfig struct {
	// EPConfig is representing endpoint configuration
	EPConfig string `json:"endpointConfig"`
}

// InjectEnv injects variables from environment to the content. It uses regex to match variables and look up them in the
// environment before processing.
// returns an error if anything happen
func InjectEnv(content string) (string, error) {
	matches := re.FindAllStringSubmatch(content, -1) // matches is [][]string

	for _, match := range matches {
		Logln("Looking for:", match[0])
		if os.Getenv(match[1]) == "" {
			return "", &ErrRequiredEnvKeyMissing{Key: match[0]}
		}
	}

	expanded := os.ExpandEnv(content)
	return expanded, nil
}

// LoadApiParams loads an configuration from a reader. It returns an error or a valid ApiParams
func LoadApiParams(r io.Reader) (*ApiParams, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	str, err := InjectEnv(string(data))
	if err != nil {
		return nil, err
	}

	apiParams := &ApiParams{}
	err = yaml.Unmarshal([]byte(str), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, nil
}

// LoadApiParamsFromFile loads a configuration YAML file located in path. It returns an error or a valid ApiParams
func LoadApiParamsFromFile(path string) (*ApiParams, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	apiConfig, err := LoadApiParams(r)
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
		if s, ok := source.(string); ok && s == "" {
			return destination
		}
		return source
	})

	return firstSourceJSON.Bytes(), nil
}

// GetEnv returns the EndpointData associated for key in the ApiParams, if not found returns nil
func (config ApiParams) GetEnv(key string) *Environment {
	for index, env := range config.Environments {
		if env.Name == key {
			return &config.Environments[index]
		}
	}
	return nil
}
