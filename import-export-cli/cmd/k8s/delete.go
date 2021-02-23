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
)

// Delete command related usage Info
const k8sDeleteCmdLiteral = "delete"
const k8sDeleteCmdShortDesc = "Delete resources related to kubernetes"
const k8sDeleteCmdLongDesc = `Delete resources by filenames, stdin, resources and names, or by resources and label selector in kubernetes mode`

const k8sDeleteCmdExamples = k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` ` + `-n petstore`

// k8sDeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:     k8sDeleteCmdLiteral,
	Short:   k8sDeleteCmdShortDesc,
	Long:    k8sDeleteCmdLongDesc,
	Example: k8sDeleteCmdExamples,
}
