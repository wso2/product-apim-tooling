package utils

import (
	"github.com/jhoonb/archivex"
	"os"
	"github.com/wso2/wum-client/utils"
)

func ZipDir(source, target string) (error){
	err := os.Chdir(source)
	if err == nil {
		utils.Logln(LogPrefixInfo + "Directory " + source + " exists")
		utils.Logln(LogPrefixInfo + "Starting Compression...")
		zip := new(archivex.ZipFile)
		zip.Create(target)
		zip.AddAll(source, true)
		zip.Close()
		utils.Logln(LogPrefixInfo + "Compression completed: Find file " + target)
		return nil
	}else{
		utils.Logln(LogPrefixError + "Compressing " + source)
		return err
	}
}
