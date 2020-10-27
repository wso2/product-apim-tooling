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
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Compose the name of the tenant specific directory to save migration artifacts
// Actual tenant name cmdResourceTenantDomain will be changed to replace '.' with '-dot-'
// If the -t option is not given in with the command, a default name 'tenant-default' will be passed
func GetMigrationExportTenantDirName(cmdResourceTenantDomain string) (resourceTenantDirName string) {
	if cmdResourceTenantDomain == "" {
		resourceTenantDirName = DefaultResourceTenantDomain
	} else {
		resourceTenantDirName = cmdResourceTenantDomain
	}

	if strings.Contains(cmdResourceTenantDomain, ".") {
		resourceTenantDirName = strings.Replace(cmdResourceTenantDomain, ".", "-dot-", -1)
	}
	return resourceTenantDirName
}

// Read the details of finally and successfully exported API into the last-succeeded-api.log file
func ReadLastSucceededAPIFileData(exportRelatedFilesPath string) API {
	var lastSucceededApiFilePath = filepath.Join(exportRelatedFilesPath, LastSucceededApiFileName)
	data, err := ioutil.ReadFile(lastSucceededApiFilePath)
	str := string(data)
	var splittedString = strings.Split(str, " ")
	var api = API{"", strings.TrimSpace(splittedString[0]), "", strings.TrimSpace(splittedString[1]), strings.TrimSpace(splittedString[2]), ""}

	if err != nil {
		HandleErrorAndExit("Error in reading file "+lastSucceededApiFilePath, err)
	}
	return api
}

// Write the last-succeeded-api.log file. It includes the meta data of the API, which was successfully exported finally
func WriteLastSuceededAPIFileData(exportRelatedFilesPath string, api API) {
	var lastSucceededApiFilePath = filepath.Join(exportRelatedFilesPath, LastSucceededApiFileName)
	var content []byte
	content = []byte(api.Name + LastSuceededContentDelimiter + api.Version + LastSuceededContentDelimiter + api.Provider)
	var error = ioutil.WriteFile(lastSucceededApiFilePath, content, 0644)

	if error != nil {
		HandleErrorAndExit("Error in writing file "+lastSucceededApiFilePath, error)
	}
}

// Read the migration-apis-export-metadata.yaml file
func (migrationApisExportMetadata *MigrationApisExportMetadata) ReadMigrationApisExportMetadataFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, migrationApisExportMetadata); err != nil {
		return err
	}
	return nil
}

// Write the migration-apis-export-metadata.yaml file. This includes the below meta data of the API export process,
// including the list of APIs exported at each iteration. (APIs are exported as multiple iterations)
// api_list_offset => offset index of list of APIs fetched from APIM server at the perticular iteration
// user => username of the user that executes the operation
// on_tenant => which tenant's APIs are exported
func WriteMigrationApisExportMetadataFile(apis []API, cmdResourceTenantDomain string,
	cmdUsername string, exportRelatedFilesPath string, apiListOffset int) {
	var exportMetaData = new(MigrationApisExportMetadata)
	exportMetaData.ApiListOffset = apiListOffset
	exportMetaData.ApiListToExport = apis
	exportMetaData.OnTenant = cmdResourceTenantDomain
	exportMetaData.User = cmdUsername

	WriteConfigFile(exportMetaData, filepath.Join(exportRelatedFilesPath, MigrationAPIsExportMetadataFileName))
}
