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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: utils.VersionCmdShortDesc,
	Long:  utils.VersionCmdLongDesc + utils.VersionCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		var version string = "0.1"
		fmt.Println("wso2apim-cli Version " + version)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
	VersionCmd.PersistentFlags().String("foo", "", "A help for foo")
	VersionCmd.PersistentFlags().StringP("full", "f", "fullValue", "show Full values")
	viper.BindPFlag("full", RootCmd.PersistentFlags().Lookup("full"))
}
