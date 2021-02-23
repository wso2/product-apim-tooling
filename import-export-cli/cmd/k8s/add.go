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

package k8s

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const K8sAddCmdLiteral = "add"
const k8sAddCmdShortDesc = "Add an API to the kubernetes cluster"
const k8sAddCmdLongDesc = `Add an API either from a Swagger file, project zip for API project to the kubernetes cluster. 
JSON, YAML, zip and API project formats are accepted.`
const k8sAddCmdExamples = utils.ProjectName + " " + K8sCmdLiteral + " " + K8sAddCmdLiteral + " " + AddApiCmdLiteral +
	" " + `-n petstore -f Swagger.json --namespace=wso2
` + utils.ProjectName + " " + K8sCmdLiteral + " " + K8sAddCmdLiteral + " " + AddApiCmdLiteral +
	" " + `-n petstore -f product-apim-tooling/import-export-cli/build/target/apictl/myapi.zip --namespace=wso2
` + utils.ProjectName + " " + K8sCmdLiteral + " " + K8sAddCmdLiteral + " " + AddApiCmdLiteral +
	" " + `-n petstore -f myapi --namespace=wso2`

// K8sAddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:     K8sAddCmdLiteral,
	Short:   k8sAddCmdShortDesc,
	Long:    k8sAddCmdLongDesc,
	Example: k8sAddCmdExamples,
}
