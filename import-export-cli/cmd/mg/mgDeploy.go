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
	mgwImportAPIFile    string
	mgDeployOverwrite   bool
	username            string
	password            string
	mgDeploySkipCleanup bool
)

const (
	mgDeployCmdLiteral   = "deploy"
	mgDeployCmdShortDesc = "Deploy an API (apictl project) in Microgateway"
	mgDeployCmdLongDesc  = "Deploy an API (apictl project) in Microgateway by " +
		"specifying the adapter host url."

	mgDeployResourcePath = "/api"
)

const mgDeployCmdExamples = utils.ProjectName + " " + mgCmdLiteral + " " + mgDeployCmdLiteral + " -h https://localhost:9095 " +
	"-f petstore -u admin -p admin" +

	"\n\nNote: The flags --host (-c), and --username (-u) are mandatory. " +
	"The password can be included via the flag --password (-p) or entered at the prompt."

//TODO: (VirajSalaka) Introduce Add environment
var MgDeployCmd = &cobra.Command{
	Use:     mgDeployCmdLiteral,
	Short:   mgDeployCmdShortDesc,
	Long:    mgDeployCmdLongDesc,
	Example: mgDeployCmdExamples,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tempMap := make(map[string]string)
		resourcePath := MgBasepath + mgDeployResourcePath

		if password == "" {
			fmt.Print("Enter Password: ")
			passwordB, err := terminal.ReadPassword(0)
			password = string(passwordB)
			fmt.Println()
			if err != nil {
				utils.HandleErrorAndExit("Error reading password", err)
			}
		}
		authToken := base64.StdEncoding.EncodeToString(
			[]byte(username + ":" + password))

		impl.DeployAPI(mgwAdapterHost+resourcePath, mgwImportAPIFile, authToken, tempMap,
			mgDeploySkipCleanup, mgDeployOverwrite)
	},
}

func init() {
	MgCmd.AddCommand(MgDeployCmd)
	//TODO: (VirajSalaka) import using just folder name
	MgDeployCmd.Flags().StringVarP(&mgwImportAPIFile, "file", "f", "", "Filepath of the apictl project to be deployed")
	MgDeployCmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "Host url for the control plane with port")
	MgDeployCmd.Flags().BoolVarP(&mgDeployOverwrite, "overwrite", "o", false, "Whether to update an existing API")
	MgDeployCmd.Flags().StringVarP(&username, "username", "u", "", "Username of the user")
	MgDeployCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the user (Can be provided at the prompt)")
	MgDeployCmd.Flags().BoolVarP(&mgDeploySkipCleanup, "skipCleanup", "", false, "Whether to keep "+
		"all temporary files created during deploy process")

	_ = MgDeployCmd.MarkFlagRequired("host")
	_ = MgDeployCmd.MarkFlagRequired("file")
	_ = MgDeployCmd.MarkFlagRequired("username")
}
