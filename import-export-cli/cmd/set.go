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
	"os"
)

var flagHttpRequestTimeout int
var flagSkipTLSVerification bool
var flagExportDirectory string

// SetCmd represents the 'set' command
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: utils.SetCmdShortDesc,
	Long: utils.SetCmdLongDesc + utils.SetCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "set called")
	},
}

func init() {
	RootCmd.AddCommand(SetCmd)

	workingDir, _ := os.Getwd()
	defaultExportDirectory := workingDir + utils.PathSeparator_ + "exported"

	SetCmd.Flags().IntVar(&flagHttpRequestTimeout, "http-request-timeout", 10000,
		"Timeout for HTTP Client")
	SetCmd.Flags().BoolVar(&flagSkipTLSVerification, "skip-tls-verification", false,
		"Skip SSL/TLS verification" )
	SetCmd.Flags().StringVar(&flagExportDirectory, "export-directory", defaultExportDirectory,
		"Path to directory where APIs should be saved")
}
