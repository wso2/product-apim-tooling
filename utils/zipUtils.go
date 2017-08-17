package utils

import (
	"github.com/jhoonb/archivex"
)

func IsDirectory(input string) bool{
	return false
}

func ZipDir(source, target string) {
	zip := new(archivex.ZipFile)
	zip.Create(target)
	zip.AddAll(source, true)
	zip.Close()
}


