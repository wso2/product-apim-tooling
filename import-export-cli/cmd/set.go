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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"fmt"
)

var flagHttpRequestTimeout int
var flagSkipTLSVerification bool
var flagExportDirectory string

// SetCmd represents the 'set' command
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: utils.SetCmdShortDesc,
	Long:  utils.SetCmdLongDesc + utils.SetCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "set called")

		// read the existing config vars
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)

		if flagHttpRequestTimeout > 0 {
			configVars.Config.HttpRequestTimeout = flagHttpRequestTimeout
		}else{
			fmt.Println("Invalid input for flag --http-request-timeout")
		}

		// boolean flag. no need to validate. default is set to false
		configVars.Config.SkipTLSVerification = flagSkipTLSVerification

		if flagExportDirectory != "" && utils.IsValid(flagExportDirectory) {
			configVars.Config.ExportDirectory = flagExportDirectory
		}else{
			fmt.Println("Invalid input for flag --export-directory")
		}

		utils.WriteConfigFile(configVars, utils.MainConfigFilePath)
	},
}

func init() {
	RootCmd.AddCommand(SetCmd)

	var defaultHttpRequestTimeout int
	var defaultExportDirectory string
	var defaultSkipTLSVerification bool

	// read current values in file to be passed into default values for flags below
	mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)

	defaultSkipTLSVerification = mainConfig.Config.SkipTLSVerification

	if mainConfig.Config.HttpRequestTimeout != 0 {
		defaultHttpRequestTimeout = mainConfig.Config.HttpRequestTimeout
	}

	if mainConfig.Config.ExportDirectory != "" {
		defaultExportDirectory = mainConfig.Config.ExportDirectory
	}

	SetCmd.Flags().IntVar(&flagHttpRequestTimeout, "http-request-timeout", defaultHttpRequestTimeout,
		"Timeout for HTTP Client")
	SetCmd.Flags().BoolVar(&flagSkipTLSVerification, "skip-tls-verification", defaultSkipTLSVerification,
		"Skip SSL/TLS verification")
	SetCmd.Flags().StringVar(&flagExportDirectory, "export-directory", defaultExportDirectory,
		"Path to directory where APIs should be saved")
}
