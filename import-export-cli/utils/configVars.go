/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"crypto/tls"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var HttpRequestTimeout = DefaultHttpRequestTimeout
var Insecure bool
var ExportDirectory string

// TLSRenegotiationMode : Defines TLS Renegotiation support mode, default is never
var TLSRenegotiationMode = tls.RenegotiateNever

// SetConfigVars
// @param mainConfigFilePath : Path to file where Configuration details are stored
// @return error
func SetConfigVars(mainConfigFilePath string) error {
	mainConfig := GetMainConfigFromFile(mainConfigFilePath)
	Logln(LogPrefixInfo + " reading '" + mainConfigFilePath + "'")

	// validate config vars
	if !(mainConfig.Config.HttpRequestTimeout >= 0) {
		Logln(LogPrefixWarning + "value of HttpRequestTimeout in '" + mainConfigFilePath + "' is less than zero")
		Logln(LogPrefixInfo + " setting HttpRequestTimeout to " + string(DefaultHttpRequestTimeout))
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

	ExportDirectory = mainConfig.Config.ExportDirectory
	Logln(LogPrefixInfo + "Setting ExportDirectory " + mainConfig.Config.ExportDirectory)

	setTLSRenegotiationMode(mainConfig)

	return nil
}

// IsValid
// @param fp : FilePath
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

func setTLSRenegotiationMode(mainConfig *MainConfig) {
	modeMap := map[string]tls.RenegotiationSupport{
		TLSRenegotiationOnce:   tls.RenegotiateOnceAsClient,
		TLSRenegotiationFreely: tls.RenegotiateFreelyAsClient,
		TLSRenegotiationNever:  tls.RenegotiateNever,
	}

	if val, ok := modeMap[mainConfig.Config.TLSRenegotiationMode]; ok {
		if ok {
			TLSRenegotiationMode = val
			Logln(LogPrefixInfo + "Setting TLSRenegotiationMode : " + mainConfig.Config.TLSRenegotiationMode)
		} else {
			Logln(LogPrefixInfo + "Setting TLSRenegotiationMode : never")
		}
	}
}
