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

package mg

import (
	"github.com/spf13/cobra"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/impl/mg"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	deployAPIDir         string
	deployAPIOverride    bool
	deployAPIEnv         string
	deployAPISkipCleanup bool
)

const (
	deployAPICmdShortDesc = "Deploy an API (apictl project) in Microgateway"
	deployAPICmdLongDesc  = "Deploy an API (apictl project) in Microgateway by " +
		"specifying the microgateway adapter environment."
)

const deployAPICmdExamples = utils.ProjectName + " " + mgCmdLiteral + " " +
	deployCmdLiteral + " " + apiCmdLiteral + " -e dev " +
	"-f petstore" +

	"\n\nNote: The flags --environment (-e), --file (-f) are mandatory. " +
	"The user needs to be logged in to use this command."

var DeployAPICmd = &cobra.Command{
	Use:     apiCmdLiteral,
	Short:   deployAPICmdShortDesc,
	Long:    deployAPICmdLongDesc,
	Example: deployAPICmdExamples,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tempMap := make(map[string]string)

		impl.DeployAPI(deployAPIEnv, deployAPIDir, tempMap,
			deployAPISkipCleanup, deployAPIOverride)
	},
}

func init() {
	DeployCmd.AddCommand(DeployAPICmd)
	DeployAPICmd.Flags().StringVarP(&deployAPIDir, "file", "f", "", "Filepath of the apictl project to be deployed")
	DeployAPICmd.Flags().StringVarP(&deployAPIEnv, "environment", "e", "", "Microgateway adapter environment to add the API")
	DeployAPICmd.Flags().BoolVarP(&deployAPIOverride, "override", "o", false, "Whether to deploy an API irrespective of its existance. Overrides when exists.")
	DeployAPICmd.Flags().BoolVarP(&deployAPISkipCleanup, "skip-cleanup", "", false, "Whether to keep "+
		"all temporary files created during deploy process")

	_ = DeployAPICmd.MarkFlagRequired("environment")
	_ = DeployAPICmd.MarkFlagRequired("file")
}
