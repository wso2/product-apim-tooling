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

package olm

import (
	"fmt"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"time"
)

// installOLM installs Operator Lifecycle Manager (OLM) with the given version
// this implements the logic in
// https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.13.0/install.sh
func InstallOLM(version string) {
	utils.Logln(utils.LogPrefixInfo + "Installing OLM")

	olmNamespace := "olm"
	csvPhaseSucceeded := "Succeeded"

	// apply OperatorHub CRDs
	if err := k8sUtils.K8sApplyFromFile(fmt.Sprintf(CrdUrlTemplate, version)); err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}

	// wait for OperatorHub CRDs
	if err := k8sUtils.K8sWaitForResourceType(10, "clusterserviceversions.operators.coreos.com", "catalogsources.operators.coreos.com", "operatorgroups.operators.coreos.com"); err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}

	// apply OperatorHub OLM
	if err := k8sUtils.K8sApplyFromFile(fmt.Sprintf(OlmUrlTemplate, version)); err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}

	// rolling out
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sRollOut, "status", "-w", "deployment/olm-operator", "-n", olmNamespace); err != nil {
		utils.HandleErrorAndExit("Error installing OLM: Rolling out deployment OLM Operator", err)
	}
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sRollOut, "status", "-w", "deployment/catalog-operator", "-n", olmNamespace); err != nil {
		utils.HandleErrorAndExit("Error installing OLM: Rolling out deployment Catalog Operator", err)
	}

	// wait max 50s to csv phase to be succeeded
	csvPhase := ""
	for i := 50; i > 0 && csvPhase != csvPhaseSucceeded; i-- {
		newCsvPhase, err := k8sUtils.GetCommandOutput(k8sUtils.Kubectl, k8sUtils.K8sGet, OperatorCsv, "-n", olmNamespace, "packageserver", "-o", `jsonpath={.status.phase}`)
		if err != nil {
			utils.HandleErrorAndExit("Error installing OLM: Getting csv phase", err)
		}

		// only print new phase
		if csvPhase != newCsvPhase {
			fmt.Println("Package server phase: " + newCsvPhase)
			csvPhase = newCsvPhase
		}

		// sleep 1 second
		time.Sleep(1e9)
	}

	if csvPhase != csvPhaseSucceeded {
		utils.HandleErrorAndExit("Error installing OLM: CSV Package Server failed to reach phase succeeded", nil)
	}
}

// InstallApiOperator installs WSO2 api-operator from Operator-Hub
func InstallApiOperator() {
	utils.Logln(utils.LogPrefixInfo + "Installing API Operator from Operator-Hub")

	err := k8sUtils.K8sApplyFromFile(ApiOperatorYamlUrl)
	if err != nil {
		utils.HandleErrorAndExit("Error installing API Operator from Operator-Hub", err)
	}
}
