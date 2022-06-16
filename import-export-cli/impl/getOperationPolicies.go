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

const (
	operationPolicyUUIDHeader              = "ID"
	operationPolicyNameHeader              = "NAME"
	operationPolicyCategoryHeader          = "CATEGORY"
	operationPolicyApplicationFlowsHeaders = "APPLICATION FLOWS"
	defaultOperationPolicyTableFormat      = "table {{.ID}}\t{{.Name}}\t{{.Category}}\t{{.Applicatio Flows}}"
)

type operationPolicy struct {
	Id                string   `json:"id"`
	DisplayName       string   `json:"displayName"`
	Category          string   `json:"category"`
	ApplicationFlows  []string `json:"applicationFlows"`
	SupportedGateways []string `json:"supportedGateways"`
}

func newOperationPolicyDefinition(a utils.OperationPolicy) *operationPolicy {
	return &operationPolicy{a.Id, a.DisplayName, a.Category, a.ApplicationFlows, a.SupportedGateways}
}

func (a operationPolicy) ID() string {
	return a.Id
}

func (a operationPolicy) PolicyDisplayName() string {
	return a.DisplayName
}

func (a operationPolicy) PolicyApplicationFlows() []string {
	return a.ApplicationFlows
}

func GetOperationPolicyListFromEnv(accessToken, environment, query string) (*resty.Response, error) {
	operationPolicyListEndpoint := utils.GetOperationPolicyListEndpointOfEnv(environment, utils.MainConfigFilePath)
	fmt.Println("Endpoint: ", operationPolicyListEndpoint)
	return getOperationPolicyList(accessToken, operationPolicyListEndpoint, query)
}

func getOperationPolicyList(accessToken string, operationPolicyListEndpoint string, query string) (*resty.Response, error) {
	url := operationPolicyListEndpoint
	queryParamString := "query=" + query
	utils.Logln(utils.LogPrefixInfo+"GetOperationPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	if query == "" {
		fmt.Println("No Query")
		resp, err := utils.InvokeGETRequest(url, headers)
		return resp, err
	} else {
		resp, err := utils.InvokeGETRequestWithQueryParamsString(url, queryParamString, headers)
		return resp, err
	}
}

// PrintOperationPolicies prints the policy list in a specific format
func PrintOperationPolicies(resp *resty.Response, format string) {
	var operationPolicyList utils.OperationPoliciesList
	err := json.Unmarshal(resp.Body(), &operationPolicyList)
	policies := operationPolicyList.List
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshalling response data", err)
	}
	if format == "" {
		format = defaultThrottlePolicyTableFormat
		// create policy context with standard output
		policyContext := formatter.NewContext(os.Stdout, format)
		// create a new renderer function which iterate collection
		renderer := func(w io.Writer, t *template.Template) error {
			for _, policy := range policies {
				if err := t.Execute(w, newOperationPolicyDefinition(policy)); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}
		// headers for table
		operationPolicyTableHeaders := map[string]string{
			"UUID":              operationPolicyUUIDHeader,
			"Name":              operationPolicyNameHeader,
			"CATEGORY":          operationPolicyCategoryHeader,
			"APPLICATION FLOWS": operationPolicyApplicationFlowsHeaders,
		}
		// execute context
		if err := policyContext.Write(renderer, operationPolicyTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else if format == utils.JsonArrayFormatType {
		utils.ListArtifactsInJsonArrayFormat(policies, utils.ProjectTypePolicy)
	}
}
