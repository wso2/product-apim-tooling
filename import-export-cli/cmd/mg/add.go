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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	addCmdLiteral   = "add"
	addCmdShortDesc = "Add Environment to Config file"
	addCmdLongDesc  = `Add new environment and its related endpoints to the config file`
)
const addCmdExamples = utils.ProjectName + " " + mgCmdLiteral + " " + addCmdLiteral + " " + envCmdLiteral +
	" prod --host  https://localhost:9443" +

	"\n\nNOTE: The flag --host (-c) is mandatory and it has to specify the microgateway adapter" +
	" url."

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:     addCmdLiteral,
	Short:   addCmdShortDesc,
	Long:    addCmdLongDesc,
	Example: addCmdExamples,
}

func init() {
	MgCmd.AddCommand(AddCmd)
}
