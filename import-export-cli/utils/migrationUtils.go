package utils

import (
	"strings"
	"io/ioutil"
	"path/filepath"
	"strconv"
)


func GetMigrationExportTenantDirName(cmdResourceTenantDomain string) (resourceTenantDirName string){
	if (cmdResourceTenantDomain == "" ){
		/*if(strings.Contains(cmdUsername, "@") ){
		   //get tenant domain by splitting the username
		   //Ok to get as this? Email user name will conflict this
		   //cmdResourceTenantDomain = strings.Split(cmdUsername, "@")[1]
	   } else {*/
		// if username doesn't contain '@' decide the tenant as 'carbon.super'
		// Only super admin can avoid passing the tenant domain. Other tenant admins must pass tenant domain with -t
		resourceTenantDirName = "carbon-dot-super"
		//}
	} else {
		resourceTenantDirName = cmdResourceTenantDomain;
	}

	if(strings.Contains(cmdResourceTenantDomain,".")) {
		resourceTenantDirName = strings.Replace(cmdResourceTenantDomain, ".", "-dot-", -1)
	}
	return resourceTenantDirName
}

func ReadLastSuceededAPIFileData(exportRelatedFilesPath string) (int, API, bool) {
	var lastSucceededApiFilePath = filepath.Join(exportRelatedFilesPath,LastSucceededApiFileName)
	data, err := ioutil.ReadFile(lastSucceededApiFilePath)
	str := string(data)
	var splittedString[] string = strings.Split(str," ")
	var api  = API{"", splittedString[1], "", splittedString[2],splittedString[3], "" }

	if err != nil {
		HandleErrorAndExit("MainConfig: File Not Found: "+lastSucceededApiFilePath, err)
	}

	var iterationNo,_ = strconv.Atoi(splittedString[0])
	var migrationExportProcessCompleted bool
	if migrationExportProcessCompleted = false ;splittedString[4] == "COMPLETED" {
		migrationExportProcessCompleted = true
	}
	return iterationNo, api, migrationExportProcessCompleted
}