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
var WSO2PublicCertificate = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10, 77, 73, 73, 68, 117, 84, 67, 67, 65, 113, 71, 103, 65, 119, 73, 66, 65, 103, 73, 85, 90, 51, 114, 89, 75, 86, 78, 90, 84, 47, 97, 84, 77, 106, 79, 67, 109, 106, 115, 66, 108, 80, 57, 108, 79, 118, 81, 119, 68, 81, 89, 74, 75, 111, 90, 73, 104, 118, 99, 78, 65, 81, 69, 76, 10, 66, 81, 65, 119, 90, 68, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 66, 104, 77, 67, 86, 86, 77, 120, 67, 122, 65, 74, 66, 103, 78, 86, 66, 65, 103, 77, 65, 107, 78, 66, 77, 82, 89, 119, 70, 65, 89, 68, 86, 81, 81, 72, 68, 65, 49, 78, 98, 51, 86, 117, 100, 71, 70, 112, 10, 98, 105, 66, 87, 97, 87, 86, 51, 77, 81, 48, 119, 67, 119, 89, 68, 86, 81, 81, 75, 68, 65, 82, 88, 85, 48, 56, 121, 77, 81, 48, 119, 67, 119, 89, 68, 86, 81, 81, 76, 68, 65, 82, 88, 85, 48, 56, 121, 77, 82, 73, 119, 69, 65, 89, 68, 86, 81, 81, 68, 68, 65, 108, 115, 10, 98, 50, 78, 104, 98, 71, 104, 118, 99, 51, 81, 119, 72, 104, 99, 78, 77, 106, 85, 119, 77, 106, 69, 122, 77, 84, 77, 119, 77, 68, 69, 120, 87, 104, 99, 78, 77, 106, 99, 119, 78, 84, 69, 53, 77, 84, 77, 119, 77, 68, 69, 120, 87, 106, 66, 107, 77, 81, 115, 119, 67, 81, 89, 68, 10, 86, 81, 81, 71, 69, 119, 74, 86, 85, 122, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 67, 65, 119, 67, 81, 48, 69, 120, 70, 106, 65, 85, 66, 103, 78, 86, 66, 65, 99, 77, 68, 85, 49, 118, 100, 87, 53, 48, 89, 87, 108, 117, 73, 70, 90, 112, 90, 88, 99, 120, 68, 84, 65, 76, 10, 66, 103, 78, 86, 66, 65, 111, 77, 66, 70, 100, 84, 84, 122, 73, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 115, 77, 66, 70, 100, 84, 84, 122, 73, 120, 69, 106, 65, 81, 66, 103, 78, 86, 66, 65, 77, 77, 67, 87, 120, 118, 89, 50, 70, 115, 97, 71, 57, 122, 100, 68, 67, 67, 10, 65, 83, 73, 119, 68, 81, 89, 74, 75, 111, 90, 73, 104, 118, 99, 78, 65, 81, 69, 66, 66, 81, 65, 68, 103, 103, 69, 80, 65, 68, 67, 67, 65, 81, 111, 67, 103, 103, 69, 66, 65, 75, 47, 84, 122, 57, 70, 118, 117, 49, 77, 122, 101, 82, 74, 57, 89, 108, 69, 80, 103, 66, 79, 115, 10, 114, 43, 111, 65, 78, 80, 121, 66, 71, 102, 72, 101, 74, 85, 121, 51, 74, 74, 118, 86, 79, 88, 104, 76, 117, 54, 76, 88, 70, 85, 112, 108, 67, 102, 80, 87, 113, 101, 104, 101, 76, 112, 77, 73, 85, 120, 78, 113, 76, 86, 100, 105, 51, 117, 101, 78, 102, 98, 113, 88, 57, 90, 105, 110, 10, 43, 65, 78, 112, 120, 53, 109, 43, 70, 116, 119, 107, 106, 53, 119, 99, 84, 80, 67, 110, 106, 68, 114, 114, 104, 110, 79, 53, 76, 84, 81, 120, 114, 111, 116, 57, 101, 116, 112, 121, 53, 49, 72, 103, 86, 87, 117, 50, 105, 85, 53, 108, 77, 101, 82, 111, 73, 52, 119, 65, 100, 105, 100, 103, 10, 119, 100, 75, 99, 90, 75, 82, 67, 69, 101, 117, 82, 121, 100, 83, 88, 101, 122, 76, 48, 67, 71, 87, 69, 112, 51, 116, 100, 65, 53, 47, 47, 115, 74, 53, 108, 105, 121, 90, 49, 120, 114, 66, 50, 56, 54, 107, 69, 74, 114, 75, 101, 71, 68, 79, 74, 53, 105, 84, 53, 104, 119, 76, 89, 10, 100, 74, 84, 99, 48, 80, 108, 100, 73, 70, 56, 72, 83, 101, 47, 98, 115, 87, 65, 108, 68, 47, 78, 89, 81, 65, 50, 111, 67, 120, 73, 70, 49, 118, 101, 47, 77, 80, 101, 79, 97, 76, 56, 107, 102, 66, 105, 116, 121, 116, 49, 54, 82, 116, 112, 55, 80, 107, 110, 105, 81, 118, 109, 55, 10, 121, 86, 87, 99, 79, 77, 99, 107, 77, 110, 115, 65, 97, 57, 56, 116, 80, 113, 109, 72, 85, 112, 52, 119, 57, 118, 117, 68, 116, 121, 67, 104, 111, 120, 117, 89, 50, 89, 120, 52, 86, 48, 101, 105, 113, 105, 82, 81, 74, 66, 100, 74, 114, 43, 105, 57, 75, 66, 90, 85, 118, 75, 77, 67, 10, 65, 119, 69, 65, 65, 97, 78, 106, 77, 71, 69, 119, 70, 65, 89, 68, 86, 82, 48, 82, 66, 65, 48, 119, 67, 52, 73, 74, 98, 71, 57, 106, 89, 87, 120, 111, 98, 51, 78, 48, 77, 66, 48, 71, 65, 49, 85, 100, 68, 103, 81, 87, 66, 66, 84, 66, 47, 98, 119, 75, 51, 89, 47, 65, 10, 117, 73, 88, 111, 78, 111, 56, 108, 78, 117, 87, 76, 52, 86, 74, 72, 66, 84, 65, 76, 66, 103, 78, 86, 72, 81, 56, 69, 66, 65, 77, 67, 66, 80, 65, 119, 72, 81, 89, 68, 86, 82, 48, 108, 66, 66, 89, 119, 70, 65, 89, 73, 75, 119, 89, 66, 66, 81, 85, 72, 65, 119, 69, 71, 10, 67, 67, 115, 71, 65, 81, 85, 70, 66, 119, 77, 67, 77, 65, 48, 71, 67, 83, 113, 71, 83, 73, 98, 51, 68, 81, 69, 66, 67, 119, 85, 65, 65, 52, 73, 66, 65, 81, 67, 108, 55, 108, 87, 102, 109, 49, 83, 78, 88, 120, 120, 122, 48, 80, 99, 99, 84, 118, 88, 98, 51, 85, 117, 116, 10, 106, 48, 89, 69, 101, 70, 107, 55, 52, 98, 82, 48, 65, 70, 107, 90, 51, 87, 84, 69, 79, 99, 104, 84, 90, 79, 97, 51, 101, 106, 74, 112, 112, 112, 76, 83, 105, 119, 65, 101, 82, 98, 87, 68, 54, 111, 54, 47, 48, 82, 48, 52, 108, 76, 103, 102, 65, 101, 55, 89, 53, 97, 107, 104, 10, 56, 88, 57, 55, 74, 51, 71, 71, 104, 111, 106, 99, 110, 57, 74, 114, 83, 43, 70, 67, 67, 106, 120, 102, 73, 84, 54, 49, 113, 119, 65, 119, 115, 97, 71, 74, 109, 103, 111, 73, 65, 112, 76, 71, 97, 57, 72, 75, 57, 49, 86, 47, 67, 103, 122, 105, 47, 108, 119, 79, 106, 88, 105, 54, 10, 82, 97, 102, 105, 48, 68, 78, 73, 57, 49, 88, 114, 73, 67, 121, 77, 71, 111, 118, 43, 86, 119, 111, 53, 121, 98, 98, 86, 121, 89, 55, 108, 97, 78, 50, 78, 86, 78, 117, 81, 107, 71, 74, 109, 118, 77, 52, 54, 119, 98, 57, 43, 50, 106, 110, 52, 83, 122, 79, 122, 89, 117, 79, 119, 10, 77, 72, 74, 105, 88, 68, 83, 104, 57, 90, 98, 117, 97, 75, 105, 78, 105, 65, 116, 98, 105, 86, 103, 89, 115, 75, 109, 102, 114, 105, 97, 115, 86, 101, 97, 119, 90, 49, 108, 81, 68, 112, 88, 74, 65, 66, 116, 43, 65, 50, 78, 110, 65, 52, 89, 114, 73, 51, 73, 74, 54, 111, 107, 54, 10, 104, 90, 106, 114, 103, 74, 113, 100, 50, 99, 53, 97, 100, 106, 68, 119, 56, 43, 76, 68, 77, 71, 89, 89, 53, 50, 48, 106, 99, 97, 118, 53, 110, 69, 117, 76, 97, 105, 98, 50, 98, 100, 70, 68, 88, 75, 54, 78, 104, 107, 66, 48, 73, 65, 87, 103, 48, 108, 83, 71, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10}

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
	APIProductDefinitionFileYaml  = "api_product.yaml"
	APIProductDefinitionFileJson  = "api_product.json"
	ApplicationDefinitionFileYaml = "application.yaml"
	ApplicationDefinitionFileJson = "application.json"
)

// project meta files
const (
	MetaFileAPI         = "api_meta.yaml"
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
