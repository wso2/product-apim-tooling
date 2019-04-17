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

type ErrRequiredEnvKeyMissing struct {
	Key string
}

func (e ErrRequiredEnvKeyMissing) Error() string {
	return fmt.Sprintf("%s is required, please set the environment variable", e.Key)
}

type Configuration struct {
	RetryTimeOut *int `yaml:"retryTimeOut" json:"retryTimeOut,string"`
	RetryDelay   *int `yaml:"retryDelay" json:"retryDelay,string"`
	Factor       *int `yaml:"factor" json:"factor,string"`
}

type Endpoint struct {
	Url    *string        `yaml:"url" json:"url"`
	Config *Configuration `yaml:"config" json:"config"`
}

type EndpointData struct {
	Production *Endpoint `yaml:"production" json:"production_endpoints"`
	Sandbox    *Endpoint `yaml:"sandbox" json:"sandbox_endpoints"`
}

type Environment struct {
	Name      string        `yaml:"name"`
	Endpoints *EndpointData `yaml:"endpoints"`
}

type APIConfig struct {
	Environments []Environment `yaml:"environments"`
}

type APIEndpointConfig struct {
	EPConfig string `json:"endpointConfig"`
}

func injectEnv(str string) (string, error) {
	matches := re.FindAllStringSubmatch(str, -1) // matches is [][]string

	for _, match := range matches {
		Logln("Looking for: ", match[0])
		if os.Getenv(match[1]) == "" {
			return "", &ErrRequiredEnvKeyMissing{Key: match[0]}
		}
	}

	expanded := os.ExpandEnv(str)
	return expanded, nil
}

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

func LoadConfigFromFile(path string) (*APIConfig, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	apiConfig, err := LoadConfig(r)
	_ = r.Close()

	return apiConfig, err
}

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

func (config APIConfig) GetEnv(key string) *EndpointData {
	for index, env := range config.Environments {
		fmt.Println("Looking", index, env.Name)
		if env.Name == key {
			fmt.Println("found", index)
			return config.Environments[index].Endpoints
		}
	}
	return nil
}
