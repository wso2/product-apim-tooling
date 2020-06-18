package params

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
)

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

// SecurityData contains the details about endpoint security from api_params.yaml
type SecurityData struct {
	// Decides whether the endpoint security is enabled
	Enabled string `yaml:"enabled" json:"enabled,omitempty"`
	// Type of the endpoint security (can be Basic or Digest)
	Type string `yaml:"type" json:"type,omitempty"`
	// Username for the endpoint
	Username string `yaml:"username" json:"username,omitempty"`
	// Password for the endpoint
	Password string `yaml:"password" json:"password,omitempty"`
}

// Cert stores certificate details
type Cert struct {
	// Host of the certificate
	Host string `yaml:"host" json:"hostName"`
	// Alias for certificate
	Alias string `yaml:"alias" json:"alias"`
	// Path for certificate file
	Path string `yaml:"path" json:"-"`
	// Certificate is used for internal purposes, it contains secret in base64
	Certificate string `json:"certificate"`
}

// Environment represents an api environment
type Environment struct {
	// Name of the environment
	Name string `yaml:"name"`
	// Endpoints contain details about endpoints in a configuration
	Endpoints *EndpointData `yaml:"endpoints"`
	// Security contains the details about endpoint security
	Security *SecurityData `yaml:"security"`
	// GatewayEnvironments contains environments that used to deploy API
	GatewayEnvironments []string `yaml:"gatewayEnvironments"`
	// Certs for environment
	Certs []Cert `yaml:"certs"`
}

// ApiParams represents environments defined in configuration file
type ApiParams struct {
	// Environments contains all environments in a configuration
	Environments []Environment   `yaml:"environments"`
	Import       APIImportParams `yaml:"import"`
}

type ApiProductParams struct {
	Import APIProductImportParams `yaml:"import"`
}

type ApplicationParams struct {
	Import ApplicationImportParams `yaml:"import"`
}

// ------------------- Structs for Import Params ----------------------------------
type APIImportParams struct {
	Update           bool `yaml:"update"`
	PreserveProvider bool `yaml:"preserveProvider"`
}

type APIProductImportParams struct {
	ImportAPIs       bool `yaml:"importApis"`
	UpdateAPIs       bool `yaml:"updateApis"`
	UpdateAPIProduct bool `yaml:"updateApiProduct"`
	PreserveProvider bool `yaml:"preserveProvider"`
}

type ApplicationImportParams struct {
	Update            bool   `yaml:"update"`
	TargetOwner       string `yaml:"targetOwner"`
	PreserveOwner     bool   `yaml:"preserveOwner"`
	SkipKeys          bool   `yaml:"skipKeys"`
	SkipSubscriptions bool   `yaml:"skipSubscriptions"`
}

type ProjectParams struct {
	Type                     string
	AbsolutePath             string
	RelativePath             string
	Name                     string
	FailedDuringPreviousPush bool
	Deleted                  bool
	ApiParams                *ApiParams
	ApiProductParams         *ApiProductParams
	ApplicationParams         *ApplicationParams
}
// ---------------- End of Structs for Project Details ---------------------------------

// APIEndpointConfig contains details about endpoints in an API
type APIEndpointConfig struct {
    // EPConfig is representing endpoint configuration
	EPConfig string `json:"endpointConfig"`
}

// loads the given file in path and substitutes environment variables that are defined as ${var} or $var in the file.
//	returns the file as string.
func getEnvSubstitutedFileContent(path string) (string, error) {
	r, err := os.Open(path)
	defer func() {
		_ = r.Close()
	}()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	str, err := utils.EnvSubstitute(string(data))
	if err != nil {
		return "", err
	}
	return str, nil
}

// LoadApiParamsFromFile loads an API Project configuration YAML file located in path.
//	It returns an error or a valid ApiParams
func LoadApiParamsFromFile(path string) (*ApiParams, error) {
	fileContent, err := getEnvSubstitutedFileContent(path)
	if err != nil {
		return nil, err
	}

	apiParams := &ApiParams{}
	err = yaml.Unmarshal([]byte(fileContent), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, err
}

// LoadApiProductParamsFromFile loads an API Product project configuration YAML file located in path.
//	It returns an error or a valid ApiProductParams
func LoadApiProductParamsFromFile(path string) (*ApiProductParams, error) {
	fileContent, err := getEnvSubstitutedFileContent(path)
	if err != nil {
		return nil, err
	}

	apiParams := &ApiProductParams{}
	err = yaml.Unmarshal([]byte(fileContent), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, err
}

// LoadApplicationParamsFromFile loads an Application project configuration YAML file located in path.
//	It returns an error or a valid ApplicationParams
func LoadApplicationParamsFromFile(path string) (*ApplicationParams, error) {
	fileContent, err := getEnvSubstitutedFileContent(path)
	if err != nil {
		return nil, err
	}

	apiParams := &ApplicationParams{}
	err = yaml.Unmarshal([]byte(fileContent), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, err
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


// GetEnv returns the EndpointData associated for key in the ApiParams, if not found returns nil
func (config ApiParams) GetEnv(key string) *Environment {
	for index, env := range config.Environments {
		if env.Name == key {
			return &config.Environments[index]
		}
	}
	return nil
}

