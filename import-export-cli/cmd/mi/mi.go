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

package mi

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	miActivateCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/activate"
	miAddCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/add"
	miDeactivateCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/deactivate"
	miDeleteCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/delete"
	miGetCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/get"
	miUpdateCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/update"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var verbose bool
var cfgFile string
var insecure bool

const miCmdShortDesc = "Micro Integrator related commands"

var miCmdLongDesc = getMiCmdLongDesc()

func getMiCmdLongDesc() string {
	if utils.GetMICmdName() == "" {
		return utils.MICmd + " is a Command Line Tool for Managing WSO2 Micro Integrator"
	}
	return "Micro Integrator related commands such as login, logout, get, add, update, delete, activate, deactivate."
}

// MICmd represents the mi command
var MICmd = &cobra.Command{
	Use:   utils.MiCmdLiteral,
	Short: miCmdShortDesc,
	Long:  miCmdLongDesc,
	// Example: miCmdExamples,
		Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + utils.MiCmdLiteral + " called")
		cmd.Help()
	},
}

func init() {
	if utils.GetMICmdName() == "" {
		createConfigFiles()

		cobra.OnInitialize(initConfig)

		cobra.EnableCommandSorting = false
		MICmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose mode")
		MICmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false,
			"Allow connections to SSL endpoints without certs")
		err := utils.SetConfigVars(utils.MainConfigFilePath)
		if err != nil {
			utils.HandleErrorAndExit("Error reading "+utils.MainConfigFilePath+".", err)
		}
	}
	MICmd.AddCommand(miGetCmd.GetCmd)
	MICmd.AddCommand(miAddCmd.AddCmd)
	MICmd.AddCommand(miDeleteCmd.DeleteCmd)
	MICmd.AddCommand(miUpdateCmd.UpdateCmd)
	MICmd.AddCommand(miActivateCmd.ActivateCmd)
	MICmd.AddCommand(miDeactivateCmd.DeactivateCmd)
}

func createConfigFiles() {
	err := utils.CreateDirIfNotExist(utils.GetConfigDirPath())
	if err != nil {
		utils.HandleErrorAndExit("Error creating config directory: "+utils.ConfigDirPath, err)
	}

	if !utils.IsFileExist(utils.MainConfigFilePath) {
		var mainConfig = new(utils.MainConfig)
		mainConfig.Config = utils.Config{HttpRequestTimeout: utils.DefaultHttpRequestTimeout,
			ExportDirectory: utils.DefaultExportDirPath}

		utils.WriteConfigFile(mainConfig, utils.MainConfigFilePath)
	}

	err = utils.CreateDirIfNotExist(utils.LocalCredentialsDirectoryPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating local directory: "+utils.LocalCredentialsDirectoryName, err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		utils.EnableVerboseMode()
		t := time.Now()
		utils.Logf("Executed ImportExportCLI (%s) on %v\n", utils.MICmd, t.Format(time.RFC1123))
	}

	utils.Logln(utils.LogPrefixInfo+"Insecure:", insecure)
	if insecure {
		utils.Insecure = true
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := MICmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
