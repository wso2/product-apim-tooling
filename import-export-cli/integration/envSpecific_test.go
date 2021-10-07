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
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func TestEnvironmentSpecificParamsEndpoint(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIEndpointParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			// Validate the imported API with the endpoint configs from the params file
			testutils.ValidateAPIsEqualWithEndpointConfigsFromParam(t, api, importedAPI, apiParams)

			// Update the endpoints in the params file and import with update and
			// validate the updated API with the endpoint configs from the updated params file
			testutils.ValidateAPIImportUpdateWithParamsEndpointConfig(t, args, apiParams, user.CtlUser.Username)
		})
	}
}

// Add an API to one environment, export it and re-import it to another environment by setting
// the configs for endpoints using the params file
func TestEnvironmentSpecificParamsEndpointConfigs(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIEndpointConfigsParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Production["url"], importedAPI.GetProductionURL())
			paramConfig := apiParams.Environments[0].Configs.Endpoints.Production["config"].(map[interface{}]interface{})

			apiEndpointConfig := importedAPI.GetProductionConfig()

			for k, v := range paramConfig {
				key := fmt.Sprintf("%v", k)
				value := fmt.Sprintf("%v", v)
				valueInImportedApi := fmt.Sprintf("%v", apiEndpointConfig[key])

				assert.Equal(t, value, valueInImportedApi)
			}

			assert.Equal(t, len(paramConfig), len(apiEndpointConfig))

			apiCopy := apim.CopyAPI(api)
			importedAPICopy := apim.CopyAPI(importedAPI)

			same := "override_with_same_value"
			apiCopy.EndpointConfig = same
			importedAPICopy.EndpointConfig = same

			testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
		})
	}
}

// Add an API to one environment, export it and re-import it to another environment
// by disabling the endpoint security using the params file
func TestEnvironmentSpecificParamsEndpointSecurityFalse(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APISecurityFalseParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			assert.Equal(t, false, importedAPI.GetProductionSecurityConfig()["enabled"])
			assert.Equal(t, false, importedAPI.GetSandboxSecurityConfig()["enabled"])

			api.EndpointConfig.(map[string]interface{})["endpoint_security"] = "override_with_the_same_value"
			importedAPI.EndpointConfig.(map[string]interface{})["endpoint_security"] = "override_with_the_same_value"

			testutils.ValidateAPIsEqual(t, api, importedAPI)
		})
	}
}

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoint security
// (with the security type digest), using the params file
func TestEnvironmentSpecificParamsEndpointSecurityDigest(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APISecurityDigestParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
		})
	}
}

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoint security
// (with the security type basic), using the params file
func TestEnvironmentSpecificParamsEndpointSecurityBasic(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APISecurityBasicParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
		})
	}
}

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoint security
// (with the security type oauth), using the params file
func TestEnvironmentSpecificParamsEndpointSecurityOauth(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APISecurityOauthParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
		})
	}
}

// Import an API with the external params file that has HTTP/REST endpoints without load balancing or failover configs
func TestHttpRestEndpointParamsWithoutLoadBalancingOrFailover(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIHttpRestEndpointWithoutLoadBalancingOrFailoverParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateHttpEndpointWithoutLoadBalancingAndFailover(t, apiParams, api, importedAPI)

		})
	}
}

// Import an API with the external params file that has HTTP/SOAP endpoints without load balancing or failover configs.
func TestHttpSoapEndpointParamsWithoutLoadBalancingOrFailover(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIHttpSoapEndpointWithoutLoadBalancingOrFailoverParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateHttpEndpointWithoutLoadBalancingAndFailover(t, apiParams, api, importedAPI)
		})
	}
}

// Import an API with the external params file that has HTTP/REST endpoints with load balancing configs
func TestHttpRestEndpointParamsWithLoadBalancing(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIHttpRestEndpointWithLoadBalancingParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateHttpEndpointWithLoadBalancing(t, apiParams, api, importedAPI)

		})
	}
}

// Import an API with the external params file that has HTTP/SOAP endpoints with load balancing config
func TestHttpSoapEndpointParamsWithLoadBalancing(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIHttpSoapEndpointWithLoadBalancingParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateHttpEndpointWithLoadBalancing(t, apiParams, api, importedAPI)
		})
	}
}

// Import an API with the external params file that has HTTP/REST endpoints with failover config
func TestHttpRestEndpointParamsWithFailover(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIHttpRestEndpointWithFailoverParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateHttpEndpointWithFailover(t, apiParams, api, importedAPI)

		})
	}
}

// Import an API with the external params file that has HTTP/SOAP endpoints with failover config
func TestHttpSoapEndpointParamsWithFailover(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIHttpSoapEndpointWithFailoverParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateHttpEndpointWithFailover(t, apiParams, api, importedAPI)

		})
	}
}

// Import an API with the external params file that has AWS Lambda Endpoint with role supplied credentials configs
func TestAwsLambdaEndpointParamsWithRoleSupplied(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIAwsRoleSuppliedCredentialsParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateAwsEndpoint(t, apiParams, api, importedAPI)

		})
	}
}

// Import an API with the external params file that has AWS Lambda Endpoint with stored credentials config
func TestAwsLambdaEndpointParamsWithStoredCred(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIAwsEndpointWithStoredCredentialsParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateAwsEndpoint(t, apiParams, api, importedAPI)
		})
	}
}

// Import an API with the external params file that has Dynamic endpoint config
func TestDynamicEndpointParams(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
				ParamsFile:  testutils.APIDynamicEndpointParamsFile,
			}

			testutils.ValidateAPIExport(t, args)

			importedAPI := testutils.GetImportedAPI(t, args)

			apiParams := testutils.ReadParams(t, args.ParamsFile)

			testutils.ValidateDynamicEndpoint(t, apiParams, api, importedAPI)
		})
	}
}

// Export an API from one environment and generate the deployment directory for that. Import it to another environment with the params
// and certificates. Validate the imported API with the used params. Again, re-export it to validate the certs.
func TestExportApiGenDeploymentDirImportSuperTenant(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExport(t, args)

			genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
				Source:      base.ConstructAPIFilePath(testutils.GetEnvAPIExportPath(dev.GetEnvName()), api.Name, api.Version),
				Destination: testutils.GetEnvAPIExportPath(dev.GetEnvName()),
			}

			testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

			// Store the deployment directory path to be provided as the params during import
			args.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination, api.Name, api.Version)
			testutils.ValidateAPIImportExportWithDeploymentDir(t, args)

		})
	}
}

// Export an API Product from one environment and generate the deployment directory for that. Import it to another environment with the params
// and certificates. Validate the imported API Product with the used params. Again, re-export it to validate the certs.
func TestExportApiProductGenDeploymentDirImport(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			// Add the first dependent API to env1
			dependentAPI1 := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, dependentAPI1.ID)

			// Add the second dependent API to env1
			dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, dependentAPI2.ID)
			os.Setenv("DEPENDENTAPI_2", dependentAPI2.Name+"-"+dependentAPI2.Version)

			// Map the real name of the API with the API
			apisList := map[string]*apim.API{
				"PizzaShackAPI":   dependentAPI1,
				"SwaggerPetstore": dependentAPI2,
			}

			// Add the API Product to env1
			apiProduct := testutils.AddAPIProductFromJSON(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, apisList)

			args := &testutils.ApiProductImportExportTestArgs{
				ApiProductProvider: testutils.Credentials{Username: user.ApiPublisher.Username, Password: user.ApiPublisher.Password},
				CtlUser:            testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				ImportApisFlag:     true,
				ApiProduct:         apiProduct,
				SrcAPIM:            dev,
				DestAPIM:           prod,
			}

			testutils.ValidateAPIProductExport(t, args)

			genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
				Source: base.ConstructAPIFilePath(testutils.GetEnvAPIProductExportPath(dev.GetEnvName()), apiProduct.Name,
					utils.DefaultApiProductVersion),
				Destination: testutils.GetEnvAPIProductExportPath(dev.GetEnvName()),
			}

			testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

			// Store the deployment directory path to be provided as the params during import
			args.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination, apiProduct.Name, utils.DefaultApiProductVersion)
			testutils.ValidateAPIProductImportExportWithDeploymentDir(t, args)

			// Validate the dependent API (SwaggerPetstore will be the only one that is in params file of the product)
			testutils.ValidateDependentAPIWithParams(t, dependentAPI2, prod, user.CtlUser.Username, user.CtlUser.Password)
		})
	}
}
