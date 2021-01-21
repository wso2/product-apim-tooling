/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	HttpRequestTimeout   int    `yaml:"http_request_timeout"`
	ExportDirectory      string `yaml:"export_directory"`
	TLSRenegotiationMode string `yaml:"tls_renegotiation_mode,omitempty"`
}

type EnvKeys struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"` // encrypted (with the user's password) and stored
	Username     string `yaml:"username"`
}

type EnvEndpoints struct {
	ApiManagerEndpoint      string `yaml:"api_manager_endpoint"`
	ApiImportExportEndpoint string `yaml:"api_import_export_endpoint"`
	ApiListEndpoint         string `yaml:"api_list_endpoint"`
	AppListEndpoint         string `yaml:"application_list_endpoint"`
	RegistrationEndpoint    string `yaml:"registration_endpoint"`
	AdminEndpoint           string `yaml:"admin_endpoint"`
	TokenEndpoint           string `yaml:"token_endpoint"`
}

// ---------------- End of Structs for YAML Config Files ---------------------------------

type API struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Context  string `json:"context"`
	Version  string `json:"version"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

type Application struct {
	ID      string `json:"applicationId"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Status  string `json:"status"`
	GroupID string `json:"groupId"`
}

type RegistrationResponse struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	ClientName   string `json:"clientName"`
	CallBackURL  string `json:"callBackURL"`
	JsonString   string `json:"jsonString"`
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

type ApplicationListResponse struct {
	Count int32         `json:"count"`
	List  []Application `json:"list"`
}

type MigrationApisExportMetadata struct {
	ApiListOffset   int    `yaml:"api_list_offset"`
	User            string `yaml:"user"`
	OnTenant        string `yaml:"on_tenant"`
	ApiListToExport []API  `yaml:"apis_to_export"`
}

type HttpErrorResponse struct {
	Code        int     `json:"code"`
	Status      string  `json:"message"`
	Description string  `json:"description"`
	MoreInfo    string  `json:"moreInfo"`
	Error       []error `json:"error"`
}
