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

package deprecated

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagAddEnvName string           // name of the environment to be added
var flagTokenEndpoint string        // token endpoint of the environment to be added
var flagPublisherEndpoint string    // Publisher endpoint of the environment to be added
var flagDevPortalEndpoint string    // DevPortal endpoint of the environment to be added
var flagRegistrationEndpoint string // registration endpoint of the environment to be added
var flagApiManagerEndpoint string   // api manager endpoint of the environment to be added
var flagAdminEndpoint string        // admin endpoint of the environment to be added

// AddEnv command related Info
const addEnvCmdLiteral = "add-env"
const addEnvCmdShortDesc = "Add Environment to Config file"
const addEnvCmdLongDesc = "Add new environment and its related endpoints to the config file"
const addEnvCmdExamples = utils.ProjectName + ` ` + addEnvCmdLiteral + ` -e production \
--apim  https://localhost:9443 

` + utils.ProjectName + ` ` + addEnvCmdLiteral + ` -e test \
--registration https://idp.com:9443 \
--publisher https://apim.com:9443 \
--devportal  https://apps.com:9443 \
--admin  https://apim.com:9443 \
--token https://gw.com:8243/token

` + utils.ProjectName + ` ` + addEnvCmdLiteral + ` -e dev \
--apim https://apim.com:9443 \
--registration https://idp.com:9443 \
--token https://gw.com:8243/token

NOTE: The flag --environment (-e) is mandatory.
You can either provide only the flag --apim , or all the other 4 flags (--registration --publisher --devportal --admin) without providing --apim flag.
If you are omitting any of --registration --publisher --devportal --admin flags, you need to specify --apim flag with the API Manager endpoint. In both of the
cases --token flag is optional and use it to specify the gateway token endpoint. This will be used for "apictl get-keys" operation.`

// addEnvCmdDeprecated represents the addEnv command
var addEnvCmdDeprecated = &cobra.Command{
	Use:        addEnvCmdLiteral,
	Short:      addEnvCmdShortDesc,
	Long:       addEnvCmdLongDesc,
	Example:    addEnvCmdExamples,
	Deprecated: "instead use \"" + cmd.AddCmdLiteral + " " + cmd.AddEnvCmdLiteral + "\".",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + addEnvCmdLiteral + " called")
		executeAddEnvCmd(utils.MainConfigFilePath)
	},
}

func executeAddEnvCmd(mainConfigFilePath string) {
	envEndpoints := new(utils.EnvEndpoints)
	envEndpoints.ApiManagerEndpoint = flagApiManagerEndpoint
	envEndpoints.RegistrationEndpoint = flagRegistrationEndpoint

	envEndpoints.PublisherEndpoint = flagPublisherEndpoint
	envEndpoints.DevPortalEndpoint = flagDevPortalEndpoint
	envEndpoints.AdminEndpoint = flagAdminEndpoint
	envEndpoints.TokenEndpoint = flagTokenEndpoint
	err := impl.AddEnv(flagAddEnvName, envEndpoints, mainConfigFilePath, addEnvCmdLiteral)
	if err != nil {
		utils.HandleErrorAndExit("Error adding environment", err)
	}
}

// init using Cobra
func init() {
	cmd.RootCmd.AddCommand(addEnvCmdDeprecated)

	addEnvCmdDeprecated.Flags().StringVarP(&flagAddEnvName, "environment", "e", "", "Name of the environment to be added")
	addEnvCmdDeprecated.Flags().StringVar(&flagApiManagerEndpoint, "apim", "", "API Manager endpoint for the environment")
	addEnvCmdDeprecated.Flags().StringVar(&flagPublisherEndpoint, "publisher", "", "Publisher endpoint for the environment")
	addEnvCmdDeprecated.Flags().StringVar(&flagDevPortalEndpoint, "devportal", "", "DevPortal endpoint for the environment")
	addEnvCmdDeprecated.Flags().StringVar(&flagTokenEndpoint, "token", "", "Token endpoint for the environment")
	addEnvCmdDeprecated.Flags().StringVar(&flagRegistrationEndpoint, "registration", "",
		"Registration endpoint for the environment")
	addEnvCmdDeprecated.Flags().StringVar(&flagAdminEndpoint, "admin", "", "Admin endpoint for the environment")
	_ = addEnvCmdDeprecated.MarkFlagRequired("environment")
}
