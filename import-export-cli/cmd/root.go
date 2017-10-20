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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"time"
)

var verbose bool
var cfgFile string
var insecure bool

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "wso2apim",
	Short: utils.RootCmdShortDesc,
	Long:  utils.RootCmdLongDesc,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	cobra.EnableCommandSorting = false
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose mode")
	RootCmd.PersistentFlags().BoolVarP(&insecure, "insecure","k", false,
		"Allow connections to SLL endpoints without certs")
	RootCmd.PersistentFlags().StringP("author", "a", "", "WSO2")

	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Init ConfigVars
	err := utils.SetConfigVars(utils.MainConfigFilePath)
	if err != nil {
		utils.HandleErrorAndExit("Error reading "+utils.MainConfigFilePath+".", err)
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		utils.EnableVerboseMode()
		t := time.Now()
		utils.Logf("Executed ImportExportCLI on %v\n", t.Format(time.RFC1123))
	}

	utils.Logln(utils.LogPrefixInfo +"Insecure:", insecure)
	if insecure {
		utils.SkipTLSVerification = true
	}

	/*
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".wso2apim-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")         // adding home directory as first search path
	viper.AutomaticEnv()                 // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	*/
}
