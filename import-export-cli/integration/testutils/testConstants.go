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

//Environment management related test constants
const ApictlInitMessage = "apictl is a Command Line Tool for Importing and Exporting APIs and Applications between " +
	"different environments of WSO2 API Manager"
const CustomTestExportDirectory = "CustomExportDirectory"
const CustomDirectoryAtInit = "CustomExportDirectoryAtInit"
const EnvVariableNameOfCustomCustomDirectoryAtInit = "APICTL_CONFIG_DIR"

const TestSwagger2DefinitionPath = "testdata/swagger2Definition.yaml"
const TestOpenAPI3DefinitionPath = "testdata/openAPI3Definition.yaml"
const TestOpenAPI3DefinitionWithoutEndpointsPath = "testdata/openAPI3DefinitionWithoutEndpoints.yaml"
const TestOpenAPISpecificationURL = "https://petstore.swagger.io/v2/swagger.json"
const TestMigrationDirectorySuffix = "/migration"

const DefaultApictlTestAppName = "default-apictl-app"

//Export test cases
const DevFirstDefaultAPIName = "SwaggerPetstoreNew"
const DevFirstSwagger2APIName = "PizzaShackAPI"
const OpenAPI3DefinitionWithoutEndpointsAPIName = "PizzaShackAPI"
const DevFirstDefaultAPIVersion = "1.0.0"
const APIYamlFilePath = "api.yaml"
const DeploymentEnvYamlFilePath = "deployment_environments.yaml"
const SampleAPIYamlFilePath = "testdata/sample-api.yaml"
const SampleRevisionedAPIYamlFilePath = "testdata/sample-revisioned-api.yaml"
const TestArtifact1Path = "testdata/TestArtifactDirectory/ArtifactSet1"
const TestCase1DocName = "/Docs/Doc1"
const TestCase1DocPath = TestArtifact1Path + "/Doc1/testDoc.pdf"
const TestCase1DestPathSuffix = TestCase1DocName + "/testDoc.pdf"
const TestCase1DocMetaDataPath = TestArtifact1Path + "/Doc1/document.yaml"
const TestCase1DestMetaDataPathSuffix = TestCase1DocName + "/document.yaml"
const TestArtifact2Path = "testdata/TestArtifactDirectory/ArtifactSet2"
const TestCase2DocName = "/Docs/Doc2"
const TestCase2DocPath = TestArtifact2Path + "/Doc2/mockPdf.pdf"
const TestCase2DestPathSuffix = TestCase2DocName + "/mockPdf.pdf"
const TestCase2DestPathSuffixForUpdate = TestCase2DocName + "/testDoc.pdf"
const TestCase2DocMetaDataPath = TestArtifact2Path + "/Doc2/document.yaml"
const TestCase2DestMetaDataPathSuffix = TestCase2DocName + "/document.yaml"
const TestCase2PngPath = TestArtifact2Path + "/icon.png"
const TestCase2DestPngPathSuffix = "/Image/icon.png"
const TestCase2JpegPath = TestArtifact2Path + "/icon.jpeg"
const TestCase2DestJpegPathSuffix = "/Image/icon.jpeg"
const TestArtifact3Path = "testdata/TestArtifactDirectory/ArtifactSet3"
const TestDefaultExtractedFileName = "/SwaggerPetstoreNew-1.0.0"
const CustomDynamicInSequenceName = "custom-header-in.xml"

//Environment specific testcase constants

// EnvParamsFilesDir : Directory that stored environment specific test resources
const EnvParamsFilesDir = "testdata/EnvParamsFiles"

// APIEndpointParamsFile : Endpoint URL api_params.yaml
const APIEndpointParamsFile = EnvParamsFilesDir + "/api_params_endpoint.yaml"

// APIEndpointConfigsParamsFile : Endpoint URL and Retry Timeout api_params.yaml
const APIEndpointConfigsParamsFile = EnvParamsFilesDir + "/api_params_endpoint_configs.yaml"

// APISecurityFalseParamsFile : Security false api_params.yaml
const APISecurityFalseParamsFile = EnvParamsFilesDir + "/api_params_security_false.yaml"

// APISecurityDigestParamsFile : Security Digest api_params.yaml
const APISecurityDigestParamsFile = EnvParamsFilesDir + "/api_params_security_digest.yaml"

// APISecurityBasicParamsFile : Security Basic api_params.yaml
const APISecurityBasicParamsFile = EnvParamsFilesDir + "/api_params_security_basic.yaml"

// APISecurityOauthParamsFile : Security Basic api_params.yaml
const APISecurityOauthParamsFile = EnvParamsFilesDir + "/api_params_security_oauth.yaml"

// APIFullParamsFile : Full api_params.yaml
const APIFullParamsFile = EnvParamsFilesDir + "/api_params_full.yaml"

// APIDynamicDataParamsFile : api_params.yaml with dynamic data
const APIDynamicDataParamsFile = EnvParamsFilesDir + "/api_params_dynamic_data.yaml"

// APIProductFullParamsFile : Full api_product_params.yaml
const APIProductFullParamsFile = EnvParamsFilesDir + "/api_product_params_full.yaml"

// CertificatesDirectoryPath : Directory path for the dummy certificates
const CertificatesDirectoryPath = "testdata/TestArtifactDirectory/certificates"

// UnlimitedPolicy : Unlimited Throttle Policy
const UnlimitedPolicy = "Unlimited"

// TenPerMinAppThrottlingPolicy : 10 per min application throttling policy
const TenPerMinAppThrottlingPolicy = "10PerMin"

// APIHttpRestEndpointWithoutLoadBalancingOrFailoverParamsFile : HTTP/REST Endpoint without Loadbalancing and Failover URLs in api_params.yaml
const APIHttpRestEndpointWithoutLoadBalancingOrFailoverParamsFile = EnvParamsFilesDir + "/api_params_http_rest_endpoint_without_lb_or_failover.yaml"

// APIHttpSoapEndpointWithoutLoadBalancingOrFailoverParamsFile : HTTP/SOAP Endpoint without Loadbalancing and Failover URLs in api_params.yaml
const APIHttpSoapEndpointWithoutLoadBalancingOrFailoverParamsFile = EnvParamsFilesDir + "/api_params_http_soap_endpoint_without_lb_or_failover.yaml"

// APIHttpRestEndpointWithLoadBalancingParamsFile : HTTP/REST Endpoint with Loadbalancing URLs in api_params.yaml
const APIHttpRestEndpointWithLoadBalancingParamsFile = EnvParamsFilesDir + "/api_params_http_rest_endpoint_with_load_balancing.yaml"

// APIHttpSoapEndpointWithLoadBalancingParamsFile : HTTP/SOAP Endpoint with Loadbalancing URLs in api_params.yaml
const APIHttpSoapEndpointWithLoadBalancingParamsFile = EnvParamsFilesDir + "/api_params_http_soap_endpoint_with_load_balancing.yaml"

// APIHttpRestEndpointWithFailoverParamsFile : HTTP/REST Endpoint with Failover URLs in api_params.yaml
const APIHttpRestEndpointWithFailoverParamsFile = EnvParamsFilesDir + "/api_params_http_rest_endpoint_with_failover.yaml"

// APIHttpSoapEndpointWithFailoverParamsFile : HTTP/SOAP Endpoint with Failover URLs in api_params.yaml
const APIHttpSoapEndpointWithFailoverParamsFile = EnvParamsFilesDir + "/api_params_http_soap_endpoint_with_failover.yaml"

// APIAwsRoleSuppliedCredentialsParamsFile : AWS Lambda Endpoint with role supplied credentials in api_params.yaml
const APIAwsRoleSuppliedCredentialsParamsFile = EnvParamsFilesDir + "/api_params_aws_lambda_endpoint_with_role_supplied_cred.yaml"

// APIAwsEndpointWithStoredCredentialsParamsFile : AWS Lambda Endpoint with stored credentials in api_params.yaml
const APIAwsEndpointWithStoredCredentialsParamsFile = EnvParamsFilesDir + "/api_params_aws_lambda_endpoint_with_stored_creds.yaml"

// APIDynamicEndpointParamsFile : Dynamic Endpoint with stored credentials in api_params.yaml
const APIDynamicEndpointParamsFile = EnvParamsFilesDir + "/api_params_dynamic_endpoint.yaml"

// API types
const APITypeREST = "HTTP"
const APITypeSoap = "SOAP"
const APITypeSoapToRest = "SOAPTOREST"
const APITypeGraphQL = "GraphQL"
const APITypeWebScoket = "WS"
const APITypeWebSub = "WEBSUB"
const APITypeSSE = "SSE"

// REST API Endpoint URL
const RESTAPIEndpoint = "https://petstore.swagger.io"

// SOAP API Endpoint URL
const SoapEndpointURL = "http://ws.cdyne.com/phoneverify/phoneverify.asmx"

// GraphQL API Endpoint URL
const GraphQLEndpoint = "http://www.mocky.io/v2/5ea84def2d0000a52d3a3ecd"

// Web Socket API Endpoint URL
const WebSocketEndpoint = "ws://echo.websocket.org:80"

// Search query types
const CustomAPIName = "Customized_API"
const CustomAPIVersion = "2.3.4"
const CustomAPIContext = "/custom"

// Endpoint security related constants
const EndpointSecurityTypeOAuth = "OAUTH"
const PasswordGrantType = "PASSWORD"
const EndpointSecurityTypeDigest = "DIGEST"
const EndpointSecurityTypeBasic = "BASIC"
