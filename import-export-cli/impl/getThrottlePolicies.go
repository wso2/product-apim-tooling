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
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io"
	"os"
	"strconv"
	"text/template"
)

const (
	policyUUIDHeader                 = "UUID"
	policyNameHeader                 = "NAME"
	isDeployedHeader                 = "IS DEPLOYED"
	policyTypeHeader                 = "TYPE"
	defaultThrottlePolicyTableFormat = "table {{.UUID}}\t{{.Name}}\t{{.Deployed}}\t{{.PolicyType}}"
)

type policy struct {
	PolicyId    int    `json:"policyId"`
	Uuid        string `json:"uuid"`
	PolicyName  string `json:"policyName"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IsDeployed  bool   `json:"isDeployed"`
	Type        string `json:"type"`
}

func newPolicyDefinition(a utils.Policy) *policy {
	return &policy{a.PolicyId, a.Uuid, a.PolicyName, a.DisplayName, a.Description,
		a.IsDeployed, a.Type}
}

func (a policy) UUID() string {
	return a.Uuid
}

func (a policy) Name() string {
	return a.PolicyName
}

func (a policy) Deployed() string {
	return strconv.FormatBool(a.IsDeployed)
}

func (a policy) PolicyType() string {
	return a.Type
}

func GETThrottlePolicyListFromEnv(accessToken, environment, query string) (*resty.Response, error) {
	adminEndpoint := utils.GetAdminEndpointOfEnv(environment, utils.MainConfigFilePath)
	throttlePolicyListEndpoint := adminEndpoint + "/throttling/policies/search"

	return GetThrottlePolicyList(accessToken, throttlePolicyListEndpoint, query)
}

func GetThrottlePolicyList(accessToken string, throttlePolicyListEndpoint string, query string) (*resty.Response, error) {
	url := throttlePolicyListEndpoint
	queryParamString := "query=" + query
	utils.Logln(utils.LogPrefixInfo+"ExportThrottlingPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	if query == "" {
		resp, err := utils.InvokeGETRequest(url, headers)
		return resp, err
	} else {
		resp, err := utils.InvokeGETRequestWithQueryParamsString(url, queryParamString, headers)
		return resp, err
	}
}

func PrintThrottlePolicies(resp *resty.Response, format string) {
	var policyList utils.PolicyList
	err := json.Unmarshal(resp.Body(), &policyList)
	policies := policyList.List
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshalling response data", err)
	}

	if format == "" {
		format = defaultThrottlePolicyTableFormat
	} else if format == utils.JsonArrayFormatType {
		utils.ListArtifactsInJsonArrayFormat(policies, utils.ProjectTypePolicy)
		return
	}

	// create policy context with standard output
	policyContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, policy := range policies {
			if err := t.Execute(w, newPolicyDefinition(policy)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}
	// headers for table
	ThrottlePolicyTableHeaders := map[string]string{
		"UUID":       policyUUIDHeader,
		"Name":       policyNameHeader,
		"Deployed":   isDeployedHeader,
		"PolicyType": policyTypeHeader,
	}
	// execute context
	if err := policyContext.Write(renderer, ThrottlePolicyTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
