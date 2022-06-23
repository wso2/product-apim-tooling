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
	operationPolicyUUIDHeader               = "ID"
	operationPolicyNameHeader               = "NAME"
	operationPolicyDisplayNameHeader        = "Display NAME"
	operationPolicyCategoryHeader           = "CATEGORY"
	operationPolicyApplicableFlowsHeaders   = "APPLICABLE FLOWS"
	operationPolicySupportedGatewaysHeaders = "SUPPORTED GATEWAYS"
	defaultOperationPolicyTableFormat       = "table {{.ID}}\t{{.Name}}\t{{.DisplayName}}\t{{.Category}}\t{{.ApplicableFlows}}\t{{.SupportedGateways}}"
)

type operationPolicy struct {
	id                string
	name              string
	displayName       string
	category          string
	applicableFlows   []string
	supportedGateways []string
}

func newOperationPolicyDefinition(a utils.OperationPolicy) *operationPolicy {
	return &operationPolicy{a.Id, a.Name, a.DisplayName, a.Category, a.ApplicableFlows, a.SupportedGateways}
}

func (a operationPolicy) ID() string {
	return a.id
}

func (a operationPolicy) Name() string {
	return a.name
}

func (a operationPolicy) DisplayName() string {
	return a.displayName
}

func (a operationPolicy) Category() string {
	return a.category
}

func (a operationPolicy) ApplicableFlows() []string {
	return a.applicableFlows
}

func (a operationPolicy) SupportedGateways() []string {
	return a.supportedGateways
}

func GetOperationPolicyListFromEnv(accessToken, environment string) (*resty.Response, error) {
	operationPolicyListEndpoint := utils.GetOperationPolicyListEndpointOfEnv(environment, utils.MainConfigFilePath)
	fmt.Println("Endpoint: ", operationPolicyListEndpoint)
	return getOperationPolicyList(accessToken, operationPolicyListEndpoint)
}

func getOperationPolicyList(accessToken string, operationPolicyListEndpoint string) (*resty.Response, error) {
	url := operationPolicyListEndpoint

	utils.Logln(utils.LogPrefixInfo+"GetOperationPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)

	return resp, err
}

// PrintOperationPolicies prints the policy list in a specific format
func PrintOperationPolicies(resp *resty.Response, format string) {
	var operationPolicyList utils.OperationPoliciesList
	err := json.Unmarshal(resp.Body(), &operationPolicyList)
	policies := operationPolicyList.List
	fmt.Println(policies[0])
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshalling response data", err)
	}
	if format == "" {
		format = defaultOperationPolicyTableFormat
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
			"ID":                operationPolicyUUIDHeader,
			"Name":              operationPolicyNameHeader,
			"DisplayName":       operationPolicyDisplayNameHeader,
			"Category":          operationPolicyCategoryHeader,
			"ApplicableFlows":   operationPolicyApplicableFlowsHeaders,
			"SupportedGateways": operationPolicySupportedGatewaysHeaders,
		}
		// execute context
		if err := policyContext.Write(renderer, operationPolicyTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else if format == utils.JsonArrayFormatType {
		utils.ListArtifactsInJsonArrayFormat(policies, utils.ProjectTypePolicy)
	}
}
