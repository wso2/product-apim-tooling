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

import "os"

const ProjectName string = "wso2apim-cli"

// File Names and Paths
const PathSeparator_ string = string(os.PathSeparator)
const ApplicationRoot string = "/home/menuka/.go/src/github" +
	".com/wso2/product-apim-tooling/import-export-cli" // TODO:: Change to a generic root
const EnvKeysAllFileName string = "env_keys_all.yaml"
const EnvKeysAllFilePath string = ApplicationRoot + PathSeparator_ + EnvKeysAllFileName
const MainConfigFileName string = "main_config.yaml"
const MainConfigFilePath string = ApplicationRoot + PathSeparator_ + MainConfigFileName
const ExportedAPIsDirectoryName string = "exported"
const ExportedAPIsDirectoryPath string = ApplicationRoot + PathSeparator_ + ExportedAPIsDirectoryName

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
