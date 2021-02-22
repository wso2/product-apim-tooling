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
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/impl/mg"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	deployAPIDir         string
	deployAPIOverride    bool
	deployAPIUsername    string
	deployAPIPassword    string
	deployAPISkipCleanup bool
)

const (
	deployAPICmdShortDesc = "Deploy an API (apictl project) in Microgateway"
	deployAPICmdLongDesc  = "Deploy an API (apictl project) in Microgateway by " +
		"specifying the adapter host url."

	mgDeployResourcePath = "/apis"
)

const deployAPICmdExamples = utils.ProjectName + " " + mgCmdLiteral + " " +
	deployCmdLiteral + " " + apiCmdLiteral + " -c https://localhost:9095 " +
	"-f petstore -u admin -p admin" +

	"\n\nNote: The flags --host (-c), and --deployAPIUsername (-u) are mandatory. " +
	"The password can be included via the flag --password (-p) or entered at the prompt."

//TODO: (VirajSalaka) Introduce Add environment
var DeployAPICmd = &cobra.Command{
	Use:     apiCmdLiteral,
	Short:   deployAPICmdShortDesc,
	Long:    deployAPICmdLongDesc,
	Example: deployAPICmdExamples,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tempMap := make(map[string]string)
		resourcePath := MgBasepath + mgDeployResourcePath

		if deployAPIPassword == "" {
			fmt.Print("Enter Password: ")
			deployAPIPasswordB, err := terminal.ReadPassword(0)
			deployAPIPassword = string(deployAPIPasswordB)
			fmt.Println()
			if err != nil {
				utils.HandleErrorAndExit("Error reading password", err)
			}
		}
		authToken := base64.StdEncoding.EncodeToString(
			[]byte(deployAPIUsername + ":" + deployAPIPassword))

		impl.DeployAPI(mgwAdapterHost+resourcePath, deployAPIDir, authToken, tempMap,
			deployAPISkipCleanup, deployAPIOverride)
	},
}

func init() {
	DeployCmd.AddCommand(DeployAPICmd)
	//TODO: (VirajSalaka) import using just folder name
	DeployAPICmd.Flags().StringVarP(&deployAPIDir, "file", "f", "", "Filepath of the apictl project to be deployed")
	DeployAPICmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "Host url for the control plane with port")
	DeployAPICmd.Flags().BoolVarP(&deployAPIOverride, "override", "o", false, "Whether to deploy an API irrespective of its existance. Overrides when exists.")
	DeployAPICmd.Flags().StringVarP(&deployAPIUsername, "username", "u", "", "Username of the user")
	DeployAPICmd.Flags().StringVarP(&deployAPIPassword, "password", "p", "", "Password of the user (Can be provided at the prompt)")
	DeployAPICmd.Flags().BoolVarP(&deployAPISkipCleanup, "skip-cleanup", "", false, "Whether to keep "+
		"all temporary files created during deploy process")

	_ = DeployAPICmd.MarkFlagRequired("host")
	_ = DeployAPICmd.MarkFlagRequired("file")
	_ = DeployAPICmd.MarkFlagRequired("username")
}
