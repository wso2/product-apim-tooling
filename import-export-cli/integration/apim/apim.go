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
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

const (
	restClientPostfix = "-integ-rest-client"
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
func (instance *Client) Setup(envName string, host string, offset int, dcrVersion string, restAPIVersion string) {
	base.Log("apim.Setup() - envName:", envName, ",host:", host, ",offset:", offset, ",dcrVersion:", dcrVersion,
		",restAPIVersion:", restAPIVersion)
	instance.apimURL = getApimURL(host, offset)
	instance.dcrURL = getDCRURL(host, offset, dcrVersion)
	instance.devPortalRestURL = getDevPortalRestURL(host, offset, restAPIVersion)
	instance.publisherRestURL = getPublisherRestURL(host, offset, restAPIVersion)
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
	api.Name = generateRandomString() + "API"
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
	api.EndpointImplementationType = "ENDPOINT"
	api.GatewayEnvironments = []string{"Production and Sandbox"}
	api.BusinessInformation = BusinessInfo{"Jane Roe", "marketing@pizzashack.com", "John Doe", "architecture@pizzashack.com"}
	api.EndpointConfig = HTTPEndpoint{"http", &URLConfig{"https://localhost:" + strconv.Itoa(9443+instance.portOffset) + "/am/sample/pizzashack/v1/api/"},
		&URLConfig{"https://localhost:" + strconv.Itoa(9443+instance.portOffset) + "/am/sample/pizzashack/v1/api/"}}
	api.Operations = generateSampleAPIOperations()

	return &api
}

func getContext(provider string) string {
	context := generateRandomString()
	if strings.Contains(provider, "@") {
		splits := strings.Split(provider, "@")
		domain := splits[len(splits)-1]
		return "/t/" + domain + "/" + context
	}

	return "/" + context
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

	// Copy GatewayEnvironments slice
	apiCopy.GatewayEnvironments = make([]string, len(apiToCopy.GatewayEnvironments))
	copy(apiCopy.GatewayEnvironments, apiToCopy.GatewayEnvironments)

	// Copy Labels slice
	apiCopy.Labels = make([]string, len(apiToCopy.Labels))
	copy(apiCopy.Labels, apiToCopy.Labels)

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

	// Sort GatewayEnvironments slice
	sort.Strings(api.GatewayEnvironments)

	// Sort Labels slice
	sort.Strings(api.Labels)

	// Sort MediationPolicies slice
	sort.Sort(ByID(api.MediationPolicies))

	// Sort SubscriptionAvailableTenants slice
	sort.Strings(api.SubscriptionAvailableTenants)

	// Sort AdditionalProperties map
	sortAdditionalProperties(api)

	// Sort AccessControlRoles slice
	sort.Strings(api.AccessControlRoles)

	// Sort Operations slice
	sort.Sort(ByTargetVerb(api.Operations))
}

func sortAdditionalProperties(api *API) {
	keys := make([]string, 0, len(api.AdditionalProperties))

	for key := range api.AdditionalProperties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	sortedMap := make(map[string]string, len(api.AdditionalProperties))

	for _, key := range keys {
		sortedMap[key] = api.AdditionalProperties[key]
	}

	api.AdditionalProperties = sortedMap
}

// GenerateSampleAppData : Generate sample Application object
func (instance *Client) GenerateSampleAppData() *Application {
	app := Application{}
	app.Name = generateRandomString() + "Application"
	app.ThrottlingPolicy = "Unlimited"
	app.Description = "Test Application"
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
func (instance *Client) AddAPI(t *testing.T, api *API, username string, password string) string {
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

	t.Cleanup(func() {
		instance.Login(username, password)
		instance.DeleteAPI(apiResponse.ID)
	})

	return apiResponse.ID
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

// DeleteAPIByName : Delete API from APIM by name
func (instance *Client) DeleteAPIByName(name string) {
	apiInfo := instance.GetAPIByName(name)
	instance.DeleteAPI(apiInfo.ID)
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

// GetAPIByName : Get API by name from APIM
func (instance *Client) GetAPIByName(name string) *APIInfo {
	apisURL := instance.publisherRestURL + "/apis"

	request := base.CreateGet(apisURL)

	base.SetDefaultRestAPIHeaders(instance.accessToken, request)

	values := url.Values{}
	values.Add("query", name)

	request.URL.RawQuery = values.Encode()

	base.LogRequest("apim.GetAPIByName()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("apim.GetAPIByName()", response, 200)

	var apiResponse APIList
	json.NewDecoder(response.Body).Decode(&apiResponse)
	return &apiResponse.List[0]
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
		subsDeleteURL := instance.devPortalRestURL + "/subscriptions/" + sub.SubscriptionID

		request = base.CreateDelete(subsDeleteURL)

		base.SetDefaultRestAPIHeaders(instance.accessToken, request)

		base.LogRequest("apim.DeleteSubscriptions() deleting Subscriptions", request)

		response = base.SendHTTPRequest(request)

		defer response.Body.Close()

		base.ValidateAndLogResponse("apim.DeleteSubscriptions() deleting Subscriptions", response, 200)
	}
}

// AddApplication : Add new Application to APIM
func (instance *Client) AddApplication(t *testing.T, application *Application, username string, password string) *Application {
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

	t.Cleanup(func() {
		instance.Login(username, password)
		instance.DeleteApplication(appResponse.ApplicationID)
	})

	return &appResponse
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
			subsDeleteURL := instance.devPortalRestURL + "/subscriptions/" + sub.SubscriptionID

			request = base.CreateDelete(subsDeleteURL)

			base.SetDefaultRestAPIHeaders(instance.accessToken, request)

			base.LogRequest("apim.DeleteAllSubscriptions() deleting Subscriptions", request)

			response = base.SendHTTPRequest(request)

			defer response.Body.Close()

			base.ValidateAndLogResponse("apim.DeleteAllSubscriptions() deleting Subscriptions", response, 200)
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

func generateSampleAPIOperations() []APIOperations {
	op1 := APIOperations{}
	op1.Target = "/order/{orderId}"
	op1.Verb = "GET"
	op1.ThrottlingPolicy = "Unlimited"
	op1.AuthType = "Application & Application User"

	op2 := APIOperations{}
	op2.Target = "/order/{orderId}"
	op2.Verb = "DELETE"
	op2.ThrottlingPolicy = "Unlimited"
	op2.AuthType = "Application & Application User"

	op3 := APIOperations{}
	op3.Target = "/order/{orderId}"
	op3.Verb = "PUT"
	op3.ThrottlingPolicy = "Unlimited"
	op3.AuthType = "Application & Application User"

	op4 := APIOperations{}
	op4.Target = "/menu"
	op4.Verb = "GET"
	op4.ThrottlingPolicy = "Unlimited"
	op4.AuthType = "Application & Application User"

	op5 := APIOperations{}
	op5.Target = "/order"
	op5.Verb = "POST"
	op5.ThrottlingPolicy = "Unlimited"
	op5.AuthType = "Application & Application User"

	return []APIOperations{op1, op2, op3, op4, op5}
}

func generateRandomString() string {
	b := make([]byte, 10)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
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
	values.Add("scope", "apim:api_view apim:api_create apim:api_publish apim:subscribe apim:api_delete")

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

func getDCRURL(host string, offset int, version string) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/client-registration/" + version + "/register"
}

func getApimURL(host string, offset int) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port)
}

func getDevPortalRestURL(host string, offset int, version string) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/api/am/store/" + version
}

func getPublisherRestURL(host string, offset int, version string) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/api/am/publisher/" + version
}

func getTokenURL(host string, offset int) string {
	port := 8243 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/token"
}
