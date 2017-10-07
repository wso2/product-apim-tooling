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
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var HttpRequestTimeout int = 2500
var SkipTLSVerification bool = true
var ExportDirectory string
var ConfigDirectory string

func SetConfigVars(mainConfigFilePath string) error {
	mainConfig := GetMainConfigFromFile(mainConfigFilePath)
	Logln(LogPrefixInfo + " reading '" + mainConfigFilePath + "'")

	// validate config vars
	if reflect.ValueOf(mainConfig.Config.HttpRequestTimeout).Kind() != reflect.Int {
		// value of httpRequestTimeout is not an int
		Logln(LogPrefixError + "value of HttpRequestTimeout in '" + mainConfigFilePath + "' is not an integer")
		return errors.New("invalid value for HttpRequestTimeout. Should be an integer")
	}
	if !(mainConfig.Config.HttpRequestTimeout >= 0) {
		Logln(LogPrefixWarning + "value of HttpRequestTimeout in '" + mainConfigFilePath + "' is less than zero")
		Logln(LogPrefixInfo + " setting HttpRequestTimeout to " + DefaultHttpRequestTimeout)
		// default it unlimited
	}
	if reflect.ValueOf(mainConfig.Config.SkipTLSVerification).Kind() != reflect.Bool {
		// value of SkipTLSVerification is not a boolean
		return errors.New("invalid value for SkipTLSVerification. Should be true/false")
	}
	if strings.TrimSpace(mainConfig.Config.ExportDirectory) == "" ||
		len(strings.TrimSpace(mainConfig.Config.ExportDirectory)) == 0 {
		return errors.New("exportDirectory cannot be blank")
	}
	if !IsValid(mainConfig.Config.ExportDirectory) {
		Logln(LogPrefixWarning + "export Directory path invalid or the user doesn't have necessary privileges")
	}

	HttpRequestTimeout = mainConfig.Config.HttpRequestTimeout
	Logln(LogPrefixInfo + "Setting HttpTimeoutRequest to " + string(mainConfig.Config.HttpRequestTimeout))

	SkipTLSVerification = mainConfig.Config.SkipTLSVerification
	Logln(LogPrefixInfo + "Setting SkipTLSVerification to " + strconv.FormatBool(mainConfig.Config.SkipTLSVerification))

	ExportDirectory = mainConfig.Config.ExportDirectory
	Logln(LogPrefixInfo + "Setting ExportDirectory " + mainConfig.Config.ExportDirectory)

	Logln(LogPrefixInfo + "Setting ConfigDirectory" + mainConfig.Config.ConfigDirectory)
	ConfigDirectory = mainConfig.Config.ConfigDirectory

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
