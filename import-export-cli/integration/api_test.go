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

import "testing"

func TestExportApiNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: apiCreator, password: apiCreatorPassword},
		ctlUser:     credentials{username: apiPublisher, password: apiPublisherPassword},
		api:         api,
		srcAPIM:     dev,
	}

	validateAPIExportFailure(t, args)
}

func TestExportImportApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: apiCreator, password: apiCreatorPassword},
		ctlUser:     credentials{username: adminUsername, password: adminPassword},
		api:         api,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAPIExportImport(t, args)
}
