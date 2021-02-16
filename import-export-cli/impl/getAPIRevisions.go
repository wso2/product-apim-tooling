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
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	revisionIdHeader          = "ID"
	revisionNameHeader        = "REVISION"
	revisionDescriptionHeader = "DESCRIPTION"
	deployedGatewaysHeader    = "GATEWAYS"

	defaultRevisionTableFormat = "table {{.Id}}\t{{.RevisionNumber}}\t{{.Description}}\t{{.Gateways}}"
)

// revisions struct holds information about an revision for outputting
type revision struct {
	id               string
	revisionNumber   string
	description      string
	deployedGateways []string
}

// creates a new revision from utils.Revisions
func newRevisionDefinitionFromRevisions(r utils.Revisions) *revision {
	return &revision{r.ID, r.RevisionNumber, r.Description, r.Gateways}
}

// Id of revision
func (r revision) Id() string {
	return r.id
}

// Revision number
func (r revision) RevisionNumber() string {
	return strings.ReplaceAll(r.revisionNumber, "Revision ", "")
}

// Revision description
func (r revision) Description() string {
	return r.description
}

// Deployed gateways of the revision
func (r revision) Gateways() []string {
	return r.deployedGateways
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (r *revision) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(r)
}

// GetRevisionListFromEnv
// @param accessToken	: Access Token for the environment
// @param environment	: Environment name to use when getting the API List
// @param apiName		: Name of the API
// @param apiVersion	: Version of the API
// @param provider		: Provider of the API
// @param query			: Query param for the filtering the revisions based on the deployed status
// @return count (no. of APIs)
// @return array of revision objects
// @return error
func GetRevisionListFromEnv(accessToken, environment, apiName, apiVersion, provider, query string) (count int32, revisions []utils.Revisions, err error) {
	apiId, err := GetAPIId(accessToken, environment, apiName, apiVersion, provider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id to list revisions ", err)
	}
	revisionListEndpoint := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	revisionListEndpoint = utils.AppendSlashToString(revisionListEndpoint)
	url := revisionListEndpoint + apiId + "/revisions"
	if query != "" {
		url += "?query=" + query
	}
	return GetRevisionsList(accessToken, url)
}

// Print Revisions in the given template
// @param revisions	Available revisions list for the API
// @param format	Format type of the output
func PrintRevisions(revisions []utils.Revisions, format string) {
	if format == "" {
		format = defaultRevisionTableFormat
	}
	// create revision Context with standard output
	revisionContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, r := range revisions {
			var gateways []string
			for _, d := range r.Deployments {
				gateways = append(gateways, d.Name)
			}
			r.Gateways = gateways
			if err := t.Execute(w, newRevisionDefinitionFromRevisions(r)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	revisionTableHeaders := map[string]string{
		"Id":             revisionIdHeader,
		"RevisionNumber": revisionNameHeader,
		"Description":    revisionDescriptionHeader,
		"Gateways":       deployedGatewaysHeader,
	}

	// execute context
	if err := revisionContext.Write(renderer, revisionTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
