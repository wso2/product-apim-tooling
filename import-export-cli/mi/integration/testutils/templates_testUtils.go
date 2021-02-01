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

// ValidateTemplateList validate ctl output with list of templates from the Management API
func ValidateTemplateList(t *testing.T, config *MiConfig, templateCmd string) {
	t.Helper()
	output, _ := ListArtifacts(t, templateCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementTemplateResource, &artifactutils.TemplateList{})
	validateTemplateListEqual(t, output, (artifactList.(*artifactutils.TemplateList)))
}

func validateTemplateListEqual(t *testing.T, templateListFromCtl string, templateList *artifactutils.TemplateList) {
	unmatchedEndpointTemplateCount := len(templateList.EndpointTemplates)
	for _, template := range templateList.EndpointTemplates {
		assert.Truef(t, strings.Contains(templateListFromCtl, template.Name), "templateListFromCtl: "+templateListFromCtl+
			" , does not contain template.Name: "+template.Name)
		unmatchedEndpointTemplateCount--
	}
	unmatchedSequenceTemplateCount := len(templateList.SequenceTemplates)
	for _, template := range templateList.SequenceTemplates {
		assert.Truef(t, strings.Contains(templateListFromCtl, template.Name), "templateListFromCtl: "+templateListFromCtl+
			" , does not contain template.Name: "+template.Name)
		unmatchedSequenceTemplateCount--
	}
	unmatchedCount := unmatchedEndpointTemplateCount + unmatchedSequenceTemplateCount
	assert.Equal(t, 0, unmatchedCount, "template lists are not equal")
}

// ValidateSpecificTemplateList validate ctl output with list of specific templates from the Management API
func ValidateSpecificTemplateList(t *testing.T, config *MiConfig, templateCmd, templateType string) {
	t.Helper()
	output, _ := GetArtifact(t, config, templateCmd, templateType)
	artifactList := config.MIClient.GetArtifactFromAPI(utils.MiManagementTemplateResource, getParamMap("type", templateType), &artifactutils.TemplateListByType{})
	validateTemplateListByTypeEqual(t, output, (artifactList.(*artifactutils.TemplateListByType)))
}

func validateTemplateListByTypeEqual(t *testing.T, templateListFromCtl string, templateList *artifactutils.TemplateListByType) {
	unmatchedCount := templateList.Count
	for _, template := range templateList.Templates {
		assert.Truef(t, strings.Contains(templateListFromCtl, template.Name), "templateListFromCtl: "+templateListFromCtl+
			" , does not contain template.Name: "+template.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "template lists are not equal")
}

// ValidateEndpointTemplate validate ctl output with the endpoint template from the Management API
func ValidateEndpointTemplate(t *testing.T, config *MiConfig, templateCmd, templateName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, templateCmd, "endpoint", templateName)
	paramMap := make(map[string]string)
	paramMap["type"] = "endpoint"
	paramMap["templateName"] = templateName
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementTemplateResource, paramMap, &artifactutils.TemplateEndpointListByName{})
	validateEndpointTemplateEqual(t, output, (artifact.(*artifactutils.TemplateEndpointListByName)))
}

func validateEndpointTemplateEqual(t *testing.T, templateFromCtl string, template *artifactutils.TemplateEndpointListByName) {
	assert.Contains(t, templateFromCtl, template.Name)
	for _, param := range template.Parameters {
		assert.Contains(t, templateFromCtl, param)
	}
}

// ValidateSequenceTemplate validate ctl output with the sequence template from the Management API
func ValidateSequenceTemplate(t *testing.T, config *MiConfig, templateCmd, templateName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, templateCmd, "sequence", templateName)
	paramMap := make(map[string]string)
	paramMap["type"] = "sequence"
	paramMap["templateName"] = templateName
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementTemplateResource, paramMap, &artifactutils.TemplateSequenceListByName{})
	validateSequenceTemplateEqual(t, output, (artifact.(*artifactutils.TemplateSequenceListByName)))
}

func validateSequenceTemplateEqual(t *testing.T, templateFromCtl string, template *artifactutils.TemplateSequenceListByName) {
	assert.Contains(t, templateFromCtl, template.Name)
	for _, param := range template.Parameters {
		assert.Contains(t, templateFromCtl, param.Name)
		assert.Contains(t, templateFromCtl, param.DefaultValue)
		assert.Contains(t, templateFromCtl, param.IsMandatory)
	}
}
