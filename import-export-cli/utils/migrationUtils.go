package utils

import (
	"strings"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

func GetMigrationExportTenantDirName(cmdResourceTenantDomain string) (resourceTenantDirName string) {
	if (cmdResourceTenantDomain == "") {
		resourceTenantDirName = DefaultResourceTenantDomain
	} else {
		resourceTenantDirName = cmdResourceTenantDomain;
	}

	if (strings.Contains(cmdResourceTenantDomain, ".")) {
		resourceTenantDirName = strings.Replace(cmdResourceTenantDomain, ".", "-dot-", -1)
	}
	return resourceTenantDirName
}

func ReadLastSucceededAPIFileData(exportRelatedFilesPath string) (API) {
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

func WriteLastSuceededAPIFileData(exportRelatedFilesPath string, api API) {
	var lastSucceededApiFilePath = filepath.Join(exportRelatedFilesPath, LastSucceededApiFileName)
	var content []byte
	content = []byte(api.Name + LastSuceededContentDelimiter + api.Version + LastSuceededContentDelimiter + api.Provider)
	var error = ioutil.WriteFile(lastSucceededApiFilePath, content, 0644)

	if (error != nil) {
		HandleErrorAndExit("Error in writing file "+lastSucceededApiFilePath, error)
	}
}

func (migrationApisExportMetadata *MigrationApisExportMetadata) ReadMigrationApisExportMetadataFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		HandleErrorAndExit("migration-apis-export-metadata.yaml: File Not Found: "+filePath, err)
	}
	if err := yaml.Unmarshal(data, migrationApisExportMetadata); err != nil {
		return err
	}
	return nil
}

func WriteMigrationApisExportMetadataFile(apis []API, cmdResourceTenantDomain string,
	cmdUsername string, exportRelatedFilesPath string, apiListOffset int) {
	var exportMetaData = new(MigrationApisExportMetadata)
	exportMetaData.ApiListOffset = apiListOffset
	exportMetaData.ApiListToExport = apis
	exportMetaData.OnTenant = cmdResourceTenantDomain
	exportMetaData.User = cmdUsername

	WriteConfigFile(exportMetaData, filepath.Join(exportRelatedFilesPath, MigrationAPIsExportMetadataFileName))
}
