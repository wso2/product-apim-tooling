/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package utils

// ------------------- Structs for YAML Config Files ----------------------------------

// For env_keys_all.yaml
// Not to be manually edited
type EnvKeysAll struct {
	Environments map[string]EnvKeys `yaml:"environments"`
}

// For main_config.yaml
// To be manually edited by the user
type MainConfig struct {
	Config       Config                  `yaml:"config"`
	Environments map[string]EnvEndpoints `yaml:"environments"`
}

type Config struct {
	HttpRequestTimeout  int    `yaml:"http_request_timeout"`
	SkipTLSVerification bool   `yaml:"skip_tls_verification"`
	ExportDirectory     string `yaml:"export_directory"`
	ConfigDirectory string `yaml:"config_directory"`
}

type EnvKeys struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"` // to be encrypted (with the user's password) and stored
	Username     string `yaml:"username"`
}

type EnvEndpoints struct {
	APIManagerEndpoint   string `yaml:"api_manager_endpoint"`
	RegistrationEndpoint string `yaml:"registration_endpoint"`
	TokenEndpoint        string `yaml:"token_endpoint"`
}

// ---------------- End of Structs for YAML Config Files ---------------------------------

type API struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Context         string `json:"context"`
	Version         string `json:"version"`
	Provider        string `json:"provider"`
	LifeCycleStatus string `json:"lifeCycleStatus"`
	WorkflowStatus  string `json:"workflowStatus"`
}

type RegistrationResponse struct {
	ClientSecretExpiresAt string `json:"client_secret_expires_at"`
	ClientID              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	ClientName            string `json:"client_name"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
}

type APIListResponse struct {
	Count int32 `json:"count"`
	List  []API `json:"list"`
}
