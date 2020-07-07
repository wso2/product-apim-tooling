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
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func addAPIProductFromJSONWithoutCleaning(t *testing.T, client *apim.Client, username string, password string, apisList map[string]*apim.API) *apim.APIProduct {
	client.Login(username, password)
	path := "testdata/SampleAPIProduct.json"
	doClean := false
	id := client.AddAPIProductFromJSON(t, path, username, password, apisList, doClean)
	apiProduct := client.GetAPIProduct(id)
	return apiProduct
}

func getAPIProduct(t *testing.T, client *apim.Client, name string, username string, password string) *apim.APIProduct {
	client.Login(username, password)
	apiProductInfo := client.GetAPIProductByName(name)
	return client.GetAPIProduct(apiProductInfo.ID)
}

func getAPIProducts(client *apim.Client, username string, password string) *apim.APIProductList {
	client.Login(username, password)
	return client.GetAPIProducts()
}

func deleteAPIProduct(t *testing.T, client *apim.Client, apiProductID string, username string, password string) {
	time.Sleep(2000 * time.Millisecond)
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
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.ApiProduct.Name, utils.DefaultApiProductVersion)

	if args.ImportApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--import-apis", "--preserve-provider=false")
	} else if args.UpdateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--update-apis", "--preserve-provider=false")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose", "--preserve-provider=false")
	}

	t.Cleanup(func() {
		args.DestAPIM.DeleteAPIProductByName(args.ApiProduct.Name)
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
	output, err := base.Execute(t, "list", "api-products", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func deleteAPIProductByCtl(t *testing.T, args *ApiProductImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "api-product", "-n", args.ApiProduct.Name, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
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

	importAPIProductPreserveProvider(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	time.Sleep(1 * time.Second)

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
	time.Sleep(1 * time.Second)

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

func ValidateAPIProductImport(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Product to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAPIProduct(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	time.Sleep(1 * time.Second)

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.DestAPIM, args.ApiProduct.Name, args.ApiProductProvider.Username, args.ApiProductProvider.Password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqualCrossTenant(t, args.ApiProduct, importedAPIProduct)
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
	time.Sleep(1 * time.Second)

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

func validateAPIProductsList(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List API Products of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, _ := listAPIProducts(t, args)

	apiProductsList := args.SrcAPIM.GetAPIProducts()

	validateListAPIProductsEqual(t, output, apiProductsList)
}

func validateListAPIProductsEqual(t *testing.T, apiProductsListFromCtl string, apiProductsList *apim.APIProductList) {

	for _, apiProduct := range apiProductsList.List {
		// If the output string contains the same API Product ID, then decrement the count
		if strings.Contains(apiProductsListFromCtl, apiProduct.ID) {
			apiProductsList.Count = apiProductsList.Count - 1
		}
	}

	// Count == 0 means that all the API Products from apiProductsList were in apiProductsListFromCtl
	assert.Equal(t, apiProductsList.Count, 0, "API Product lists are not equal")
}

func validateAPIProductDelete(t *testing.T, args *ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Delete an API Product of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	time.Sleep(1 * time.Second)
	apiProductsListBeforeDelete := args.SrcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	apiProductsListAfterDelete := args.SrcAPIM.GetAPIProducts()
	time.Sleep(1 * time.Second)

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

	time.Sleep(1 * time.Second)
	apiProductsListBeforeDelete := args.SrcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	time.Sleep(1 * time.Second)
	apiProductsListAfterDelete := args.SrcAPIM.GetAPIProducts()

	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count, "API Product delete is successful")
}
