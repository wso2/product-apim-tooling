package utils

import (
	"testing"
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
)

func TestEnvExistsInEndpointsFile(t *testing.T) {
	env := "dev"
	envExistsInEndpointsFile := EnvExistsInEndpointsFile(env)
	data, err := ioutil.ReadFile("./env_endpoints_all.yaml")
	if err != nil {
		fmt.Println("Error:", err)
	}

	envEndpointsAll := &EnvEndpointsAll{}
	_ = yaml.Unmarshal(data, envEndpointsAll)

	for key, _ := range envEndpointsAll.Environments {
		if  key == env && envExistsInEndpointsFile{
			return
		}else if key == env && !envExistsInEndpointsFile {
			t.Errorf("Error in method EnvExistInEndpointsFile()")
		}else{
			continue
		}
	}
}

func TestEnvExistsInKeysFile(t *testing.T) {
	env := "dev"
	envExistsInKeysFile := EnvExistsInKeysFile(env)
	data, err := ioutil.ReadFile("./env_keys_all.yaml")

	if err != nil {
		fmt.Println("Error:", err)
	}

	envKeysAll := &EnvKeysAll{}
	_ = yaml.Unmarshal(data, envKeysAll)

	for key, _ := range envKeysAll.Environments {
		if  key == env && envExistsInKeysFile{
			return
		}else if key == env && !envExistsInKeysFile {
			t.Errorf("Error in method EnvExistInEndpointsFile()")
		}else{
			continue
		}
	}

}