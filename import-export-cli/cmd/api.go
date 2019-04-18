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
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/go-resty/resty"
	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var apiUsername = ""
var apiPassword = ""
var apiEnvironment = ""
var apiLifCycleAction = ""
var apiName = ""
var apiVersion = ""
var apiID = ""

// API command related usage Info
const apiCmdLiteral = "api"
const apiCmdShortDesc = "Control API related stuffs"

var apiCmdLongDesc = dedent.Dedent(`
			allows to control an API given by id or name/version
	`)

var apiCmdExamples = dedent.Dedent(`
		Examples:
		apimcli api lifecycle -e dev --action="Publish" -u admin -p admin -k
	`)

// Status command related stuffs
const lifecycleCmdLiteral = "lifecycle"
const lifecycleCmdShortDesc = "Change lifecycle of an API"

var lifecycleCmdLongDesc = dedent.Dedent(`
			allows to control an API given by id or name/version
	`)

var lifecycleCmdExamples = dedent.Dedent(`
		Examples:
		apimcli api lifecycle -e dev --id="your-api-id" --action="Deploy as a Prototype" -u admin -p admin -k
		apimcli api lifecycle -e dev --name="your-api-name" --version="api-version" --action="Deploy as a Prototype" -u admin -p admin -k
	`)

// lifecycleCmd represents lifecycle command
var lifecycleCmd = &cobra.Command{
	Use:   lifecycleCmdLiteral + " --action=<action to take> -e <environment> [flags]",
	Short: lifecycleCmdShortDesc,
	Long:  lifecycleCmdLongDesc + lifecycleCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + lifecycleCmdLiteral + " called")
		err := changeLifecycle()
		if err != nil {
			utils.HandleErrorAndExit("Error", err)
		}
	},
}

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   apiCmdLiteral,
	Short: apiCmdShortDesc,
	Long:  apiCmdLongDesc + apiCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apiCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	lifecycleCmd.Flags().StringVarP(&apiEnvironment, "environment", "e",
		utils.DefaultEnvironmentName, "Environment from the which the API should be imported")
	lifecycleCmd.Flags().StringVarP(&apiUsername, "username", "u", "", "Username")
	lifecycleCmd.Flags().StringVarP(&apiPassword, "password", "p", "", "Password")
	lifecycleCmd.Flags().StringVarP(&apiName, "name", "", "", "API name i.e. PizzaShackAPI")
	lifecycleCmd.Flags().StringVarP(&apiID, "id", "", "", "API ID")
	lifecycleCmd.Flags().StringVarP(&apiVersion, "version", "", "", "API version i.e. 1.0.0")
	lifecycleCmd.Flags().StringVarP(&apiLifCycleAction, "action", "", "", "Action to be "+
		"taken. Supported actions: Publish, Deploy as a Prototype, Demote to Created, Demote to Prototyped, Block, Deprecate, "+
		"Re-Publish, Retire")

	apiCmd.AddCommand(lifecycleCmd)
	RootCmd.AddCommand(apiCmd)
}

// change lifecycle of an api
func changeLifecycle() error {
	// get access token
	accessOAuthToken, err :=
		utils.ExecutePreCommandWithOAuth(apiEnvironment, apiUsername, apiPassword,
			utils.MainConfigFilePath, utils.EnvKeysAllFilePath)
	if err != nil {
		return err
	}

	// check for lifecycle
	if apiLifCycleAction == "" {
		return errors.New("api lifecycle action is not defined")
	}
	// check for apiID
	if apiID == "" {
		// if not provided check whether name and version is there
		if apiName == "" || apiVersion == "" {
			return errors.New("either set apiID or name and version")
		}
		// get apiID from api manager
		apiID, err = getApiID(apiName, apiVersion, "", apiEnvironment, accessOAuthToken)
		if err != nil {
			return err
		}
	}

	utils.Logln(utils.LogPrefixInfo + "Changing api lifecycle to " + apiLifCycleAction)
	err = changeAPIStatusByID(apiID, apiLifCycleAction, apiEnvironment, accessOAuthToken)
	if err != nil {
		return err
	}
	fmt.Println("Successfully changed lifecycle to " + apiLifCycleAction)
	return nil
}

// changeAPIStatusByID changes status of an api given by apiID for the environment
func changeAPIStatusByID(apiID, status, environment, accessToken string) error {
	if apiID == "" {
		return errors.New("api ID is not set")
	}
	if status == "" {
		return errors.New("status is not set")
	}

	// get endpoint for API
	endpointString := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	endpoint, err := url.Parse(endpointString)
	if err != nil {
		return err
	}
	endpoint.Path = path.Join(endpoint.Path, "/change-lifecycle")

	resty.SetTimeout(time.Duration(utils.HttpRequestTimeout) * time.Millisecond)
	if utils.Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	}
	resp, err := resty.R().SetQueryParams(map[string]string{
		"apiId":  apiID,
		"action": status,
	}).SetAuthToken(accessToken).Post(endpoint.String())
	if err != nil {
		return err
	}

	// check if response is 200
	if resp.StatusCode() == 200 {
		return nil
	}
	return errors.New(string(resp.Body()))
}
