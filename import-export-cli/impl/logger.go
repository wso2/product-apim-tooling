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

package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
)

const (
	loggingApiIdHeader       = "API_ID"
	loggingApiContextHeader  = "API_CONTEXT"
	loggingApiLogLevelHeader  = "LOG_LEVEL"

	defaultLoggingApiTableFormat = "table {{.Id}}\t{{.Context}}\t{{.LogLevel}}"
)

// LoggingApi holds information about an API for outputting
type loggingApi struct {
	id              string
	context         string
	logLevel        string
}

// Creates a new api from utils.APILogger
func newLoggingApiDefinitionFromAPI(a utils.APILogger) *loggingApi {
	return &loggingApi{a.ID, a.Context, a.LogLevel}
}

// Id of loggingApi
func (a loggingApi) Id() string {
	return a.id
}

// Context of loggingApi
func (a loggingApi) Context() string {
	return a.context
}

// LogLevel of loggingApi
func (a loggingApi) LogLevel() string {
	return a.logLevel
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (a *loggingApi) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

var tenantDomain = utils.DefaultTenantDomain

// Add new APILogger object and use it in the array
func GetPerAPILoggingListFromEnv(credential credentials.Credential, environment string) (apis []utils.APILogger, err error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	if strings.Contains(credential.Username, "@") {
		splits := strings.Split(credential.Username, "@")
		tenantDomain = splits[len(splits)-1]
	}
	apiListEndpoint := utils.GetAPILoggingListEndpointOfEnv(environment, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", apiListEndpoint)
	resp, err := utils.InvokeGETRequest(apiListEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to " + apiListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiLoggerListResponse := &utils.APILoggerListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiLoggerListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return apiLoggerListResponse.Apis, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}

// PrintAPILoggers
func PrintAPILoggers(apis []utils.APILogger, format string) {
	if format == "" {
		format = defaultLoggingApiTableFormat
	}
	// Create API context with standard output
	apiContext := formatter.NewContext(os.Stdout, format)

	// Create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apis {
			if err := t.Execute(w, newLoggingApiDefinitionFromAPI(a)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// Headers for table
	apiLoggerTableHeaders := map[string]string{
		"Id":              loggingApiIdHeader,
		"Context":         loggingApiContextHeader,
		"LogLevel":        loggingApiLogLevelHeader,
	}

	// Execute context
	if err := apiContext.Write(renderer, apiLoggerTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}

// Check how to send the API details in a single object instead of an array
func GetPerAPILoggingDetailsFromEnv(credential credentials.Credential, environment, apiId string) (apis []utils.APILogger, err error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	if strings.Contains(credential.Username, "@") {
		splits := strings.Split(credential.Username, "@")
		tenantDomain = splits[len(splits)-1]
	}
	apiDetailsEndpoint := utils.GetAPILoggingDetailsEndpointOfEnv(environment, apiId, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", apiDetailsEndpoint)
	resp, err := utils.InvokeGETRequest(apiDetailsEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiDetailsEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiLoggerListResponse := &utils.APILoggerListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiLoggerListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return apiLoggerListResponse.Apis, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}	
}

// SetAPILoggingLevel
func SetAPILoggingLevel(credential credentials.Credential, environment, apiId, logLevel string) (*resty.Response, error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	if strings.Contains(credential.Username, "@") {
		splits := strings.Split(credential.Username, "@")
		tenantDomain = splits[len(splits)-1]
	}
	apiSetEndpoint := utils.GetAPILoggingSetEndpointOfEnv(environment, apiId, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", apiSetEndpoint)
	body := `{"logLevel":"` + logLevel + `"}`
	resp, err := utils.InvokePutRequest(nil, apiSetEndpoint, headers, body)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiSetEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		return resp, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}
