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

// File Names and Paths
var CurrentDir, _ = os.Getwd()

const ConfigDirName = ".wso2apictl"

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

var ConfigDirPath = filepath.Join(HomeDirectory, ConfigDirName)

const LocalCredentialsDirectoryName = ".wso2apictl.local"
const EnvKeysAllFileName = "env_keys_all.yaml"
const MainConfigFileName = "main_config.yaml"
const SampleMainConfigFileName = "main_config.yaml.sample"
const DefaultAPISpecFileName = "default_api.yaml"

var LocalCredentialsDirectoryPath = filepath.Join(HomeDirectory, LocalCredentialsDirectoryName)
var EnvKeysAllFilePath = filepath.Join(LocalCredentialsDirectoryPath, EnvKeysAllFileName)
var MainConfigFilePath = filepath.Join(ConfigDirPath, MainConfigFileName)
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

var DefaultExportDirPath = filepath.Join(ConfigDirPath, DefaultExportDirName)
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
var WSO2PublicCertificate = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 13, 10, 77, 73, 73, 68, 113, 84, 67, 67, 65, 112, 71, 103, 65, 119, 73, 66, 65, 103, 73, 69, 89, 47, 90, 97, 65, 122, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 115, 70, 65, 68, 66, 107, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 71, 69, 119, 74, 86, 13, 10, 85, 122, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 67, 65, 119, 67, 81, 48, 69, 120, 70, 106, 65, 85, 66, 103, 78, 86, 66, 65, 99, 77, 68, 85, 49, 118, 100, 87, 53, 48, 89, 87, 108, 117, 73, 70, 90, 112, 90, 88, 99, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 111, 77, 13, 10, 66, 70, 100, 84, 84, 122, 73, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 115, 77, 66, 70, 100, 84, 84, 122, 73, 120, 69, 106, 65, 81, 66, 103, 78, 86, 66, 65, 77, 77, 67, 87, 120, 118, 89, 50, 70, 115, 97, 71, 57, 122, 100, 68, 65, 101, 70, 119, 48, 121, 77, 122, 65, 121, 13, 10, 77, 106, 73, 120, 79, 68, 65, 52, 77, 68, 78, 97, 70, 119, 48, 121, 78, 84, 65, 49, 77, 106, 99, 120, 79, 68, 65, 52, 77, 68, 78, 97, 77, 71, 81, 120, 67, 122, 65, 74, 66, 103, 78, 86, 66, 65, 89, 84, 65, 108, 86, 84, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 73, 13, 10, 68, 65, 74, 68, 81, 84, 69, 87, 77, 66, 81, 71, 65, 49, 85, 69, 66, 119, 119, 78, 84, 87, 57, 49, 98, 110, 82, 104, 97, 87, 52, 103, 86, 109, 108, 108, 100, 122, 69, 78, 77, 65, 115, 71, 65, 49, 85, 69, 67, 103, 119, 69, 86, 49, 78, 80, 77, 106, 69, 78, 77, 65, 115, 71, 13, 10, 65, 49, 85, 69, 67, 119, 119, 69, 86, 49, 78, 80, 77, 106, 69, 83, 77, 66, 65, 71, 65, 49, 85, 69, 65, 119, 119, 74, 98, 71, 57, 106, 89, 87, 120, 111, 98, 51, 78, 48, 77, 73, 73, 66, 73, 106, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 69, 70, 13, 10, 65, 65, 79, 67, 65, 81, 56, 65, 77, 73, 73, 66, 67, 103, 75, 67, 65, 81, 69, 65, 122, 78, 113, 106, 102, 66, 52, 121, 112, 89, 48, 81, 107, 77, 52, 81, 110, 114, 116, 99, 70, 108, 82, 109, 70, 106, 100, 74, 89, 84, 122, 101, 84, 69, 122, 90, 106, 51, 80, 97, 74, 116, 78, 109, 13, 10, 99, 79, 67, 120, 86, 78, 119, 111, 109, 85, 119, 98, 107, 73, 83, 111, 103, 118, 52, 79, 52, 74, 48, 108, 68, 86, 65, 121, 113, 54, 97, 97, 112, 68, 78, 89, 57, 74, 122, 120, 115, 111, 76, 101, 104, 111, 115, 113, 121, 117, 75, 97, 114, 51, 73, 71, 83, 74, 104, 109, 56, 73, 77, 56, 13, 10, 78, 55, 117, 86, 102, 84, 48, 109, 76, 81, 43, 114, 104, 111, 51, 122, 88, 69, 55, 47, 70, 97, 104, 83, 43, 114, 119, 73, 112, 43, 79, 85, 80, 113, 74, 118, 82, 72, 56, 101, 110, 99, 50, 109, 112, 102, 103, 104, 71, 56, 99, 118, 66, 120, 52, 113, 113, 54, 86, 122, 77, 83, 51, 66, 13, 10, 55, 50, 67, 102, 78, 65, 80, 121, 69, 101, 70, 74, 119, 105, 52, 82, 52, 90, 103, 88, 122, 115, 108, 98, 114, 47, 111, 71, 77, 74, 66, 72, 83, 81, 68, 104, 86, 69, 111, 65, 56, 117, 107, 81, 122, 85, 115, 76, 97, 102, 98, 116, 51, 115, 70, 77, 68, 86, 121, 48, 116, 53, 75, 78, 13, 10, 83, 97, 122, 82, 66, 76, 99, 72, 83, 80, 108, 120, 53, 66, 48, 87, 52, 74, 83, 87, 102, 47, 118, 118, 49, 65, 47, 99, 49, 118, 57, 65, 74, 83, 107, 115, 120, 114, 83, 115, 82, 113, 82, 106, 111, 97, 72, 103, 112, 51, 65, 122, 104, 90, 119, 55, 76, 68, 115, 111, 119, 102, 120, 113, 13, 10, 72, 90, 53, 98, 70, 48, 84, 104, 120, 113, 52, 79, 88, 69, 79, 115, 107, 47, 114, 83, 80, 73, 72, 118, 82, 105, 108, 107, 55, 79, 80, 43, 101, 112, 84, 88, 119, 90, 112, 119, 52, 81, 73, 68, 65, 81, 65, 66, 111, 50, 77, 119, 89, 84, 65, 85, 66, 103, 78, 86, 72, 82, 69, 69, 13, 10, 68, 84, 65, 76, 103, 103, 108, 115, 98, 50, 78, 104, 98, 71, 104, 118, 99, 51, 81, 119, 72, 81, 89, 68, 86, 82, 48, 79, 66, 66, 89, 69, 70, 68, 121, 115, 106, 85, 85, 104, 103, 47, 56, 65, 115, 72, 54, 48, 114, 74, 118, 114, 65, 110, 57, 87, 97, 121, 49, 47, 77, 66, 48, 71, 13, 10, 65, 49, 85, 100, 74, 81, 81, 87, 77, 66, 81, 71, 67, 67, 115, 71, 65, 81, 85, 70, 66, 119, 77, 66, 66, 103, 103, 114, 66, 103, 69, 70, 66, 81, 99, 68, 65, 106, 65, 76, 66, 103, 78, 86, 72, 81, 56, 69, 66, 65, 77, 67, 66, 80, 65, 119, 68, 81, 89, 74, 75, 111, 90, 73, 13, 10, 104, 118, 99, 78, 65, 81, 69, 76, 66, 81, 65, 68, 103, 103, 69, 66, 65, 75, 82, 65, 78, 43, 68, 65, 51, 66, 72, 53, 80, 120, 86, 47, 48, 51, 83, 122, 76, 77, 110, 103, 120, 110, 47, 119, 56, 104, 77, 97, 98, 78, 88, 80, 53, 70, 112, 116, 118, 66, 114, 80, 100, 70, 90, 104, 13, 10, 54, 56, 75, 49, 118, 82, 43, 80, 75, 111, 112, 81, 66, 89, 73, 82, 78, 106, 111, 108, 122, 109, 114, 102, 76, 73, 110, 78, 122, 48, 113, 108, 117, 89, 103, 122, 56, 78, 108, 100, 111, 87, 43, 83, 75, 113, 52, 107, 106, 102, 71, 98, 78, 49, 82, 105, 74, 86, 87, 52, 104, 102, 89, 68, 13, 10, 98, 56, 85, 112, 82, 119, 111, 121, 87, 67, 121, 98, 80, 65, 67, 67, 74, 117, 122, 115, 105, 74, 86, 43, 86, 70, 88, 66, 117, 57, 104, 71, 107, 53, 49, 78, 90, 72, 119, 89, 54, 81, 103, 74, 68, 70, 107, 72, 99, 102, 51, 100, 79, 102, 106, 77, 69, 48, 53, 105, 121, 67, 87, 52, 13, 10, 99, 43, 120, 89, 87, 99, 74, 66, 55, 78, 113, 110, 48, 110, 51, 69, 71, 116, 75, 49, 68, 111, 105, 105, 69, 78, 112, 84, 67, 47, 52, 80, 75, 89, 49, 89, 71, 118, 69, 81, 88, 122, 68, 113, 73, 55, 51, 103, 74, 110, 103, 75, 78, 90, 55, 56, 115, 120, 78, 120, 48, 87, 97, 81, 13, 10, 73, 102, 78, 51, 51, 72, 111, 118, 48, 88, 116, 47, 99, 76, 77, 107, 49, 122, 122, 108, 74, 118, 68, 43, 78, 56, 82, 56, 110, 55, 67, 122, 83, 43, 106, 114, 49, 83, 110, 67, 111, 74, 109, 67, 90, 70, 74, 122, 70, 106, 113, 80, 78, 120, 122, 115, 84, 115, 102, 70, 84, 105, 114, 116, 13, 10, 106, 113, 121, 89, 114, 80, 87, 72, 104, 85, 56, 68, 57, 77, 99, 51, 105, 113, 116, 103, 97, 117, 82, 70, 52, 85, 50, 107, 112, 70, 115, 100, 112, 75, 112, 87, 79, 104, 81, 61, 13, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 13, 10}

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
