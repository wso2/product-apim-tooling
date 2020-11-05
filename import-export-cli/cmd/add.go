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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const AddCmdLiteral = "add"
const AddCmdShortDesc = "Add Environment to Config file"
const AddCmdLongDesc = `Add new environment and its related endpoints to the config file`
const addCmdExamples = utils.ProjectName + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` production \
--apim  https://localhost:9443 

` + utils.ProjectName + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` test \
--registration https://idp.com:9443 \
--publisher https://apim.com:9443 \
--devportal  https://apps.com:9443 \
--admin  https://apim.com:9443 \
--token https://gw.com:8243/token

` + utils.ProjectName + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` dev \
--apim https://apim.com:9443 \
--registration https://idp.com:9443 \
--token https://gw.com:8243/token

NOTE: The flag --environment (-e) is mandatory.
You can either provide only the flag --apim , or all the other 4 flags (--registration --publisher --devportal --admin) without providing --apim flag.
If you are omitting any of --registration --publisher --devportal --admin flags, you need to specify --apim flag with the API Manager endpoint. In both of the
cases --token flag is optional and use it to specify the gateway token endpoint. This will be used for "apictl get-keys" operation.`

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:     AddCmdLiteral,
	Short:   AddCmdShortDesc,
	Long:    AddCmdLongDesc,
	Example: addCmdExamples,
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
