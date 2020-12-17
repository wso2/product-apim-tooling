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
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	mgwControlPlaneHost string
	mgwImportAPIFile    string
	username            string
	password            string
	mgDeploySkipCleanup bool
)

const (
	mgDeployCmdLiteral   = "deploy"
	mgDeployCmdShortDesc = "Deploy API"
	mgDeployCmdLongDesc  = "Deploy the API (apictl project) in Microgateway"

	mgDeployResourcePath = "/apis"
)

const mgDeployCmdExamples = utils.ProjectName + " " + mgCmdLiteral + " " + mgDeployCmdLiteral + " -h https://localhost:9095 " +
	"-f qa/TwitterAPI.zip -u admin -p admin\n" +
	"cat ~/.mypassword | " + utils.ProjectName + mgCmdLiteral + " " + " " + mgDeployCmdLiteral + " -h https://localhost:9095 " +
	"-f qa/TwitterAPI.zip -u admin"

type MgwResponse struct {
	Message string
}

//TODO: (VirajSalaka) Introduce Add environment
var MgDeployCmd = &cobra.Command{
	Use: mgDeployCmdLiteral + " --host [control plane url] --file [file name] " +
		"--username [username] --password [password]",
	Short:   "Deploy apictl project.",
	Long:    "Deploy the apictl project in Microgateway",
	Example: mgDeployCmdExamples,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var tempMap map[string]string
		resourcePath := MgBasepath + mgDeployResourcePath
		if password == "" {
			fmt.Printf("Provide the password for the user: %v \n", username)
			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			password = strings.TrimRight(strings.TrimSuffix(string(data), "\n"), "\r")
		}
		authToken := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		err := impl.ImportAPIToMGW(mgwControlPlaneHost+resourcePath, mgwImportAPIFile, authToken, tempMap, mgDeploySkipCleanup)
		if err != nil {
			utils.HandleErrorAndExit("Error adding swagger to microgateway", err)
		}
	},
}

func init() {
	MgCmd.AddCommand(MgDeployCmd)
	//TODO: (VirajSalaka) import using just folder name
	MgDeployCmd.Flags().StringVarP(&mgwImportAPIFile, "file", "f", "",
		"Provide the filepath of the apictl project to be imported")
	MgDeployCmd.Flags().StringVarP(&mgwControlPlaneHost, "host", "c", "",
		"Provide the host url for the control plane with port")
	MgDeployCmd.Flags().StringVarP(&username, "username", "u", "",
		"Provide the username")
	MgDeployCmd.Flags().StringVarP(&password, "password", "p", "",
		"Provide the password")
	MgDeployCmd.Flags().BoolVarP(&mgDeploySkipCleanup, "skipCleanup", "", false, "Leave "+
		"all temporary files created during import process")

	_ = MgDeployCmd.MarkFlagRequired("host")
	_ = MgDeployCmd.MarkFlagRequired("file")
	_ = MgDeployCmd.MarkFlagRequired("username")
}
