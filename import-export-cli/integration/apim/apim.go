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

package apim

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/adminservices"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	restClientPostfix = "-integ-rest-client"

	ApplicationThrottlePolicyType  = "application"
	CustomThrottlePolicyType       = "custom"
	AdvancedThrottlePolicyType     = "advanced"
	SubscriptionThrottlePolicyType = "subscription"

	applicationThrottlePolicyQuery  = "app"
	customThrottlePolicyQuery       = "global"
	advancedThrottlePolicyQuery     = "api"
	subscriptionThrottlePolicyQuery = "sub"
)

// Client : Enables interacting with an instance of APIM
type Client struct {
	portOffset       int
	host             string
	dcrURL           string
	restClientName   string
	tokenURL         string
	apimURL          string
	accessToken      string
	publisherRestURL string
	devPortalRestURL string
	adminRestURL     string
	devopsRestURL    string
	EnvName          string
}

type dcrRequest struct {
	CallbackURL         string `json:"callbackUrl"`
	ClientName          string `json:"clientName"`
	Owner               string `json:"owner"`
	SupportedGrantTypes string `json:"grantType"`
	IsSaaSApp           bool   `json:"saasApp"`
}

type jsonString struct {
	UserName     string `json:"username"`
	ClientName   string `json:"client_name"`
	RedirectURIs string `json:"redirect_uris"`
	GrantTypes   string `json:"grant_types"`
}

type dcrResponse struct {
	CallbackURL  string      `json:"callBackURL"`
	ClientName   string      `json:"clientName"`
	JSONString   *jsonString `json:"jsonString"`
	ClientID     string      `json:"clientId"`
	ClientSecret string      `json:"clientSecret"`
	IsSaaSApp    bool        `json:"isSaasApplication"`
	Owner        string      `json:"appOwner"`
}

type tokenResponse struct {
	Scope        string `json:"scope"`
	TokeType     string `json:"token_type"`
	ValidTime    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

// GetHost : Get Host of APIM
func (instance *Client) GetHost() string {
	return instance.host
}

// GetPortOffset : Get Port Offset of APIM
func (instance *Client) GetPortOffset() int {
	return instance.portOffset
}

// Login : Login to APIM Client instance
func (instance *Client) Login(username string, password string) {
	instance.restClientName = username + restClientPostfix
	instance.accessToken = instance.getToken(username, password)
	base.Log("apim.Login() - username:", username, ",password", password)
}

// Setup : Setup APIM Client config
func (instance *Client) Setup(envName string, host string, offset int, dcrVersion, adminRestAPIVersion,
	devportalRestAPIVersion, publisherRestAPIVersion, devopsRestAPIVersion string) {
	base.Log("apim.Setup() - envName:", envName, ",host:", host, ",offset:", offset, ",dcrVersion:", dcrVersion,
		",adminRestAPIVersion:", adminRestAPIVersion, ",devportalRestAPIVersion:", devportalRestAPIVersion,
		",publisherRestAPIVersion:", publisherRestAPIVersion)
	instance.apimURL = getApimURL(host, offset)
	instance.dcrURL = getDCRURL(host, dcrVersion, offset)
	instance.devPortalRestURL = getDevPortalRestURL(host, devportalRestAPIVersion, offset)
	instance.publisherRestURL = getPublisherRestURL(host, publisherRestAPIVersion, offset)
	instance.adminRestURL = getAdminRestURL(host, adminRestAPIVersion, offset)
	instance.devopsRestURL = getDevOpsRestURL(host, devopsRestAPIVersion, offset)
	instance.portOffset = offset
	instance.tokenURL = getTokenURL(host, offset)
	instance.host = host
	instance.EnvName = envName
}

// GetEnvName : Get environment name
func (instance *Client) GetEnvName() string {
	return instance.EnvName
}

// GetApimURL : Get APIM URL
func (instance *Client) GetApimURL() string {
	return instance.apimURL
}

// GetTokenURL : Get Token URL
func (instance *Client) GetTokenURL() string {
	return instance.tokenURL
}

// GenerateSampleAPIData : Generate sample Pizzashack API object
func (instance *Client) GenerateSampleAPIData(provider string) *API {
	api := API{}
	api.Name = base.GenerateRandomString() + "API"
	api.Description = "This is a simple API for Pizza Shack online pizza delivery store."
	api.Context = getContext(provider)
	api.Version = "1.0.0"
	api.Provider = provider
	api.Transport = []string{"http", "https"}
	api.Tags = []string{"pizza"}
	api.Policies = []string{"Unlimited"}
	api.APIThrottlingPolicy = "Unlimited"
	api.SecurityScheme = []string{"oauth2"}
	api.Visibility = "PUBLIC"
	api.Type = "HTTP"
	api.SubscriptionAvailability = "CURRENT_TENANT"
	api.AccessControl = "NONE"
	api.AuthorizationHeader = "Authorization"
	api.EndpointImplementationType = "ENDPOINT"
	api.BusinessInformation = BusinessInfo{"Jane Roe", "marketing@pizzashack.com", "John Doe", "architecture@pizzashack.com"}
	api.EndpointConfig = HTTPEndpoint{"http", &URLConfig{"https://localhost:" + strconv.Itoa(9443+instance.portOffset) + "/am/sample/pizzashack/v1/api/"},
		&URLConfig{"https://localhost:" + strconv.Itoa(9443+instance.portOffset) + "/am/sample/pizzashack/v1/api/"}}
	api.Operations = generateSampleAPIOperations()
	api.AdvertiseInformation = AdvertiseInfo{Advertised: false}
	api.GatewayType = "wso2/synapse"

	return &api
}

// GenerateSampleStreamingAPIData : Generate sample Streaming API object
func (instance *Client) GenerateSampleStreamingAPIData(provider, apiType string) *API {
	api := API{}
	api.Name = base.GenerateRandomString() + "API"
	api.Description = "This is a simple Streaming API."
	api.Context = getContext(provider)
	api.Version = "1.0.0"
	api.Provider = provider
	api.Type = apiType
	if strings.EqualFold(apiType, "WS") {
		api.Policies = []string{"AsyncUnlimited"}
		api.EndpointConfig = HTTPEndpoint{"ws", &URLConfig{"ws://echo.websocket.org:" + strconv.Itoa(80+instance.portOffset)},
			&URLConfig{"ws://echo.websocket.org:" + strconv.Itoa(80+instance.portOffset)}}

	}
	if strings.EqualFold(apiType, "WEBSUB") {
		api.Policies = []string{"AsyncWHUnlimited"}
	}
	if strings.EqualFold(apiType, "SSE") {
		api.Policies = []string{"AsyncUnlimited"}
		api.EndpointConfig = HTTPEndpoint{"http", &URLConfig{"http://localhost:8080"}, &URLConfig{"http://localhost:8080"}}
	}
	return &api
}

func GenerateAdvertiseOnlyProperties(api *API, originalDevportalUrl, productionEp, sandboxEp string) {
	api.AdvertiseInformation.Advertised = true
	api.AdvertiseInformation.ApiOwner = api.Provider
	api.AdvertiseInformation.OriginalDevPortalUrl = originalDevportalUrl
	api.AdvertiseInformation.ApiExternalProductionEndpoint = productionEp
	api.AdvertiseInformation.ApiExternalSandboxEndpoint = sandboxEp
	api.AdvertiseInformation.Vendor = "WSO2"
}

func getContext(provider string) string {
	context := base.GenerateRandomString()
	if strings.Contains(provider, "@") {
		splits := strings.Split(provider, "@")
		domain := splits[len(splits)-1]
		return "/t/" + domain + "/" + context
	}

	return "/" + context
}

// GenerateAdditionalProperties : Generate additional properties to create an API from swagger
func (instance *Client) GenerateAdditionalProperties(provider, endpointUrl, apiType string, operations []interface{}) string {
	additionalProperties := `{"name":"` + base.GenerateRandomString() + `",
	"version":"1.0.5",
	"context":"` + getContext(provider) + `",
	"policies":[
	   "Unlimited"
	],
	`
	if len(operations) > 0 {
		operationsData, _ := json.Marshal(operations)
		additionalProperties += ` "operations": ` + string(operationsData) + `, `
	}

	if strings.EqualFold(apiType, "SOAPTOREST") {
		additionalProperties = additionalProperties +
			`"endpointConfig": {   
				"endpoint_type":"address",
					"sandbox_endpoints":{
						"type": "address",
						"url":"` + endpointUrl + `"
					},
					"production_endpoints":{
						"type": "address",
						"url":"` + endpointUrl + `"
					}
			}
		}`
	} else if strings.EqualFold(apiType, "WS") {
		additionalProperties = additionalProperties + `"type":"` + apiType + `",
			"endpointConfig": {   
				"endpoint_type":"ws",
					"sandbox_endpoints":{
						"url":"` + endpointUrl + `"
					},
					"production_endpoints":{
						"url":"` + endpointUrl + `"
					}
			}
		}`
	} else if strings.EqualFold(apiType, "ASYNC") {
		api := API{}
		api.Provider = provider
		GenerateAdvertiseOnlyProperties(&api, "https://localhost:9443/devportal", "amqp://production-ep:9000",
			"amqp://sandbox-ep:9000")
		advertiseInfo, _ := json.Marshal(api.AdvertiseInformation)
		additionalProperties = additionalProperties + `"type":"` + apiType + `",
		"advertiseInfo": ` + string(advertiseInfo) + `}`
	} else {
		additionalProperties = additionalProperties +
			`"endpointConfig": {   
				"endpoint_type":"http",
					"sandbox_endpoints":{
						"url":"` + endpointUrl + `"
					},
					"production_endpoints":{
						"url":"` + endpointUrl + `"
					}
			}
		}`
	}
	return additionalProperties
}

// CopyAPI : Create a deep copy of an API object
func CopyAPI(apiToCopy *API) API {
	apiCopy := *apiToCopy

	// Copy Transport slice
	apiCopy.Transport = make([]string, len(apiToCopy.Transport))
	copy(apiCopy.Transport, apiToCopy.Transport)

	// Copy Tags slice
	apiCopy.Tags = make([]string, len(apiToCopy.Tags))
	copy(apiCopy.Tags, apiToCopy.Tags)

	// Copy Policies slice
	apiCopy.Policies = make([]string, len(apiToCopy.Policies))
	copy(apiCopy.Policies, apiToCopy.Policies)

	// Copy SecurityScheme slice
	apiCopy.SecurityScheme = make([]string, len(apiToCopy.SecurityScheme))
	copy(apiCopy.SecurityScheme, apiToCopy.SecurityScheme)

	// Copy VisibleRoles slice
	apiCopy.VisibleRoles = make([]string, len(apiToCopy.VisibleRoles))
	copy(apiCopy.VisibleRoles, apiToCopy.VisibleRoles)

	// Copy VisibleTenants slice
	apiCopy.VisibleTenants = make([]string, len(apiToCopy.VisibleTenants))
	copy(apiCopy.VisibleTenants, apiToCopy.VisibleTenants)

	// Copy MediationPolicies slice
	apiCopy.MediationPolicies = make([]MediationPolicy, len(apiToCopy.MediationPolicies))
	copy(apiCopy.MediationPolicies, apiToCopy.MediationPolicies)

	// Copy SubscriptionAvailableTenants slice
	apiCopy.SubscriptionAvailableTenants = make([]string, len(apiToCopy.SubscriptionAvailableTenants))
	copy(apiCopy.SubscriptionAvailableTenants, apiToCopy.SubscriptionAvailableTenants)

	// Copy AdditionalProperties
	for key, value := range apiToCopy.AdditionalProperties {
		apiCopy.AdditionalProperties[key] = value
	}

	// Copy AccessControlRoles slice
	apiCopy.AccessControlRoles = make([]string, len(apiToCopy.AccessControlRoles))
	copy(apiCopy.AccessControlRoles, apiToCopy.AccessControlRoles)

	// Copy Operations slice
	apiCopy.Operations = make([]APIOperations, len(apiToCopy.Operations))
	copy(apiCopy.Operations, apiToCopy.Operations)

	return apiCopy
}

// CopyAPIProduct : Create a deep copy of an API Product object
func CopyAPIProduct(apiProductToCopy *APIProduct) APIProduct {
	apiProductCopy := *apiProductToCopy

	// Copy Transport slice
	apiProductCopy.Transport = make([]string, len(apiProductToCopy.Transport))
	copy(apiProductCopy.Transport, apiProductToCopy.Transport)

	// Copy Tags slice
	apiProductCopy.Tags = make([]string, len(apiProductToCopy.Tags))
	copy(apiProductCopy.Tags, apiProductToCopy.Tags)

	// Copy Policies slice
	apiProductCopy.Policies = make([]string, len(apiProductToCopy.Policies))
	copy(apiProductCopy.Policies, apiProductToCopy.Policies)

	// Copy SecurityScheme slice
	apiProductCopy.SecurityScheme = make([]string, len(apiProductToCopy.SecurityScheme))
	copy(apiProductCopy.SecurityScheme, apiProductToCopy.SecurityScheme)

	// Copy VisibleRoles slice
	apiProductCopy.VisibleRoles = make([]string, len(apiProductToCopy.VisibleRoles))
	copy(apiProductCopy.VisibleRoles, apiProductToCopy.VisibleRoles)

	// Copy VisibleTenants slice
	apiProductCopy.VisibleTenants = make([]string, len(apiProductToCopy.VisibleTenants))
	copy(apiProductCopy.VisibleTenants, apiProductToCopy.VisibleTenants)

	// Copy SubscriptionAvailableTenants slice
	apiProductCopy.SubscriptionAvailableTenants = make([]string, len(apiProductToCopy.SubscriptionAvailableTenants))
	copy(apiProductCopy.SubscriptionAvailableTenants, apiProductToCopy.SubscriptionAvailableTenants)

	// Copy AdditionalProperties
	for key, value := range apiProductToCopy.AdditionalProperties {
		apiProductCopy.AdditionalProperties[key] = value
	}

	// Copy AccessControlRoles slice
	apiProductCopy.AccessControlRoles = make([]string, len(apiProductToCopy.AccessControlRoles))
	copy(apiProductCopy.AccessControlRoles, apiProductToCopy.AccessControlRoles)

	// Copy APIs slice
	apiProductCopy.APIs = make([]interface{}, len(apiProductToCopy.APIs))
	copy(apiProductCopy.APIs, apiProductToCopy.APIs)

	return apiProductCopy
}

// SortAPIMembers : Sort API object members such as slices
func SortAPIMembers(api *API) {
	// Sort Transport slice
	sort.Strings(api.Transport)

	// Sort Tags slice
	sort.Strings(api.Tags)

	// Sort Policies slice
	sort.Strings(api.Policies)

	// Sort SecurityScheme slice
	sort.Strings(api.SecurityScheme)

	// Sort VisibleRoles slice
	sort.Strings(api.VisibleRoles)

	// Sort VisibleTenants slice
	sort.Strings(api.VisibleTenants)

	// Sort MediationPolicies slice
	sort.Sort(ByID(api.MediationPolicies))

	// Sort SubscriptionAvailableTenants slice
	sort.Strings(api.SubscriptionAvailableTenants)

	// Sort AccessControlRoles slice
	sort.Strings(api.AccessControlRoles)

	// Sort Operations slice
	sort.Sort(ByTargetVerb(api.Operations))
}

// SortAPIProductMembers : Sort API Product object members such as slices
func SortAPIProductMembers(apiProduct *APIProduct) {
	// Sort Transport slice
	sort.Strings(apiProduct.Transport)

	// Sort Tags slice
	sort.Strings(apiProduct.Tags)

	// Sort Policies slice
	sort.Strings(apiProduct.Policies)

	// Sort SecurityScheme slice
	sort.Strings(apiProduct.SecurityScheme)

	// Sort VisibleRoles slice
	sort.Strings(apiProduct.VisibleRoles)

	// Sort VisibleTenants slice
	sort.Strings(apiProduct.VisibleTenants)

	// Sort SubscriptionAvailableTenants slice
	sort.Strings(apiProduct.SubscriptionAvailableTenants)

	// Sort AccessControlRoles slice
	sort.Strings(apiProduct.AccessControlRoles)

	// Sort APIs slice
	// sort.Sort(ByTargetVerb(apiProduct.APIs))
}

// GenerateSampleAppData : Generate sample Application object
func (instance *Client) GenerateSampleAppData() *Application {
	app := Application{}
	app.Name = base.GenerateRandomString() + "Application"
	app.ThrottlingPolicy = "Unlimited"
	app.Description = "Test Application"
	app.TokenType = "JWT"
	return &app
}

// GenerateSampleAppData : Generate sample Application object with space in the application Name
func (instance *Client) GenerateSampleAppWithNameInSpaceData() *Application {
	app := Application{}
	app.Name = base.GenerateRandomString() + "Test Application"
	app.ThrottlingPolicy = "Unlimited"
	app.Description = "Test Application with space in the name"
	app.TokenType = "JWT"
	return &app
}

// CopyApp : Create a deep copy of an Application object
func CopyApp(appToCopy *Application) Application {
	appCopy := Application{}
	appCopy = *appToCopy

	// Copy Groups slice
	appCopy.Groups = make([]string, len(appToCopy.Groups))
	copy(appCopy.Groups, appToCopy.Groups)

	// Copy Keys slice
	appCopy.Keys = make([]ApplicationKey, len(appToCopy.Keys))
	copy(appCopy.Keys, appToCopy.Keys)

	// Copy SubscriptionScopes slice
	appCopy.SubscriptionScopes = make([]string, len(appToCopy.SubscriptionScopes))
	copy(appCopy.SubscriptionScopes, appToCopy.SubscriptionScopes)

	return appCopy
}

// AddAPI : Add new API to APIM
func (instance *Client) AddAPI(t *testing.T, api *API, username string, password string, doClean bool) string {
	apisURL := instance.publisherRestURL + "/apis"

	data, err := json.Marshal(api)

	if err != nil {
		t.Fatal(err)
	}

	request := base.CreatePost(apisURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.AddAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddAPI()", response, 201)

	var apiResponse API
	json.NewDecoder(response.Body).Decode(&apiResponse)

	if doClean {
		t.Cleanup(func() {
			username, password := RetrieveAdminCredentialsInsteadCreator(username, password)
			instance.Login(username, password)
			instance.DeleteAPI(apiResponse.ID)
		})
	}

	return apiResponse.ID
}

// UpdateAPI : Update API in APIM
func (instance *Client) UpdateAPI(t *testing.T, api *API, username string, password string) string {
	apisURL := instance.publisherRestURL + "/apis/" + api.ID

	data, err := json.Marshal(api)

	if err != nil {
		t.Fatal(err)
	}

	request := base.CreatePut(apisURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.UpdateAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.UpdateAPI()", response, 200)

	var apiResponse API
	json.NewDecoder(response.Body).Decode(&apiResponse)

	return apiResponse.ID
}

// AddSoapAPI : Add new SOAP API to APIM
func (instance *Client) AddSoapAPI(t *testing.T, path, additionalProperties, username, password, apiType string) string {
	apisURL := instance.publisherRestURL + "/apis/import-wsdl"

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	part, err = writer.CreateFormField("additionalProperties")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(additionalProperties))

	part, err = writer.CreateFormField("implementationType")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(apiType))

	err = writer.Close()

	request := base.CreatePost(apisURL, body)

	base.SetDefaultRestAPIHeadersToConsumeFormData(instance.accessToken, request)

	base.LogRequest("apim.AddSoapAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddSoapAPI()", response, 201)

	var apiResponse API
	json.NewDecoder(response.Body).Decode(&apiResponse)

	t.Cleanup(func() {
		username, password := RetrieveAdminCredentialsInsteadCreator(username, password)
		instance.Login(username, password)
		instance.DeleteAPI(apiResponse.ID)
	})

	return apiResponse.ID
}

// AddAPIFromOpenAPIDefinition : Add Petstore API using an OpenAPI Definition to APIM
func (instance *Client) AddAPIFromOpenAPIDefinition(t *testing.T, path string, additionalProperties string, username string, password string) string {
	apisURL := instance.publisherRestURL + "/apis/import-openapi"

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	part, err = writer.CreateFormField("additionalProperties")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(additionalProperties))

	err = writer.Close()

	request := base.CreatePost(apisURL, body)

	base.SetDefaultRestAPIHeadersToConsumeFormData(instance.accessToken, request)

	base.LogRequest("apim.AddAPIFromOpenAPIDefinition()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddAPIFromOpenAPIDefinition()", response, 201)

	var apiResponse API
	json.NewDecoder(response.Body).Decode(&apiResponse)

	t.Cleanup(func() {
		username, password := RetrieveAdminCredentialsInsteadCreator(username, password)
		instance.Login(username, password)
		instance.DeleteAPI(apiResponse.ID)
	})

	return apiResponse.ID
}

// AddGraphQLAPI : Add new GraphQL API to APIM
func (instance *Client) AddGraphQLAPI(t *testing.T, path, additionalProperties, username, password string) string {
	apisURL := instance.publisherRestURL + "/apis/import-graphql-schema"

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	part, err = writer.CreateFormField("additionalProperties")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(additionalProperties))

	part, err = writer.CreateFormField("implementationType")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("GraphQL"))

	err = writer.Close()

	request := base.CreatePost(apisURL, body)

	base.SetDefaultRestAPIHeadersToConsumeFormData(instance.accessToken, request)

	base.LogRequest("apim.AddGraphQLAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddGraphQLAPI()", response, 201)

	var apiResponse API
	json.NewDecoder(response.Body).Decode(&apiResponse)

	t.Cleanup(func() {
		username, password := RetrieveAdminCredentialsInsteadCreator(username, password)
		instance.Login(username, password)
		instance.DeleteAPI(apiResponse.ID)
	})

	return apiResponse.ID
}

// ValidateGraphQLSchema : Validate the GraphQL schema
func (instance *Client) ValidateGraphQLSchema(t *testing.T, path, username, password string) GraphQLValidationResponseDTO {
	apisURL := instance.publisherRestURL + "/apis/validate-graphql-schema"

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()

	request := base.CreatePost(apisURL, body)

	base.SetDefaultRestAPIHeadersToConsumeFormData(instance.accessToken, request)

	base.LogRequest("apim.ValidateGraphQLSchema()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.ValidateGraphQLSchema()", response, 200)

	var validationResponse GraphQLValidationResponseDTO
	json.NewDecoder(response.Body).Decode(&validationResponse)

	return validationResponse
}

// AddStreamingAPI : Add new Streaming API to APIM from definition
func (instance *Client) AddStreamingAPI(t *testing.T, path, additionalProperties, username, password string) string {
	apisURL := instance.publisherRestURL + "/apis/import-asyncapi"

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	part, err = writer.CreateFormField("additionalProperties")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(additionalProperties))

	err = writer.Close()

	request := base.CreatePost(apisURL, body)

	base.SetDefaultRestAPIHeadersToConsumeFormData(instance.accessToken, request)

	base.LogRequest("apim.AddStreamingAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddStreamingAPI()", response, 201)

	var apiResponse API
	json.NewDecoder(response.Body).Decode(&apiResponse)

	t.Cleanup(func() {
		username, password := RetrieveAdminCredentialsInsteadCreator(username, password)
		instance.Login(username, password)
		instance.DeleteAPI(apiResponse.ID)
	})

	return apiResponse.ID
}

// AddAPIProductFromJSON : Add SampleAPIProduct using using a json file
func (instance *Client) AddAPIProductFromJSON(t *testing.T, path string, username string, password string, apisList map[string]*API, doClean bool) string {
	apiProductsURL := instance.publisherRestURL + "/api-products"

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var apiProductData interface{}
	err = json.Unmarshal([]byte(data), &apiProductData)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a random string as the API Product name and context
	apiProductData.(map[string]interface{})["name"] = base.GenerateRandomString()
	apiProductData.(map[string]interface{})["context"] = getContext(username)

	// Retrieve the APIs in the json file of the API Product
	apisInAPIProduct := apiProductData.(map[string]interface{})["apis"]

	// Iterate through the APIs in the apis array
	for _, apiFromJSONFile := range apisInAPIProduct.([]interface{}) {
		// Iterate through the realAPIName:API map
		for apiNameKey, dependentAPI := range apisList {
			// If the name of the apiFromJSONFile matches with the real API name in the map
			if apiFromJSONFile.(map[string]interface{})["name"] == apiNameKey {
				// Replace the real API name with the random string name generated when importing the API
				apiFromJSONFile.(map[string]interface{})["name"] = dependentAPI.Name
				// Replace the apiId witht the ID in the APIM for that particular API
				apiFromJSONFile.(map[string]interface{})["apiId"] = dependentAPI.ID
			}
		}
	}

	data, err = json.Marshal(apiProductData)
	if err != nil {
		t.Fatal(err)
	}

	request := base.CreatePost(apiProductsURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.AddAPIProductFromJSON()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddAPIProductFromYaml()", response, 201)

	var apiProductResponse APIProduct
	json.NewDecoder(response.Body).Decode(&apiProductResponse)

	if doClean {
		t.Cleanup(func() {
			instance.Login(username, password)
			instance.DeleteAPIProduct(apiProductResponse.ID)
		})
	}

	return apiProductResponse.ID
}

// PublishAPIProduct : Publish API Product from APIM
func (instance *Client) PublishAPIProduct(apiProductID string) {
	lifeCycleURL := instance.publisherRestURL + "/api-products/change-lifecycle"

	request := base.CreatePostEmptyBody(lifeCycleURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("action", "Publish")
	values.Add("apiProductId", apiProductID)

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.PublishAPIProduct()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.PublishAPIProduct()", response, 200)
}

// GetAPIRevisions : Get API revisions
func (instance *Client) GetAPIRevisions(apiID, query string) *APIRevisionList {
	revisioningURL := instance.publisherRestURL + "/apis/" + apiID + "/revisions"

	if query != "" {
		revisioningURL += "?query=" + query
	}

	request := base.CreateGet(revisioningURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetAPIRevisions()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPIRevisions()", response, 200)

	var revisionsList APIRevisionList
	json.NewDecoder(response.Body).Decode(&revisionsList)
	return &revisionsList
}

// CreateAPIRevision : Create API revision
func (instance *Client) CreateAPIRevision(apiID string) *APIRevision {
	revisioningURL := instance.publisherRestURL + "/apis/" + apiID + "/revisions"

	request := base.CreatePost(revisioningURL, bytes.NewBuffer([]byte("{ \"description\": \"\" }")))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.CreateAPIRevision()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.CreateAPIRevision()", response, 201)

	var revision APIRevision
	json.NewDecoder(response.Body).Decode(&revision)

	return &revision
}

// DeployAPIRevision : Deploy API revision
func (instance *Client) DeployAPIRevision(t *testing.T, apiID, deploymentName, vhost, revisionID string) {
	deployURL := instance.publisherRestURL + "/apis/" + apiID + "/deploy-revision"

	deploymentInfoArray := []APIRevisionDeployment{}
	deploymentInfo := APIRevisionDeployment{}
	deploymentInfo.RevisionUUID = revisionID

	if deploymentName == "" {
		deploymentInfo.Name = "Default"
	} else {
		deploymentInfo.Name = deploymentName
	}
	if vhost == "" {
		deploymentInfo.VHost = "localhost"
	} else {
		deploymentInfo.VHost = vhost
	}

	deploymentInfo.DisplayOnDevportal = true
	deploymentInfoArray = append(deploymentInfoArray, deploymentInfo)

	data, err := json.Marshal(deploymentInfoArray)

	if err != nil {
		t.Fatal(err)
	}

	request := base.CreatePost(deployURL, bytes.NewBuffer(data))

	values := url.Values{}
	values.Add("revisionId", revisionID)

	request.URL.RawQuery = values.Encode()

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeployAPIRevision()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeployAPIRevision()", response, 201)
}

// GetAPIProductRevisions : Get API Product revisions
func (instance *Client) GetAPIProductRevisions(apiID, query string) *APIRevisionList {
	revisioningURL := instance.publisherRestURL + "/api-products/" + apiID + "/revisions"

	if query != "" {
		revisioningURL += "?query=" + query
	}

	request := base.CreateGet(revisioningURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetAPIProductRevisions()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPIProductRevisions()", response, 200)

	var revisionsList APIRevisionList
	json.NewDecoder(response.Body).Decode(&revisionsList)
	return &revisionsList
}

// CreateAPIProductRevision : Create API Product revision
func (instance *Client) CreateAPIProductRevision(apiProductID string) *APIRevision {
	revisioningURL := instance.publisherRestURL + "/api-products/" + apiProductID + "/revisions"

	request := base.CreatePost(revisioningURL, bytes.NewBuffer([]byte("{ \"description\": \"\" }")))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.CreateAPIProductRevision()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.CreateAPIProductRevision()", response, 201)

	var revision APIRevision
	json.NewDecoder(response.Body).Decode(&revision)

	return &revision
}

// DeployAPIProductRevision : Deploy API Product revision
func (instance *Client) DeployAPIProductRevision(t *testing.T, apiProductID, deploymentName, vhost, revisionID string) {
	deployURL := instance.publisherRestURL + "/api-products/" + apiProductID + "/deploy-revision"

	deploymentInfoArray := []APIRevisionDeployment{}
	deploymentInfo := APIRevisionDeployment{}
	deploymentInfo.RevisionUUID = revisionID

	if deploymentName == "" {
		deploymentInfo.Name = "Default"
	} else {
		deploymentInfo.Name = deploymentName
	}
	if vhost == "" {
		deploymentInfo.VHost = "localhost"
	} else {
		deploymentInfo.VHost = vhost
	}

	deploymentInfo.DisplayOnDevportal = true
	deploymentInfoArray = append(deploymentInfoArray, deploymentInfo)

	data, err := json.Marshal(deploymentInfoArray)

	if err != nil {
		t.Fatal(err)
	}

	request := base.CreatePost(deployURL, bytes.NewBuffer(data))

	values := url.Values{}
	values.Add("revisionId", revisionID)

	request.URL.RawQuery = values.Encode()

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeployAPIProductRevision()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeployAPIProductRevision()", response, 201)
}

// DeleteAPI : Delete API from APIM
func (instance *Client) DeleteAPI(apiID string) {
	apisURL := instance.publisherRestURL + "/apis/" + apiID

	request := base.CreateDelete(apisURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteAPI()", response, 200)
}

// DeleteAPIProduct : Delete API Product from APIM
func (instance *Client) DeleteAPIProduct(apiProductID string) {
	apiProductsURL := instance.publisherRestURL + "/api-products/" + apiProductID

	request := base.CreateDelete(apiProductsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteAPIProduct()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteAPIProduct()", response, 200)
}

// AddGatewayEnv : Add a gateway environment
func (instance *Client) AddGatewayEnv(t *testing.T, environment Environment, username, password string) *Environment {
	environmentsURL := instance.adminRestURL + "/environments"

	data, err := json.Marshal(environment)

	if err != nil {
		t.Fatal(err)
	}

	request := base.CreatePost(environmentsURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.AddGatewayEnv()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddGatewayEnv()", response, 201)

	var envResponse Environment
	json.NewDecoder(response.Body).Decode(&envResponse)

	t.Cleanup(func() {
		username, password := RetrieveAdminCredentialsInsteadCreator(username, password)
		instance.Login(username, password)
		instance.DeleteGatewayEnv(envResponse.ID)
	})

	return &envResponse
}

// DeleteGatewayEnv : Delete a gateway environment
func (instance *Client) DeleteGatewayEnv(envID string) {
	environmentsURL := instance.adminRestURL + "/environments/" + envID

	request := base.CreateDelete(environmentsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteGatewayEnv()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteGatewayEnv()", response, 200)
}

// DeleteAPIByName : Delete API from APIM by name
func (instance *Client) DeleteAPIByName(name string) error {
	apiInfo, err := instance.GetAPIByName(name)

	if err == nil {
		instance.DeleteAPI(apiInfo.ID)
	}

	return err
}

// DeleteAPIProductByName : Delete API from APIM by name
func (instance *Client) DeleteAPIProductByName(name string) {
	apiProductInfo := instance.GetAPIProductByName(name)
	instance.DeleteAPIProduct(apiProductInfo.ID)
}

// GetAPI : Get API from APIM
func (instance *Client) GetAPI(apiID string) *API {
	apisURL := instance.publisherRestURL + "/apis/" + apiID

	request := base.CreateGet(apisURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPI()", response, 200)

	var apiResponse API
	json.NewDecoder(response.Body).Decode(&apiResponse)
	return &apiResponse
}

// GetAPIs : Get APIs from APIM
func (instance *Client) GetAPIs() *APIList {
	apisURL := instance.publisherRestURL + "/apis"

	request := base.CreateGet(apisURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPIs()", response, 200)

	var apisResponse APIList
	json.NewDecoder(response.Body).Decode(&apisResponse)
	return &apisResponse
}

// GetAPIProduct : Get API Product from APIM
func (instance *Client) GetAPIProduct(apiProductID string) *APIProduct {
	apisURL := instance.publisherRestURL + "/api-products/" + apiProductID

	request := base.CreateGet(apisURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetAPIProduct()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPIProduct()", response, 200)

	var apiProductResponse APIProduct
	json.NewDecoder(response.Body).Decode(&apiProductResponse)
	return &apiProductResponse
}

// GetAPIProducts : Get API Products from APIM
func (instance *Client) GetAPIProducts() *APIProductList {
	apisURL := instance.publisherRestURL + "/api-products"

	request := base.CreateGet(apisURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetAPIProducts()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPIProducts()", response, 200)

	var apiProductsResponse APIProductList
	json.NewDecoder(response.Body).Decode(&apiProductsResponse)
	return &apiProductsResponse
}

// GetAPIByName : Get API by name from APIM
func (instance *Client) GetAPIByName(name string) (*APIInfo, error) {

	apisURL := instance.publisherRestURL + "/apis"

	request := base.CreateGet(apisURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("query", name)

	request.URL.RawQuery = values.Encode()

	attempts := 0

	base.LogRequest("apim.GetAPIByName()", request)

	for attempts != base.GetMaxInvocationAttempts() {

		base.Log("apim.GetAPIByName() attempts = ", attempts)

		response := base.SendHTTPRequest(request)

		defer response.Body.Close()

		base.ValidateAndLogResponse("apim.GetAPIByName()", response, 200)

		var apiResponse APIList
		json.NewDecoder(response.Body).Decode(&apiResponse)

		if len(apiResponse.List) > 0 {
			return &apiResponse.List[0], nil
		}

		base.WaitForIndexing()

		attempts++
	}

	return nil, errors.New("apim.GetAPIByName() did not return result for: " + name +
		", it is possible that sufficient time is not allowed for solr indexing." +
		"Consider the user of base.WaitForIndexing() in the execution flow where appropriate or " +
		"increasing the `indexing-delay` value in the integration test config.yaml")
}

// GetAPIProductByName : Get API Product by name from APIM
func (instance *Client) GetAPIProductByName(name string) *APIProductInfo {
	apiProductsURL := instance.publisherRestURL + "/api-products"

	request := base.CreateGet(apiProductsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("query", name)

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.GetAPIProductByName()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPIProductByName()", response, 200)

	var apiProductResponse APIProductList
	json.NewDecoder(response.Body).Decode(&apiProductResponse)
	return &apiProductResponse.List[0]
}

// PublishAPI : Publish API from APIM
func (instance *Client) PublishAPI(apiID string) {
	lifeCycleURL := instance.publisherRestURL + "/apis/change-lifecycle"

	request := base.CreatePostEmptyBody(lifeCycleURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("action", "Publish")
	values.Add("apiId", apiID)

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.PublishAPI()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.PublishAPI()", response, 200)
}

// GetApplicationSubscriptions : Get subscriptions of an application
func (instance *Client) GetApplicationSubscriptions(appID string) *SubscriptionList {
	appsURL := instance.devPortalRestURL + "/subscriptions"

	request := base.CreateGet(appsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("applicationId", appID)

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.GetApplicationSubscriptions()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetApplicationSubscriptions()", response, 200)

	var subscriptionsList SubscriptionList
	json.NewDecoder(response.Body).Decode(&subscriptionsList)
	return &subscriptionsList
}

// AddSubscription : Subscribe an App to a given API in APIM
func (instance *Client) AddSubscription(t *testing.T, apiID string, appID string, throttlePolicy string, username string, password string) {
	subscriptionURL := instance.devPortalRestURL + "/subscriptions"

	subscription := Subscription{}
	subscription.APIID = apiID
	subscription.ApplicationID = appID
	subscription.ThrottlingPolicy = throttlePolicy

	data, err := json.Marshal(subscription)

	if err != nil {
		base.Fatal(err)
	}

	request := base.CreatePost(subscriptionURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.AddSubscription()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddSubscription()", response, 201)

	var subsResponse Subscription
	json.NewDecoder(response.Body).Decode(&subsResponse)

	t.Cleanup(func() {
		instance.Login(username, password)
		instance.deleteSubscription(subsResponse.SubscriptionID)
	})

}

// DeleteSubscriptions : Delete Subscriptions for an API from APIM
func (instance *Client) DeleteSubscriptions(apiID string) {
	subsGetURL := instance.devPortalRestURL + "/subscriptions"

	request := base.CreateGet(subsGetURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("apiId", apiID)

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.DeleteSubscriptions() getting Subscriptions", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteSubscriptions() getting Subscriptions", response, 200)

	var subsResponse SubscriptionList
	json.NewDecoder(response.Body).Decode(&subsResponse)

	for _, sub := range subsResponse.List {
		instance.deleteSubscription(sub.SubscriptionID)
	}
}

func (instance *Client) deleteSubscription(subsID string) {
	subsDeleteURL := instance.devPortalRestURL + "/subscriptions/" + subsID

	request := base.CreateDelete(subsDeleteURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.deleteSubscription() deleting Subscription", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.deleteSubscription() deleting Subscription", response, 200)
}

// AddApplication : Add new Application to APIM
func (instance *Client) AddApplication(t *testing.T, application *Application, username string, password string, doClean bool) *Application {
	appsURL := instance.devPortalRestURL + "/applications"

	data, err := json.Marshal(application)

	if err != nil {
		base.Fatal(err)
	}

	request := base.CreatePost(appsURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.AddApplication()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddApplication()", response, 201)

	var appResponse Application
	json.NewDecoder(response.Body).Decode(&appResponse)

	if doClean {
		t.Cleanup(func() {
			instance.Login(username, password)
			instance.DeleteApplication(appResponse.ApplicationID)
		})
	}

	return &appResponse
}

// GenerateKeys : Generate keys for an application
func (instance *Client) GenerateKeys(t *testing.T, keyGenRequest utils.KeygenRequest, appId string) ApplicationKey {
	appsURL := instance.devPortalRestURL + "/applications/" + appId + "/generate-keys"

	data, err := json.Marshal(keyGenRequest)

	if err != nil {
		base.Fatal(err)
	}

	request := base.CreatePost(appsURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GenerateKeys()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GenerateKeys()", response, 200)

	var keyGenResponse ApplicationKey
	json.NewDecoder(response.Body).Decode(&keyGenResponse)

	return keyGenResponse
}

// GetOauthKeys : Get Oauth keys of an application
func (instance *Client) GetOauthKeys(t *testing.T, application *Application) *ApplicationKeysList {
	appsURL := instance.devPortalRestURL + "/applications/" + application.ApplicationID + "/oauth-keys"

	request := base.CreateGet(appsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetOauthKeys()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetOauthKeys()", response, 200)

	var applicationKeysList ApplicationKeysList
	json.NewDecoder(response.Body).Decode(&applicationKeysList)

	if len(applicationKeysList.List) > 0 {
		return &applicationKeysList
	} else {
		return &ApplicationKeysList{}
	}
}

// DeleteApplication : Delete Application from APIM
func (instance *Client) DeleteApplication(appID string) {
	appsURL := instance.devPortalRestURL + "/applications/" + appID

	request := base.CreateDelete(appsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteApplication()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteApplication()", response, 200)
}

// DeleteApplicationByName : Delete Application from App by name
func (instance *Client) DeleteApplicationByName(name string) {
	appInfo := instance.GetApplicationByName(name)

	instance.DeleteApplication(appInfo.ApplicationID)
}

// GetApplication : Get Application from APIM
func (instance *Client) GetApplication(appID string) *Application {
	appsURL := instance.devPortalRestURL + "/applications/" + appID

	request := base.CreateGet(appsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetApplication()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetApplication()", response, 200)

	var appResponse Application
	json.NewDecoder(response.Body).Decode(&appResponse)
	return &appResponse
}

// GetApplications : Get Applications list from APIM
func (instance *Client) GetApplications() *ApplicationList {
	appsURL := instance.devPortalRestURL + "/applications"
	request := base.CreateGet(appsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetApplications()", request)

	response := base.SendHTTPRequest(request)
	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetApplications()", response, 200)

	var appResponse ApplicationList
	json.NewDecoder(response.Body).Decode(&appResponse)
	return &appResponse
}

// GetApplicationByName : Get Application from APIM by name
func (instance *Client) GetApplicationByName(name string) *ApplicationInfo {
	appsURL := instance.devPortalRestURL + "/applications"

	request := base.CreateGet(appsURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("query", name)

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.GetApplicationByName()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetApplicationByName()", response, 200)

	var appsResponse ApplicationList
	json.NewDecoder(response.Body).Decode(&appsResponse)

	return &appsResponse.List[0]
}

// DeleteAllSubscriptions : Delete All Subscriptions from APIM
func (instance *Client) DeleteAllSubscriptions() {
	apisGetURL := instance.devPortalRestURL + "/apis"

	request := base.CreateGet(apisGetURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteAllSubscriptions() getting APIs", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteAllSubscriptions() getting APIs", response, 200)

	var apisResponse APIList
	json.NewDecoder(response.Body).Decode(&apisResponse)

	for _, api := range apisResponse.List {

		subsGetURL := instance.devPortalRestURL + "/subscriptions"

		request = base.CreateGet(subsGetURL)

		base.SetDefaultRestAPIHeaders(instance.accessToken, request)

		values := url.Values{}
		values.Add("apiId", api.ID)

		request.URL.RawQuery = values.Encode()

		base.LogRequest("apim.DeleteAllSubscriptions() getting Subscriptions", request)

		response = base.SendHTTPRequest(request)

		defer response.Body.Close()

		base.ValidateAndLogResponse("apim.DeleteAllSubscriptions() getting Subscriptions", response, 200)

		var subsResponse SubscriptionList
		json.NewDecoder(response.Body).Decode(&subsResponse)

		for _, sub := range subsResponse.List {
			instance.deleteSubscription(sub.SubscriptionID)
		}
	}
}

// DeleteAllApplications : Delete All Applications from APIM
func (instance *Client) DeleteAllApplications() {
	appsGetURL := instance.devPortalRestURL + "/applications"

	request := base.CreateGet(appsGetURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteAllApplications()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteAllApplications() getting Applications", response, 200)

	var appsResponse ApplicationList
	json.NewDecoder(response.Body).Decode(&appsResponse)

	for _, app := range appsResponse.List {
		appsDeleteURL := instance.devPortalRestURL + "/applications/" + app.ApplicationID

		request = base.CreateDelete(appsDeleteURL)

		base.SetDefaultRestAPIHeaders(instance.accessToken, request)

		response = base.SendHTTPRequest(request)

		defer response.Body.Close()

		base.ValidateAndLogResponse("apim.DeleteAllApplications() deleting Applications", response, 200)
	}
}

// DeleteAllAPIs : Delete All APIs from APIM
func (instance *Client) DeleteAllAPIs() {
	apisGetURL := instance.publisherRestURL + "/apis"

	request := base.CreateGet(apisGetURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteAllAPIs() getting APIs", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteAllAPIs() getting APIs", response, 200)

	var apisResponse APIList
	json.NewDecoder(response.Body).Decode(&apisResponse)

	for _, api := range apisResponse.List {
		apisDeleteURL := instance.publisherRestURL + "/apis/" + api.ID
		request = base.CreateDelete(apisDeleteURL)

		base.SetDefaultRestAPIHeaders(instance.accessToken, request)

		base.LogRequest("apim.DeleteAllAPIs() deleting APIs", request)

		response = base.SendHTTPRequest(request)
		defer response.Body.Close()

		base.ValidateAndLogResponse("apim.DeleteAllAPIs() deleting APIs", response, 200)
	}
}

// DeleteAllAPIProducts : Delete All API Products from APIM
func (instance *Client) DeleteAllAPIProducts() {
	apiProductsGetURL := instance.publisherRestURL + "/api-products"

	request := base.CreateGet(apiProductsGetURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteAllAPIProducts() getting API Products", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteAllAPIProducts() getting API Products", response, 200)

	var apiProductsResponse APIProductList
	json.NewDecoder(response.Body).Decode(&apiProductsResponse)

	for _, apiProduct := range apiProductsResponse.List {
		apiProductsDeleteURL := instance.publisherRestURL + "/api-products/" + apiProduct.ID

		request = base.CreateDelete(apiProductsDeleteURL)

		base.SetDefaultRestAPIHeaders(instance.accessToken, request)

		base.LogRequest("apim.DeleteAllAPIProducts() deleting API Products", request)

		response = base.SendHTTPRequest(request)

		defer response.Body.Close()

		base.ValidateAndLogResponse("apim.DeleteAllAPIProducts() deleting API Products", response, 200)
	}
}

// RemoveAllEndpointCerts : Remove All Endpoint Certs from the Truststore
func (instance *Client) RemoveAllEndpointCerts() {
	apisGetURL := instance.publisherRestURL + "/endpoint-certificates"

	request := base.CreateGet(apisGetURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.RemoveAllEndpointCerts() getting Certs", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.RemoveAllEndpointCerts() getting Certs", response, 200)

	var certificatesResponse Certificates
	json.NewDecoder(response.Body).Decode(&certificatesResponse)

	for _, certificate := range certificatesResponse.List {
		instance.RemoveEndpointCert(certificate.Alias)
	}
}

// RemoveEndpointCert : Remove Endpoint Cert from the Truststore
func (instance *Client) RemoveEndpointCert(alias string) {
	certificatesDeleteURL := instance.publisherRestURL + "/endpoint-certificates/" + alias
	request := base.CreateDelete(certificatesDeleteURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.RemoveEndpointCert() deleting Cert", request)

	response := base.SendHTTPRequest(request)
	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.RemoveEndpointCert() deleting Cert", response, 200)
}

// Retrieve admin credentials
func RetrieveAdminCredentialsInsteadCreator(username, password string) (string, string) {
	newUsername := username
	newPassword := password
	if strings.EqualFold(adminservices.CreatorUsername, username) {
		newUsername = adminservices.AdminUsername
		newPassword = adminservices.AdminPassword
	}
	if strings.EqualFold(adminservices.CreatorUsername+"@"+adminservices.Tenant1, username) {
		newUsername = adminservices.AdminUsername + "@" + adminservices.Tenant1
		newPassword = adminservices.AdminPassword
	}
	return newUsername, newPassword
}

// Get log level of APIs
func (instance *Client) GetAPILogLevel(username, password, tenantDomain, apiId string) (*APILogLevelList, error) {
	url := instance.devopsRestURL + "/tenant-logs/" + tenantDomain + "/apis/" + apiId
	request := base.CreateGet(url)
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	request.Header.Set("Authorization", "Basic "+encoded)
	request.Header.Set("Content-Type", "application/json")

	base.LogRequest("apim.GetAPILogLevel()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	if response.StatusCode == 200 {
		var apiLogLevelList APILogLevelList
		json.NewDecoder(response.Body).Decode(&apiLogLevelList)
		return &apiLogLevelList, nil
	}

	return nil, errors.New("Error with status code " + response.Status + "while retrieving log levels.")
}

// Set log level of an API
func (instance *Client) SetAPILogLevel(username, password, tenantDomain, apiId, logLevel string) (*APILogLevel, error) {
	url := instance.devopsRestURL + "/tenant-logs/" + tenantDomain + "/apis/" + apiId
	data, err := json.Marshal(APILogLevel{LogLevel: logLevel})

	if err != nil {
		return nil, errors.New("Error while building request payload for setting log level of API " + apiId + ".")
	}

	request := base.CreatePut(url, bytes.NewBuffer(data))
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	request.Header.Set("Authorization", "Basic "+encoded)
	request.Header.Set("Content-Type", "application/json")

	base.LogRequest("apim.SetAPILogLevel()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	if response.StatusCode == 200 {
		var apiLogLevel APILogLevel
		json.NewDecoder(response.Body).Decode(&apiLogLevel)
		return &apiLogLevel, nil
	}

	return nil, errors.New("Error with status code " + response.Status + "while setting log level of API " + apiId + ".")
}

func generateSampleAPIOperations() []APIOperations {

	op1 := APIOperations{}
	op1.Target = "/order/{orderId}"
	op1.Verb = "GET"
	op1.ThrottlingPolicy = "Unlimited"
	op1.AuthType = "Application & Application User"
	op1.OperationPolicies = OperationPolicies{[]string{}, []string{}, []string{}}

	op2 := APIOperations{}
	op2.Target = "/order/{orderId}"
	op2.Verb = "DELETE"
	op2.ThrottlingPolicy = "Unlimited"
	op2.AuthType = "Application & Application User"
	op2.OperationPolicies = OperationPolicies{[]string{}, []string{}, []string{}}

	op3 := APIOperations{}
	op3.Target = "/order/{orderId}"
	op3.Verb = "PUT"
	op3.ThrottlingPolicy = "Unlimited"
	op3.AuthType = "Application & Application User"
	op3.OperationPolicies = OperationPolicies{[]string{}, []string{}, []string{}}

	op4 := APIOperations{}
	op4.Target = "/menu"
	op4.Verb = "GET"
	op4.ThrottlingPolicy = "Unlimited"
	op4.AuthType = "Application & Application User"
	op4.OperationPolicies = OperationPolicies{[]string{}, []string{}, []string{}}

	op5 := APIOperations{}
	op5.Target = "/order"
	op5.Verb = "POST"
	op5.ThrottlingPolicy = "Unlimited"
	op5.AuthType = "Application & Application User"
	op5.OperationPolicies = OperationPolicies{[]string{}, []string{}, []string{}}

	return []APIOperations{op1, op2, op3, op4, op5}
}

func (instance *Client) getToken(username string, password string) string {
	registrationResponse := instance.registerClient(username, password)

	request := base.CreatePostEmptyBody(instance.tokenURL)
	request.SetBasicAuth(registrationResponse.ClientID, registrationResponse.ClientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	values := url.Values{}
	values.Add("grant_type", "password")
	values.Add("username", username)
	values.Add("password", password)
	values.Add("scope",
		"apim:admin apim:api_view apim:api_create apim:api_publish apim:subscribe apim:api_delete "+
			"apim:app_import_export apim:api_import_export apim:api_product_import_export apim:app_manage apim:sub_manage")

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.getToken()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.getToken()", response, 200)

	var jsonResp tokenResponse
	json.NewDecoder(response.Body).Decode(&jsonResp)
	return jsonResp.AccessToken
}

func (instance *Client) registerClient(username string, password string) dcrResponse {
	dcrPayload := dcrRequest{}

	dcrPayload.CallbackURL = "http://localhost"
	dcrPayload.ClientName = instance.restClientName
	dcrPayload.IsSaaSApp = true
	dcrPayload.Owner = username
	dcrPayload.SupportedGrantTypes = "password refresh_token"

	data, err := json.Marshal(dcrPayload)

	if err != nil {
		base.Fatal(err)
	}

	request := base.CreatePost(instance.dcrURL, bytes.NewBuffer(data))

	request.SetBasicAuth(username, password)
	request.Header.Set("Content-Type", "application/json")

	base.LogRequest("apim.registerClient()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	contentType := response.Header["Content-Type"][0]

	if response.StatusCode != 200 {
		base.FatalStatusCodeResponse("apim.registerClient()", response)
	}

	// If DCR endpoint being invoked is invalid, an HTML error page will be returned.
	// We cannot rely on checking the response code since it will always be 200.
	// Therefore need to validate the Content Type of the response to detect this condition.
	if contentType != "application/json" {
		base.FatalContentTypeResponse("apim.registerClient()", response)
	}

	base.LogResponse("apim.registerClient()", response)

	var jsonResp dcrResponse
	json.NewDecoder(response.Body).Decode(&jsonResp)

	return jsonResp
}

// AddThrottlePolicy : Add new Throttle Policy of different policy types to APIM
func (instance *Client) AddThrottlePolicy(t *testing.T, policy interface{}, username, password, policyType string, doClean bool) map[string]interface{} {
	var throttlePolicyResponse map[string]interface{}

	throttlePolicyURL := instance.adminRestURL + "/throttling/policies/" + policyType

	data, err := json.Marshal(policy)
	if err != nil {
		t.Fatal(err)
	}

	request := base.CreatePost(throttlePolicyURL, bytes.NewBuffer(data))

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.AddThrottlePolicy()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.AddThrottlePolicy()", response, 201)

	json.NewDecoder(response.Body).Decode(&throttlePolicyResponse)
	policyId := fmt.Sprintf("%v", throttlePolicyResponse["policyId"])
	if doClean {
		t.Cleanup(func() {
			instance.Login(username, password)
			instance.DeleteThrottlePolicy(policyId, policyType)
		})
	}
	return throttlePolicyResponse
}

// DeleteThrottlePolicy : Deletes Throttling Policy from APIM using UUID
func (instance *Client) DeleteThrottlePolicy(policyID, policyType string) {

	policiesURL := instance.adminRestURL + "/throttling/policies/" + policyType + "/" + policyID

	request := base.CreateDelete(policiesURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.DeleteThrottlePolicy()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.DeleteThrottlePolicy()", response, 200)
}

// GetThrottlePolicy : Get Throttle Policy from APIM using UUID
func (instance *Client) GetThrottlePolicy(policyID, policyType string) map[string]interface{} {
	var policyResponse map[string]interface{}

	policiesURL := instance.adminRestURL + "/throttling/policies/" + policyType + "/" + policyID

	request := base.CreateGet(policiesURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetThrottlePolicy()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetThrottlePolicy()", response, 200)

	json.NewDecoder(response.Body).Decode(&policyResponse)

	return policyResponse
}

// DeleteThrottlePolicyByName : Deletes Throttling Policy from APIM using policy name
func (instance *Client) DeleteThrottlePolicyByName(t *testing.T, policyName, policyType string) {

	policyID := instance.GetThrottlePolicyID(t, policyName, policyType)
	instance.DeleteThrottlePolicy(policyID, policyType)
}

// GetThrottlePolicyID : Get Throttle Policy UUID using policy name from APIM
func (instance *Client) GetThrottlePolicyID(t *testing.T, policyName, policyType string) string {
	var policyListResponse utils.ThrottlingPoliciesDetailsList
	var uuid string

	queryType := ""
	switch policyType {
	case ApplicationThrottlePolicyType:
		queryType = applicationThrottlePolicyQuery
	case CustomThrottlePolicyType:
		queryType = customThrottlePolicyQuery
	case AdvancedThrottlePolicyType:
		queryType = advancedThrottlePolicyQuery
	case SubscriptionThrottlePolicyType:
		queryType = subscriptionThrottlePolicyQuery
	}

	getPoliciesURL := instance.adminRestURL + "/throttling/policies/search/?query=type:" + queryType

	request := base.CreateGet(getPoliciesURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetThrottlePolicyID()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetThrottlePolicyID()", response, 200)

	err := json.NewDecoder(response.Body).Decode(&policyListResponse)
	if err != nil {
		t.Fatal(err)
	}
	throttlePolicyList := policyListResponse.List
	for _, policy := range throttlePolicyList {
		if policy.PolicyName == policyName {
			uuid = policy.Uuid
			break
		}
	}
	return uuid
}

// GetThrottlePolicies : Get Throttle Policies list of all types from APIM
func (instance *Client) GetThrottlePolicies(t *testing.T) *utils.ThrottlingPoliciesDetailsList {
	var policyListResponse *utils.ThrottlingPoliciesDetailsList

	getPoliciesURL := instance.adminRestURL + "/throttling/policies/search/?query=type:all"

	request := base.CreateGet(getPoliciesURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	base.LogRequest("apim.GetThrottlePolicies()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetThrottlePolicies()", response, 200)

	err := json.NewDecoder(response.Body).Decode(&policyListResponse)
	if err != nil {
		t.Fatal(err)
	}

	return policyListResponse
}

// GenerateSampleThrottlePolicyData : Generate sample ThrottlePolicy of a specific throttling policy type
func (instance *Client) GenerateSampleThrottlePolicyData(policyType string) interface{} {
	const (
		policyString      = "Policy"
		requestCountLimit = "REQUESTCOUNTLIMIT"
		timeUnit          = "min"
		unitTime          = 10
		requestCount      = 5
	)

	switch policyType {
	case ApplicationThrottlePolicyType:
		policy := ApplicationThrottlePolicy{}
		policy.PolicyName = base.GenerateRandomString() + policyString
		policy.DisplayName = "This is a Test Application Policy"
		policy.Description = "This is a Test Application Policy"
		policy.IsDeployed = false
		policy.Type = "ApplicationThrottlePolicy"
		policy.DefaultLimit.Type = requestCountLimit
		policy.DefaultLimit.RequestCount.TimeUnit = timeUnit
		policy.DefaultLimit.RequestCount.UnitTime = unitTime
		policy.DefaultLimit.RequestCount.RequestCount = requestCount
		return &policy
	case CustomThrottlePolicyType:
		policy := CustomThrottlePolicy{}
		policy.PolicyName = base.GenerateRandomString() + policyString
		policy.Description = "This is a Test Custom Policy"
		policy.IsDeployed = false
		policy.Type = "CustomRule"
		policy.SiddhiQuery = "FROM RequestStream\\nSELECT userId, ( userId == 'admin@carbon.super' ) AS isEligible , " +
			"str:concat('admin@carbon.super','') as throttleKey\\nINSERT INTO EligibilityStream; \\n\\nFROM " +
			"EligibilityStream[isEligible==true]#throttler:timeBatch(1 min) \\nSELECT throttleKey, (count(userId) >= 10) " +
			"as isThrottled, expiryTimeStamp group by throttleKey \\nINSERT ALL EVENTS into ResultStream;\n"
		policy.KeyTemplate = "$userId"
		return &policy
	case AdvancedThrottlePolicyType:
		policy := AdvancedThrottlePolicy{}
		conditionalGroup := AdvancedPolicyConditionalGroup{}
		condition := AdvancedPolicyCondition{}
		policy.PolicyName = base.GenerateRandomString() + policyString
		policy.Description = "This is a Test Advanced Policy"
		policy.IsDeployed = false
		policy.Type = "AdvancedThrottlePolicy"
		policy.DefaultLimit.Type = requestCountLimit
		policy.DefaultLimit.RequestCount.TimeUnit = timeUnit
		policy.DefaultLimit.RequestCount.UnitTime = unitTime
		policy.DefaultLimit.RequestCount.RequestCount = requestCount
		conditionalGroup.Description = "Sample description about condition group"
		condition.Type = "HEADERCONDITION"
		condition.HeaderCondition.HeaderName = "Test"
		condition.HeaderCondition.HeaderValue = "TestValue"
		conditionalGroup.Conditions = []AdvancedPolicyCondition{condition}
		conditionalGroup.Limit.Type = requestCountLimit
		conditionalGroup.Limit.RequestCount.TimeUnit = timeUnit
		conditionalGroup.Limit.RequestCount.UnitTime = unitTime
		conditionalGroup.Limit.RequestCount.RequestCount = requestCount
		policy.ConditionalGroups = []AdvancedPolicyConditionalGroup{conditionalGroup}
		return &policy
	case SubscriptionThrottlePolicyType:
		policy := SubscriptionThrottlePolicy{}
		policy.PolicyName = base.GenerateRandomString() + policyString
		policy.Description = "This is a Test Subscription Policy"
		policy.IsDeployed = false
		policy.Type = "SubscriptionThrottlePolicy"
		policy.DefaultLimit.Type = requestCountLimit
		policy.DefaultLimit.RequestCount.TimeUnit = timeUnit
		policy.DefaultLimit.RequestCount.UnitTime = unitTime
		policy.DefaultLimit.RequestCount.RequestCount = requestCount
		policy.Permissions.PermissionType = "ALLOW"
		policy.Permissions.Roles = []string{"admin"}
		return &policy
	}
	return nil
}

func getDCRURL(host, version string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/client-registration/" + version + "/register"
}

func getApimURL(host string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port)
}

func getDevPortalRestURL(host, version string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/api/am/devportal/" + version
}

func getPublisherRestURL(host, version string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/api/am/publisher/" + version
}

func getAdminRestURL(host, version string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/api/am/admin/" + version
}

func getDevOpsRestURL(host, version string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/api/am/devops/" + version
}

func getTokenURL(host string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/oauth2/token"
}
