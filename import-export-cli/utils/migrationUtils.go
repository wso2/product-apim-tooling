package utils

import (
	"strings"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"gopkg.in/yaml.v2"
)

func GetMigrationExportTenantDirName(cmdResourceTenantDomain string) (resourceTenantDirName string) {
	if (cmdResourceTenantDomain == "") {
		/*if(strings.Contains(cmdUsername, "@") ){
		   //get tenant domain by splitting the username
		   //Ok to get as this? Email user name will conflict this
		   //cmdResourceTenantDomain = strings.Split(cmdUsername, "@")[1]
	   } else {*/
		// if username doesn't contain '@' decide the tenant as 'carbon.super'
		// Only super admin can avoid passing the tenant domain. Other tenant admins must pass tenant domain with -t
		resourceTenantDirName = DefaultResourceTenantDomain
		//}
	} else {
		resourceTenantDirName = cmdResourceTenantDomain;
	}

	if (strings.Contains(cmdResourceTenantDomain, ".")) {
		resourceTenantDirName = strings.Replace(cmdResourceTenantDomain, ".", "-dot-", -1)
	}
	return resourceTenantDirName
}

func ReadLastSuceededAPIFileData(exportRelatedFilesPath string) (int, API) {
	var lastSucceededApiFilePath = filepath.Join(exportRelatedFilesPath, LastSucceededApiFileName)
	data, err := ioutil.ReadFile(lastSucceededApiFilePath)
	str := string(data)
	var splittedString = strings.Split(str, " ")
	var api = API{"", strings.TrimSpace(splittedString[1]), "", strings.TrimSpace(splittedString[2]), strings.TrimSpace(splittedString[3]), ""}

	if err != nil {
		HandleErrorAndExit("Error in reading file "+lastSucceededApiFilePath, err)
	}

	var iterationNo, _ = strconv.Atoi(splittedString[0])
	return iterationNo, api
}

func WriteLastSuceededAPIFileData(exportRelatedFilesPath string, iterationNo int, api API) {
	var lastSucceededApiFilePath = filepath.Join(exportRelatedFilesPath, LastSucceededApiFileName)
	var content []byte
	content = []byte(strconv.Itoa(iterationNo) + LastSuceededContentDelimiter + api.Name + LastSuceededContentDelimiter + api.Version + LastSuceededContentDelimiter + api.Provider)
	var error = ioutil.WriteFile(lastSucceededApiFilePath, content, 0644)

	if (error != nil) {
		HandleErrorAndExit("Error in writing file "+lastSucceededApiFilePath, error)
	}
}

func (migrationApisExportMetadata *MigrationApisExportMetadata) ReadMigrationApisExportMetadataFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		HandleErrorAndExit("MainConfig: File Not Found: "+filePath, err)
	}
	if err := yaml.Unmarshal(data, migrationApisExportMetadata); err != nil {
		return err
	}
	return nil
}

func WriteMigrationApisExportMetadataFile(apis []API, cmdResourceTenantDomain string,
	cmdUsername string, exportRelatedFilesPath string, iterationNo int) {
	var exportMetaData = new(MigrationApisExportMetadata)
	exportMetaData.IterationNo = iterationNo
	exportMetaData.ApiListToExport = apis
	exportMetaData.OnTenant = cmdResourceTenantDomain
	exportMetaData.User = cmdUsername

	WriteConfigFile(exportMetaData, filepath.Join(exportRelatedFilesPath, MigrationAPIsExportMetadataFileName))
}
