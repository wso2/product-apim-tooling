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

package utils

import (
	"os"
	"os/user"
	"path/filepath"
)

const ProjectName = "apictl"

var MICmd = "apictl"

func GetMICmdName() string {
	if MICmd == "mi" {
		return ""
	}
	envProjName := os.Getenv("MICmd")
	if envProjName == "mi" {
		MICmd = envProjName
		return ""
	}
	return MICmd
}

// File Names and Paths
var CurrentDir, _ = os.Getwd()

const ConfigDirName = ".wso2apictl"

const MIConfigDirName = ".wso2mi"

var HomeDirectory = getConfigHomeDir()

func getConfigHomeDir() string {
	value := os.Getenv("APICTL_CONFIG_DIR")
	if len(value) == 0 {
		value, err := os.UserHomeDir()
		if len(value) == 0 || err != nil {
			current, err := user.Current()
			if err != nil || current == nil {
				HandleErrorAndExit("User's HOME folder location couldn't be identified", nil)
				return ""
			}
			return current.HomeDir
		}
		return value
	}
	return value
}

func GetConfigDirPath() string {
	if MICmd == "mi" {
		return filepath.Join(HomeDirectory, MIConfigDirName)
	}
	return filepath.Join(HomeDirectory, ConfigDirName)
}

func getLocalCredentialsDirectoryName() string {
	if MICmd == "mi" {
		return filepath.Join(HomeDirectory, MILocalCredentialsDirectoryName)
	}
	return filepath.Join(HomeDirectory, LocalCredentialsDirectoryName)
}

var ConfigDirPath = filepath.Join(HomeDirectory, ConfigDirName)

const LocalCredentialsDirectoryName = ".wso2apictl.local"
const MILocalCredentialsDirectoryName = ".wso2mi.local"
const EnvKeysAllFileName = "env_keys_all.yaml"
const MainConfigFileName = "main_config.yaml"
const SampleMainConfigFileName = "main_config.yaml.sample"
const DefaultAPISpecFileName = "default_api.yaml"

var LocalCredentialsDirectoryPath = getLocalCredentialsDirectoryName()
var EnvKeysAllFilePath = filepath.Join(LocalCredentialsDirectoryPath, EnvKeysAllFileName)
var MainConfigFilePath = filepath.Join(GetConfigDirPath(), MainConfigFileName)
var SampleMainConfigFilePath = filepath.Join(ConfigDirPath, SampleMainConfigFileName)
var DefaultAPISpecFilePath = filepath.Join(ConfigDirPath, DefaultAPISpecFileName)

const DefaultExportDirName = "exported"
const ExportedApisDirName = "apis"
const ExportedMCPServersDirName = "mcp-servers"
const ExportedPoliciesDirName = "policies"
const ExportedThrottlePoliciesDirName = "rate-limiting"
const ExportedAPIPoliciesDirName = "api"
const ExportedApiProductsDirName = "api-products"
const ExportedAppsDirName = "apps"
const ExportedMigrationArtifactsDirName = "migration"
const CertificatesDirName = "certs"

const (
	InitProjectDefinitions              = "Definitions"
	InitProjectDefinitionsSwagger       = InitProjectDefinitions + string(os.PathSeparator) + "swagger.yaml"
	InitProjectDefinitionsGraphQLSchema = InitProjectDefinitions + string(os.PathSeparator) + "schema.graphql"
	InitProjectDefinitionsAsyncAPI      = InitProjectDefinitions + string(os.PathSeparator) + "asyncapi.yaml"
	InitProjectImage                    = "Image"
	InitProjectDocs                     = "Docs"
	InitProjectSequences                = "Policies"
	InitProjectClientCertificates       = "Client-certificates"
	InitProjectEndpointCertificates     = "Endpoint-certificates"
	InitProjectInterceptors             = "Interceptors"
	InitProjectLibs                     = "libs"
	InitProjectWSDL                     = "WSDL"
)

const DeploymentDirPrefix = "DeploymentArtifacts_"
const DeploymentCertificatesDirectory = "certificates"

var DefaultExportDirPath = filepath.Join(GetConfigDirPath(), DefaultExportDirName)
var DefaultCertDirPath = filepath.Join(ConfigDirPath, CertificatesDirName)

const defaultApiApplicationImportExportSuffix = "api/am/admin/v4"
const defaultPublisherApiImportExportSuffix = "api/am/publisher/v4"
const defaultApiListEndpointSuffix = "api/am/publisher/v4/apis"
const defaultMcpServerListEndpointSuffix = "api/am/publisher/v4/mcp-servers"
const defaultAPIPolicyListEndpointSuffix = "api/am/publisher/v4/operation-policies"
const defaultApiProductListEndpointSuffix = "api/am/publisher/v4/api-products"
const defaultUnifiedSearchEndpointSuffix = "api/am/publisher/v4/search"
const defaultAdminApplicationListEndpointSuffix = "api/am/admin/v4/applications"
const defaultDevPortalApplicationListEndpointSuffix = "api/am/devportal/v3/applications"
const defaultDevPortalThrottlingPoliciesEndpointSuffix = "api/am/devportal/v3/throttling-policies"
const defaultClientRegistrationEndpointSuffix = "client-registration/v0.17/register"
const defaultTokenEndPoint = "oauth2/token"
const defaultRevokeEndpointSuffix = "oauth2/revoke"
const defaultAPILoggingBaseEndpoint = "api/am/devops/v0/tenant-logs"
const defaultAPILoggingApisEndpoint = "apis"
const defaultCorrelationLoggingEndpoint = "api/am/devops/v0/config/correlation"
const defaultAIServiceEndpoint = "https://dev-tools.wso2.com/apim-ai-service/v2"
const defaultAITokenServiceEndpoint = "https://api.asgardeo.io/t/wso2devtools/oauth2/token"

const DefaultEnvironmentName = "default"
const DefaultTenantDomain = "carbon.super"

// API Product related constants
const DefaultApiProductVersion = "1.0.0"
const DefaultApiProductType = "APIProduct"

// MCP Server related constants
const DefaultMcpServerType = "MCP"

// Application keys related constants
const ProductionKeyType = "PRODUCTION"
const SandboxKeyType = "SANDBOX"

var GrantTypesToBeSupported = []string{"refresh_token", "password", "client_credentials"}

// WSO2PublicCertificate : wso2 public certificate in PEM format
var WSO2PublicCertificate = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10, 77, 73, 73, 68, 117, 84, 67, 67, 65, 113, 71, 103, 65, 119, 73, 66, 65, 103, 73, 85, 89, 72, 57, 51, 108, 65, 84, 111, 112, 116, 100, 119, 53, 113, 122, 53, 117, 57, 77, 66, 109, 75, 76, 111, 74, 57, 99, 119, 68, 81, 89, 74, 75, 111, 90, 73, 104, 118, 99, 78, 65, 81, 69, 76, 13, 10, 66, 81, 65, 119, 90, 68, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 66, 104, 77, 67, 86, 86, 77, 120, 67, 122, 65, 74, 66, 103, 78, 86, 66, 65, 103, 77, 65, 107, 78, 66, 77, 82, 89, 119, 70, 65, 89, 68, 86, 81, 81, 72, 68, 65, 49, 78, 98, 51, 86, 117, 100, 71, 70, 112, 13, 10, 98, 105, 66, 87, 97, 87, 86, 51, 77, 81, 48, 119, 67, 119, 89, 68, 86, 81, 81, 75, 68, 65, 82, 88, 85, 48, 56, 121, 77, 81, 48, 119, 67, 119, 89, 68, 86, 81, 81, 76, 68, 65, 82, 88, 85, 48, 56, 121, 77, 82, 73, 119, 69, 65, 89, 68, 86, 81, 81, 68, 68, 65, 108, 115, 13, 10, 98, 50, 78, 104, 98, 71, 104, 118, 99, 51, 81, 119, 72, 104, 99, 78, 77, 106, 85, 119, 79, 68, 69, 122, 77, 84, 65, 120, 78, 84, 65, 119, 87, 104, 99, 78, 77, 106, 99, 120, 77, 106, 73, 50, 77, 84, 65, 120, 78, 84, 65, 119, 87, 106, 66, 107, 77, 81, 115, 119, 67, 81, 89, 68, 13, 10, 86, 81, 81, 71, 69, 119, 74, 86, 85, 122, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 67, 65, 119, 67, 81, 48, 69, 120, 70, 106, 65, 85, 66, 103, 78, 86, 66, 65, 99, 77, 68, 85, 49, 118, 100, 87, 53, 48, 89, 87, 108, 117, 73, 70, 90, 112, 90, 88, 99, 120, 68, 84, 65, 76, 13, 10, 66, 103, 78, 86, 66, 65, 111, 77, 66, 70, 100, 84, 84, 122, 73, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 115, 77, 66, 70, 100, 84, 84, 122, 73, 120, 69, 106, 65, 81, 66, 103, 78, 86, 66, 65, 77, 77, 67, 87, 120, 118, 89, 50, 70, 115, 97, 71, 57, 122, 100, 68, 67, 67, 13, 10, 65, 83, 73, 119, 68, 81, 89, 74, 75, 111, 90, 73, 104, 118, 99, 78, 65, 81, 69, 66, 66, 81, 65, 68, 103, 103, 69, 80, 65, 68, 67, 67, 65, 81, 111, 67, 103, 103, 69, 66, 65, 76, 85, 122, 85, 112, 101, 115, 50, 109, 100, 54, 73, 77, 72, 101, 122, 65, 108, 43, 48, 71, 67, 82, 13, 10, 47, 56, 113, 74, 90, 108, 71, 49, 67, 112, 86, 119, 57, 80, 70, 122, 107, 107, 77, 74, 47, 98, 88, 43, 43, 104, 76, 52, 111, 49, 52, 99, 48, 52, 119, 98, 70, 79, 119, 88, 90, 69, 65, 119, 43, 72, 112, 72, 118, 106, 65, 108, 54, 50, 56, 88, 97, 108, 116, 71, 97, 118, 107, 56, 13, 10, 50, 103, 120, 105, 104, 111, 76, 114, 55, 104, 102, 74, 118, 71, 108, 78, 57, 54, 120, 104, 83, 83, 49, 110, 72, 106, 88, 120, 56, 117, 90, 82, 107, 87, 66, 68, 79, 89, 49, 104, 117, 118, 118, 65, 49, 120, 107, 121, 84, 112, 52, 119, 113, 110, 56, 56, 100, 73, 121, 105, 71, 114, 76, 43, 13, 10, 68, 74, 118, 83, 101, 102, 88, 115, 109, 57, 112, 49, 99, 120, 116, 83, 55, 118, 114, 109, 101, 118, 98, 69, 52, 118, 105, 53, 106, 79, 90, 77, 120, 122, 118, 66, 53, 69, 84, 116, 100, 50, 79, 116, 119, 88, 56, 51, 67, 71, 114, 66, 100, 73, 117, 118, 99, 56, 101, 50, 67, 111, 67, 52, 13, 10, 56, 90, 90, 49, 121, 100, 50, 120, 89, 49, 106, 98, 78, 101, 47, 111, 69, 117, 65, 48, 116, 101, 49, 56, 67, 107, 66, 77, 105, 79, 115, 56, 84, 121, 83, 66, 102, 56, 81, 76, 106, 121, 107, 80, 55, 101, 111, 78, 104, 73, 53, 74, 110, 110, 49, 118, 84, 98, 49, 55, 74, 116, 56, 88, 13, 10, 55, 122, 76, 81, 111, 106, 70, 73, 80, 105, 75, 109, 117, 100, 100, 118, 119, 105, 43, 50, 73, 67, 117, 67, 75, 52, 99, 98, 57, 106, 80, 108, 88, 102, 72, 115, 70, 73, 74, 109, 90, 77, 81, 103, 108, 47, 112, 50, 87, 71, 117, 100, 117, 105, 99, 69, 107, 119, 89, 47, 99, 99, 99, 67, 13, 10, 65, 119, 69, 65, 65, 97, 78, 106, 77, 71, 69, 119, 70, 65, 89, 68, 86, 82, 48, 82, 66, 65, 48, 119, 67, 52, 73, 74, 98, 71, 57, 106, 89, 87, 120, 111, 98, 51, 78, 48, 77, 66, 48, 71, 65, 49, 85, 100, 68, 103, 81, 87, 66, 66, 82, 78, 71, 80, 118, 75, 56, 114, 53, 56, 13, 10, 120, 98, 111, 47, 111, 79, 75, 118, 53, 67, 50, 77, 50, 86, 75, 108, 97, 106, 65, 76, 66, 103, 78, 86, 72, 81, 56, 69, 66, 65, 77, 67, 66, 80, 65, 119, 72, 81, 89, 68, 86, 82, 48, 108, 66, 66, 89, 119, 70, 65, 89, 73, 75, 119, 89, 66, 66, 81, 85, 72, 65, 119, 69, 71, 13, 10, 67, 67, 115, 71, 65, 81, 85, 70, 66, 119, 77, 67, 77, 65, 48, 71, 67, 83, 113, 71, 83, 73, 98, 51, 68, 81, 69, 66, 67, 119, 85, 65, 65, 52, 73, 66, 65, 81, 66, 76, 75, 88, 112, 72, 116, 120, 50, 112, 115, 65, 98, 90, 75, 55, 119, 49, 80, 47, 120, 110, 49, 68, 47, 78, 13, 10, 108, 83, 43, 121, 122, 116, 57, 87, 73, 82, 107, 57, 118, 105, 49, 65, 78, 80, 73, 112, 104, 111, 54, 71, 116, 81, 86, 67, 74, 105, 87, 90, 50, 43, 51, 89, 53, 105, 76, 105, 56, 103, 107, 51, 98, 76, 101, 47, 43, 114, 117, 84, 74, 90, 55, 68, 50, 66, 72, 109, 115, 51, 104, 101, 13, 10, 110, 86, 106, 107, 47, 102, 79, 114, 57, 74, 65, 109, 73, 66, 48, 113, 112, 120, 99, 86, 106, 79, 77, 56, 56, 51, 78, 76, 52, 81, 121, 65, 69, 83, 74, 88, 72, 48, 121, 112, 65, 68, 86, 72, 70, 109, 112, 86, 81, 57, 66, 101, 103, 98, 78, 117, 50, 65, 75, 89, 119, 115, 98, 53, 13, 10, 54, 66, 106, 49, 119, 68, 70, 100, 52, 54, 88, 67, 99, 114, 97, 53, 114, 101, 104, 89, 118, 79, 43, 101, 70, 77, 78, 80, 122, 82, 53, 87, 90, 67, 70, 103, 108, 111, 111, 122, 109, 81, 66, 88, 68, 113, 49, 56, 84, 47, 80, 112, 56, 52, 106, 56, 103, 75, 79, 97, 72, 81, 106, 70, 13, 10, 81, 73, 65, 90, 47, 75, 84, 80, 81, 109, 69, 79, 81, 55, 43, 105, 97, 99, 66, 104, 100, 67, 102, 115, 122, 48, 120, 78, 90, 49, 52, 90, 104, 74, 100, 87, 77, 118, 66, 100, 86, 80, 47, 116, 70, 43, 112, 85, 79, 51, 103, 43, 48, 70, 117, 106, 73, 97, 122, 57, 50, 119, 98, 70, 13, 10, 82, 122, 122, 66, 106, 47, 82, 99, 88, 117, 50, 56, 55, 43, 110, 69, 117, 84, 70, 70, 120, 74, 100, 119, 121, 117, 56, 56, 67, 72, 69, 71, 56, 86, 77, 90, 56, 70, 71, 88, 101, 52, 117, 101, 49, 122, 66, 109, 55, 97, 84, 53, 78, 76, 105, 88, 50, 114, 106, 78, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10}

// Headers and Header Values
const HeaderAuthorization = "Authorization"
const HeaderContentType = "Content-Type"
const HeaderConnection = "Connection"
const HeaderAccept = "Accept"
const HeaderProduces = "Produces"
const HeaderConsumes = "Consumes"
const HeaderContentEncoding = "Content-Encoding"
const HeaderTransferEncoding = "transfer-encoding"
const HeaderValueChunked = "chunked"
const HeaderValueGZIP = "gzip"
const HeaderValueKeepAlive = "keep-alive"
const HeaderValueApplicationZip = "application/zip"
const HeaderValueApplicationJSON = "application/json"
const HeaderValueXWWWFormUrlEncoded = "application/x-www-form-urlencoded"
const HeaderValueAuthBearerPrefix = "Bearer"
const HeaderValueAuthBasicPrefix = "Basic"
const HeaderValueMultiPartFormData = "multipart/form-data"
const HeaderToken = "token="
const TokenTypeForRevocation = "&token_type_hint=access_token"

// Logging Prefixes
const LogPrefixInfo = "[INFO]: "
const LogPrefixWarning = "[WARN]: "
const LogPrefixError = "[ERROR]: "

// String Constants
const SearchAndTag = "&"

// Other
const DefaultTokenValidityPeriod = 3600
const DefaultHttpRequestTimeout = 10000

// AI
const DefaultAIThreadCount = 3
const DefaultAIEndpoint = "https://e95488c8-8511-4882-967f-ec3ae2a0f86f-prod.e1-us-east-azure.choreoapis.dev/lgpt/interceptor-service/interceptor-service-be2/v1.0"

// TLSRenegotiationNever : never negotiate
const TLSRenegotiationNever = "never"

// TLSRenegotiationOnce : negotiate once
const TLSRenegotiationOnce = "once"

// TLSRenegotiationFreely : negotiate freely
const TLSRenegotiationFreely = "freely"

// Migration export
const MaxAPIsToExportOnce = 20
const MaxAppsToExportOnce = 20
const MaxMCPServersToExportOnce = 20
const MigrationAPIsExportMetadataFileName = "migration-apis-export-metadata.yaml"
const MigrationAppsExportMetadataFileName = "migration-apps-export-metadata.yaml"
const MigrationMCPServersExportMetadataFileName = "migration-mcp-servers-export-metadata.yaml"
const LastSucceededApiFileName = "last-succeeded-api.log"
const LastSucceededAppFileName = "last-succeeded-app.log"
const LastSucceededMCPServerFileName = "last_succeeded_mcp_server.log"
const LastSuceededContentDelimiter = " " // space
const DefaultResourceTenantDomain = "tenant-default"
const ApplicationId = "applicationId"
const ApiId = "apiId"
const APIProductId = "apiProductId"
const DefaultCliApp = "default-apictl-app"
const DefaultTokenType = "JWT"

const LifeCycleAction = "action"

var ValidInitialStates = []string{"CREATED", "PUBLISHED"}

// The list of repos and directories that can be used when replcing env variables
var EnvReplaceFilePaths = []string{
	"Policies",
}

// The list of file extensions when replcing env variables related to Policies
var EnvReplacePoliciesFileExtensions = []string{
	"j2",
	"gotmpl",
}

// project types
const (
	ProjectTypeNone        = "None"
	ProjectTypeApi         = "API"
	ProjectTypeMcpServer   = "MCP Server"
	ProjectTypeApiProduct  = "API Product"
	ProjectTypeApplication = "Application"
	ProjectTypeRevision    = "Revision"
	ProjectTypePolicy      = "Policy"
	ProjectTypeAPIPolicy   = "API Policy"
)

// project param files
const ParamFile = "params.yaml"
const ParamsIntermediateFile = "intermediate_params.yaml"

const (
	APIDefinitionFileYaml         = "api.yaml"
	APIDefinitionFileJson         = "api.json"
	MCPServerDefinitionFileYaml   = "mcp_server.yaml"
	MCPServerDefinitionFileJson   = "mcp_server.json"
	APIProductDefinitionFileYaml  = "api_product.yaml"
	APIProductDefinitionFileJson  = "api_product.json"
	ApplicationDefinitionFileYaml = "application.yaml"
	ApplicationDefinitionFileJson = "application.json"
)

// project meta files
const (
	MetaFileAPI         = "api_meta.yaml"
	MetaFileMCPServer   = "mcp_server_meta.yaml"
	MetaFileAPIProduct  = "api_product_meta.yaml"
	MetaFileApplication = "application_meta.yaml"
)

// Constants related to meta file structs
const DeployImportRotateRevision = "deploy.import.rotateRevision"
const DeployImportSkipSubscriptions = "deploy.import.skipSubscriptions"

const DeploymentEnvFile = "deployment_environments.yaml"
const PrivateJetModeConst = "privateJet"
const SidecarModeConst = "sidecar"

// Default values for Help commands
const DefaultApisDisplayLimit = 25
const DefaultApiProductsDisplayLimit = 25
const DefaultAppsDisplayLimit = 25
const DefaultExportFormat = "YAML"
const DefaultPoliciesDisplayLimit = 25

const InitDirName = string(os.PathSeparator) + "init" + string(os.PathSeparator)

// AWS API security document constants
const DefaultAWSDocFileName = "document.yaml"

const ResourcePolicyDocName = "resource_policy_doc"
const ResourcePolicyDocDisplayName = "Resource Policy"
const ResourcePolicyDocSummary = "This document contains details related to AWS resource policies"

const CognitoUserPoolDocName = "cognito_userpool_doc"
const CognitoDocDisplayName = "Cognito Userpool"
const CognitoDocSummary = "This document contains details related to AWS cognito user pools"

const AWSAPIKeyDocName = "aws_apikey_doc"
const ApiKeysDocDisplayName = "AWS APIKeys"
const ApiKeysDocSummary = "This document contains details related to AWS API keys"

const AWSSigV4DocName = "aws_sigv4_doc"
const AWSSigV4DocDisplayName = "AWS Signature Version4"
const AWSSigV4DocSummary = "This document contains details related to AWS signature version 4"

// MiCmdLiteral denote the alias for micro integrator related commands
const MiCmdLiteral = "mi"

// MiManagementAPIContext
const MiManagementAPIContext = "management"

// Mi Management Resource paths
const MiManagementCarbonAppResource = "applications"
const MiManagementServiceResource = "services"
const MiManagementAPIResource = "apis"
const MiManagementProxyServiceResource = "proxy-services"
const MiManagementInboundEndpointResource = "inbound-endpoints"
const MiManagementEndpointResource = "endpoints"
const MiManagementMessageProcessorResource = "message-processors"
const MiManagementTemplateResource = "templates"
const MiManagementConnectorResource = "connectors"
const MiManagementMessageStoreResource = "message-stores"
const MiManagementLocalEntrieResource = "local-entries"
const MiManagementSequenceResource = "sequences"
const MiManagementTaskResource = "tasks"
const MiManagementLogResource = "logs"
const MiManagementLoggingResource = "logging"
const MiManagementServerResource = "server"
const MiManagementDataServiceResource = "data-services"
const MiManagementMiLoginResource = "login"
const MiManagementMiLogoutResource = "logout"
const MiManagementUserResource = "users"
const MiManagementTransactionResource = "transactions"
const MiManagementTransactionCountResource = "count"
const MiManagementTransactionReportResource = "report"
const MiManagementExternalVaultsResource = "external-vaults"
const MiManagementExternalVaultHashiCorpResource = "hashicorp"
const MiManagementRoleResource = "roles"

const ZipFileSuffix = ".zip"

// Output format types
const JsonArrayFormatType = "jsonArray"

const ThrottlingPolicyTypeSub = "subscription"
const ThrottlingPolicyTypeApp = "application"
const ThrottlingPolicyTypeAdv = "advanced"
const ThrottlingPolicyTypeCus = "custom"
