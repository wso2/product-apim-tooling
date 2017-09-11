package utils

import (
)


/*
func TestEnvExistsInEndpointsFile(t *testing.T) {
	env := "dev"
	workingDir, _ := os.Getwd()
	envExistsInEndpointsFile := EnvExistsInEndpointsFile(env)
	data, err := ioutil.ReadFile(workingDir +"/" + EnvEndpointsAllFileName)
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
*/

/*
func TestEnvExistsInKeysFile(t *testing.T) {
	env := "dev"
	envExistsInKeysFile := EnvExistsInKeysFile(env)
	data, err := ioutil.ReadFile("../env_keys_all.yaml")

	if err != nil {
		fmt.Println("Error:", err)
	}

	envKeysAll := &EnvKeysAll{}
	_ = yaml.Unmarshal(data, envKeysAll)

	for key, _ := range envKeysAll.Environments {
		if key == env {
			// env actually exists
			if envExistsInKeysFile{
				// function says env exists
				return
			}else{
				// functions says env doesn't exist
				t.Errorf("Error in method EnvExistInKeysFile()")
			}
		}else{
			// env doesn't actually exist
			if envExistsInKeysFile{
				// function says env exists
				t.Errorf("Error in method EnvExistInKeysFile()")
			}else{
				// functions says env doesn't exist
				continue
			}
		}
	}

}
*/