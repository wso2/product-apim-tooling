// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/menuka94/wso2apim-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var verbose bool
var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "wso2apim-cli",
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
	RootCmd.PersistentFlags().StringP("author", "a", "", "WSO2")

	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wso2apim-cli.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Temporary code
	//rootDir := "/home/menuka/.go/src/github.com/menuka94/wso2apim-cli/"
	//err := utils.ZipDir(rootDir + "exported/hogwarts", rootDir + "exported/hogwarts.zip" )
	//if err == nil {
	//	fmt.Println("Directory Zipped Successfully")
	//}else{
	//	fmt.Println("Error: ", err)
	//	os.Exit(1)
	//}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		utils.EnableVerboseMode()
	}

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
}

func importAPI(cmd *cobra.Command, args []string) {
	log.Println("importAPI command is executed")
}
