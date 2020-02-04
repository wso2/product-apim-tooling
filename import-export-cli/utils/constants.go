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

const defaultApiApplicationImportExportSuffix = "api/am/admin/v0.15"
const defaultApiListEndpointSuffix = "api/am/publisher/v0.15/apis"
const defaultApplicationListEndpointSuffix = "api/am/admin/v0.15/applications"

const DefaultEnvironmentName = "default"

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

// Kubernetes Constants
const DefaultKubernetesMode = false
const Kubectl = "kubectl"
const Create = "create"
const K8sApply = "apply"
const K8sDelete = "delete"
const K8sRollOut = "rollout"
const K8sGet = "get"

// WSO2 API Operator constats
const ApiOpControllerConfigMap = "controller-config"
const ApiOpWso2Namespace = "wso2-system"

// Regex Validation
const UsernameValidationRegex = `^[\w\d\-]+$`
const UrlValidationRegex = `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`

// Operator Hub Constants
const OlmCrdUrlTemplate = "https://github.com/operator-framework/operator-lifecycle-manager/releases/download/%s/crds.yaml"
const OlmOlmUrlTemplate = "https://github.com/operator-framework/operator-lifecycle-manager/releases/download/%s/olm.yaml"
const OperatorYamlUrl = "https://operatorhub.io/install/api-operator.yaml"
const DockerRegistryUrl = "https://index.docker.io/v2/"
const OperatorCsv = "csv"
const OlmVersion = "0.13.0"

// Other
const DefaultTokenValidityPeriod = "3600"
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
