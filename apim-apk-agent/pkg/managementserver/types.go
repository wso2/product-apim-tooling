/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package managementserver

// Subscription for struct subscription
type Subscription struct {
	SubStatus     string         `json:"subStatus,omitempty"`
	UUID          string         `json:"uuid,omitempty"`
	Organization  string         `json:"organization,omitempty"`
	SubscribedAPI *SubscribedAPI `json:"subscribedApi,omitempty"`
	TimeStamp     int64          `json:"timeStamp,omitempty"`
	RateLimit     string         `json:"rateLimit,omitempty"`
}

// SubscriptionList for struct list of applications
type SubscriptionList struct {
	List []Subscription `json:"list"`
}

// SubscribedAPI for struct subscribedAPI
type SubscribedAPI struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// Application for struct application
type Application struct {
	UUID         string            `json:"uuid,omitempty"`
	Name         string            `json:"name,omitempty"`
	Owner        string            `json:"owner,omitempty"`
	Organization string            `json:"organization,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
	TimeStamp    int64             `json:"timeStamp,omitempty"`
}

// ApplicationList for struct list of application
type ApplicationList struct {
	List []Application `json:"list"`
}

// ResolvedApplicationList for struct list of resolved application
type ResolvedApplicationList struct {
	List []ResolvedApplication `json:"list"`
}

// ResolvedApplication for struct resolvedApplication
type ResolvedApplication struct {
	UUID            string            `json:"uuid,omitempty"`
	Name            string            `json:"name,omitempty"`
	Owner           string            `json:"owner,omitempty"`
	Organization    string            `json:"organization,omitempty"`
	Attributes      map[string]string `json:"attributes,omitempty"`
	TimeStamp       int64             `json:"timeStamp,omitempty"`
	SecuritySchemes []SecurityScheme  `json:"securitySchemes,omitempty"`
}

// SecurityScheme for struct securityScheme
type SecurityScheme struct {
	SecurityScheme        string `json:"securityScheme,omitempty"`
	ApplicationIdentifier string `json:"applicationIdentifier,omitempty"`
	KeyType               string `json:"keyType,omitempty"`
	EnvID                 string `json:"envID,omitempty"`
}

// ApplicationKeyMapping for struct applicationKeyMapping
type ApplicationKeyMapping struct {
	ApplicationUUID       string `json:"applicationUUID,omitempty"`
	SecurityScheme        string `json:"securityScheme,omitempty"`
	ApplicationIdentifier string `json:"applicationIdentifier,omitempty"`
	KeyType               string `json:"keyType,omitempty"`
	EnvID                 string `json:"envID,omitempty"`
	Timestamp             int64  `json:"timestamp,omitempty"`
	Organization          string `json:"organization,omitempty"`
}

// ApplicationKeyMappingList for struct list of applicationKeyMapping
type ApplicationKeyMappingList struct {
	List []ApplicationKeyMapping `json:"list"`
}

// ApplicationMapping for struct applicationMapping
type ApplicationMapping struct {
	UUID            string `json:"uuid,omitempty"`
	ApplicationRef  string `json:"applicationRef,omitempty"`
	SubscriptionRef string `json:"subscriptionRef,omitempty"`
	Organization    string `json:"organization,omitempty"`
}

// ApplicationMappingList for struct list of applicationMapping
type ApplicationMappingList struct {
	List []ApplicationMapping `json:"list"`
}

// APICPEvent holds data of a specific API event from adapter
type APICPEvent struct {
	Event EventType `json:"event"`
	API   API       `json:"payload"`
}

// EventType is the type of api event. One of (CREATE, UPDATE, DELETE)
type EventType string

const (
	// CreateEvent is create api event
	CreateEvent EventType = "CREATE"
	// DeleteEvent is delete api event
	DeleteEvent EventType = "DELETE"
)

// API holds the api data from adapter api event
type API struct {
	APIUUID          string            `json:"apiUUID"`
	APIName          string            `json:"apiName"`
	APIVersion       string            `json:"apiVersion"`
	IsDefaultVersion bool              `json:"isDefaultVersion"`
	Definition       string            `json:"definition"`
	APIType          string            `json:"apiType"`
	APISubType       string            `json:"apiSubType"`
	BasePath         string            `json:"basePath"`
	Organization     string            `json:"organization"`
	SystemAPI        bool              `json:"systemAPI"`
	APIProperties    map[string]string `json:"apiProperties,omitempty"`
	Environment      string            `json:"environment,omitempty"`
	RevisionID       string            `json:"revisionID"`
	SandEndpoint     string            `json:"sandEndpoint"`
	ProdEndpoint     string            `json:"prodEndpoint"`
	EndpointProtocol string            `json:"endpointProtocol"`
	CORSPolicy       *CORSPolicy       `json:"cORSPolicy"`
	Vhost            string            `json:"vhost"`
	SandVhost        string            `json:"sandVhost"`
	SecurityScheme   []string          `json:"securityScheme"`
	AuthHeader       string            `json:"authHeader"`
	APIKeyHeader     string            `json:"apiKeyHeader"`
	Operations       []OperationFromDP `json:"operations"`
	SandAIRL         *AIRL             `json:"sandAIRL"`
	ProdAIRL         *AIRL             `json:"prodAIRL"`
	AIConfiguration  AIConfiguration   `json:"aiConfiguration"`
}

// AIRL holds AI ratelimit related data
type AIRL struct {
	PromptTokenCount     *uint32 `json:"promptTokenCount"`
	CompletionTokenCount *uint32 `json:"CompletionTokenCount"`
	TotalTokenCount      *uint32 `json:"totalTokenCount"`
	TimeUnit             string  `json:"timeUnit"`
	RequestCount         *uint32 `json:"requestCount"`
}

// AIConfiguration holds the AI configuration
type AIConfiguration struct {
	LLMProviderID         string `json:"llmProviderID"`
	LLMProviderName       string `json:"llmProviderName"`
	LLMProviderAPIVersion string `json:"llmProviderApiVersion"`
}

// APKHeaders contains the request and response header modifier information
type APKHeaders struct {
	Policy
	RequestHeaders  APKHeaderModifier `json:"requestHeaders"`
	ResponseHeaders APKHeaderModifier `json:"responseHeaders"`
}

// APKHeaderModifier contains header modifier values
type APKHeaderModifier struct {
	AddHeaders    []APKHeader `json:"addHeaders"`
	RemoveHeaders []string    `json:"removeHeaders"`
}

// APKHeader contains the header information
type APKHeader struct {
	Name  string `json:"headerName" yaml:"headerName"`
	Value string `json:"headerValue,omitempty" yaml:"headerValue,omitempty"`
}

// OperationFromDP holds the path, verb, throttling and interceptor policy
type OperationFromDP struct {
	Path    string   `json:"path"`
	Verb    string   `json:"verb"`
	Scopes  []string `json:"scopes"`
	Filters []Filter `json:"filters"`
}

// Policy holds the policy name and version
type Policy struct {
	PolicyName    string `json:"policyName"`
	PolicyVersion string `json:"policyVersion"`
}

// Filter interface is used to define the type of parameters that can be used in an operation policy
type Filter interface {
	GetPolicyName() string
	GetPolicyVersion() string
	isFilter()
}

// GetPolicyName returns the name of the policy sent to the APIM
func (p *Policy) GetPolicyName() string {
	return p.PolicyName
}

// GetPolicyVersion returns the version of the policy sent to the APIM
func (p *Policy) GetPolicyVersion() string {
	return p.PolicyVersion
}

func (h APKHeaders) isFilter() {}

// APKRedirectRequest defines the parameters of a redirect request policy sent from the APK
type APKRedirectRequest struct {
	Policy
	URL string `json:"url"`
}

func (r APKRedirectRequest) isFilter() {}

// APKMirrorRequest defines the parameters of a mirror request policy sent from the APK
type APKMirrorRequest struct {
	Policy
	URLs []string `json:"urls"`
}

func (m APKMirrorRequest) isFilter() {}

// CORSPolicy hold cors configs
type CORSPolicy struct {
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty"`
	AccessControlExposeHeaders    []string `json:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAge           *int     `json:"accessControlMaxAge,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty"`
}
