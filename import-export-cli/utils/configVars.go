/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/

package utils

import (
	"errors"
	"reflect"
	"os"
	"io/ioutil"
	"fmt"
)

var HttpRequestTimeout int = 2500
var SkipTLSVerification bool = true
var ExportDirectory string


func SetConfigVars() (error){
	mainConfig := GetMainConfigFromFile(MainConfigFilePath)
	fmt.Println("httprequesttimeout:", mainConfig.Config.HttpRequestTimeout)
	fmt.Println("skiptlsverification:", mainConfig.Config.SkipTLSVerification)
	fmt.Println("exportdirectory:", mainConfig.Config.ExportDirectory)

	// validate config vars
	if reflect.ValueOf(mainConfig.Config.HttpRequestTimeout).Kind() != reflect.Int {
		return errors.New("invalid value for HttpRequestTimeout. Should be an integer")
	}
	if !(mainConfig.Config.HttpRequestTimeout >= 0)  {
		return errors.New("invalid HttpRequestTimeout")
	}
	if reflect.ValueOf(mainConfig.Config.SkipTLSVerification).Kind() != reflect.Bool {
		return errors.New("invalid value for SkipTLSVerification. Should be true/false")
	}
	fmt.Println("Test 5")
	if mainConfig.Config.ExportDirectory == "" || len(mainConfig.Config.ExportDirectory) == 0{
		errors.New("exportDirectory cannot be blank")
	}
	if !IsValid(mainConfig.Config.ExportDirectory) {
		errors.New("export Directory path in valid or the user doesn't have necessary privileges")
	}


	HttpRequestTimeout = mainConfig.Config.HttpRequestTimeout
	SkipTLSVerification = mainConfig.Config.SkipTLSVerification
	ExportDirectory = mainConfig.Config.ExportDirectory

	return nil
}

// Attempt to create a file and delete it right after
func IsValid(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}

