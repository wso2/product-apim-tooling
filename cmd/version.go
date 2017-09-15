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

	"github.com/menuka94/wso2apim-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: utils.VersionCmdShortDesc,
	Long:  utils.VersionCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		var version string = "0.1"
		fmt.Println("wso2apim-cli Version " + version)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	VersionCmd.PersistentFlags().String("foo", "", "A help for foo")
	VersionCmd.PersistentFlags().StringP("full", "f", "fullValue", "show Full values")
	viper.BindPFlag("full", RootCmd.PersistentFlags().Lookup("full"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// VersionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
