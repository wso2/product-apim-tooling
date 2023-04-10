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
	"os"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getLogCmdEnvironment string
var getLogCmdFormat string
var logFileDownloadPath string

const getLogCmdLiteral = "logs [file-name]"

const getLogCmdShortDesc = "List all the available log files"
const getLogCmdLongDesc = "Download a log file by providing the file name and download location,\n" +
	"if not provided, list all the log files of the Micro Integrator in the environment specified by the flag --environment, -e"

var getLogCmdExamples = "Example:\n" +
	"To list all the log files\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getLogCmdLiteral) + " -e dev\n" +
	"To download a selected log file\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getLogCmdLiteral) + " [file-name] -p [download-location] -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var getLogCmd = &cobra.Command{
	Use:     getLogCmdLiteral,
	Short:   getLogCmdShortDesc,
	Long:    getLogCmdLongDesc,
	Example: getLogCmdExamples,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetLogCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getLogCmd)
	setEnvFlag(getLogCmd, &getLogCmdEnvironment)
	setFormatFlag(getLogCmd, &getLogCmdFormat)
	getLogCmd.Flags().StringVarP(&logFileDownloadPath, "path", "p", "", "Path the file should be downloaded")
}

func handleGetLogCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getLogCmdLiteral))
	credentials.HandleMissingCredentials(getLogCmdEnvironment)
	if len(args) == 1 {
		var logFileName = args[0]
		if isEmptyOrCurrentDir(logFileDownloadPath) {
			logFileDownloadPath, _ = os.Getwd()
		}
		executeDownloadLogFile(logFileDownloadPath, logFileName)
	} else {
		executeListLogFiles()
	}
}

func executeListLogFiles() {
	fileList, err := impl.GetLogFileList(getLogCmdEnvironment)
	if err == nil {
		logFileList := impl.FilterOnlyLogFiles(fileList)
		impl.PrintLogFileList(logFileList, getLogCmdFormat)
	} else {
		printErrorForArtifactList("log files", err)
	}
}

func executeDownloadLogFile(targetDirectory, logFileName string) {
	logFile, err := impl.GetLogFile(getLogCmdEnvironment, logFileName)
	if err == nil {
		impl.WriteLogFile(logFile, targetDirectory+"/"+logFileName)
	} else {
		printErrorForArtifact("log file", logFileName, err)
	}
}
