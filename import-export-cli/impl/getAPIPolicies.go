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
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// TODO: Add Policy Version
const (
	apiPolicyUUIDHeader               = "ID"
	apiPolicyNameHeader               = "NAME"
	apiPolicyDisplayNameHeader        = "DISPLAY NAME"
	apiPolicyVersion                  = "VERSION"
	apiPolicyCategoryHeader           = "CATEGORY"
	apiPolicyApplicableFlowsHeaders   = "APPLICABLE FLOWS"
	apiPolicySupportedGatewaysHeaders = "SUPPORTED GATEWAYS"
	defaultAPIPolicyTableFormat       = "table {{.ID}}\t{{.Name}}\t{{.DisplayName}}\t{{.Version}}\t{{.Category}}\t{{.ApplicableFlows}}\t{{.SupportedGateways}}"
)

type apiPolicy struct {
	id                string
	name              string
	displayName       string
	version           string
	category          string
	applicableFlows   []string
	supportedGateways []string
}

func newAPIPolicyDefinition(a utils.APIPolicy) *apiPolicy {
	return &apiPolicy{a.Id, a.Name, a.DisplayName, a.Version, a.Category, a.ApplicableFlows, a.SupportedGateways}
}

func (a apiPolicy) ID() string {
	return a.id
}

func (a apiPolicy) Name() string {
	return a.name
}

func (a apiPolicy) DisplayName() string {
	return a.displayName
}

func (a apiPolicy) Version() string {
	return a.version
}

func (a apiPolicy) Category() string {
	return a.category
}

func (a apiPolicy) ApplicableFlows() []string {
	return a.applicableFlows
}

func (a apiPolicy) SupportedGateways() []string {
	return a.supportedGateways
}

func GetAPIPolicyListFromEnv(accessToken, environment, limit string) (*resty.Response, error) {
	apiPolicyListEndpoint := utils.GetPublisherEndpointOfEnv(environment, utils.MainConfigFilePath)
	return getAPIPolicyList(accessToken, apiPolicyListEndpoint, limit)
}

func getAPIPolicyList(accessToken, apiPolicyListEndpoint, limit string) (*resty.Response, error) {
	apiPolicyListEndpoint = utils.AppendSlashToString(apiPolicyListEndpoint)
	apiPolicyResource := "operation-policies"

	if limit != "" {
		query := `?limit=` + limit
		apiPolicyResource += query
	}

	url := apiPolicyListEndpoint + apiPolicyResource

	utils.Logln(utils.LogPrefixInfo+"GetAPIPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)

	return resp, err
}

// PrintAPIPolicies prints the policy list in a specific format
func PrintAPIPolicies(resp *resty.Response, format string) {
	var apiPolicyList utils.APIPoliciesList
	err := json.Unmarshal(resp.Body(), &apiPolicyList)
	policies := apiPolicyList.List

	if err != nil {
		utils.HandleErrorAndExit("Error unmarshalling response data", err)
	}
	if format == "" {
		format = defaultAPIPolicyTableFormat
		// create policy context with standard output
		policyContext := formatter.NewContext(os.Stdout, format)
		// create a new renderer function which iterate collection
		renderer := func(w io.Writer, t *template.Template) error {
			for _, policy := range policies {
				if err := t.Execute(w, newAPIPolicyDefinition(policy)); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}

		// headers for table
		apiPolicyTableHeaders := map[string]string{
			"ID":                apiPolicyUUIDHeader,
			"Name":              apiPolicyNameHeader,
			"DisplayName":       apiPolicyDisplayNameHeader,
			"Version":           apiPolicyVersion,
			"Category":          apiPolicyCategoryHeader,
			"ApplicableFlows":   apiPolicyApplicableFlowsHeaders,
			"SupportedGateways": apiPolicySupportedGatewaysHeaders,
		}
		// execute context
		if err := policyContext.Write(renderer, apiPolicyTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else if format == utils.JsonArrayFormatType {
		utils.ListArtifactsInJsonArrayFormat(policies, utils.ProjectTypeAPIPolicy)
	}
}
