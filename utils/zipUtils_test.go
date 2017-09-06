package utils

import (
	"testing"
	"os"
)

func TestZipDirError(t *testing.T) {
	err := ZipDir("", "")
	if err == nil {
		t.Errorf("ZipDir() didn't return an error for invalid source and destination")
	}
}

func TestZipDirOK(t *testing.T) {
	directoryName := "wso2apimZipTest"
	workingDir, _ := os.Getwd()

	directoryPath := workingDir + "/" + directoryName
	fileName := "test.txt"
	filePath := directoryPath + "/" + fileName

	os.Mkdir(directoryPath, os.ModePerm)

	// check if directory exists
	var _, err = os.Stat(directoryPath)
	if err != nil {
		t.Errorf("Error opening directory")
	}

	// create directory if it doesn't already exist
	if os.IsNotExist(err) {
		var file, err = os.Create(directoryPath)
		if err != nil {
			t.Errorf("Error creating sample directory for compressing: %s\n", err)
		}

		defer file.Close()
	}

	// check if file exists
	_, err = os.Stat(filePath)

	// create file if it doesn't already exist
	if os.IsNotExist(err) {
		var file, err = os.Create(filePath)
		if err != nil {
			t.Errorf("Error creating sample file for compressing: %s\n", err)
		}
		defer file.Close()
	}

	// Open file using READ & WRITE permissions
	var file, err1 = os.OpenFile(filePath, os.O_RDWR, 0644)
	if err1 != nil {
		t.Errorf("Error opening sample file: %s\n", err1)
	}
	defer file.Close()

	// Write content to file
	_, err = file.WriteString("Test line\n")
	if err != nil {
		t.Errorf("Error writing content to file: %s\n", err)
	}

	// Save changes
	err = file.Sync()
	if err != nil {
		t.Errorf("Error saving file: %s\n", err)
	}

	zipFile := directoryPath + "/testZip.zip"

	// now try compressing
	err = ZipDir(directoryPath, zipFile)

	if err != nil {
		t.Errorf("Error compressing directory: %s\n", err)
	}

	// delete file
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error deleting file: %s\n", err)
	}

	// delete zip file
	err = os.Remove(zipFile)
	if err != nil {
		t.Errorf("Error deleting file: %s\n", err)
	}

	// delete directory
	err = os.Remove(directoryPath)
	if err != nil {
		t.Errorf("Error deleting directory: %s\n", err)
	}
}
