package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetConfigVarsCorrect(t *testing.T) {
	WriteCorrectMainConfig()
	err := SetConfigVars(testMainConfigFilePath)

	if err != nil {
		t.Errorf("Error in setting Config Vars: %s\n", err.Error())
	}

	defer os.Remove(testMainConfigFilePath)
}

func TestSetConfigVarsIncorrect1(t *testing.T) {
	testIncorrectMainConfig := new(MainConfig)
	testIncorrectMainConfigFileName := "test_incorrect_main_config.yaml"
	testIncorrectMainConfigFilePath := filepath.Join(ConfigDirPath, testIncorrectMainConfigFileName)
	testIncorrectMainConfig.Config = Config{0, ""}
	WriteConfigFile(testIncorrectMainConfig, testIncorrectMainConfigFilePath)

	err := SetConfigVars(testIncorrectMainConfigFilePath)

	if err == nil {
		t.Errorf("Expected error, got nil\n")
	}

	defer os.Remove(testIncorrectMainConfigFilePath)
}

func TestSetConfigVarsIncorrect2(t *testing.T) {
	testIncorrectMainConfig := new(MainConfig)
	testIncorrectMainConfigFileName := "test_incorrect_main_config.yaml"
	testIncorrectMainConfigFilePath := filepath.Join(ConfigDirPath, testIncorrectMainConfigFileName)
	testIncorrectMainConfig.Config = Config{-10, ""}
	WriteConfigFile(testIncorrectMainConfig, testIncorrectMainConfigFilePath)

	err := SetConfigVars(testIncorrectMainConfigFilePath)

	if err == nil {
		t.Errorf("Expected error, got nil\n")
	}

	defer os.Remove(testIncorrectMainConfigFilePath)
}

// TestIsValid1 - Create new file
func TestIsValid1(t *testing.T) {
	filePath := filepath.Join(CurrentDir, "test.txt")
	IsValid(filePath)
}

// TestIsValid2 - Create new file
func TestIsValid2(t *testing.T) {
	fileName := "test.txt"
	os.Create(fileName)
	filePath := filepath.Join(CurrentDir, fileName)
	IsValid(filePath)
	os.Remove(filePath)
}
