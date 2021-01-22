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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagHttpRequestTimeout int
var flagExportDirectory string
var flagTLSRenegotiationMode string

// Set command related Info
const setCmdLiteral = "set"
const setCmdShortDesc = "Set configuration"

const setCmdLongDesc = `Set configuration parameters. Use at least one of the following flags
* --http-request-timeout <time-in-milli-seconds>
* --export-directory <path-to-directory-where-apis-should-be-saved>
* --tls_renegotiation_mode <never|once|freely>`

const setCmdExamples = utils.ProjectName + ` ` + setCmdLiteral + ` --http-request-timeout 3600 --export-directory /home/user/exported-apis
` + utils.ProjectName + ` ` + setCmdLiteral + ` --http-request-timeout 5000 --export-directory C:\Documents\exported
` + utils.ProjectName + ` ` + setCmdLiteral + ` --http-request-timeout 5000
` + utils.ProjectName + ` ` + setCmdLiteral + ` --tls_renegotiation_mode freely`

// SetCmd represents the 'set' command
var SetCmd = &cobra.Command{
	Use:     "set",
	Short:   setCmdShortDesc,
	Long:    setCmdLongDesc,
	Example: setCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + setCmdLiteral + " called")
		executeSetCmd(utils.MainConfigFilePath, utils.ExportDirectory)
	},
}

func executeSetCmd(mainConfigFilePath, exportDirectory string) {
	// read the existing config vars
	configVars := utils.GetMainConfigFromFile(mainConfigFilePath)
	if flagHttpRequestTimeout > 0 {
		configVars.Config.HttpRequestTimeout = flagHttpRequestTimeout
	} else {
		fmt.Println("Invalid input for flag --http-request-timeout")
	}
	if flagExportDirectory != "" && utils.IsValid(flagExportDirectory) {
		configVars.Config.ExportDirectory = flagExportDirectory
	} else {
		fmt.Println("Invalid input for flag --export-directory")
	}
	if flagTLSRenegotiationMode == utils.TLSRenegotiationNever ||
		flagTLSRenegotiationMode == utils.TLSRenegotiationOnce ||
		flagTLSRenegotiationMode == utils.TLSRenegotiationFreely {
		configVars.Config.TLSRenegotiationMode = flagTLSRenegotiationMode
	} else {
		fmt.Println("Invalid input for flag --tls_renegotiation_mode")
	}
	utils.WriteConfigFile(configVars, mainConfigFilePath)
}

// init using Cobra
func init() {
	RootCmd.AddCommand(SetCmd)

	var defaultHttpRequestTimeout int
	var defaultExportDirectory string

	// read current values in file to be passed into default values for flags below
	mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)

	if mainConfig.Config.HttpRequestTimeout != 0 {
		defaultHttpRequestTimeout = mainConfig.Config.HttpRequestTimeout
	}

	if mainConfig.Config.ExportDirectory != "" {
		defaultExportDirectory = mainConfig.Config.ExportDirectory
	}

	SetCmd.Flags().IntVar(&flagHttpRequestTimeout, "http-request-timeout", defaultHttpRequestTimeout,
		"Timeout for HTTP Client")
	SetCmd.Flags().StringVar(&flagExportDirectory, "export-directory", defaultExportDirectory,
		"Path to directory where APIs should be saved")
	SetCmd.Flags().StringVar(&flagTLSRenegotiationMode, "tls_renegotiation_mode", utils.TLSRenegotiationNever,
		"Supported TLS renegotiation mode")
}
