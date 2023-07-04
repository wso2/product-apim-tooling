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
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getTemplateCmdEnvironment string
var getTemplateCmdFormat string

const artifactTemplates = "templates"
const getTemplateCmdLiteral = "templates [template-type] [template-name]"

const getTemplateCmdShortDesc = "Get information about templates deployed in a Micro Integrator"

const getTemplateCmdLongDesc = "Get information about the template specified by command line arguments [template-type] and [template-name]\n" +
	"If not specified, list all the templates in the environment specified by the flag --environment, -e"

var getTemplateCmdExamples = "To list all the " + artifactTemplates + "\n" +
	utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTemplateCmdLiteral) + " -e dev\n" +
	"To get details about a specific template type\n" +
	utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTemplateCmdLiteral) + " TemplateType\n" +
	"To get details about a specific template\n" +
	utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTemplateCmdLiteral) + " TemplateType TemplateName -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

const endpointKey string = "endpoint"
const sequenceKey string = "sequence"

var getTemplateCmd = &cobra.Command{
	Use:     getTemplateCmdLiteral,
	Short:   getTemplateCmdShortDesc,
	Long:    getTemplateCmdLongDesc,
	Example: getTemplateCmdExamples,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			var errMessage = "accepts at most 2 arg(s), received " + fmt.Sprint(len(args))
			return errors.New(errMessage)
		} else if len(args) > 0 {
			if !isValidTemplateType(args[0]) {
				var errMessage = "accepts " + endpointKey + " or " + sequenceKey + " as template-type, invalid template type " + args[0]
				return errors.New(errMessage)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleGetTemplateCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getTemplateCmd)
	setEnvFlag(getTemplateCmd, &getTemplateCmdEnvironment)
	setFormatFlag(getTemplateCmd, &getTemplateCmdFormat)
}

func handleGetTemplateCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getTemplateCmdLiteral))
	credentials.HandleMissingCredentials(getTemplateCmdEnvironment)
	if len(args) == 2 {
		var templateType = args[0]
		var templateName = args[1]
		executeGetTemplateByNameCmd(templateType, templateName)
	} else if len(args) == 1 {
		var templateType = args[0]
		executeGetTemplateByTypeCmd(templateType)
	} else {
		executeListTemplates()
	}
}

func executeListTemplates() {
	templateList, err := impl.GetTemplateList(getTemplateCmdEnvironment)
	if err == nil {
		impl.PrintTemplateList(templateList, getTemplateCmdFormat)
	} else {
		printErrorForArtifactList(artifactTemplates, err)
	}
}

func executeGetTemplateByTypeCmd(templateType string) {
	templateList, err := impl.GetTemplatesByType(getTemplateCmdEnvironment, templateType)
	if err == nil {
		impl.PrintTemplatesByType(templateList, getTemplateCmdFormat)
	} else {
		printErrorForArtifact(artifactTemplates, templateType, err)
	}
}

func executeGetTemplateByNameCmd(templateType, templateName string) {
	if templateType == sequenceKey {
		sequenceTemplate, err := impl.GetSequenceTemplate(getTemplateCmdEnvironment, templateName)
		if err == nil {
			impl.PrintSequenceTemplateDetails(sequenceTemplate, getTemplateCmdFormat)
		} else {
			printErrorForArtifact(artifactTemplates, templateName, err)
		}
	}
	if templateType == endpointKey {
		endpointTemplate, err := impl.GetEndpointTemplate(getTemplateCmdEnvironment, templateName)
		if err == nil {
			impl.PrintEndpointTemplateDetails(endpointTemplate, getTemplateCmdFormat)
		} else {
			printErrorForArtifact(artifactTemplates, templateName, err)
		}
	}
}

func isValidTemplateType(templateType string) bool {
	templateType = strings.ToLower(templateType)
	return templateType == endpointKey || templateType == sequenceKey
}
