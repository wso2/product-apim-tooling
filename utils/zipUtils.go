package utils

import (
	"github.com/jhoonb/archivex"
	"os"
	"fmt"
)

func IsDirectory(input string) bool {
	return false
}

func ZipDir(source, target string) (error){
	err := os.Chdir(source)
	if err == nil {
		fmt.Println("Directory available...")
		zip := new(archivex.ZipFile)
		zip.Create(target)
		zip.AddAll(source, true)
		zip.Close()
		return nil
	}else{
		return err
	}
}
