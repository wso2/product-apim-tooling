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
	Config         Config                  `yaml:"config"`
	Environments   map[string]EnvEndpoints `yaml:"environments"`
	MgwAdapterEnvs map[string]MgwEndpoints `yaml:"mgw-clusters"`
}

type Config struct {
	HttpRequestTimeout    int    `yaml:"http_request_timeout"`
	ExportDirectory       string `yaml:"export_directory"`
	KubernetesMode        bool   `yaml:"kubernetes_mode"`
	TokenType             string `yaml:"token_type"`
	VCSDeletionEnabled    bool   `yaml:"vcs_deletion_enabled"`
	VCSConfigFilePath     string `yaml:"vcs_config_file_path"`
	VCSSourceRepoPath     string `yaml:"vcs_source_repo_path"`
	VCSDeploymentRepoPath string `yaml:"vcs_deployment_repo_path"`
	TLSRenegotiationMode  string `yaml:"tls-renegotiation-mode"`
}

type EnvKeys struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"` // encrypted (with the user's password) and stored
	Username     string `yaml:"username"`
}

type EnvEndpoints struct {
	ApiManagerEndpoint   string `yaml:"apim"`
	PublisherEndpoint    string `yaml:"publisher"`
	DevPortalEndpoint    string `yaml:"devportal"`
	RegistrationEndpoint string `yaml:"registration"`
	AdminEndpoint        string `yaml:"admin"`
	TokenEndpoint        string `yaml:"token"`
	MiManagementEndpoint string `yaml:"mi"`
}

type MgwEndpoints struct {
	AdapterEndpoint string `yaml:"adapter"`
}

// ---------------- End of Structs for YAML Config Files ---------------------------------

type API struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Context         string `json:"context"`
	Version         string `json:"version"`
	Provider        string `json:"provider"`
	LifeCycleStatus string `json:"lifeCycleStatus"`
}

type APIProduct struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Context         string `json:"context"`
	Provider        string `json:"provider"`
	LifeCycleStatus string `json:"status"`
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

type APIProductListResponse struct {
	Count int32        `json:"count"`
	List  []APIProduct `json:"list"`
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

//Key generation request
type KeygenRequest struct {
	KeyType                 string   `json:"keyType"`
	GrantTypesToBeSupported []string `json:"grantTypesToBeSupported"`
	ValidityTime            int      `json:"validityTime"`
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

//Application Keys
type AppKeyList struct {
	Count int              `json:"count"`
	List  []ApplicationKey `json:"list"`
}

// Consumer Secret regeneration response
type ConsumerSecretRegenResponse struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
}

//Applications get response structure
type AppList struct {
	Count int `json:"count"`
	List  []struct {
		ApplicationID     string        `json:"applicationId"`
		Name              string        `json:"name"`
		Owner             string        `json:"owner"`
		ThrottlingPolicy  string        `json:"throttlingPolicy"`
		Description       interface{}   `json:"description"`
		Status            string        `json:"status"`
		Groups            []interface{} `json:"groups"`
		SubscriptionCount int           `json:"subscriptionCount"`
		Attributes        struct {
		} `json:"attributes"`
	} `json:"list"`
	Pagination struct {
		Offset   int    `json:"offset"`
		Limit    int    `json:"limit"`
		Total    int    `json:"total"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"pagination"`
}

//Specific application details structure
type AppDetails struct {
	ApplicationID     string        `json:"applicationId"`
	Name              string        `json:"name"`
	ThrottlingPolicy  string        `json:"throttlingPolicy"`
	Description       interface{}   `json:"description"`
	TokenType         string        `json:"tokenType"`
	Status            string        `json:"status"`
	Groups            []interface{} `json:"groups"`
	SubscriptionCount int           `json:"subscriptionCount"`
	Keys              []ApplicationKey
	Attributes        struct {
	} `json:"attributes"`
	SubscriptionScopes []struct {
		Key         string   `json:"key"`
		Name        string   `json:"name"`
		Roles       []string `json:"roles"`
		Description string   `json:"description"`
	} `json:"subscriptionScopes"`
	Owner       string `json:"owner"`
	HashEnabled bool   `json:"hashEnabled"`
}

// Application key details
type ApplicationKey struct {
	ConsumerKey         string      `json:"consumerKey"`
	ConsumerSecret      string      `json:"consumerSecret"`
	SupportedGrantTypes []string    `json:"supportedGrantTypes"`
	CallbackURL         interface{} `json:"callbackUrl"`
	KeyState            string      `json:"keyState"`
	KeyType             string      `json:"keyType"`
}

// Application creation request
type AppCreateRequest struct {
	Name             string `json:"name"`
	ThrottlingPolicy string `json:"throttlingPolicy"`
	Description      string `json:"description"`
	TokenType        string `json:"tokenType"`
}

//Subscriptions List response struct
type SubscriptionList struct {
	Count      int            `json:"count"`
	List       []Subscription `json:"list"`
	Pagination interface{}    `json:"pagination"`
}

//Subscription
type Subscription struct {
	SubscriptionID string `json:"subscriptionId"`
	ApplicationID  string `json:"applicationId"`
	APIID          string `json:"apiId"`
	APIInfo        struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Context         string `json:"context"`
		Version         string `json:"version"`
		Provider        string `json:"provider"`
		LifeCycleStatus string `json:"lifeCycleStatus"`
	} `json:"apiInfo"`
	ApplicationInfo struct {
		ApplicationID string        `json:"applicationId"`
		Name          string        `json:"name"`
		Status        string        `json:"status"`
		Groups        []interface{} `json:"groups"`
		Owner         string        `json:"owner"`
	} `json:"applicationInfo"`
	ThrottlingPolicy  string      `json:"throttlingPolicy"`
	Status            string      `json:"status"`
	RedirectionParams interface{} `json:"redirectionParams"`
}

//Throttling Policies List response struct
type ThrottlingPoliciesList struct {
	Count      int                `json:"count"`
	List       []ThrottlingPolicy `json:"list"`
	Pagination interface{}        `json:"pagination"`
}

//ThrottlingPolicy
type ThrottlingPolicy struct {
	Name                        string      `json:"name"`
	Description                 string      `json:"description"`
	PolicyLevel                 string      `json:"policyLevel"`
	Attributes                  interface{} `json:"attributes"`
	RequestCount                int         `json:"requestCount"`
	UnitTime                    int         `json:"unitTime"`
	TierPlan                    string      `json:"tierPlan"`
	StopOnQuotaReach            bool        `json:"stopOnQuotaReach"`
	MonetizationAttributes      interface{}
	throttlingPolicyPermissions interface{}
}

//Subscription creation request
type SubscriptionCreateRequest struct {
	ApplicationID    string `json:"applicationId"`
	APIID            string `json:"apiId"`
	ThrottlingPolicy string `json:"throttlingPolicy"`
}

//API Search response struct. This includes common attributes for both store and publisher REST API search
type ApiSearch struct {
	Count int `json:"count"`
	List  []struct {
		ID              string      `json:"id"`
		Name            string      `json:"name"`
		Description     interface{} `json:"description"`
		Context         string      `json:"context"`
		Version         string      `json:"version"`
		Provider        string      `json:"provider"`
		Type            string      `json:"type"`
		LifeCycleStatus string      `json:"lifeCycleStatus"`
	} `json:"list"`
	Pagination struct {
		Offset   int    `json:"offset"`
		Limit    int    `json:"limit"`
		Total    int    `json:"total"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"pagination"`
}

//get detailed API response
type APIData struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Context             string      `json:"context"`
	Version             string      `json:"version"`
	Provider            string      `json:"provider"`
	LifeCycleStatus     string      `json:"lifeCycleStatus"`
	HasThumbnail        interface{} `json:"hasThumbnail"`
	Policies            []string    `json:"policies"`
	BusinessInformation struct {
		BusinessOwner       string `json:"businessOwner"`
		BusinessOwnerEmail  string `json:"businessOwnerEmail"`
		TechnicalOwner      string `json:"technicalOwner"`
		TechnicalOwnerEmail string `json:"technicalOwnerEmail"`
	} `json:"businessInformation"`
	WsdlInfo         interface{} `json:"wsdlInfo"`
	WsdlURL          interface{} `json:"wsdlUrl"`
	IsDefaultVersion bool        `json:"isDefaultVersion"`
	EndpointConfig   struct {
		EndpointType     string `json:"endpoint_type"`
		SandboxEndpoints struct {
			URL string `json:"url"`
		} `json:"sandbox_endpoints"`
		ProductionEndpoints struct {
			URL string `json:"url"`
		} `json:"production_endpoints"`
	} `json:"endpointConfig"`
	Transport []string `json:"transport"`
	Tags      []string `json:"tags"`
}

// Project MetaData struct
type MetaData struct {
	Name         string       `json:"name,omitempty" yaml:"name,omitempty"`
	Version      string       `json:"version,omitempty" yaml:"version,omitempty"`
	Owner        string       `json:"owner,omitempty" yaml:"owner,omitempty"`
	DeployConfig DeployConfig `json:"deploy,omitempty" yaml:"deploy,omitempty"`
}

type DeployConfig struct {
	Import ImportConfig `json:"import,omitempty" yaml:"import,omitempty"`
}

type ImportConfig struct {
	Update            bool `json:"update,omitempty" yaml:"update,omitempty"`
	PreserveProvider  bool `json:"preserveProvider,omitempty" yaml:"preserveProvider,omitempty"`
	RotateRevision    bool `json:"rotateRevision" yaml:"rotateRevision"`
	ImportAPIs        bool `json:"importApis,omitempty" yaml:"importApis,omitempty"`
	UpdateAPIProduct  bool `json:"updateApiProduct,omitempty" yaml:"updateApiProduct,omitempty"`
	UpdateAPIs        bool `json:"updateApis,omitempty" yaml:"updateApis,omitempty"`
	PreserveOwner     bool `json:"preserveOwner,omitempty" yaml:"preserveOwner,omitempty"`
	SkipSubscriptions bool `json:"skipSubscriptions,omitempty" yaml:"skipSubscriptions,omitempty"`
	SkipKeys          bool `json:"skipKeys,omitempty" yaml:"skipKeys,omitempty"`
}

type RevisionListResponse struct {
	Count int32       `json:"count"`
	List  []Revisions `json:"list"`
}

type Revisions struct {
	ID             string       `json:"id"`
	RevisionNumber string       `json:"displayName"`
	Description    string       `json:"description"`
	Deployments    []Deployment `json:"deploymentInfo"`
	GatewayEnvs    []string
}

type Deployment struct {
	Name               string `json:"name"`
	DisplayOnDevportal bool   `json:"displayOnDevportal"`
}

// APIEntry Api List Entry struct to support  different formats of output in the list command
type APIEntry struct {
	Id              string
	Name            string
	Context         string
	Version         string
	LifeCycleStatus string
	Provider        string
}

// APIProductEntry Api Product List Entry struct to support  different formats of output in the list command
type APIProductEntry struct {
	Id              string
	Name            string
	Context         string
	LifeCycleStatus string
	Provider        string
}

// ApplicationEntry Application List Entry struct to support  different formats of output in the list command
type ApplicationEntry struct {
	Id      string
	Name    string
	Status  string
	Owner   string
	GroupId string
}
