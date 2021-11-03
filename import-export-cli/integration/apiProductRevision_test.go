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

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

func TestExportInvalidApiProductRevision(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add the first dependent API to env1
			dependentAPI1 := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, dependentAPI1.ID)

			// Add the second dependent API to env1
			dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, dependentAPI2.ID)

			// Map the real name of the API with the API
			apisList := map[string]*apim.API{
				"PizzaShackAPI":   dependentAPI1,
				"SwaggerPetstore": dependentAPI2,
			}

			// Add the API Product to env1
			apiProduct := testutils.AddAPIProductFromJSON(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, apisList)

			testutils.CreateAndDeployAPIProductRevision(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, apiProduct.ID)

			// Export an invalid revision
			args := &testutils.ApiProductImportExportTestArgs{
				CtlUser:    user.CtlUser,
				ApiProduct: apiProduct,
				SrcAPIM:    dev,
				Revision:   "100",
			}

			testutils.ValidateExportedAPIProductRevisionFailure(t, args)
		})
	}
}
