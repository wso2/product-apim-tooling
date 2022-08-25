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

package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func ImportAPIPolicyToEnv(accessOAuthToken, importEnvironment, importPath string) error {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	if _, err := os.Stat(importPath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	publisherEndpoint = utils.AppendSlashToString(publisherEndpoint)
	uri := publisherEndpoint + "operation-policies/import"
	err := importAPIPolicy(uri, importPath, accessOAuthToken, true)
	return err
}

func importAPIPolicy(endpoint string, importPath string, accessToken string, isOauth bool) error {
	exportDirectory := filepath.Join(utils.ExportDirectory, utils.ExportedPoliciesDirName, utils.ExportedAPIPoliciesDirName)

	resolvedPolicyFilePath, err := resolvePolicyImportFilePath(importPath, exportDirectory)
	if err != nil {
		return err
	}

	utils.Logln(utils.LogPrefixInfo+"Policy Location:", resolvedPolicyFilePath)

	utils.Logln(utils.LogPrefixInfo + "Creating workspace")
	tmpPath, err := utils.GetTempCloneFromDirOrZip(resolvedPolicyFilePath)

	if err != nil {
		return err
	}

	defer func() {
		utils.Logln(utils.LogPrefixInfo+"Deleting", tmpPath)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			utils.Logln(utils.LogPrefixError + err.Error())
		}
	}()

	files, err := ioutil.ReadDir(tmpPath)

	if err != nil {
		return err
	}

	policyPaths := strings.Split(tmpPath, "/")
	policyName := policyPaths[len(policyPaths)-1]

	for _, file := range files {
		originalFilePath := tmpPath + "/" + file.Name()
		ext := filepath.Ext(originalFilePath)

		expectedPolicyFileName := policyName + ext

		isExt := ext == ".yaml" || ext == ".yml" || ext == ".json" || ext == ".j2" || ext == ".gotmpl"

		if isExt && expectedPolicyFileName != file.Name() {
			errorTxt := file.Name() + " should be equivalent to the policy name " + policyName
			utils.HandleErrorAndExit("Policy Directory name and policy files are not consistent", errors.New(errorTxt))
		}
	}

	utils.Logln(utils.LogPrefixInfo + "Substituting environment variables in API Policy files...")
	err = replaceEnvVariablesInPolicies(tmpPath)
	if err != nil {
		return err
	}

	// if policyFilePath contains a directory, zip it. Otherwise, leave it as it is.
	policyFilePath, err, cleanupFunc := utils.CreateZipFileFromProject(tmpPath, false)
	if err != nil {
		return err
	}

	//cleanup the temporary artifacts once consuming the zip file
	if cleanupFunc != nil {
		defer cleanupFunc()
	}

	resp, err := executeAPIPolicyImportRequest(endpoint, policyFilePath, accessToken, isOauth)
	utils.Logf("Response : %v", resp)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}

	var errorResponse utils.HttpErrorResponse

	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		// 201 Created or 200 OK
		fmt.Println("Successfully Imported API Policy.")
		return nil
	} else if resp.StatusCode() == http.StatusConflict {

		err := json.Unmarshal(resp.Body(), &errorResponse)

		if err != nil {
			return err
		}

		fmt.Println("Error importing API Policy due to: ", errorResponse.Description)
		fmt.Println("Please change the Policy name and re-import")

		if err != nil {
			return err
		}

		return errors.New(errorResponse.Status)
	} else {
		fmt.Println("Error importing API Policy.")
		fmt.Println("Status: " + resp.Status())
		fmt.Println("Response:", resp.IsSuccess())

		err := json.Unmarshal(resp.Body(), &errorResponse)

		if err != nil {
			return err
		}

		return errors.New(errorResponse.Status)
	}

	return nil
}

func executeAPIPolicyImportRequest(uri string, importPath string, accessToken string, isOAuthToken bool) (*resty.Response, error) {
	fileParamName := "file"

	headers := make(map[string]string)
	if isOAuthToken {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	} else {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + accessToken
	}
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive
	return utils.InvokePOSTRequestWithFile(uri, headers, fileParamName, importPath)
}

// resolveImportFilePath resolves the archive/directory for importing API policy
// First will resolve in given path, if not found will try to load from exported directory
func resolvePolicyImportFilePath(file, defaultExportDirectory string) (string, error) {
	// check current path
	utils.Logln(utils.LogPrefixInfo + "Resolving for Policy path...")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// if the file not in given path it might be inside exported directory
		utils.Logln(utils.LogPrefixInfo+"Looking for Policy in", defaultExportDirectory)
		file = filepath.Join(defaultExportDirectory, file)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return "", err
		}
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// Substitutes environment variables in the project files.
func replaceEnvVariablesInPolicies(policyFilePath string) error {
	for _, replacePath := range utils.EnvReplaceFilePaths {
		absFile := filepath.Join(policyFilePath, replacePath)
		// check if the path exists. If exists, proceed with processing. Otherwise, continue with the next items
		if fi, err := os.Stat(absFile); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				utils.Logln(utils.LogPrefixInfo+"Substituting env variables of files in folder path: ", absFile)
				if strings.EqualFold(replacePath, utils.InitProjectSequences) {
					err = utils.EnvSubstituteInFolder(absFile, utils.EnvReplacePoliciesFileExtensions)
				} else {
					err = utils.EnvSubstituteInFolder(absFile, nil)
				}
			case mode.IsRegular():
				utils.Logln(utils.LogPrefixInfo+"Substituting env of file: ", absFile)
				err = utils.EnvSubstituteInFile(absFile, nil)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetAPIPolicyId Get the ID of an API Policy if available
// @param accessToken : Token to call the Publisher Rest API
// @param environment : Environment where API policy needs to be located
// @param policyName : Name of the API policy
// @param policyVersion : Version of the API policy
// @return apiId, error
func GetAPIPolicyId(accessToken, environment, policyName, policyVersion string) (string, error) {
	apiPolicyEndpoint := utils.GetPublisherEndpointOfEnv(environment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	queryParams := `query=name:` + policyName + ` version:` + policyVersion

	apiPolicyEndpoint = utils.AppendSlashToString(apiPolicyEndpoint)

	apiPolicyResource := "operation-policies"

	url := apiPolicyEndpoint + apiPolicyResource

	utils.Logln(utils.LogPrefixInfo+"GetAPIPolicy: URL:", url)

	resp, err := utils.InvokeGETRequestWithQueryParamsString(url, queryParams, headers)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() == http.StatusOK {
		policyData := &utils.APIPoliciesList{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &policyData)
		if policyData.List[0].Id != "" {
			return policyData.List[0].Id, err
		}

		return "", errors.New("Requested API Policy is not available in the Publisher. Policy: " + policyName +
			" Version: " + policyVersion)
	} else if resp.StatusCode() == http.StatusNotFound {
		var errorResponse utils.HttpErrorResponse
		err := json.Unmarshal(resp.Body(), &errorResponse)

		if err != nil {
			return "", err
		}
		return "", errors.New(errorResponse.Description)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while getting API Policy " + policyName)
		}
		return "", errors.New("Request didn't respond 200 OK for getting API Policy. Status: " + resp.Status())
	}
}
