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
var WSO2PublicCertificate = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10, 77, 73, 73, 68, 113, 84, 67, 67, 65, 112, 71, 103, 65, 119, 73, 66, 65, 103, 73, 69, 89, 102, 69, 86, 83, 106, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 115, 70, 65, 68, 66, 107, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 71, 69, 119, 74, 86, 10, 85, 122, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 67, 65, 119, 67, 81, 48, 69, 120, 70, 106, 65, 85, 66, 103, 78, 86, 66, 65, 99, 77, 68, 85, 49, 118, 100, 87, 53, 48, 89, 87, 108, 117, 73, 70, 90, 112, 90, 88, 99, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 111, 77, 10, 66, 70, 100, 84, 84, 122, 73, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 115, 77, 66, 70, 100, 84, 84, 122, 73, 120, 69, 106, 65, 81, 66, 103, 78, 86, 66, 65, 77, 77, 67, 87, 120, 118, 89, 50, 70, 115, 97, 71, 57, 122, 100, 68, 65, 101, 70, 119, 48, 121, 77, 106, 65, 120, 10, 77, 106, 89, 119, 79, 84, 77, 121, 78, 84, 104, 97, 70, 119, 48, 121, 78, 68, 65, 48, 77, 122, 65, 119, 79, 84, 77, 121, 78, 84, 104, 97, 77, 71, 81, 120, 67, 122, 65, 74, 66, 103, 78, 86, 66, 65, 89, 84, 65, 108, 86, 84, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 73, 10, 68, 65, 74, 68, 81, 84, 69, 87, 77, 66, 81, 71, 65, 49, 85, 69, 66, 119, 119, 78, 84, 87, 57, 49, 98, 110, 82, 104, 97, 87, 52, 103, 86, 109, 108, 108, 100, 122, 69, 78, 77, 65, 115, 71, 65, 49, 85, 69, 67, 103, 119, 69, 86, 49, 78, 80, 77, 106, 69, 78, 77, 65, 115, 71, 10, 65, 49, 85, 69, 67, 119, 119, 69, 86, 49, 78, 80, 77, 106, 69, 83, 77, 66, 65, 71, 65, 49, 85, 69, 65, 119, 119, 74, 98, 71, 57, 106, 89, 87, 120, 111, 98, 51, 78, 48, 77, 73, 73, 66, 73, 106, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 69, 70, 10, 65, 65, 79, 67, 65, 81, 56, 65, 77, 73, 73, 66, 67, 103, 75, 67, 65, 81, 69, 65, 107, 100, 103, 110, 99, 111, 67, 114, 122, 54, 53, 53, 76, 113, 56, 112, 84, 100, 88, 48, 55, 101, 111, 86, 66, 106, 100, 90, 68, 67, 85, 69, 54, 117, 101, 66, 100, 48, 68, 49, 104, 112, 74, 48, 10, 47, 122, 69, 51, 120, 51, 65, 122, 54, 116, 108, 118, 122, 115, 57, 56, 80, 115, 80, 117, 71, 122, 97, 81, 79, 77, 109, 117, 76, 97, 52, 113, 120, 78, 74, 43, 79, 75, 120, 74, 109, 117, 116, 68, 85, 108, 67, 108, 112, 117, 118, 120, 117, 102, 43, 106, 121, 113, 52, 103, 67, 86, 53, 116, 10, 69, 73, 73, 76, 87, 82, 77, 66, 106, 108, 66, 69, 112, 74, 102, 87, 109, 54, 51, 43, 86, 75, 75, 85, 52, 110, 118, 66, 87, 78, 74, 55, 75, 102, 104, 87, 106, 108, 56, 43, 68, 85, 100, 78, 83, 104, 50, 112, 67, 68, 76, 112, 85, 79, 98, 109, 98, 57, 75, 113, 117, 113, 99, 49, 10, 120, 52, 66, 103, 116, 116, 106, 78, 52, 114, 120, 47, 80, 43, 51, 47, 118, 43, 49, 106, 69, 84, 88, 122, 73, 80, 49, 76, 52, 52, 121, 72, 116, 112, 81, 78, 118, 48, 107, 104, 89, 102, 52, 106, 47, 97, 72, 106, 99, 69, 114, 105, 57, 121, 107, 118, 112, 122, 49, 109, 116, 100, 97, 99, 10, 98, 114, 75, 75, 50, 53, 78, 52, 86, 49, 72, 72, 82, 119, 68, 113, 90, 105, 74, 122, 79, 67, 67, 73, 83, 88, 68, 117, 113, 66, 54, 119, 103, 117, 89, 47, 118, 52, 110, 48, 108, 49, 88, 116, 114, 69, 115, 55, 105, 67, 121, 102, 82, 70, 119, 78, 83, 75, 78, 114, 76, 113, 114, 50, 10, 51, 116, 82, 49, 67, 115, 99, 109, 76, 102, 98, 72, 54, 90, 76, 103, 53, 67, 89, 74, 84, 68, 43, 49, 117, 80, 83, 120, 48, 72, 77, 79, 66, 52, 87, 118, 53, 49, 80, 98, 87, 119, 73, 68, 65, 81, 65, 66, 111, 50, 77, 119, 89, 84, 65, 85, 66, 103, 78, 86, 72, 82, 69, 69, 10, 68, 84, 65, 76, 103, 103, 108, 115, 98, 50, 78, 104, 98, 71, 104, 118, 99, 51, 81, 119, 72, 81, 89, 68, 86, 82, 48, 79, 66, 66, 89, 69, 70, 72, 48, 75, 81, 51, 89, 84, 90, 74, 120, 84, 115, 78, 115, 80, 121, 114, 90, 79, 83, 70, 103, 88, 88, 104, 71, 43, 77, 66, 48, 71, 10, 65, 49, 85, 100, 74, 81, 81, 87, 77, 66, 81, 71, 67, 67, 115, 71, 65, 81, 85, 70, 66, 119, 77, 66, 66, 103, 103, 114, 66, 103, 69, 70, 66, 81, 99, 68, 65, 106, 65, 76, 66, 103, 78, 86, 72, 81, 56, 69, 66, 65, 77, 67, 66, 80, 65, 119, 68, 81, 89, 74, 75, 111, 90, 73, 10, 104, 118, 99, 78, 65, 81, 69, 76, 66, 81, 65, 68, 103, 103, 69, 66, 65, 70, 78, 74, 51, 52, 67, 73, 105, 73, 108, 67, 120, 109, 121, 112, 50, 55, 43, 75, 65, 50, 50, 52, 76, 97, 72, 86, 116, 76, 53, 68, 117, 99, 70, 75, 48, 80, 50, 50, 70, 81, 43, 81, 75, 107, 79, 78, 10, 105, 85, 119, 79, 55, 48, 75, 111, 86, 70, 114, 101, 66, 72, 49, 83, 109, 120, 117, 52, 101, 80, 87, 107, 54, 114, 77, 90, 70, 79, 77, 53, 111, 76, 56, 72, 88, 89, 103, 51, 116, 119, 121, 43, 53, 101, 71, 99, 76, 51, 80, 81, 100, 55, 88, 53, 100, 119, 65, 113, 108, 86, 105, 118, 10, 122, 111, 107, 111, 105, 54, 83, 68, 97, 65, 47, 98, 73, 71, 54, 74, 47, 79, 49, 85, 57, 81, 100, 52, 88, 69, 86, 74, 100, 86, 117, 76, 113, 106, 107, 49, 43, 99, 112, 55, 48, 65, 76, 116, 48, 88, 54, 66, 55, 115, 78, 76, 102, 106, 70, 99, 98, 122, 51, 106, 81, 85, 76, 78, 10, 110, 75, 56, 72, 78, 118, 113, 98, 110, 55, 122, 81, 117, 80, 49, 48, 115, 56, 112, 53, 121, 50, 113, 86, 107, 80, 66, 65, 47, 112, 106, 105, 103, 82, 68, 115, 73, 87, 82, 54, 112, 55, 56, 81, 69, 83, 70, 43, 84, 97, 72, 70, 106, 120, 102, 99, 68, 54, 102, 57, 99, 110, 89, 105, 10, 101, 43, 121, 69, 72, 69, 82, 116, 71, 56, 107, 56, 120, 53, 106, 76, 70, 101, 43, 111, 100, 73, 49, 47, 81, 71, 90, 80, 56, 70, 121, 48, 111, 75, 84, 43, 69, 47, 84, 74, 49, 70, 66, 104, 52, 114, 66, 49, 70, 116, 75, 121, 108, 113, 71, 101, 97, 117, 80, 117, 56, 57, 68, 110, 10, 97, 74, 57, 43, 107, 118, 112, 78, 81, 57, 52, 121, 70, 109, 69, 117, 104, 116, 68, 66, 121, 118, 68, 105, 106, 120, 65, 113, 118, 108, 105, 110, 51, 84, 80, 73, 102, 121, 56, 61, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10}

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

//Default values for Help commands
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
