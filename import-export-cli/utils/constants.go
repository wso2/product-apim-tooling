/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"path/filepath"
)

const ProjectName string = "wso2apim"

// File Names and Paths

//var ApplicationRoot, _ = os.Getwd()
var ApplicationRoot  = "/home/menuka/.go/src/github.com/wso2/product-apim-tooling/import-export-cli"

const EnvKeysAllFileName string = "env_keys_all.yaml"

var EnvKeysAllFilePath string = filepath.Join(ApplicationRoot, EnvKeysAllFileName)

const MainConfigFileName string = "main_config.yaml"

var MainConfigFilePath string = filepath.Join(ApplicationRoot, MainConfigFileName)

const ExportedAPIsDirectoryName string = "exported"

var ExportedAPIsDirectoryPath string = filepath.Join(ApplicationRoot, ExportedAPIsDirectoryName)

const DefaultEnvironmentName string = "default"

// Headers and Header Values
const HeaderAuthorization string = "Authorization"
const HeaderContentType string = "Content-Type"
const HeaderConnection string = "Connection"
const HeaderAccept string = "Accept"
const HeaderProduces string = "Produces"
const HeaderConsumes string = "Consumes"
const HeaderContentEncoding string = "Content-Encoding"
const HeaderTransferEncoding string = "transfer-encoding"
const HeaderValueChunked string = "chunked"
const HeaderValueGZIP string = "gzip"
const HeaderValueKeepAlive string = "keep-alive"
const HeaderValueApplicationZip = "application/zip"
const HeaderValueApplicationJSON string = "application/json"
const HeaderValueXWWWFormUrlEncoded string = "application/x-www-form-urlencoded"
const HeaderValueAuthBearerPrefix string = "Bearer"
const HeaderValueAuthBasicPrefix string = "Basic"
const HeaderValueMultiPartFormData string = "multipart/form-data"

// Logging Prefixes
const LogPrefixInfo = "[INFO]: "
const LogPrefixWarning = "[WARN]: "
const LogPrefixError = "[ERROR]: "

// Other
const DefaultTokenValidityPeriod string = "3600"
const DefaultHttpRequestTimeout int = 10000
