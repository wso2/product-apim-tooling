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
	"path/filepath"
)

const ProjectName = "apictl"

// File Names and Paths
var CurrentDir, _ = os.Getwd()

const ConfigDirName = ".wso2apictl"

var HomeDirectory = os.Getenv("HOME")

var ConfigDirPath = filepath.Join(HomeDirectory, ConfigDirName)

const EnvKeysAllFileName = "env_keys_all.yaml"

var EnvKeysAllFilePath = filepath.Join(ConfigDirPath, EnvKeysAllFileName)

const MainConfigFileName = "main_config.yaml"
const SampleMainConfigFileName = "main_config.yaml.sample"
const DefaultAPISpecFileName = "default_api.yaml"

var MainConfigFilePath = filepath.Join(ConfigDirPath, MainConfigFileName)
var SampleMainConfigFilePath = filepath.Join(ConfigDirPath, SampleMainConfigFileName)
var DefaultAPISpecFilePath = filepath.Join(ConfigDirPath, DefaultAPISpecFileName)

const DefaultExportDirName = "exported"
const ExportedApisDirName = "apis"
const ExportedAppsDirName = "apps"
const ExportedMigrationArtifactsDirName = "migration"

var DefaultExportDirPath = filepath.Join(ConfigDirPath, DefaultExportDirName)

const defaultApiApplicationImportExportSuffix = "api/am/admin/v1"
const defaultApiListEndpointSuffix = "api/am/publisher/v1/apis"
const defaultApiProductListEndpointSuffix = "api/am/publisher/v1/api-products"
const defaultUnifiedSearchEndpointSuffix = "api/am/publisher/v1/search"
const defaultAdminApplicationListEndpointSuffix = "api/am/admin/v1/applications"
const defaultDevPortalApplicationListEndpointSuffix = "api/am/store/v1/applications"
const defaultDevPortalThrottlingPoliciesEndpointSuffix = "api/am/store/v1/throttling-policies"
const defaultClientRegistrationEndpointSuffix = "client-registration/v0.16/register"

const DefaultEnvironmentName = "default"

// API Product related constants
const DefaultApiProductVersion = "1.0.0"
const DefaultApiProductType = "APIProduct"

// WSO2PublicCertificate : wso2 public certificate
var WSO2PublicCertificate = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10, 77, 73, 73, 68, 113, 84, 67, 67, 65, 112, 71, 103, 65, 119, 73, 66, 65, 103, 73, 69, 88, 98, 65, 66, 111, 122, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 115, 70, 65, 68, 66, 107, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 71, 69, 119, 74, 86, 13, 10, 85, 122, 69, 76, 77, 65, 107, 71, 65, 49, 85, 69, 67, 65, 119, 67, 81, 48, 69, 120, 70, 106, 65, 85, 66, 103, 78, 86, 66, 65, 99, 77, 68, 85, 49, 118, 100, 87, 53, 48, 89, 87, 108, 117, 73, 70, 90, 112, 90, 88, 99, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 111, 77, 13, 10, 66, 70, 100, 84, 84, 122, 73, 120, 68, 84, 65, 76, 66, 103, 78, 86, 66, 65, 115, 77, 66, 70, 100, 84, 84, 122, 73, 120, 69, 106, 65, 81, 66, 103, 78, 86, 66, 65, 77, 77, 67, 87, 120, 118, 89, 50, 70, 115, 97, 71, 57, 122, 100, 68, 65, 101, 70, 119, 48, 120, 79, 84, 69, 119, 13, 10, 77, 106, 77, 119, 78, 122, 77, 119, 78, 68, 78, 97, 70, 119, 48, 121, 77, 106, 65, 120, 77, 106, 85, 119, 78, 122, 77, 119, 78, 68, 78, 97, 77, 71, 81, 120, 67, 122, 65, 74, 66, 103, 78, 86, 66, 65, 89, 84, 65, 108, 86, 84, 77, 81, 115, 119, 67, 81, 89, 68, 86, 81, 81, 73, 13, 10, 68, 65, 74, 68, 81, 84, 69, 87, 77, 66, 81, 71, 65, 49, 85, 69, 66, 119, 119, 78, 84, 87, 57, 49, 98, 110, 82, 104, 97, 87, 52, 103, 86, 109, 108, 108, 100, 122, 69, 78, 77, 65, 115, 71, 65, 49, 85, 69, 67, 103, 119, 69, 86, 49, 78, 80, 77, 106, 69, 78, 77, 65, 115, 71, 13, 10, 65, 49, 85, 69, 67, 119, 119, 69, 86, 49, 78, 80, 77, 106, 69, 83, 77, 66, 65, 71, 65, 49, 85, 69, 65, 119, 119, 74, 98, 71, 57, 106, 89, 87, 120, 111, 98, 51, 78, 48, 77, 73, 73, 66, 73, 106, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 69, 70, 13, 10, 65, 65, 79, 67, 65, 81, 56, 65, 77, 73, 73, 66, 67, 103, 75, 67, 65, 81, 69, 65, 120, 101, 113, 111, 90, 89, 98, 81, 47, 83, 114, 56, 68, 79, 70, 81, 43, 47, 113, 98, 69, 98, 67, 112, 54, 86, 122, 98, 53, 104, 122, 72, 55, 111, 97, 51, 104, 102, 50, 70, 90, 120, 82, 75, 13, 10, 70, 48, 72, 54, 98, 56, 67, 79, 77, 122, 122, 56, 43, 48, 109, 118, 69, 100, 89, 86, 118, 98, 47, 51, 49, 106, 77, 69, 76, 50, 67, 73, 81, 104, 107, 81, 82, 111, 108, 49, 73, 114, 117, 68, 54, 110, 66, 79, 109, 107, 106, 117, 88, 74, 83, 66, 102, 105, 99, 107, 108, 77, 97, 74, 13, 10, 90, 79, 82, 104, 117, 67, 114, 66, 52, 114, 111, 72, 120, 122, 111, 71, 49, 57, 97, 87, 109, 115, 99, 65, 48, 103, 110, 102, 66, 75, 111, 50, 111, 71, 88, 83, 106, 74, 109, 110, 90, 120, 73, 104, 43, 50, 88, 54, 115, 121, 72, 67, 102, 121, 77, 90, 90, 48, 48, 76, 122, 68, 121, 114, 13, 10, 103, 111, 88, 87, 81, 88, 121, 70, 118, 67, 65, 50, 97, 120, 53, 52, 115, 55, 115, 75, 105, 72, 79, 77, 51, 80, 52, 65, 57, 87, 52, 81, 85, 119, 109, 111, 69, 105, 52, 72, 81, 109, 80, 103, 74, 106, 73, 77, 52, 101, 71, 86, 80, 104, 48, 71, 116, 73, 65, 78, 78, 43, 66, 79, 13, 10, 81, 49, 75, 107, 85, 73, 55, 79, 122, 116, 101, 72, 67, 84, 76, 117, 51, 86, 106, 120, 77, 48, 115, 119, 56, 81, 82, 97, 121, 90, 100, 104, 110, 105, 80, 70, 43, 85, 57, 110, 51, 102, 97, 49, 109, 79, 52, 75, 76, 66, 115, 87, 52, 109, 68, 76, 106, 103, 56, 82, 47, 74, 117, 65, 13, 10, 71, 84, 88, 47, 83, 69, 69, 71, 106, 48, 66, 53, 72, 87, 81, 65, 80, 54, 109, 121, 120, 75, 70, 122, 50, 120, 119, 68, 97, 67, 71, 118, 84, 43, 114, 100, 118, 107, 107, 116, 79, 119, 73, 68, 65, 81, 65, 66, 111, 50, 77, 119, 89, 84, 65, 85, 66, 103, 78, 86, 72, 82, 69, 69, 13, 10, 68, 84, 65, 76, 103, 103, 108, 115, 98, 50, 78, 104, 98, 71, 104, 118, 99, 51, 81, 119, 72, 81, 89, 68, 86, 82, 48, 79, 66, 66, 89, 69, 70, 69, 68, 112, 76, 66, 52, 80, 68, 103, 122, 115, 100, 120, 68, 50, 70, 86, 51, 114, 86, 110, 79, 114, 47, 65, 48, 68, 77, 66, 48, 71, 13, 10, 65, 49, 85, 100, 74, 81, 81, 87, 77, 66, 81, 71, 67, 67, 115, 71, 65, 81, 85, 70, 66, 119, 77, 66, 66, 103, 103, 114, 66, 103, 69, 70, 66, 81, 99, 68, 65, 106, 65, 76, 66, 103, 78, 86, 72, 81, 56, 69, 66, 65, 77, 67, 66, 80, 65, 119, 68, 81, 89, 74, 75, 111, 90, 73, 13, 10, 104, 118, 99, 78, 65, 81, 69, 76, 66, 81, 65, 68, 103, 103, 69, 66, 65, 69, 56, 72, 47, 97, 120, 65, 103, 88, 106, 116, 57, 51, 72, 71, 67, 89, 71, 117, 109, 85, 76, 87, 50, 108, 75, 107, 103, 113, 69, 118, 88, 114, 121, 80, 50, 81, 107, 82, 112, 98, 121, 81, 83, 115, 84, 89, 13, 10, 99, 76, 55, 90, 76, 83, 86, 66, 55, 77, 86, 86, 72, 116, 73, 115, 72, 104, 56, 102, 49, 67, 52, 88, 113, 54, 81, 117, 56, 78, 85, 114, 113, 117, 53, 90, 76, 67, 49, 112, 85, 66, 121, 97, 113, 82, 50, 90, 73, 122, 99, 106, 47, 79, 87, 76, 71, 89, 82, 106, 83, 84, 72, 83, 13, 10, 86, 109, 86, 73, 113, 57, 81, 113, 66, 113, 49, 106, 55, 114, 54, 102, 51, 66, 87, 113, 97, 79, 73, 105, 107, 110, 109, 84, 122, 69, 117, 113, 73, 86, 108, 79, 84, 89, 48, 103, 79, 43, 83, 72, 100, 83, 54, 50, 118, 114, 50, 70, 67, 122, 52, 121, 79, 114, 66, 69, 117, 108, 71, 65, 13, 10, 118, 111, 109, 115, 85, 56, 115, 113, 103, 52, 80, 104, 70, 110, 107, 104, 120, 73, 52, 77, 57, 49, 50, 76, 121, 43, 50, 82, 103, 78, 57, 76, 55, 65, 107, 104, 122, 75, 43, 69, 122, 88, 89, 49, 47, 81, 116, 108, 73, 47, 86, 121, 115, 78, 102, 83, 54, 122, 114, 72, 97, 115, 75, 122, 13, 10, 54, 67, 114, 75, 75, 67, 71, 113, 81, 110, 66, 110, 83, 118, 83, 84, 121, 70, 57, 79, 82, 53, 75, 70, 72, 110, 107, 65, 119, 69, 57, 57, 53, 73, 90, 114, 99, 83, 81, 105, 99, 77, 120, 115, 76, 104, 84, 77, 85, 72, 68, 76, 81, 47, 103, 82, 121, 121, 55, 86, 47, 90, 112, 68, 13, 10, 77, 102, 65, 87, 82, 43, 53, 79, 101, 81, 105, 78, 65, 112, 47, 98, 71, 52, 102, 106, 74, 111, 84, 100, 111, 113, 107, 117, 108, 53, 49, 43, 50, 98, 72, 72, 86, 114, 85, 61, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10}

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

// Logging Prefixes
const LogPrefixInfo = "[INFO]: "
const LogPrefixWarning = "[WARN]: "
const LogPrefixError = "[ERROR]: "

// String Constants
const SearchAndTag = "&"

// Regex Validation
const UsernameValidRegex = `^[\w\d\-]*$`
const PositiveNoValidRegex = `^[1-9]\d*$`
const UrlValidRegex = `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`

// Other
const DefaultTokenValidityPeriod = 3600
const DefaultHttpRequestTimeout = 10000

// Migration export
const MaxAPIsToExportOnce = 20
const MigrationAPIsExportMetadataFileName = "migration-apis-export-metadata.yaml"
const LastSucceededApiFileName = "last-succeeded-api.log"
const LastSuceededContentDelimiter = " " // space
const DefaultResourceTenantDomain = "tenant-default"
const ApplicationId = "applicationId"
const ApiId = "apiId"
const DefaultCliApp = "default-apictl-app"
const DefaultTokenType = "JWT"

var ValidInitialStates = []string{"CREATED", "PUBLISHED"}

var EnvReplaceFilePaths = []string{
	"Docs" + string(os.PathSeparator) + "docs.yaml",
	"Docs" + string(os.PathSeparator) + "InlineContents",
	"Meta-information",
	"WSDL",
	"Sequences",
	"SoapToRest",
}

const PrivateJetModeConst = "privateJet"
const SidecarModeConst = "sidecar"
