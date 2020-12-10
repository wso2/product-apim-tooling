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

package get

import (
	"fmt"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func generateGetCmdShortDescForArtifact(artifact string) string {
	return "Get information about " + artifact + " deployed in a Micro Integrator"
}

func generateGetCmdLongDescForArtifact(artifact, argument string) string {
	return "Get information about the " + artifact + " specified by command line argument [" + argument + "]\nIf not specified, list all the " + artifact + " deployed in a Micro Integrator in the environment specified by the flag --environment, -e"
}

func generateGetCmdExamplesForArtifact(resourceType, cmdLiteral, sampleResourceName string) string {
	return "To list all the " + resourceType + "\n" +
		"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + cmdLiteral + " -e dev\n" +
		"To get details about a specific " + resourceType + "\n" +
		"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + cmdLiteral + " " + sampleResourceName + " -e dev\n" +
		"NOTE: The flag (--environment (-e)) is mandatory"
}

func getTrimmedCmdLiteral(cmd string) string {
	cmdParts := strings.Fields(cmd)
	return cmdParts[0]
}

func generateFormatFlagUsage(resource string) string {
	return "Pretty-print " + resource + " using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields"
}

func printErrorForArtifact(artifactType, artifactName string, err error) {
	fmt.Println(utils.LogPrefixError+"Getting Information of "+artifactType+" [ "+artifactName+" ] ", err)
}

func printErrorForArtifactList(artifactType string, err error) {
	fmt.Println(utils.LogPrefixError+"Getting List of "+artifactType, err)
}

func printGetCmdVerboseLogForArtifact(artifactType string) {
	utils.Logln(utils.LogPrefixInfo + GetCmdLiteral + " " + artifactType + " called")
}

func isEmptyOrCurrentDir(path string) bool {
	return path == "" || path == "."
}
