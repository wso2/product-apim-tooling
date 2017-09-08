package utils

import (
	"github.com/jhoonb/archivex"
	"os"
)

func ZipDir(source, target string) (error){
	err := os.Chdir(source)
	if err == nil {
		Logln(LogPrefixInfo + "Directory " + source + " exists")
		Logln(LogPrefixInfo + "Starting Compression...")
		zip := new(archivex.ZipFile)
		zip.Create(target)
		zip.AddAll(source, true)
		zip.Close()
		Logln(LogPrefixInfo + "Compression completed: Find file " + target)
		return nil
	}else{
		Logln(LogPrefixError + "Compressing " + source)
		return err
	}
}
