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
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/adminservices"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func AddAPIProductFromJSONWithoutCleaning(t *testing.T, client *apim.Client, username string, password string, apisList map[string]*apim.API) *apim.APIProduct {
	client.Login(username, password)
	path := "testdata/SampleAPIProduct.json"
	doClean := false
	id := client.AddAPIProductFromJSON(t, path, username, password, apisList, doClean)
	revision := client.CreateAPIProductRevision(id)
	client.DeployAPIProductRevision(t, id, "", "", revision.ID)
	apiProduct := client.GetAPIProduct(id)
	return apiProduct
}

func AddAPIProductFromJSON(t *testing.T, client *apim.Client, username string, password string, apisList map[string]*apim.API) *apim.APIProduct {
	client.Login(username, password)
	path := "testdata/SampleAPIProduct.json"
	doClean := true
	id := client.AddAPIProductFromJSON(t, path, username, password, apisList, doClean)

	base.WaitForIndexing()

	apiProduct := client.GetAPIProduct(id)
	return apiProduct
}

func CreateAndDeployAPIProductRevision(t *testing.T, client *apim.Client, username, password, apiProductID string) string {
	client.Login(username, password)
	revision := client.CreateAPIProductRevision(apiProductID)
	client.DeployAPIProductRevision(t, apiProductID, "", "", revision.ID)
	return revision.ID
}

func getAPIProduct(t *testing.T, client *apim.Client, name string, username string, password string) *apim.APIProduct {
	if username == adminservices.DevopsUsername {
		client.Login(adminservices.AdminUsername, adminservices.AdminPassword)
	} else if username == adminservices.DevopsUsername+"@"+adminservices.Tenant1 {
		client.Login(adminservices.AdminUsername+"@"+adminservices.Tenant1, adminservices.AdminPassword)
	} else {
		client.Login(username, password)
	}
	apiProductInfo := client.GetAPIProductByName(name)
	return client.GetAPIProduct(apiProductInfo.ID)
}

func getAPIProducts(client *apim.Client, username string, password string) *apim.APIProductList {
	client.Login(username, password)
	return client.GetAPIProducts()
}

func deleteAPIProduct(t *testing.T, client *apim.Client, apiProductID string, username string, password string) {
	base.WaitForIndexing()
	client.Login(username, password)
	client.DeleteAPIProduct(apiProductID)
}

func getResourceURLForAPIProduct(apim *apim.Client, apiProduct *apim.APIProduct) string {
	port := 8280 + apim.GetPortOffset()
	return "http://" + apim.GetHost() + ":" + strconv.Itoa(port) + apiProduct.Context + "/menu"
}

func getEnvAPIProductExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedApiProductsDirName, envName)
}

func flagAPIsAddedViaProductImportForRemoval(t *testing.T, client *apim.Client, apiProviders *map[string]Credentials) {
	if len(*apiProviders) > 0 {
		t.Cleanup(func() {
			for name, credentials := range *apiProviders {
				username, password := apim.RetrieveAdminCredentialsInsteadCreator(credentials.Username, credentials.Password)
				client.Login(username, password)
				err := client.DeleteAPIByName(name)

				if err != nil {
					t.Fatal(err)
				}
				base.WaitForIndexing()
			}
		})
	}
}

func exportAPIProduct(t *testing.T, name string, version string, env string) (string, error) {
	output, err := base.Execute(t, "export", "api-product", "-n", name, "-e", env, "-k", "--verbose")

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, getEnvAPIProductExportPath(env), name, version)
	})

	return output, err
}

func importAPIProductPreserveProvider(t *testing.T, args *ApiProductImportExportTestArgs) (string, error) {
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.ApiProduct.Name, utils.DefaultApiProductVersion)

	if args.ImportApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--import-apis")
	} else if args.UpdateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--update-apis")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose")
	}

	t.Cleanup(func() {
		args.DestAPIM.DeleteAPIProductByName(args.ApiProduct.Name)
		base.WaitForIndexing()
	})

	return output, err
}

func importUpdateAPIProductPreserveProvider(t *testing.T, args *ApiProductImportExportTestArgs) (string, error) {
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.ApiProduct.Name, utils.DefaultApiProductVersion)

	if args.UpdateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--update-apis")
	} else if args.UpdateApiProductFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--update-api-product")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose")
	}

	return output, err
}

func importAPIProduct(t *testing.T, args *ApiProductImportExportTestArgs) (string, error) {

	fileName := base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.ApiProduct.Name, utils.DefaultApiProductVersion)

	params := []string{"import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--preserve-provider=false"}

	if args.ImportApisFlag {
		params = append(params, "--import-apis")
	}
	if args.UpdateApisFlag {
		params = append(params, "--update-apis")
	}
	if args.ParamsFile != "" {
		params = append(params, "--params", args.ParamsFile)
	}

	output, err := base.Execute(t, params...)

	t.Cleanup(func() {
		args.DestAPIM.DeleteAPIProductByName(args.ApiProduct.Name)
		base.WaitForIndexing()
	})

	return output, err
}

func importUpdateAPIProduct(t *testing.T, args *ApiProductImportExportTestArgs) (string, error) {
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.ApiProduct.Name, utils.DefaultApiProductVersion)

	if args.UpdateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--update-apis", "--preserve-provider=false")
	} else if args.UpdateApiProductFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--update-api-product", "--preserve-provider=false")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--preserve-provider=false")
	}

	return output, err
}

func listAPIProducts(t *testing.T, args *ApiProductImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "get", "api-products", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func deleteAPIProductByCtl(t *testing.T, args *ApiProductImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "api-product", "-n", args.ApiProduct.Name, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

// AddAPIProductWithTwoDependentAPIs : Helper function for adding and API Product along with two dependent APIs to an env
//
func AddAPIProductWithTwoDependentAPIs(t *testing.T, client *apim.Client, apiCreator *Credentials, apiPublisher *Credentials) *ApiProductImportExportTestArgs {
	t.Helper()

	// Add the first dependent API to env1
	dependentAPI1 := AddAPI(t, client, apiCreator.Username, apiCreator.Password)
	PublishAPI(client, apiPublisher.Username, apiPublisher.Password, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := AddAPIFromOpenAPIDefinition(t, client, apiCreator.Username, apiCreator.Password)
	PublishAPI(client, apiPublisher.Username, apiPublisher.Password, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := AddAPIProductFromJSON(t, client, apiPublisher.Username, apiPublisher.Password, apisList)

	apiProviders := map[string]Credentials{}
	apiProviders[dependentAPI1.Name] = *apiCreator
	apiProviders[dependentAPI2.Name] = *apiCreator

	return &ApiProductImportExportTestArgs{
		ApiProviders:       apiProviders,
		ApiProductProvider: *apiPublisher,
		ApiProduct:         apiProduct,
		SrcAPIM:            client,
	}
}

func ValidateAPIProductExportFailure(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting API Product from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPIProduct(t, args.ApiProduct.Name, utils.DefaultApiProductVersion, args.SrcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsAPIArchiveExists(t, getEnvAPIProductExportPath(args.SrcAPIM.GetEnvName()),
		args.ApiProduct.Name, utils.DefaultApiProductVersion))
}

func ValidateAPIProductExportImportPreserveProvider(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export API Product from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPIProduct(t, args.ApiProduct.Name, utils.DefaultApiProductVersion, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIProductExportPath(args.SrcAPIM.GetEnvName()),
		args.ApiProduct.Name, utils.DefaultApiProductVersion))

	// Import API Product to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// If any APIs were added via the API Product import, flag them for removal during cleanup
	flagAPIsAddedViaProductImportForRemoval(t, args.DestAPIM, &args.ApiProviders)

	importAPIProductPreserveProvider(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	base.WaitForIndexing()

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.DestAPIM, args.ApiProduct.Name, args.ApiProductProvider.Username, args.ApiProductProvider.Password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqual(t, args.ApiProduct, importedAPIProduct)
}

func ValidateAPIProductImportUpdatePreserveProvider(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// This is used when you have previously imported an API Product (with preserving the provider) and validated it.
	// So when doing the cleaning you do not need to clean twice. For that, importUpdateAPIProductPreserveProvider will not be doing cleaning again.
	importUpdateAPIProductPreserveProvider(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	base.WaitForIndexing()

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.DestAPIM, args.ApiProduct.Name, args.ApiProductProvider.Username, args.ApiProductProvider.Password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqual(t, args.ApiProduct, importedAPIProduct)
}

func ValidateAPIProductExport(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export API Product from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPIProduct(t, args.ApiProduct.Name, utils.DefaultApiProductVersion, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIProductExportPath(args.SrcAPIM.GetEnvName()),
		args.ApiProduct.Name, utils.DefaultApiProductVersion))
}

func ValidateAPIProductImport(t *testing.T, args *ApiProductImportExportTestArgs, skipvalidateAPIProductsEqual bool) *apim.APIProduct {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Product to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Since --preserve-provider=false, the API Provider of the APIs and API Products that will be created
	// will be considered to be the apictl user.
	for key := range args.ApiProviders {
		args.ApiProviders[key] = args.CtlUser
	}

	args.ApiProductProvider = args.CtlUser

	// If any APIs were added via the API Product import, flag them for removal during cleanup
	flagAPIsAddedViaProductImportForRemoval(t, args.DestAPIM, &args.ApiProviders)

	importAPIProduct(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	base.WaitForIndexing()

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.DestAPIM, args.ApiProduct.Name, args.ApiProductProvider.Username, args.ApiProductProvider.Password)

	if !skipvalidateAPIProductsEqual {
		// Validate env 1 and env 2 API Products are equal
		validateAPIProductsEqualCrossTenant(t, args.ApiProduct, importedAPIProduct)
	}

	return importedAPIProduct
}

func ValidateAPIProductImportUpdate(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Product to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// This is used when you have previously imported an API Product and validated it.
	// So when doing the cleaning you do not need to clean twice. For that, importUpdateAPIProduct will not be doing cleaning again.
	importUpdateAPIProduct(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	base.WaitForIndexing()

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.DestAPIM, args.ApiProduct.Name, args.ApiProductProvider.Username, args.ApiProductProvider.Password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqualCrossTenant(t, args.ApiProduct, importedAPIProduct)
}

func validateAPIProductsEqual(t *testing.T, apiProduct1 *apim.APIProduct, apiProduct2 *apim.APIProduct) {
	t.Helper()

	apiProduct1Copy := apim.CopyAPIProduct(apiProduct1)
	apiProduct2Copy := apim.CopyAPIProduct(apiProduct2)

	same := "override_with_same_value"
	// Since the API Products are from too different envs, their respective ID will defer.
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	apiProduct1Copy.ID = same
	apiProduct2Copy.ID = same

	apiProduct1Copy.CreatedTime = same
	apiProduct2Copy.CreatedTime = same

	apiProduct1Copy.LastUpdatedTime = same
	apiProduct2Copy.LastUpdatedTime = same

	// Check the validity of the operations in the APIs array of the two API Products
	err := validateOperations(&apiProduct1Copy, &apiProduct2Copy)
	if err != nil {
		t.Fatal(err)
	}

	// Sort member collections to make equality check possible
	apim.SortAPIProductMembers(&apiProduct1Copy)
	apim.SortAPIProductMembers(&apiProduct2Copy)

	assert.Equal(t, apiProduct1Copy, apiProduct2Copy, "API Product obejcts are not equal")
}

func validateAPIProductsEqualCrossTenant(t *testing.T, apiProduct1 *apim.APIProduct, apiProduct2 *apim.APIProduct) {
	t.Helper()

	apiProduct1Copy := apim.CopyAPIProduct(apiProduct1)
	apiProduct2Copy := apim.CopyAPIProduct(apiProduct2)

	same := "override_with_same_value"
	// Since the API Products are from too different envs, their respective ID will defer.
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	apiProduct1Copy.ID = same
	apiProduct2Copy.ID = same

	apiProduct1Copy.CreatedTime = same
	apiProduct2Copy.CreatedTime = same

	apiProduct1Copy.LastUpdatedTime = same
	apiProduct2Copy.LastUpdatedTime = same

	// The contexts and providers will differ since this is a cross tenant import
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	apiProduct1Copy.Context = same
	apiProduct2Copy.Context = same

	apiProduct1Copy.Provider = same
	apiProduct2Copy.Provider = same

	// Check the validity of the operations in the APIs array of the two API Products
	err := validateOperations(&apiProduct1Copy, &apiProduct2Copy)
	if err != nil {
		t.Fatal(err)
	}

	// Sort member collections to make equality check possible
	apim.SortAPIProductMembers(&apiProduct1Copy)
	apim.SortAPIProductMembers(&apiProduct2Copy)

	assert.Equal(t, apiProduct1Copy, apiProduct2Copy, "API Product obejcts are not equal")
}

func validateOperations(apiProduct1Copy, apiProduct2Copy *apim.APIProduct) error {

	// To store the validity of each dependent API operation
	var isOperationsValid []bool

	// Iterate thorugh the APIs array of API Product 1
	for index, apiInProduct1 := range apiProduct1Copy.APIs {
		// Iterate thorugh the APIs array of API Product 2
		for _, apiInProduct2 := range apiProduct2Copy.APIs {
			// If the name of the APIs in the two API Products are same, those should be compared
			if apiInProduct1.(map[string]interface{})["name"] == apiInProduct2.(map[string]interface{})["name"] {

				// Convert the maps to APIOperations array structs (so that the structs can be compared easily)
				var operationsList []apim.APIOperations
				operationsInApiInProduct1, _ := json.Marshal(apiInProduct1.(map[string]interface{})["operations"])
				err := json.Unmarshal(operationsInApiInProduct1, &operationsList)
				if err != nil {
					return err
				}
				operationsInApiInProduct2, _ := json.Marshal(apiInProduct2.(map[string]interface{})["operations"])
				err = json.Unmarshal(operationsInApiInProduct2, &operationsList)
				if err != nil {
					return err
				}

				// Compare the two APIOperations array structs, whether they are equal
				if cmp.Equal(operationsInApiInProduct1, operationsInApiInProduct2) {
					// If the operations are equal, it is valid
					isOperationsValid = append(isOperationsValid, true)

					// Since the apiIds of the dependent APIs in the environments differ, those will be assigned integer values (index value in the loop)
					// Same APIs in both apiInProduct2 and apiInProduct2 will be assigned same integer values based on the order they are in APIs array
					apiInProduct2.(map[string]interface{})["apiId"] = index
					apiInProduct1.(map[string]interface{})["apiId"] = index
					break
				}
				// If the operations are not equal, it is not valid
				isOperationsValid = append(isOperationsValid, false)
			}
		}
	}

	// To store the overall result of whether the operations are valid.
	isAllOperationsValid := true
	for _, value := range isOperationsValid {
		// If any of the value in isOperationsValid array is false, the overall result should be false
		if value == false {
			isAllOperationsValid = false
		}
	}

	if isAllOperationsValid {
		// If all the operations in both of the APIs arrays are equal, make those two APIs arrays equals
		apiProduct2Copy.APIs = apiProduct1Copy.APIs
	}
	return nil
}

func ValidateAPIProductsList(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List API Products of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIProducts(t, args)

	apiProductsList := args.SrcAPIM.GetAPIProducts()

	ValidateListAPIProductsEqual(t, output, apiProductsList)
}

func ValidateListAPIProductsEqual(t *testing.T, apiProductsListFromCtl string, apiProductsList *apim.APIProductList) {
	unmatchedCount := apiProductsList.Count
	for _, apiProduct := range apiProductsList.List {
		// If the output string contains the same API Product ID, then decrement the count
		assert.Truef(t, strings.Contains(apiProductsListFromCtl, apiProduct.ID), "apiProductsListFromCtl: "+apiProductsListFromCtl+
			" , does not contain apiProduct.ID: "+apiProduct.ID)
		unmatchedCount--
	}

	// Count == 0 means that all the API Products from apiProductsList were in apiProductsListFromCtl
	assert.Equal(t, 0, unmatchedCount, "API Product lists are not equal")
}

func ValidateAPIProductDelete(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Delete an API Product of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()
	apiProductsListBeforeDelete := args.SrcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	apiProductsListAfterDelete := args.SrcAPIM.GetAPIProducts()
	base.WaitForIndexing()

	// Validate whether the expected number of API Product count is there
	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count+1, "Expected number of API Products not deleted")

	// Validate that the delete is a success
	validateAPIProductIsDeleted(t, args.ApiProduct, apiProductsListAfterDelete)
}

func validateAPIProductIsDeleted(t *testing.T, apiProduct *apim.APIProduct, apiProductsListAfterDelete *apim.APIProductList) {
	for _, existingAPIProduct := range apiProductsListAfterDelete.List {
		assert.NotEqual(t, existingAPIProduct.ID, apiProduct.ID, "API Product delete is not successful")
	}
}

func ValidateAPIProductDeleteFailure(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Delete an API Product of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()
	apiProductsListBeforeDelete := args.SrcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	base.WaitForIndexing()
	apiProductsListAfterDelete := args.SrcAPIM.GetAPIProducts()

	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count, "API Product delete is successful")
}

func ValidateAPIProductDeleteFailureWithExistingEnv(t *testing.T, args *ApiProductImportExportTestArgs) {

	apiProductsListBeforeDelete := args.SrcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	base.WaitForIndexing()
	apiProductsListAfterDelete := args.SrcAPIM.GetAPIProducts()

	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count, "API Product delete is successful")

	//Remove subscription and remove Api-Product for cleanup
	t.Cleanup(func() {
		UnsubscribeAPI(args.SrcAPIM, args.CtlUser.Username, args.CtlUser.Password, args.ApiProduct.ID)
		deleteAPIProductByCtl(t, args)
	})
}
