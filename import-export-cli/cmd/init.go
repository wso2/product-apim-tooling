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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	initCmdOutputDir         string
	initCmdSwaggerPath       string
	initCmdApiDefinitionPath string
	initCmdInitialState      string
	initCmdForced            bool
)

const initCmdExample = `apictl init myapi --oas petstore.yaml
apictl init Petstore --oas https://petstore.swagger.io/v2/swagger.json
apictl init Petstore --oas https://petstore.swagger.io/v2/swagger.json --initial-state=PUBLISHED
apictl init MyAwesomeAPI --oas ./swagger.yaml -d definition.yaml`

var InitCommand = &cobra.Command{
	Use:     "init [project path]",
	Short:   "Initialize a new project in given path",
	Long:    "Initialize a new project in given path. If a OpenAPI specification provided API will be populated with details from it",
	Example: initCmdExample,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "init called")
		initCmdOutputDir = args[0]

		// check for dir existence, if so stop it unless forced flag is present
		if stat, err := os.Stat(initCmdOutputDir); !os.IsNotExist(err) {
			fmt.Printf("%s already exists\n", initCmdOutputDir)
			if !stat.IsDir() {
				fmt.Printf("%s is not a directory\n", initCmdOutputDir)
				os.Exit(1)
			}
			if !initCmdForced {
				fmt.Println("Run with -f or --force to overwrite directory and create project")
				os.Exit(1)
			}
			fmt.Println("Running command in forced mode")
		}

		// check the validity of initial-state before initializing
		if initCmdInitialState != "" {
			validState := false
			for _, state := range utils.ValidInitialStates {
				if initCmdInitialState == state {
					validState = true
					break
				}
			}
			if !validState {
				utils.HandleErrorAndExit(fmt.Sprintf(
					"Invalid initial API state: %s\nValid initial states: %v",
					initCmdInitialState, utils.ValidInitialStates,
				), nil)
			}
		}

		err := impl.InitAPIProject(initCmdOutputDir, initCmdInitialState, initCmdSwaggerPath, initCmdApiDefinitionPath, false)
		if err != nil {
			utils.HandleErrorAndContinue("Error initializing project", err)
			// Remove the already created project with its content since it is partially created and wrong
			dir, err := filepath.Abs(initCmdOutputDir)
			if err != nil {
				utils.HandleErrorAndExit("Error retrieving file path of the project", err)
			}
			fmt.Println("Removing the project directory " + dir + " with its content")
			err = os.RemoveAll(dir)
			if err != nil {
				utils.HandleErrorAndExit("Error removing project directory", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(InitCommand)
	InitCommand.Flags().StringVarP(&initCmdApiDefinitionPath, "definition", "d", "", "Provide a "+
		"YAML definition of API")
	InitCommand.Flags().StringVarP(&initCmdSwaggerPath, "oas", "", "", "Provide an OpenAPI "+
		"specification file for the API")
	InitCommand.Flags().StringVar(&initCmdInitialState, "initial-state", "", fmt.Sprintf("Provide the initial state "+
		"of the API; Valid states: %v", utils.ValidInitialStates))
	InitCommand.Flags().BoolVarP(&initCmdForced, "force", "f", false, "Force create project")
}
