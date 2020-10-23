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
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	appIdHeader      = "ID"
	appNameHeader    = "NAME"
	appOwnerHeader   = "OWNER"
	appStatusHeader  = "STATUS"
	appGroupIdHeader = "GROUP ID"

	defaultAppTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Owner}}\t{{.Status}}\t{{.GroupId}}"
)

// app contains information about util.Application
type app struct {
	id      string
	name    string
	owner   string
	status  string
	groupId string
}

// creates a new app definition from utils.Application
func newAppDefinitionFromApplication(a utils.Application) *app {
	return &app{a.ID, a.Name, a.Owner, a.Status, a.GroupID}
}

// Id of application
func (a app) Id() string {
	return a.id
}

// Name of application
func (a app) Name() string {
	return a.name
}

// Owner of application
func (a app) Owner() string {
	return a.owner
}

// Status of application
func (a app) Status() string {
	return a.status
}

// GroupId of application
func (a app) GroupId() string {
	return a.groupId
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (a *app) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

//GetApplicationListFromEnv
// @param accessToken : Access Token for the environment
// @param environment : Environment to get the list of applications
// @param appOwner : Owner of the applications
// @param limit : Max number of results to return
// @return count (no. of Applications)
// @return array of Application objects
// @return error
func GetApplicationListFromEnv(accessToken, environment, appOwner, limit string) (count int32, apps []utils.Application, err error) {
	applicationListEndpoint := utils.GetAdminApplicationListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return GetApplicationList(accessToken, applicationListEndpoint, appOwner, limit)
}

//Get Application List
// @param accessToken : Access Token for the environment
// @param applicationListEndpoint : Endpoint to use for listing applications
// @param appOwner : Owner of the applications
// @param limit : Max number of results to return
// @return count (no. of Applications)
// @return array of Application objects
// @return error
func GetApplicationList(accessToken, applicationListEndpoint, appOwner, limit string) (count int32, apps []utils.Application,
	err error) {

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	if limit != "" {
		applicationListEndpoint += "?limit=" + limit
	}

	var resp *resty.Response
	if appOwner == "" {
		resp, err = utils.InvokeGETRequest(applicationListEndpoint, headers)
	} else {
		resp, err = utils.InvokeGETRequestWithQueryParam("user", appOwner, applicationListEndpoint, headers)
	}
	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+applicationListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		appListResponse := &utils.ApplicationListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &appListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return appListResponse.Count, appListResponse.List, nil

	} else {
		return 0, nil, errors.New(resp.Status())
	}
}

// extractAppDefinition extracts ApplicationDefinition from jsonContent
func extractAppDefinition(jsonContent []byte) (*v2.ApplicationDefinition, error) {
	application := &v2.ApplicationDefinition{}
	err := json.Unmarshal(jsonContent, &application)
	if err != nil {
		return nil, err
	}

	return application, nil
}

// GetApplicationDefinition scans filePath and returns ApplicationDefinition or an error
func GetApplicationDefinition(filePath string) (*v2.ApplicationDefinition, []byte, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}

	var buffer []byte
	if info.IsDir() {
		_, content, err := resolveYamlOrJSON(path.Join(filePath, filepath.Base(filePath)))
		if err != nil {
			return nil, nil, err
		}
		buffer = content
	} else {
		return nil, nil, fmt.Errorf("looking for directory, found %s", info.Name())
	}
	api, err := extractAppDefinition(buffer)
	if err != nil {
		return nil, nil, err
	}
	return api, buffer, nil
}

// PrintApps
func PrintApps(apps []utils.Application, format string) {
	if format == "" {
		format = defaultAppTableFormat
	}
	// create new app context with standard output
	appContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection of apps
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apps {
			if err := t.Execute(w, newAppDefinitionFromApplication(a)); err != nil {
				return err
			}
			// write a new line after executing template
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	appTableHeaders := map[string]string{
		"Id":      appIdHeader,
		"Name":    appNameHeader,
		"Status":  appStatusHeader,
		"Owner":   appOwnerHeader,
		"GroupId": appGroupIdHeader,
	}

	// execute context
	if err := appContext.Write(renderer, appTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
