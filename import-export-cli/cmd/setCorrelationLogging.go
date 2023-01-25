/*
*  Copyright (c) WSO2 LLC (http://www.wso2.com) All Rights Reserved.
*
*  WSO2 LLC licenses this file to you under the Apache License,
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
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var setCorrelationLoggingEnvironment string
var setCorrelationLoggingComponentName string
var setCorrelationLoggingDeniedThreads string
var setCorrelationLoggingEnabled string

const SetCorrelationLoggingCmdLiteral = "correlation-logging"
const setCorrelationLoggingCmdShortDesc = "Set the correlation configs for a correlation logging component in an environment"
const setCorrelationLoggingCmdLongDesc = `Set the correlation configs for a correlation logging component in the environment specified
NOTE: The flags (--component-name (-i), --enable and --environment (-e)) are mandatory.`

var setCorrelationLoggingCmdExamples = utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetCorrelationLoggingCmdLiteral + ` --component-name http --enable true -e dev
` + utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetCorrelationLoggingCmdLiteral + ` --component-name jdbc --enable true --denied-threads MessageDeliveryTaskThreadPool,HumanTaskServer,BPELServer -e dev`

var setCorrelationLoggingCmd = &cobra.Command{
	Use:     SetCorrelationLoggingCmdLiteral,
	Short:   setCorrelationLoggingCmdShortDesc,
	Long:    setCorrelationLoggingCmdLongDesc,
	Example: setCorrelationLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + SetCmdLiteral + " " + SetCorrelationLoggingCmdLiteral + " called")
		cred, err := GetCredentials(setCorrelationLoggingEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeSetCorrelationLoggingCmd(cred)
	},
}

func executeSetCorrelationLoggingCmd(credential credentials.Credential) {
	resp, err := impl.SetCorrelationLoggingComponent(credential, setCorrelationLoggingEnvironment, setCorrelationLoggingComponentName, setCorrelationLoggingEnabled, setCorrelationLoggingDeniedThreads)
	if err != nil {
		utils.Logln(utils.LogPrefixError+"Setting the correlation log component configs", err)
		utils.HandleErrorAndExit("Error while setting the correlation log component configs", err)
	}
	// Print info on response
	utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp)
	if resp.StatusCode() == http.StatusOK {
		// 200 OK
		if strings.ToLower(setCorrelationLoggingEnabled) == "true" {
			fmt.Println("Correlation component " + setCorrelationLoggingComponentName + " is successfully enabled.")
		} else {
			fmt.Println("Correlation component " + setCorrelationLoggingComponentName + " is successfully disabled.")
		}
	} else {
		fmt.Println("Setting the correlation components : ", resp.Status(), "\n", string(resp.Body()))
	}
}

func init() {
	SetCmd.AddCommand(setCorrelationLoggingCmd)

	setCorrelationLoggingCmd.Flags().StringVarP(&setCorrelationLoggingComponentName, "component-name", "i",
		"", "Component Name")
	setCorrelationLoggingCmd.Flags().StringVarP(&setCorrelationLoggingEnabled, "enable", "",
		"", "Enable - true or false")
	setCorrelationLoggingCmd.Flags().StringVarP(&setCorrelationLoggingDeniedThreads, "denied-threads", "",
		"", "Denied Threads")
	setCorrelationLoggingCmd.Flags().StringVarP(&setCorrelationLoggingEnvironment, "environment", "e",
		"", "Environment where the correlation component configuration should be set")
	_ = setCorrelationLoggingCmd.MarkFlagRequired("environment")
	_ = setCorrelationLoggingCmd.MarkFlagRequired("component-name")
	_ = setCorrelationLoggingCmd.MarkFlagRequired("enabled")
}
