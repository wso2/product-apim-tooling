/*
*  Copyright (c) 2025 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 LLC. licenses this file to you under the Apache License,
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

package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// importMCPServer imports an MCP Server to the API manager
func importMCPServer(endpoint, filePath, accessToken string, extraParams map[string]string, isOauth bool, dryRun bool,
	mcpServerLoggingCmdFormat string) error {
	resp, err := ExecuteNewFileUploadRequest(endpoint, extraParams, "file",
		filePath, accessToken, isOauth)
	utils.Logf("Response : %v", resp)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}
	if dryRun {
		if resp.StatusCode() == http.StatusOK && resp.String() != "" {
			// 200 OK
			var data JSONData
			err := json.Unmarshal([]byte(resp.String()), &data)
			if err != nil {
				utils.Logln(utils.LogPrefixError, err)
				fmt.Println("Error occurred while validating MCP Server")
				return errors.New(resp.Status())
			}
			if data.ComplianceCheck.Result == "fail" {
				PrintViolations(data.ComplianceCheck.Violations, mcpServerLoggingCmdFormat)
			} else if resp.StatusCode() == http.StatusOK {
				fmt.Printf("No violations found for the MCP Server")
			}
		} else {
			// We have an HTTP error
			utils.Logln(utils.LogPrefixError, err)
			fmt.Println("Error occurred while validating MCP Server")
			return errors.New(resp.Status())
		}
	} else {
		if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
			// 201 Created or 200 OK
			fmt.Println("Successfully imported MCP Server.")
		} else {
			// We have an HTTP error
			utils.Logln(utils.LogPrefixError, err)
			fmt.Println("Status: " + resp.Status())
			fmt.Println("Response:", resp)
			return errors.New(resp.Status())
		}
	}
	return nil
}

// ImportMCPServerToEnv function is used with import mcp-server command
func ImportMCPServerToEnv(accessOAuthToken, environment, importFile, paramsFile string, update,
	preserveProvider, skipCleanup, rotateRevision, skipDeployments, dryRun bool,
	loggingCmdFormat string) error {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(environment, utils.MainConfigFilePath)
	return ImportMCPServer(accessOAuthToken, publisherEndpoint, environment, importFile, paramsFile, update,
		preserveProvider, skipCleanup, rotateRevision, skipDeployments, dryRun, loggingCmdFormat)
}

// ImportMCPServer function is used with import mcp-server command
func ImportMCPServer(accessOAuthToken, publisherEndpoint, importEnvironment, importPath, mcpServerParamsPath string, importMCPServerUpdate,
	preserveProvider, importMCPServerSkipCleanup, importMCPServerRotateRevision, importMCPServerSkipDeployments bool,
	dryRun bool, mcpServerLoggingCmdFormat string) error {
	exportDirectory := filepath.Join(utils.ExportDirectory, utils.ExportedMCPServersDirName)
	resolvedMCPServerFilePath, err := resolveImportFilePath(importPath, exportDirectory)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"MCP Server Location: ", resolvedMCPServerFilePath)

	utils.Logln(utils.LogPrefixInfo + "Creating workspace")
	tmpPath, err := utils.GetTempCloneFromDirOrZip(resolvedMCPServerFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if importMCPServerSkipCleanup {
			utils.Logln(utils.LogPrefixInfo+"Leaving", tmpPath)
			return
		}
		utils.Logln(utils.LogPrefixInfo+"Deleting", tmpPath)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			utils.Logln(utils.LogPrefixError + err.Error())
		}
	}()
	mcpServerFilePath := tmpPath

	utils.Logln(utils.LogPrefixInfo + "Substituting environment variables in MCP Server files...")
	err = replaceEnvVariables(mcpServerFilePath)
	if err != nil {
		return err
	}

	if importMCPServerSkipDeployments {
		//If skip deployments flag used, deployment_environments files will be removed from import artifacts
		loc := filepath.Join(mcpServerFilePath, utils.DeploymentEnvFile)
		utils.Logln(utils.LogPrefixInfo + "Removing the deployment environments file from " + loc)
		err := utils.RemoveFileIfExists(loc)
		if err != nil {
			return err
		}
	}

	if mcpServerParamsPath != "" {
		//Reading params file of the MCP Server and add configurations into temp artifact
		err := handleCustomizedParameters(mcpServerFilePath, mcpServerParamsPath, importEnvironment)
		if err != nil {
			return err
		}
	}

	// if mcpServerFilePath contains a directory, zip it. Otherwise, leave it as it is.
	mcpServerFilePath, err, cleanupFunc := utils.CreateZipFileFromProject(mcpServerFilePath, importMCPServerSkipCleanup)
	if err != nil {
		return err
	}

	//cleanup the temporary artifacts once consuming the zip file
	if cleanupFunc != nil {
		defer cleanupFunc()
	}

	extraParams := map[string]string{}
	publisherEndpoint = utils.AppendSlashToString(publisherEndpoint)
	publisherEndpoint += "mcp-servers/import"
	if importMCPServerUpdate {
		publisherEndpoint += "?overwrite=" + strconv.FormatBool(true) + "&preserveProvider=" +
			strconv.FormatBool(preserveProvider) + "&rotateRevision=" + strconv.FormatBool(importMCPServerRotateRevision)
	} else {
		publisherEndpoint += "?preserveProvider=" + strconv.FormatBool(preserveProvider) + "&rotateRevision=" +
			strconv.FormatBool(importMCPServerRotateRevision)
	}

	if dryRun {
		publisherEndpoint += "&dryRun=" + strconv.FormatBool(true)
	}
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + publisherEndpoint)

	err = importMCPServer(publisherEndpoint, mcpServerFilePath, accessOAuthToken, extraParams, true, dryRun, mcpServerLoggingCmdFormat)
	return err
}
