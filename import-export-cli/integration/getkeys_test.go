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
)

func TestGetKeysNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, api.ID)

	args := &apiGetKeyTestArgs{
		ctlUser: credentials{username: apiPublisher, password: apiPublisherPassword},
		api:     api,
		apim:    dev,
	}

	validateGetKeysFailure(t, args)
}

func TestGetKeysNonAdminTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, api.ID)

	args := &apiGetKeyTestArgs{
		ctlUser: credentials{username: apiPublisher, password: apiPublisherPassword},
		api:     api,
		apim:    dev,
	}

	validateGetKeysFailure(t, args)
}

func TestGetKeysAdminSuperTenantUser(t *testing.T) {
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	publishAPI(dev, apiPublisher, apiPublisherPassword, api.ID)

	args := &apiGetKeyTestArgs{
		ctlUser: credentials{username: adminUser, password: adminPassword},
		api:     api,
		apim:    dev,
	}

	validateGetKeys(t, args)
}

func TestGetKeysAdminTenantUser(t *testing.T) {
	adminUser := superAdminUser + "@" + TENANT1
	adminPassword := superAdminPassword

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	publishAPI(dev, apiPublisher, apiPublisherPassword, api.ID)

	args := &apiGetKeyTestArgs{
		ctlUser: credentials{username: adminUser, password: adminPassword},
		api:     api,
		apim:    dev,
	}

	validateGetKeys(t, args)
}

/*
TODO: Uncomment these when secondary user store automation is supported
func TestGetKeysSecondaryUserStoreAdminSuperTenantUser(t *testing.T) {
	username := "SECOND.COM/super"
	password := "admin"

	provider := "creator"
	name := "PizzaShackAPI"
	version := "1.0.0"
	apiResourceURL := "http://localhost:8280/pizzashack/1.0.0/menu"
	base.SetupEnv(t, devEnv, devApim, devTokenEP)
	base.Login(t, devEnv, username, password)
	result, err := getKeys(t, provider, name, version, devEnv)

	assert.Nil(t, err, "Error while getting key")

	invokeAPI(t, apiResourceURL, base.GetValueOfUniformResponse(result), 200)

	provider = "SECOND.COM/super"
	name = "seconds"
	version = "1.0.0"
	apiResourceURL = "http://localhost:8280/seconds/1.0.0/menu"

	result, err = getKeys(t,provider, name, version, devEnv)

	assert.Nil(t, err, "Error while getting key")

	invokeAPI(t, apiResourceURL, base.GetValueOfUniformResponse(result), 200)

}
*/

func TestGetKeysNonPublishedAPI(t *testing.T) {
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &apiGetKeyTestArgs{
		ctlUser: credentials{username: adminUser, password: adminPassword},
		api:     api,
		apim:    dev,
	}

	validateGetKeysFailure(t, args)
}

func TestGetKeysForAPIProductNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the first dependent API to env1
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiGetKeyTestArgs{
		ctlUser:    credentials{username: apiPublisher, password: apiPublisherPassword},
		apiProduct: apiProduct,
		apim:       dev,
	}

	validateGetKeysFailure(t, args)
}

func TestGetKeysForAPIProductNonAdminTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the first dependent API to env1
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiGetKeyTestArgs{
		ctlUser:    credentials{username: apiPublisher, password: apiPublisherPassword},
		apiProduct: apiProduct,
		apim:       dev,
	}

	validateGetKeysFailure(t, args)
}

func TestGetKeysForAPIProductAdminSuperTenantUser(t *testing.T) {
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the first dependent API to env1
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiGetKeyTestArgs{
		ctlUser:    credentials{username: adminUser, password: adminPassword},
		apiProduct: apiProduct,
		apim:       dev,
	}

	validateGetKeys(t, args)
}

func TestGetKeysForAPIProductAdminTenantUser(t *testing.T) {
	adminUser := superAdminUser + "@" + TENANT1
	adminPassword := superAdminPassword

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the first dependent API to env1
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiGetKeyTestArgs{
		ctlUser:    credentials{username: adminUser, password: adminPassword},
		apiProduct: apiProduct,
		apim:       dev,
	}

	validateGetKeys(t, args)
}
