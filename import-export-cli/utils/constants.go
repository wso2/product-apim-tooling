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

const DefaultEnvironmentName = "default"
const DefaultTenantDomain = "carbon.super"

// API Product related constants
const DefaultApiProductVersion = "1.0.0"
const DefaultApiProductType = "APIProduct"

// Application keys related constants
const ProductionKeyType = "PRODUCTION"
const SandboxKeyType = "SANDBOX"

var GrantTypesToBeSupported = []string{"refresh_token", "password", "client_credentials"}

// WSO2PublicCertificate : wso2 public certificate in PEM format
var WSO2PublicCertificate = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10, 77, 73, 73, 68, 117, 84, 67, 67, 65, 113, 71, 103, 65, 119, 73, 66, 65, 103, 73, 85, 89, 121, 83, 43, 98, 99, 122, 115, 56, 71, 83, 119, 99, 120, 104, 81, 69, 50, 89, 66, 106, 100, 65, 69, 106, 76, 56, 119, 68, 81, 89, 74, 75, 111, 90, 73, 104, 118, 99, 78, 65, 81, 69, 76, 10, 66, 81, 65, 119, 90, 68, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 66, 104, 77, 67, 86, 86, 77, 120, 67, 122, 65, 74, 66, 103, 78, 86, 66, 65, 103, 77, 65, 107, 78, 66, 77, 82, 89, 119, 70, 65, 89, 68, 86, 81, 81, 72, 68, 65, 49, 78, 98, 51, 86, 117, 100, 71, 70, 112, 10, 98, 105, 66, 87, 97, 87, 86, 51, 77, 81, 48, 119, 67, 119, 89, 68, 86, 81, 81, 75, 68, 65, 82, 88, 85, 48, 56, 121, 77, 81, 48, 119, 67, 119, 89, 68, 86, 81, 81, 76, 68, 65, 82, 88, 85, 48, 56, 121, 77, 82, 73, 119, 69, 65, 89, 68, 86, 81, 81, 68, 68, 65, 108, 115, 10, 98, 50, 78, 104, 98, 71, 104, 118, 99, 51, 81, 119, 72, 104, 99, 78, 77, 106, 81, 119, 77, 106, 73, 121, 77, 68, 103, 49, 77, 106, 85, 49, 87, 104, 99, 78, 77, 106, 89, 119, 78, 84, 73, 51, 77, 68, 103, 49, 77, 106, 85, 49, 87, 106, 66, 107, 77, 81, 115, 119, 67, 81, 89, 68, 10, 86, 81, 81, 71, 69, 119, 74, 86, 85, 122, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 67, 65, 119, 67, 81, 48, 69, 120, 70, 106, 65, 85, 66, 103, 78, 86, 66, 65, 99, 77, 68, 85, 49, 118, 100, 87, 53, 48, 89, 87, 108, 117, 73, 70, 90, 112, 90, 88, 99, 120, 68, 84, 65, 76, 10, 66, 103, 78, 86, 66, 65, 111, 77, 66, 70, 100, 84, 84, 122, 73, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 115, 77, 66, 70, 100, 84, 84, 122, 73, 120, 69, 106, 65, 81, 66, 103, 78, 86, 66, 65, 77, 77, 67, 87, 120, 118, 89, 50, 70, 115, 97, 71, 57, 122, 100, 68, 67, 67, 10, 65, 83, 73, 119, 68, 81, 89, 74, 75, 111, 90, 73, 104, 118, 99, 78, 65, 81, 69, 66, 66, 81, 65, 68, 103, 103, 69, 80, 65, 68, 67, 67, 65, 81, 111, 67, 103, 103, 69, 66, 65, 79, 98, 72, 101, 52, 114, 114, 50, 49, 70, 84, 55, 88, 117, 111, 121, 122, 51, 67, 83, 48, 114, 43, 10, 121, 56, 81, 76, 104, 112, 121, 111, 103, 70, 108, 88, 119, 48, 76, 67, 73, 80, 70, 106, 86, 113, 76, 85, 109, 65, 66, 73, 57, 83, 70, 73, 122, 55, 112, 86, 122, 74, 112, 113, 84, 80, 108, 52, 66, 108, 97, 85, 113, 71, 55, 55, 66, 55, 50, 50, 83, 73, 115, 107, 66, 110, 76, 82, 10, 71, 115, 53, 71, 87, 53, 65, 72, 71, 117, 106, 72, 85, 102, 106, 114, 82, 85, 67, 47, 76, 69, 48, 104, 119, 97, 97, 55, 65, 53, 107, 83, 84, 110, 120, 112, 118, 113, 76, 85, 57, 103, 122, 80, 108, 53, 71, 81, 50, 56, 119, 98, 111, 112, 105, 112, 119, 48, 89, 50, 67, 78, 71, 113, 10, 108, 89, 114, 84, 74, 119, 76, 48, 47, 83, 50, 55, 66, 67, 101, 79, 98, 89, 107, 82, 120, 90, 107, 52, 74, 47, 70, 86, 98, 75, 101, 77, 86, 110, 82, 65, 85, 89, 101, 65, 55, 82, 54, 73, 82, 56, 119, 66, 117, 54, 119, 97, 80, 83, 53, 67, 78, 79, 84, 107, 105, 47, 90, 119, 10, 53, 48, 65, 111, 106, 71, 97, 100, 86, 102, 109, 71, 121, 54, 81, 67, 54, 103, 86, 97, 105, 67, 56, 73, 97, 109, 66, 90, 75, 107, 118, 79, 106, 99, 47, 116, 67, 97, 117, 80, 57, 66, 78, 107, 72, 103, 112, 65, 110, 65, 97, 88, 79, 90, 116, 83, 114, 69, 88, 54, 105, 47, 113, 74, 10, 71, 47, 90, 88, 109, 79, 99, 74, 83, 70, 71, 72, 75, 122, 106, 72, 79, 117, 70, 105, 54, 43, 56, 98, 80, 43, 75, 119, 115, 51, 85, 77, 47, 118, 70, 102, 100, 52, 71, 104, 50, 97, 83, 50, 98, 54, 103, 79, 119, 50, 53, 116, 47, 82, 48, 108, 68, 74, 86, 43, 119, 116, 69, 67, 10, 65, 119, 69, 65, 65, 97, 78, 106, 77, 71, 69, 119, 70, 65, 89, 68, 86, 82, 48, 82, 66, 65, 48, 119, 67, 52, 73, 74, 98, 71, 57, 106, 89, 87, 120, 111, 98, 51, 78, 48, 77, 66, 48, 71, 65, 49, 85, 100, 68, 103, 81, 87, 66, 66, 83, 72, 65, 120, 75, 111, 75, 101, 90, 115, 10, 90, 115, 49, 110, 105, 43, 69, 88, 89, 87, 55, 53, 102, 77, 55, 48, 98, 106, 65, 100, 66, 103, 78, 86, 72, 83, 85, 69, 70, 106, 65, 85, 66, 103, 103, 114, 66, 103, 69, 70, 66, 81, 99, 68, 65, 81, 89, 73, 75, 119, 89, 66, 66, 81, 85, 72, 65, 119, 73, 119, 67, 119, 89, 68, 10, 86, 82, 48, 80, 66, 65, 81, 68, 65, 103, 84, 119, 77, 65, 48, 71, 67, 83, 113, 71, 83, 73, 98, 51, 68, 81, 69, 66, 67, 119, 85, 65, 65, 52, 73, 66, 65, 81, 68, 71, 76, 57, 83, 83, 65, 99, 56, 122, 101, 90, 112, 85, 111, 105, 104, 83, 97, 49, 86, 122, 48, 100, 99, 87, 10, 110, 69, 71, 89, 100, 116, 55, 100, 66, 112, 120, 83, 122, 105, 113, 82, 86, 72, 53, 79, 54, 50, 113, 85, 79, 118, 100, 56, 109, 108, 104, 121, 51, 72, 120, 55, 112, 118, 90, 109, 50, 43, 106, 74, 98, 53, 76, 120, 71, 74, 51, 99, 47, 87, 56, 70, 51, 109, 117, 116, 49, 83, 116, 106, 10, 65, 83, 69, 69, 117, 108, 57, 97, 120, 117, 120, 111, 53, 102, 110, 122, 85, 87, 107, 54, 65, 87, 109, 50, 112, 97, 104, 108, 75, 109, 101, 48, 98, 74, 79, 47, 117, 66, 114, 99, 77, 116, 90, 85, 54, 77, 119, 66, 118, 69, 55, 54, 65, 51, 102, 104, 43, 107, 90, 122, 86, 84, 68, 108, 10, 76, 110, 50, 84, 84, 111, 114, 74, 88, 114, 116, 47, 73, 66, 110, 74, 117, 86, 71, 79, 101, 51, 88, 108, 70, 78, 83, 50, 48, 74, 55, 108, 75, 69, 82, 107, 119, 80, 118, 68, 102, 88, 100, 103, 110, 90, 73, 105, 84, 98, 71, 48, 119, 57, 98, 65, 76, 85, 74, 106, 81, 56, 85, 50, 10, 65, 65, 80, 57, 97, 107, 103, 117, 55, 112, 76, 87, 76, 55, 71, 53, 53, 115, 83, 86, 111, 98, 70, 114, 66, 76, 74, 69, 76, 103, 113, 113, 83, 85, 84, 85, 120, 65, 57, 117, 114, 111, 54, 48, 70, 49, 81, 48, 85, 97, 88, 112, 78, 119, 104, 116, 49, 66, 119, 65, 49, 85, 80, 66, 10, 115, 43, 53, 79, 118, 102, 87, 56, 103, 49, 104, 100, 86, 103, 86, 118, 99, 102, 101, 70, 72, 66, 79, 104, 73, 50, 78, 107, 84, 109, 122, 100, 70, 70, 87, 74, 101, 48, 69, 120, 79, 67, 111, 121, 98, 71, 87, 73, 102, 70, 100, 98, 86, 109, 106, 85, 89, 47, 98, 82, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10}

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
const DefaultMarketplaceAssistantThreadCount = 3
const DefaultAIEndpoint = "https://dev-tools.wso2.com/apim-ai-service"

// TLSRenegotiationNever : never negotiate
const TLSRenegotiationNever = "never"

// TLSRenegotiationOnce : negotiate once
const TLSRenegotiationOnce = "once"

// TLSRenegotiationFreely : negotiate freely
const TLSRenegotiationFreely = "freely"

// Migration export
const MaxAPIsToExportOnce = 20
const MigrationAPIsExportMetadataFileName = "migration-apis-export-metadata.yaml"
const LastSucceededApiFileName = "last-succeeded-api.log"
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
