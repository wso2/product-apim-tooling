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

const ProjectName = "apimcli"
const ConfigDirectoryName = "wso2-apim-cli-config"

// File Names and Paths

var CurrentDir, _ = os.Getwd()
var ApplicationRoot = CurrentDir

var PathSeparator_ = string(os.PathSeparator)

//var ApplicationRoot = "/home/menuka/.go/src/github.com/wso2/product-apim-tooling/import-export-cli"

const EnvKeysAllFileName = "env_keys_all.yaml"

var EnvKeysAllFilePath = filepath.Join(ApplicationRoot, EnvKeysAllFileName)

const MainConfigFileName = "main_config.yaml"

var MainConfigFilePath = filepath.Join(ApplicationRoot, MainConfigFileName)

const ExportedAPIsDirectoryName = "exported"

var ExportedAPIsDirectoryPath = filepath.Join(ApplicationRoot, ExportedAPIsDirectoryName)

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
const LogPrefixInfo = "[INFO] "
const LogPrefixWarning = "[WARN] "
const LogPrefixError = "[ERROR] "

// Other
const DefaultTokenValidityPeriod = "3600"
const DefaultHttpRequestTimeout = 100000
