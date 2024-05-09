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

package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagHttpRequestTimeout int
var flagMarketplaceAssistantThreadCount int
var flagExportDirectory string
var flagKubernetesMode string
var flagTLSRenegotiationMode string

var flagVCSDeletionEnabled bool
var flagVCSConfigPath string
var flagVCSSourceRepoPath string
var flagVCSDeploymentRepoPath string

const flagVCSConfigPathName = "vcs-config-path"
const flagVCSSourceRepoPathName = "vcs-source-repo-path"
const flagVCSDeploymentRepoPathName = "vcs-deployment-repo-path"

// Set command related Info
const SetCmdLiteral = "set"
const setCmdShortDesc = "Set configuration parameters, per API log levels or correlation component configurations"

const setCmdLongDesc = `Set configuration parameters. You can use one of the following flags
* --http-request-timeout <time-in-milli-seconds>
* --marketplace-assistant-thread-count <number-of-threads>
* --tls-renegotiation-mode <never|once|freely>
* --export-directory <path-to-directory-where-apis-should-be-saved>
* --vcs-deletion-enabled <enable-or-disable-project-deletion-via-vcs>
* --vcs-config-path <path-to-custom-vcs-config-file>
* --vcs-deployment-repo-path <path-to-deployment-repo-for-vcs>
* --vcs-source-repo-path <path-to-source-repo-for-vcs>`

const setCmdExamples = utils.ProjectName + ` ` + SetCmdLiteral + ` --http-request-timeout 3600 --export-directory /home/user/exported-apis
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --http-request-timeout 5000 --export-directory C:\Documents\exported
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --http-request-timeout 5000
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --marketplace-assistant-thread-count 5
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --tls-renegotiation-mode freely
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --vcs-deletion-enabled=true
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --vcs-config-path /home/user/custom/vcs-config.yaml
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --vcs-deployment-repo-path /home/user/custom/deployment
` + utils.ProjectName + ` ` + SetCmdLiteral + ` --vcs-source-repo-path /home/user/custom/source
` + utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetApiLoggingCmdLiteral + ` --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level full -e dev --tenant-domain carbon.super
` + utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetCorrelationLoggingCmdLiteral + ` --component-name http --enable true -e dev`

// SetCmd represents the 'set' command
var SetCmd = &cobra.Command{
	Use:     "set",
	Short:   setCmdShortDesc,
	Long:    setCmdLongDesc,
	Example: setCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + SetCmdLiteral + " called")
		executeSetCmd(utils.MainConfigFilePath, cmd)
	},
}

func executeSetCmd(mainConfigFilePath string, cmd *cobra.Command) {
	// read the existing config vars
	configVars := utils.GetMainConfigFromFile(mainConfigFilePath)
	//Change Http Request timeout
	if flagHttpRequestTimeout > 0 {
		//Check whether the provided Http time out value is not equal to default value
		if flagHttpRequestTimeout != configVars.Config.HttpRequestTimeout {
			fmt.Println("Http Request Timeout is set to : ", flagHttpRequestTimeout)
		}
		configVars.Config.HttpRequestTimeout = flagHttpRequestTimeout
	} else {
		fmt.Println("Invalid input for flag --http-request-timeout")
	}

	if flagMarketplaceAssistantThreadCount > 0 {
		//Check whether the provided Http time out value is not equal to default value
		if flagMarketplaceAssistantThreadCount != configVars.Config.MarketplaceAssistantThreadCount {
			fmt.Println("Marketplace Assistant Thread Size is set to : ", flagMarketplaceAssistantThreadCount)
		}
		configVars.Config.MarketplaceAssistantThreadCount = flagMarketplaceAssistantThreadCount
	} else {
		fmt.Println("Invalid input for flag --marketplace-assistant-thread-count")
	}

	//Change Export Directory path
	if flagExportDirectory != "" && utils.IsValid(flagExportDirectory) {
		//Check whether the provided export directory is not equal to default value
		if flagExportDirectory != configVars.Config.ExportDirectory {
			fmt.Println("Export Directory is set to  : ", flagExportDirectory)
		}
		configVars.Config.ExportDirectory = flagExportDirectory
	} else {
		fmt.Println("Invalid input for flag --export-directory")
	}

	//Change Mode
	if flagKubernetesMode != "" {
		if strings.EqualFold(flagKubernetesMode, "kubernetes") || strings.EqualFold(flagKubernetesMode, "k8s") {
			//Check whether the provided mode value is not equal to default value
			if true != configVars.Config.KubernetesMode {
				fmt.Println("Mode is set to : ", flagKubernetesMode)
			}
			configVars.Config.KubernetesMode = true
		} else if strings.EqualFold(flagKubernetesMode, "default") {
			if false != configVars.Config.KubernetesMode {
				fmt.Println("Mode is set to : ", flagKubernetesMode)
			}
			configVars.Config.KubernetesMode = false
		} else {
			utils.HandleErrorAndExit("Error changing mode ",
				errors.New("mode should be set to either kubernetes or none"))
		}
	}

	// TLS renegotiation mode
	if flagTLSRenegotiationMode == utils.TLSRenegotiationNever ||
		flagTLSRenegotiationMode == utils.TLSRenegotiationOnce ||
		flagTLSRenegotiationMode == utils.TLSRenegotiationFreely {
		configVars.Config.TLSRenegotiationMode = flagTLSRenegotiationMode
	} else {
		fmt.Println("Invalid input for flag --tls-renegotiation-mode")
	}

	//VCS configs
	if configVars.Config.VCSDeletionEnabled != flagVCSDeletionEnabled {
		if flagVCSDeletionEnabled {
			fmt.Println("Project deletion is enabled in VCS")
		} else {
			fmt.Println("Project deletion is disabled in VCS")
		}
		configVars.Config.VCSDeletionEnabled = flagVCSDeletionEnabled
	}
	if cmd.Flags().Changed(flagVCSConfigPathName) {
		configVars.Config.VCSConfigFilePath = flagVCSConfigPath
		fmt.Println("VCS config file path is set to : " + flagVCSConfigPath)
	}
	if cmd.Flags().Changed(flagVCSSourceRepoPathName) {
		configVars.Config.VCSSourceRepoPath = flagVCSSourceRepoPath
		fmt.Println("VCS source repo path is set to : " + flagVCSSourceRepoPath)
	}
	if cmd.Flags().Changed(flagVCSDeploymentRepoPathName) {
		configVars.Config.VCSDeploymentRepoPath = flagVCSDeploymentRepoPath
		fmt.Println("VCS deployment repo path is set to : " + flagVCSDeploymentRepoPath)
	}

	utils.WriteConfigFile(configVars, mainConfigFilePath)
}

// init using Cobra
func init() {
	RootCmd.AddCommand(SetCmd)

	var defaultHttpRequestTimeout int
	var defaultMarketplaceAssistantThreadCount int
	var defaultExportDirectory string

	// read current values in file to be passed into default values for flags below
	mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)

	if mainConfig.Config.HttpRequestTimeout != 0 {
		defaultHttpRequestTimeout = mainConfig.Config.HttpRequestTimeout
	}

	if mainConfig.Config.MarketplaceAssistantThreadCount != 0 {
		defaultMarketplaceAssistantThreadCount = mainConfig.Config.MarketplaceAssistantThreadCount
	}

	if mainConfig.Config.ExportDirectory != "" {
		defaultExportDirectory = mainConfig.Config.ExportDirectory
	}

	SetCmd.Flags().IntVar(&flagHttpRequestTimeout, "http-request-timeout", defaultHttpRequestTimeout,
		"Timeout for HTTP Client")
	SetCmd.Flags().IntVar(&flagMarketplaceAssistantThreadCount, "marketplace-assistant-thread-count", defaultMarketplaceAssistantThreadCount,
		"No of threads to be used by Marketplace Assistant for parallel processing")
	SetCmd.Flags().StringVar(&flagExportDirectory, "export-directory", defaultExportDirectory,
		"Path to directory where APIs should be saved")
	SetCmd.Flags().StringVar(&flagTLSRenegotiationMode, "tls-renegotiation-mode", utils.TLSRenegotiationNever,
		"Supported TLS renegotiation mode")
	SetCmd.Flags().StringVarP(&flagKubernetesMode, "mode", "m", utils.DefaultEnvironmentName,
		"If mode is set to \"k8s\", apictl is capable of executing Kubectl commands. For example \"apictl get "+
			"pods\" -> \"kubectl get pods\". To go back to the default mode, set the mode to \"default\"")
	SetCmd.Flags().MarkDeprecated("mode", "if you want to run the commands in kubernetes mode use 'k8s' after 'apictl' (eg: apictl k8s add)")
	SetCmd.Flags().BoolVar(&flagVCSDeletionEnabled, "vcs-deletion-enabled", false,
		"Specifies whether project deletion is allowed during deployment.")
	SetCmd.Flags().StringVar(&flagVCSConfigPath, flagVCSConfigPathName, "",
		"Path to the VCS Configuration yaml file which keeps the VCS meta data")
	SetCmd.Flags().StringVar(&flagVCSSourceRepoPath, flagVCSSourceRepoPathName, "",
		"Path to the source repository to be considered during VCS deploy")
	SetCmd.Flags().StringVar(&flagVCSDeploymentRepoPath, flagVCSDeploymentRepoPathName, "",
		"Path to the deoployment repository to be considered during VCS deploy")
}
