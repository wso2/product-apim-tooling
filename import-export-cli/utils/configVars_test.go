package utils

import (
	"testing"
	"os"
	"path/filepath"
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
	testIncorrectMainConfigFilePath := filepath.Join(ApplicationRoot, testIncorrectMainConfigFileName)
	testIncorrectMainConfig.Config = Config{0, ""}
	WriteConfigFile(testIncorrectMainConfig, testIncorrectMainConfigFilePath)

	err := SetConfigVars(testIncorrectMainConfigFilePath)

	if err == nil {
		t.Errorf("Expected error, got nil\n")
	}

	os.Remove(testIncorrectMainConfigFilePath)
}

func TestSetConfigVarsIncorrect2(t *testing.T) {
	testIncorrectMainConfig := new(MainConfig)
	testIncorrectMainConfigFileName := "test_incorrect_main_config.yaml"
	testIncorrectMainConfigFilePath := filepath.Join(ApplicationRoot, testIncorrectMainConfigFileName)
	testIncorrectMainConfig.Config = Config{-10, ""}
	WriteConfigFile(testIncorrectMainConfig, testIncorrectMainConfigFilePath)

	err := SetConfigVars(testIncorrectMainConfigFilePath)

	if err == nil {
		t.Errorf("Expected error, got nil\n")
	}

	os.Remove(testIncorrectMainConfigFilePath)
}


