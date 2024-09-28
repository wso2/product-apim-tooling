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
var WSO2PublicCertificate = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10, 77, 73, 73, 68, 113, 84, 67, 67, 65, 112, 71, 103, 65, 119, 73, 66, 65, 103, 73, 69, 90, 116, 43, 57, 56, 106, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 115, 70, 65, 68, 66, 107, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 71, 69, 119, 74, 86, 10, 85, 122, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 67, 65, 119, 67, 81, 48, 69, 120, 70, 106, 65, 85, 66, 103, 78, 86, 66, 65, 99, 77, 68, 85, 49, 118, 100, 87, 53, 48, 89, 87, 108, 117, 73, 70, 90, 112, 90, 88, 99, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 111, 77, 10, 66, 70, 100, 84, 84, 122, 73, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 115, 77, 66, 70, 100, 84, 84, 122, 73, 120, 69, 106, 65, 81, 66, 103, 78, 86, 66, 65, 77, 77, 67, 87, 120, 118, 89, 50, 70, 115, 97, 71, 57, 122, 100, 68, 65, 101, 70, 119, 48, 121, 78, 68, 65, 53, 10, 77, 84, 65, 119, 77, 122, 77, 122, 77, 68, 90, 97, 70, 119, 48, 121, 78, 106, 69, 121, 77, 84, 81, 119, 77, 122, 77, 122, 77, 68, 90, 97, 77, 71, 81, 120, 67, 122, 65, 74, 66, 103, 78, 86, 66, 65, 89, 84, 65, 108, 86, 84, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 73, 10, 68, 65, 74, 68, 81, 84, 69, 87, 77, 66, 81, 71, 65, 49, 85, 69, 66, 119, 119, 78, 84, 87, 57, 49, 98, 110, 82, 104, 97, 87, 52, 103, 86, 109, 108, 108, 100, 122, 69, 78, 77, 65, 115, 71, 65, 49, 85, 69, 67, 103, 119, 69, 86, 49, 78, 80, 77, 106, 69, 78, 77, 65, 115, 71, 10, 65, 49, 85, 69, 67, 119, 119, 69, 86, 49, 78, 80, 77, 106, 69, 83, 77, 66, 65, 71, 65, 49, 85, 69, 65, 119, 119, 74, 98, 71, 57, 106, 89, 87, 120, 111, 98, 51, 78, 48, 77, 73, 73, 66, 73, 106, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 69, 70, 10, 65, 65, 79, 67, 65, 81, 56, 65, 77, 73, 73, 66, 67, 103, 75, 67, 65, 81, 69, 65, 117, 72, 115, 80, 102, 76, 106, 109, 66, 88, 50, 67, 75, 104, 101, 50, 120, 68, 80, 70, 72, 53, 98, 108, 105, 118, 97, 112, 109, 79, 101, 73, 43, 71, 99, 68, 101, 75, 74, 68, 79, 83, 110, 104, 10, 78, 115, 53, 120, 111, 101, 85, 43, 79, 82, 81, 109, 84, 105, 80, 48, 103, 84, 65, 51, 72, 97, 79, 86, 51, 90, 107, 68, 114, 114, 115, 54, 74, 108, 104, 103, 48, 50, 122, 70, 97, 115, 114, 117, 48, 111, 90, 87, 116, 76, 102, 113, 106, 99, 78, 101, 110, 43, 119, 53, 112, 79, 108, 86, 10, 103, 118, 105, 50, 51, 83, 112, 57, 73, 81, 109, 54, 108, 110, 102, 86, 80, 103, 73, 79, 56, 112, 90, 98, 106, 97, 43, 114, 86, 100, 86, 53, 74, 78, 55, 85, 88, 99, 117, 111, 111, 100, 112, 108, 121, 68, 97, 110, 65, 79, 74, 56, 90, 115, 101, 57, 110, 67, 43, 80, 55, 74, 57, 88, 10, 84, 105, 102, 101, 83, 99, 114, 99, 107, 112, 109, 78, 106, 111, 103, 80, 85, 101, 77, 79, 97, 50, 49, 103, 108, 43, 119, 110, 89, 68, 79, 117, 111, 86, 65, 80, 72, 43, 73, 104, 120, 57, 47, 74, 90, 117, 69, 66, 89, 99, 79, 65, 76, 86, 114, 54, 107, 57, 119, 51, 70, 119, 118, 83, 10, 57, 50, 72, 90, 56, 70, 115, 76, 97, 82, 118, 102, 53, 50, 52, 120, 68, 103, 88, 53, 108, 68, 112, 103, 82, 98, 54, 47, 122, 56, 120, 121, 117, 66, 102, 83, 68, 120, 55, 80, 69, 87, 85, 66, 119, 55, 109, 109, 57, 54, 82, 84, 100, 115, 74, 85, 103, 79, 81, 74, 48, 88, 98, 106, 10, 78, 71, 100, 57, 107, 72, 97, 51, 49, 86, 47, 71, 82, 70, 48, 106, 97, 90, 70, 83, 76, 102, 79, 82, 68, 97, 106, 85, 56, 101, 78, 120, 79, 87, 122, 118, 52, 49, 77, 90, 117, 119, 73, 68, 65, 81, 65, 66, 111, 50, 77, 119, 89, 84, 65, 85, 66, 103, 78, 86, 72, 82, 69, 69, 10, 68, 84, 65, 76, 103, 103, 108, 115, 98, 50, 78, 104, 98, 71, 104, 118, 99, 51, 81, 119, 72, 81, 89, 68, 86, 82, 48, 79, 66, 66, 89, 69, 70, 67, 103, 74, 51, 71, 72, 107, 79, 117, 87, 65, 47, 102, 49, 113, 113, 66, 112, 105, 77, 53, 104, 51, 88, 79, 114, 115, 77, 66, 48, 71, 10, 65, 49, 85, 100, 74, 81, 81, 87, 77, 66, 81, 71, 67, 67, 115, 71, 65, 81, 85, 70, 66, 119, 77, 66, 66, 103, 103, 114, 66, 103, 69, 70, 66, 81, 99, 68, 65, 106, 65, 76, 66, 103, 78, 86, 72, 81, 56, 69, 66, 65, 77, 67, 66, 80, 65, 119, 68, 81, 89, 74, 75, 111, 90, 73, 10, 104, 118, 99, 78, 65, 81, 69, 76, 66, 81, 65, 68, 103, 103, 69, 66, 65, 66, 110, 104, 88, 86, 97, 98, 118, 74, 99, 80, 117, 121, 53, 73, 99, 98, 71, 57, 106, 57, 47, 120, 119, 90, 76, 52, 77, 106, 52, 75, 116, 53, 75, 106, 121, 110, 98, 50, 67, 115, 89, 111, 111, 50, 88, 89, 10, 77, 84, 52, 55, 106, 75, 117, 84, 101, 80, 50, 66, 80, 102, 79, 75, 112, 113, 52, 43, 82, 89, 86, 80, 69, 50, 67, 85, 79, 115, 114, 81, 118, 68, 106, 81, 75, 115, 99, 102, 90, 54, 78, 77, 109, 107, 88, 47, 76, 117, 105, 73, 66, 78, 81, 89, 116, 120, 90, 69, 66, 79, 110, 75, 10, 101, 85, 107, 111, 100, 72, 53, 105, 97, 99, 70, 87, 111, 85, 88, 103, 100, 66, 83, 105, 72, 109, 105, 104, 99, 55, 77, 49, 97, 65, 88, 52, 97, 68, 48, 65, 113, 98, 75, 54, 56, 98, 118, 122, 104, 67, 108, 100, 113, 119, 66, 87, 67, 101, 109, 76, 43, 90, 104, 113, 115, 72, 99, 57, 10, 102, 71, 113, 106, 101, 109, 71, 52, 47, 52, 108, 55, 75, 83, 53, 99, 111, 114, 53, 104, 119, 47, 108, 76, 72, 106, 103, 118, 109, 54, 83, 67, 80, 120, 57, 85, 82, 76, 90, 111, 97, 87, 83, 68, 88, 65, 113, 102, 109, 97, 88, 43, 122, 70, 119, 83, 80, 71, 86, 47, 72, 88, 109, 114, 10, 88, 90, 72, 74, 114, 72, 54, 79, 53, 67, 54, 53, 71, 70, 119, 56, 113, 50, 122, 110, 101, 66, 112, 106, 114, 86, 56, 115, 56, 48, 68, 52, 107, 119, 89, 68, 72, 82, 108, 77, 87, 86, 113, 103, 87, 99, 88, 100, 88, 57, 110, 120, 89, 104, 85, 121, 80, 69, 112, 67, 102, 57, 76, 112, 10, 83, 116, 97, 53, 97, 81, 83, 78, 49, 108, 111, 102, 84, 90, 103, 68, 77, 111, 118, 89, 72, 111, 83, 103, 75, 79, 87, 115, 88, 50, 66, 121, 120, 65, 102, 82, 110, 119, 69, 61, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10}

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
