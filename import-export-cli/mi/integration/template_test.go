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

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/integration/testutils"
)

const validEndpointTemplateName = "sample_template"
const validSequenceTemplateName = "sample_seq_template"
const invalidTemplateName = "abc-template"
const templateCmd = "templates"
const endpointTemplateCmd = "endpoint"
const sequenceTemplateCmd = "sequence"
const invalidTemplatetype = "abc"

func TestGetTemplates(t *testing.T) {
	testutils.ValidateTemplateList(t, config, templateCmd)
}

func TestGetTemplatesByEndpointType(t *testing.T) {
	testutils.ValidateSpecificTemplateList(t, config, templateCmd, endpointTemplateCmd)
}

func TestGetTemplatesBySequenceType(t *testing.T) {
	testutils.ValidateSpecificTemplateList(t, config, templateCmd, sequenceTemplateCmd)
}

func TestGetEndpointTemplateByName(t *testing.T) {
	testutils.ValidateEndpointTemplate(t, config, templateCmd, validEndpointTemplateName)
}

func TestGetSequenceTemplateByName(t *testing.T) {
	testutils.ValidateSequenceTemplate(t, config, sequenceTemplateCmd, validSequenceTemplateName)
}

func TestGetNonExistingEndpointTemplateByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, config, templateCmd, endpointTemplateCmd, invalidTemplateName)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of templates [ "+invalidTemplateName+" ]  404 Not Found")
}

func TestGetNonExistingSequenceTemplateByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, config, templateCmd, sequenceTemplateCmd, invalidTemplateName)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of templates [ "+invalidTemplateName+" ]  404 Not Found")
}

func TestGetTemplatesWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, templateCmd)
}

func TestGetTemplatesWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, templateCmd, config)
}

func TestGetTemplatesWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, templateCmd, config)
}

func TestGetTemplatesWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 2, 3, false, templateCmd, validEndpointTemplateName, invalidTemplateName, "abc123")
}

func TestGetTemplatesWithInvalidType(t *testing.T) {
	response, _ := base.Execute(t, "mi", "get", templateCmd, invalidTemplatetype, "-k")
	base.Log(response)
	expected := "accepts endpoint or sequence as template-type, invalid template type " + invalidTemplatetype
	assert.Contains(t, response, expected)
}
