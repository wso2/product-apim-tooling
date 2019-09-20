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
	HttpRequestTimeout int    `yaml:"http_request_timeout"`
	ExportDirectory    string `yaml:"export_directory"`
	KubernetesMode     bool   `yaml:"kubernetes_mode"`
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

//Key generation response
type KeygenResponse struct {
	CallbackURL         interface{} `json:"callbackUrl"`
	ConsumerKey         string      `json:"consumerKey"`
	ConsumerSecret      string      `json:"consumerSecret"`
	GroupID             interface{} `json:"groupId"`
	KeyState            string      `json:"keyState"`
	KeyType             string      `json:"keyType"`
	SupportedGrantTypes []string    `json:"supportedGrantTypes"`
	Token               struct {
		AccessToken  string   `json:"accessToken"`
		TokenScopes  []string `json:"tokenScopes"`
		ValidityTime int      `json:"expires_in"`
	} `json:"token"`
}

//Applications get response structure
type AppData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	List     []struct {
		ApplicationID  string `json:"applicationId"`
		Name           string `json:"name"`
		Subscriber     string `json:"subscriber"`
		ThrottlingTier string `json:"throttlingTier"`
		Description    string `json:"description"`
		Status         string `json:"status"`
		GroupID        string `json:"groupId"`
		Attributes     struct {
		} `json:"attributes"`
	} `json:"list"`
}

//Specific application details structure
type AppDetails struct {
	GroupID        string      `json:"groupId"`
	CallbackURL    interface{} `json:"callbackUrl"`
	Subscriber     string      `json:"subscriber"`
	ThrottlingTier string      `json:"throttlingTier"`
	ApplicationID  string      `json:"applicationId"`
	Description    interface{} `json:"description"`
	Status         string      `json:"status"`
	Name           string      `json:"name"`
	Keys           []struct {
		ConsumerKey         string      `json:"consumerKey"`
		ConsumerSecret      string      `json:"consumerSecret"`
		KeyState            string      `json:"keyState"`
		KeyType             string      `json:"keyType"`
		SupportedGrantTypes interface{} `json:"supportedGrantTypes"`
		Token               struct {
			ValidityTime int      `json:"validityTime"`
			AccessToken  string   `json:"accessToken"`
			TokenScopes  []string `json:"tokenScopes"`
		} `json:"token"`
	} `json:"keys"`
}

type App struct {
	ApplicationID     string        `json:"applicationId"`
	Name              string        `json:"name"`
	ThrottlingPolicy  string        `json:"throttlingPolicy"`
	Description       string        `json:"description"`
	TokenType         string        `json:"tokenType"`
	Status            string        `json:"status"`
	Groups            []interface{} `json:"groups"`
	SubscriptionCount int           `json:"subscriptionCount"`
	Keys              []interface{} `json:"keys"`
	Attributes        struct {
	} `json:"attributes"`
	SubscriptionScopes []interface{} `json:"subscriptionScopes"`
	Owner              string        `json:"owner"`
}
//Specific subscription details
type Subscription struct {
	Tier           string `json:"tier"`
	SubscriptionID string `json:"subscriptionId"`
	APIIdentifier  string `json:"apiIdentifier"`
	ApplicationID  string `json:"applicationId"`
	Status         string `json:"status"`
}


//API Search response struct
type ApiSearch struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	List     []struct {
		ID           string      `json:"id"`
		Name         string      `json:"name"`
		Description  interface{} `json:"description"`
		Context      string      `json:"context"`
		Version      string      `json:"version"`
		Provider     string      `json:"provider"`
		Status       string      `json:"status"`
		ThumbnailURI interface{} `json:"thumbnailUri"`
	} `json:"list"`
	Pagination struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

//Subscriptions details response struct
type SubscriptionDetail struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	List     []struct {
		SubscriptionID string `json:"subscriptionId"`
		ApplicationID  string `json:"applicationId"`
		APIIdentifier  string `json:"apiIdentifier"`
		Tier           string `json:"tier"`
		Status         string `json:"status"`
	} `json:"list"`
}


//Scope details response struct
type Scopes struct {
	List []struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Roles       string `json:"roles"`
		Description string `json:"description"`
	} `json:"list"`
}
