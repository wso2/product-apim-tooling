package params

import (
	"encoding/json"
	"fmt"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Configuration represents endpoint config
type Configuration struct {
	// RetryTimeOut for endpoint
	RetryTimeOut *int `yaml:"retryTimeOut,omitempty" json:"retryTimeOut,omitempty"`
	// RetryDelay for endpoint
	RetryDelay *int `yaml:"retryDelay,omitempty" json:"retryDelay,omitempty"`
	// Factor used for config
	Factor *int `yaml:"factor,omitempty" json:"factor,omitempty"`
}

// Endpoint details
type Endpoint struct {
	// Type of the endpoints
	EndpointType string `json:"endpoint_type,omitempty"`
	// Url of the endpoint
	Url *string `yaml:"url" json:"url"`
	// Config of endpoint
	Config *Configuration `yaml:"config,omitempty" json:"config,omitempty"`
}

// EndpointData contains details about endpoints
type EndpointData struct {
	// Type of the endpoints
	EndpointType string `json:"endpoint_type"`
	// Production endpoint
	Production *Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Sandbox endpoint
	Sandbox *Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
}

// LoadBalanceEndpointsData contains details about endpoints mainly to be used in load balancing
type LoadBalanceEndpointsData struct {
	// Type of the endpoints
	EndpointType string `json:"endpoint_type"`
	// Production endpoints list for load balancing
	Production []Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Sandbox endpoints list for load balancing
	Sandbox []Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
	// Session management method from the load balancing group. Values can be "none", "transport" (by default), "soap", "simpleClientSession" (Client ID)
	SessionManagement string `yaml:"sessionManagement" json:"sessionManagement,omitempty"`
	// Session timeout means the number of milliseconds after which the session would time out
	SessionTimeout int `yaml:"sessionTimeOut" json:"sessionTimeOut,omitempty"`
	// Class name for algorithm to be used if load balancing should be done
	AlgorithmClassName string `yaml:"algoClassName" json:"algoClassName,omitempty"`
}

// FailoverEndpointsData contains details about endpoints mainly to be used in failover scenario
type FailoverEndpointsData struct {
	// Type of the endpoints
	EndpointType string `json:"endpoint_type"`
	// Primary production endpoint for failover
	Production *Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Production failover endpoints list for failover
	ProductionFailovers []Endpoint `yaml:"productionFailovers" json:"production_failovers,omitempty"`
	// Primary sandbox endpoint for failover
	Sandbox *Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
	// Production failover endpoints list for failover endpoint types
	SandboxFailovers []Endpoint `yaml:"sandboxFailovers" json:"sandbox_failovers,omitempty"`
	// To enable failover endpoints
	Failover bool `json:"failOver,omitempty"`
}

// AWSLambdaEndpointsData contains details about endpoints to be used with AWS Lambda endpoints
type AWSLambdaEndpointsData struct {
	// Type of the endpoints
	EndpointType string `json:"endpoint_type"`
	// Access method for endpoint. Values can be "role-supplied" (Using IAM role-supplied temporary AWS credentials) and "stored" (Using stored AWS credentials)
	AccessMethod string `yaml:"accessMethod" json:"access_method,omitempty"`
	// Region where endpoint located (Regions list https://docs.aws.amazon.com/general/latest/gr/rande.html)
	AmznRegion string `yaml:"amznRegion" json:"amznRegion,omitempty"`
	// Access Key for endpoint
	AmznAccessKey string `yaml:"amznAccessKey" json:"amznAccessKey,omitempty"`
	// Access Secret for endpoint
	AmznSecretKey string `yaml:"amznSecretKey" json:"amznSecretKey,omitempty"`
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
	Host string `yaml:"hostName" json:"hostName"`
	// Alias for certificate
	Alias string `yaml:"alias" json:"alias"`
	// Path for certificate file
	Path string `yaml:"path" json:"-"`
	// Certificate is used for internal purposes, it contains secret in base64
	Certificate string `json:"certificate"`
}

// MutualSslCert stores mutualssl certificate details
type MutualSslCert struct {
	// TierName of the certificate (eg:- Unlimited, Gold, Silver, Bronze)
	TierName string `yaml:"tierName" json:"tierName"`
	// Alias for certificate
	Alias string `yaml:"alias" json:"alias"`
	// Path for certificate file
	Path string `yaml:"path" json:"-"`
	// Certificate is used for internal purposes, it contains secret in base64
	Certificate string `json:"certificate"`
	// ApiIdentifier is used for internal purposes, it contains details of the API to be stored in client_certificates file
	APIIdentifier APIIdentifier `json:"apiIdentifier"`
}

// ApiIdentifier stores API Identifier details
type APIIdentifier struct {
	// Name of the provider of the API
	ProviderName string `json:"providerName"`
	// Name of the API
	APIName string `json:"apiName"`
	// Version of the API
	Version string `json:"version"`
}

type Environment struct {
	Name string `yaml:"name"`
	Config map[string]interface{} `yaml:"configs"`
}

// ApiParams represents environments defined in configuration file
type ApiParams struct {
	// Environments contains all environments in a configuration
	Environments []Environment `yaml:"environments"`
	Deploy       APIVCSParams  `yaml:"deploy"`
}

type ApiProductParams struct {
	Deploy ApiProductVCSParams `yaml:"deploy"`
}

type ApplicationParams struct {
	Deploy ApplicationVCSParams `yaml:"deploy"`
}

// ------------------- Structs for VCS Import Params ----------------------------------

type ApplicationVCSParams struct {
	Import ApplicationImportParams `yaml:"import"`
}

type APIVCSParams struct {
	Import APIImportParams `yaml:"import"`
}

type ApiProductVCSParams struct {
	Import APIProductImportParams `yaml:"import"`
}

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
	Type                       string             `yaml:"type"`
	AbsolutePath               string             `yaml:"absolutePath,omitempty"`
	RelativePath               string             `yaml:"relativePath,omitempty"`
	NickName                   string             `yaml:"nickName,omitempty"`
	FailedDuringPreviousDeploy bool               `yaml:"failedDuringPreviousDeploy,omitempty"`
	Deleted                    bool               `yaml:"deleted,omitempty"`
	ProjectInfo                ProjectInfo        `yaml:"projectInfo,omitempty"`
	ApiParams                  *ApiParams         `yaml:"apiParams,omitempty"`
	ApiProductParams           *ApiProductParams  `yaml:"apiProductParams,omitempty"`
	ApplicationParams          *ApplicationParams `yaml:"applicationParams,omitempty"`
}

type ProjectInfo struct {
	Owner   string `yaml:"owner,omitempty"`
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
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

// LoadApiParamsFromDirectory loads an API Project configuration YAML file located in path when the root
// directory is provided instead of yaml file.
//	It returns an error or a valid ApiParams
func LoadApiParamsFromDirectory(path string) (*ApiParams, error) {
	paramsFilePath:= filepath.Join(path,utils.ParamFileAPI)
	fmt.Println(paramsFilePath)
	fileContent, err := getEnvSubstitutedFileContent(paramsFilePath)
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

// LoadApiProductParamsFromFile loads an API Product project configuration YAML file located in path when the root
// directory is provided instead of yaml file.
//	It returns an error or a valid ApiProductParams
func LoadApiProductParamsFromDirectory(path string) (*ApiProductParams, error) {
	paramsFilePath:= filepath.Join(path,utils.ParamFileAPIProduct)
	fileContent, err := getEnvSubstitutedFileContent(paramsFilePath)
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

// LoadApplicationParamsFromDirectory loads an Application project configuration YAML file located in path when the root
// directory is provided instead of yaml file.
//	It returns an error or a valid ApplicationParams
func LoadApplicationParamsFromDirectory(path string) (*ApplicationParams, error) {
	paramsFilePath:= filepath.Join(path,utils.ParamFileApplication)
	fileContent, err := getEnvSubstitutedFileContent(paramsFilePath)
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
