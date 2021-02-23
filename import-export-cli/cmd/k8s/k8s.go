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
	"bytes"
	"github.com/spf13/cobra"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io"
	"os"
	"os/exec"
)

// K8s command related usage Info
const K8sCmdLiteral = "k8s"
const k8sCmdShortDesc = "Kubernetes mode based commands"

const k8sCmdLongDesc = `Kubernetes mode based commands such as add, update and delete API`

const k8sCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sAddCmdLiteral + ` ` + AddApiCmdLiteral + ` ` +
	`-n petstore -f Swagger.json --namespace=wso2
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sUpdateCmdLiteral + ` ` + AddApiCmdLiteral + ` ` +
	`-n petstore -f Swagger.json --namespace=wso2
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` ` +
	`-n petstore`

// K8sCmd represents the import command
var Cmd = &cobra.Command{
	Use:     K8sCmdLiteral,
	Short:   k8sCmdShortDesc,
	Long:    k8sCmdLongDesc,
	Example: k8sCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + K8sCmdLiteral + " called")
		ExecuteKubernetes(args...)
	},
}

//execute kubernetes commands
func ExecuteKubernetes(arg ...string) {
	cmd := exec.Command(
		k8sUtils.Kubectl,
		arg...,
	)
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	err := cmd.Run()
	if err != nil {
		utils.HandleErrorAndExit("Error executing kubernetes commands ", err)
	}
}

// init using Cobra
func init() {
	Cmd.AddCommand(AddCmd)
	Cmd.AddCommand(GenCmd)
	Cmd.AddCommand(DeleteCmd)
	Cmd.AddCommand(UpdateCmd)
}
