package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

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
	Url    *string        `yaml:"url"`
	Config *Configuration `yaml:"config"`
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

func injectEnv(str string) (string, error) {
	matches := re.FindAllStringSubmatch(str, -1) // matches is [][]string

	for _, match := range matches {
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

	return LoadConfig(r)
}
