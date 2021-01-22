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

package testutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ValidateConnectorList validate ctl output with list of connectors from the Management API
func ValidateConnectorList(t *testing.T, connectorCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, connectorCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementConnectorResource, &artifactutils.ConnectorList{})
	validateConnectorListEqual(t, output, (artifactList.(*artifactutils.ConnectorList)))
}

func validateConnectorListEqual(t *testing.T, connectorListFromCtl string, connectorList *artifactutils.ConnectorList) {
	unmatchedCount := connectorList.Count
	for _, connector := range connectorList.Connectors {
		assert.Truef(t, strings.Contains(connectorListFromCtl, connector.Name), "connectorListFromCtl: "+connectorListFromCtl+
			" , does not contain connector.Name: "+connector.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Connector lists are not equal")
}
