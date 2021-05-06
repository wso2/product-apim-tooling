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

package integration

import (
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

const revisionNumber = "1"

// Undeploy an API revision from one gateway
func TestUndeployAPIRevisionSingleGateway(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API to env
			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			// Add a new gateway env by the respective admin user
			gatewayEnv := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIRevision(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)

			// Deploy the same revision in another gateway
			dev.DeployAPIRevision(t, api.ID, gatewayEnv.Name, gatewayEnv.VHosts[0].Host, revisionId)
			base.WaitForIndexing()

			args := &testutils.UndeployTestArgs{
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				RevisionNo:  revisionNumber,
				GatewayEnvs: []string{gatewayEnv.Name},
			}

			// Validate the undeploy command
			testutils.ValidateAPIUndeploy(t, args, "", revisionId)
		})
	}
}

// Undeploy an API revision from multiple gateways but not all
func TestUndeployAPIRevisionMulitpleGateways(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API to env
			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			// Add two new gateway envs by the respective admin user
			gatewayEnv1 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)
			gatewayEnv2 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIRevision(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)

			// Deploy the same revision in other gateways
			dev.DeployAPIRevision(t, api.ID, gatewayEnv1.Name, gatewayEnv1.VHosts[0].Host, revisionId)
			base.WaitForIndexing()
			dev.DeployAPIRevision(t, api.ID, gatewayEnv2.Name, gatewayEnv2.VHosts[0].Host, revisionId)
			base.WaitForIndexing()

			args := &testutils.UndeployTestArgs{
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				RevisionNo:  revisionNumber,
				GatewayEnvs: []string{gatewayEnv1.Name, gatewayEnv2.Name},
			}

			// Validate the undeploy command
			testutils.ValidateAPIUndeploy(t, args, "", revisionId)
		})
	}
}

// Undeploy an API revision from all the gateways
func TestUndeployAPIRevisionAllGateways(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API to env
			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			// Add two new gateway envs by the respective admin user
			gatewayEnv1 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)
			gatewayEnv2 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIRevision(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)

			// Deploy the same revision in other gateways
			dev.DeployAPIRevision(t, api.ID, gatewayEnv1.Name, gatewayEnv1.VHosts[0].Host, revisionId)
			base.WaitForIndexing()
			dev.DeployAPIRevision(t, api.ID, gatewayEnv2.Name, gatewayEnv2.VHosts[0].Host, revisionId)
			base.WaitForIndexing()

			args := &testutils.UndeployTestArgs{
				CtlUser:    testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:        api,
				SrcAPIM:    dev,
				RevisionNo: revisionNumber,
			}

			// Validate the undeploy command
			testutils.ValidateAPIUndeploy(t, args, "", revisionId)
		})
	}
}

// Undeploy an API revision from a gateway that does not exist
func TestUndeployAPIRevisionFailure(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API to env
			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIRevision(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)

			args := &testutils.UndeployTestArgs{
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				RevisionNo:  revisionNumber,
				GatewayEnvs: []string{base.GenerateRandomString()},
			}

			// Validate the undeploy command failure
			testutils.ValidateAPIUndeployFailure(t, args, "", revisionId)
		})
	}
}

// Undeploy an API Product revision from one gateway
func TestUndeployAPIProductRevisionSingleGateway(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API Product to env
			apiProductArgs := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &user.ApiCreator, &user.ApiPublisher)

			// Add a new gateway env by the respective admin user
			gatewayEnv := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIProductRevision(t, dev, user.ApiPublisher.Username,
				user.ApiPublisher.Password, apiProductArgs.ApiProduct.ID)

			// Deploy the same revision in another gateway
			dev.DeployAPIProductRevision(t, apiProductArgs.ApiProduct.ID, gatewayEnv.Name, gatewayEnv.VHosts[0].Host, revisionId)
			base.WaitForIndexing()

			args := &testutils.UndeployTestArgs{
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				ApiProduct:  apiProductArgs.ApiProduct,
				SrcAPIM:     dev,
				RevisionNo:  revisionNumber,
				GatewayEnvs: []string{gatewayEnv.Name},
			}

			// Validate the undeploy command
			testutils.ValidateAPIProductUndeploy(t, args, "", revisionId)
		})
	}
}

// Undeploy an API Product revision from multiple specified gateways but not all
func TestUndeployAPIProductRevisionMulitpleGateways(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API Product to env
			apiProductArgs := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &user.ApiCreator, &user.ApiPublisher)

			// Add two new gateway envs by the respective admin user
			gatewayEnv1 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)
			gatewayEnv2 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIProductRevision(t, dev, user.ApiPublisher.Username,
				user.ApiPublisher.Password, apiProductArgs.ApiProduct.ID)

			// Deploy the same revision in other gateways
			dev.DeployAPIProductRevision(t, apiProductArgs.ApiProduct.ID, gatewayEnv1.Name, gatewayEnv1.VHosts[0].Host, revisionId)
			base.WaitForIndexing()
			dev.DeployAPIProductRevision(t, apiProductArgs.ApiProduct.ID, gatewayEnv2.Name, gatewayEnv2.VHosts[0].Host, revisionId)
			base.WaitForIndexing()

			args := &testutils.UndeployTestArgs{
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				ApiProduct:  apiProductArgs.ApiProduct,
				SrcAPIM:     dev,
				RevisionNo:  revisionNumber,
				GatewayEnvs: []string{gatewayEnv1.Name, gatewayEnv2.Name},
			}

			// Validate the undeploy command
			testutils.ValidateAPIProductUndeploy(t, args, "", revisionId)
		})
	}
}

// Undeploy an API Product revision from all the gateways
func TestUndeployAPIProductRevisionAllGateways(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API Product to env
			apiProductArgs := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &user.ApiCreator, &user.ApiPublisher)

			// Add two new gateway envs by the respective admin user
			gatewayEnv1 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)
			gatewayEnv2 := testutils.AddGatewayEnv(t, dev, user.Admin.Username, user.Admin.Password)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIProductRevision(t, dev, user.ApiPublisher.Username,
				user.ApiPublisher.Password, apiProductArgs.ApiProduct.ID)

			// Deploy the same revision in other gateways
			dev.DeployAPIProductRevision(t, apiProductArgs.ApiProduct.ID, gatewayEnv1.Name, gatewayEnv1.VHosts[0].Host, revisionId)
			base.WaitForIndexing()
			dev.DeployAPIProductRevision(t, apiProductArgs.ApiProduct.ID, gatewayEnv2.Name, gatewayEnv2.VHosts[0].Host, revisionId)
			base.WaitForIndexing()

			args := &testutils.UndeployTestArgs{
				CtlUser:    testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				ApiProduct: apiProductArgs.ApiProduct,
				SrcAPIM:    dev,
				RevisionNo: revisionNumber,
			}

			// Validate the undeploy command
			testutils.ValidateAPIProductUndeploy(t, args, "", revisionId)
		})
	}
}

// Undeploy an API Product revision from a gateway that does not exist
func TestUndeployAPIProductRevisionFailure(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the API Product to env
			apiProductArgs := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &user.ApiCreator, &user.ApiPublisher)

			// Create and Deploy Revision of the above API to the default gateway
			revisionId := testutils.CreateAndDeployAPIProductRevision(t, dev, user.ApiPublisher.Username,
				user.ApiPublisher.Password, apiProductArgs.ApiProduct.ID)

			args := &testutils.UndeployTestArgs{
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				ApiProduct:  apiProductArgs.ApiProduct,
				SrcAPIM:     dev,
				RevisionNo:  revisionNumber,
				GatewayEnvs: []string{base.GenerateRandomString()},
			}

			// Validate the undeploy command failure
			testutils.ValidateAPIProductUndeployFailure(t, args, "", revisionId)
		})
	}
}
