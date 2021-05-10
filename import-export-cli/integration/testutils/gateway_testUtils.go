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

package testutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

func AddGatewayEnv(t *testing.T, client *apim.Client, username, password string) *apim.Environment {
	client.Login(username, password)
	env := GenerateSampleGatewayData()
	addedEnv := client.AddGatewayEnv(t, env, username, password)
	return addedEnv
}

func GenerateSampleGatewayData() apim.Environment {
	env := apim.Environment{}
	env.Name = base.GenerateRandomString() + "-Gateway"
	env.DisplayName = env.Name + "-Display Name"
	env.Description = "Gateway environment for testing purposes"
	vhost := apim.VHost{
		Host: env.Name + ".com",
	}
	env.VHosts = []apim.VHost{}
	env.VHosts = append(env.VHosts, vhost)

	return env
}

func undeployAPI(t *testing.T, args *UndeployTestArgs, provider string) (string, error) {
	params := []string{"undeploy", "api", "-n", args.Api.Name, "-v", args.Api.Version,
		"--rev", args.RevisionNo, "-e", args.SrcAPIM.GetEnvName(), "-k", "--verbose"}

	if provider != "" {
		params = append(params, "-r", provider)
	}

	if len(args.GatewayEnvs) > 0 {
		for _, gatewayEnv := range args.GatewayEnvs {
			params = append(params, "-g", "\""+gatewayEnv+"\"")
		}
	}

	output, err := base.Execute(t, params...)

	return output, err
}

func undeployAPIProduct(t *testing.T, args *UndeployTestArgs, provider string) (string, error) {
	params := []string{"undeploy", "api-product", "-n", args.ApiProduct.Name,
		"--rev", args.RevisionNo, "-e", args.SrcAPIM.GetEnvName(), "-k", "--verbose"}

	if provider != "" {
		params = append(params, "-r", provider)
	}

	if len(args.GatewayEnvs) > 0 {
		for _, gatewayEnv := range args.GatewayEnvs {
			params = append(params, "-g", "\""+gatewayEnv+"\"")
		}
	}

	output, err := base.Execute(t, params...)

	return output, err
}

func ValidateAPIUndeploy(t *testing.T, args *UndeployTestArgs, provider, revisionId string) {
	t.Helper()

	deployedAPIRevisionsBeforeUndeploy := args.SrcAPIM.GetAPIRevisions(args.Api.ID, "deployed:true")

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Execute undeploy API command
	result, err := undeployAPI(t, args, provider)

	assert.Nil(t, err, "Should return nil error")
	assert.Contains(t, result, "Revision "+args.RevisionNo+" of API "+args.Api.Name+"_"+
		args.Api.Version+" successfully undeployed")

	deployedAPIRevisionsAfterUndeploy := args.SrcAPIM.GetAPIRevisions(args.Api.ID, "deployed:true")

	if len(args.GatewayEnvs) > 0 {
		// Validate the deployed gateways before and after executing the undeploy command
		ValidateDeployedGateways(t, deployedAPIRevisionsBeforeUndeploy, deployedAPIRevisionsAfterUndeploy, args, revisionId)
	} else {
		// This scenario is that the API revision is undeployed from all the gateways
		assert.Equal(t, len(deployedAPIRevisionsAfterUndeploy.List), 0)
	}
}

func ValidateAPIUndeployFailure(t *testing.T, args *UndeployTestArgs, provider, revisionId string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Execute undeploy API command
	result, err := undeployAPI(t, args, provider)

	assert.NotNil(t, err, "Should not return nil error")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1",
		"Test failed because API was undeployed successfully")
}

func ValidateDeployedGateways(t *testing.T, deployedRevisionsBeforeUndeploy *apim.APIRevisionList,
	deployedRevisionsAfterUndeploy *apim.APIRevisionList, args *UndeployTestArgs, revisionId string) {

	// Validate whether the deployed gateways contain the gateway envs before undeploying
	assert.True(t, containsGatewaysInDeployment(deployedRevisionsBeforeUndeploy, args.GatewayEnvs, revisionId))

	// Validate whether the deployed gateways do not contain the gateway envs after undeploying
	assert.False(t, containsGatewaysInDeployment(deployedRevisionsAfterUndeploy, args.GatewayEnvs, revisionId))
}

func containsGatewaysInDeployment(deployedRevisions *apim.APIRevisionList, gatewayEnvs []string,
	revisionId string) bool {
	for _, deployedRevision := range deployedRevisions.List {
		if strings.EqualFold(deployedRevision.ID, revisionId) {
			containsGatewaysInDeployment := []bool{}

			// Check whether the passed gateway labels to the command are there in the deployed list
			// If so mark it as true, else false
			for _, gatewayEnv := range gatewayEnvs {
				containsGateway := false
				for _, deployment := range deployedRevision.DeploymentInfo {
					if strings.EqualFold(deployment.Name, gatewayEnv) {
						containsGateway = true
						break
					}
				}
				containsGatewaysInDeployment = append(containsGatewaysInDeployment, containsGateway)
			}

			// If any of the gateways are not inside the deployed list, false should be returned
			for _, containsGatewayInDeployment := range containsGatewaysInDeployment {
				if !containsGatewayInDeployment {
					return false
				}
			}
			return true
		}
	}
	return false
}

func ValidateAPIProductUndeploy(t *testing.T, args *UndeployTestArgs, provider, revisionId string) {
	t.Helper()

	deployedAPIProductRevisionsBeforeUndeploy := args.SrcAPIM.GetAPIProductRevisions(args.ApiProduct.ID,
		"deployed:true")

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Execute undeploy API Product command
	result, err := undeployAPIProduct(t, args, provider)

	assert.Nil(t, err, "Should return nil error")
	assert.Contains(t, result, "Revision "+args.RevisionNo+" of API Product "+
		args.ApiProduct.Name+" successfully undeployed")

	deployedAPIProductRevisionsAfterUndeploy := args.SrcAPIM.GetAPIProductRevisions(args.ApiProduct.ID,
		"deployed:true")

	if len(args.GatewayEnvs) > 0 {
		// Validate the deployed gateways before and after executing the undeploy command
		ValidateDeployedGateways(t, deployedAPIProductRevisionsBeforeUndeploy,
			deployedAPIProductRevisionsAfterUndeploy, args, revisionId)
	} else {
		// This scenario is that the API revision is undeployed from all the gateways
		assert.Equal(t, len(deployedAPIProductRevisionsAfterUndeploy.List), 0)
	}
}

func ValidateAPIProductUndeployFailure(t *testing.T, args *UndeployTestArgs, provider, revisionId string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Execute undeploy API command
	result, _ := undeployAPIProduct(t, args, provider)

	assert.Contains(t, result, "400", "Test failed because API Product was undeployed successfully")
}
