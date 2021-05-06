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
	"strconv"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

func TestExportApiNonDeloyedRevision(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			apiRevisions := testutils.CreateAndDeploySeriesOfAPIRevisions(t, dev, api, &user.ApiCreator, &user.ApiPublisher)

			revNumber := len(apiRevisions) - 1
			apiRevision := apiRevisions[revNumber]

			// Export middle revision and check if it matches the corresponding API
			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: user.ApiCreator,
				CtlUser:     user.CtlUser,
				Api:         apiRevision,
				SrcAPIM:     dev,
				Revision:    strconv.Itoa(revNumber),
			}

			testutils.ValidateExportedAPIRevisionStructure(t, args)
		})
	}
}

func TestExportApiDeloyedRevision(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			apiRevisions := testutils.CreateAndDeploySeriesOfAPIRevisions(t, dev, api, &user.ApiCreator, &user.ApiPublisher)

			finalDeployedRevision := len(apiRevisions)

			// Export final revision and see if it matches the corresponding API
			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: user.ApiCreator,
				CtlUser:     user.CtlUser,
				Api:         apiRevisions[finalDeployedRevision],
				SrcAPIM:     dev,
				Revision:    strconv.Itoa(finalDeployedRevision),
				IsDeployed:  true,
			}

			testutils.ValidateExportedAPIRevisionStructure(t, args)
		})
	}
}

func TestExportApiWorkingCopy(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			testutils.CreateAndDeploySeriesOfAPIRevisions(t, dev, api, &user.ApiCreator, &user.ApiPublisher)

			// Export final revision and see if it matches the corresponding API
			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: user.ApiCreator,
				CtlUser:     user.CtlUser,
				Api:         api,
				SrcAPIM:     dev,
			}

			testutils.ValidateExportedAPIStructure(t, args)
		})
	}
}

func TestExportApiLatestRevision(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			apiRevisions := testutils.CreateAndDeploySeriesOfAPIRevisions(t, dev, api, &user.ApiCreator, &user.ApiPublisher)

			finalDeployedRevision := len(apiRevisions)

			// Export final revision and see if it matches the corresponding API
			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: user.ApiCreator,
				CtlUser:     user.CtlUser,
				Api:         apiRevisions[finalDeployedRevision],
				SrcAPIM:     dev,
				Revision:    strconv.Itoa(finalDeployedRevision),
				IsDeployed:  true,
				IsLatest:    true,
			}

			testutils.ValidateExportedAPIRevisionStructure(t, args)
		})
	}
}

func TestExportImportApiSameGWEnv(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			revision := testutils.CreateAndDeployAPIRevision(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)

			api = testutils.GetAPIById(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, revision)

			// Export final revision and see if it matches the corresponding API
			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: user.ApiCreator,
				CtlUser:     user.CtlUser,
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				Revision:    "1",
				IsDeployed:  true,
			}

			testutils.ValidateExportedAPIRevisionStructure(t, args)

			testutils.ValidateAPIRevisionExportImport(t, args, testutils.APITypeREST)
		})
		break
	}
}
