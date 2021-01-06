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

package credentials

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

// MiCredential for storing mi user details
type MiCredential struct {
	// Username of mi user
	Username string `json:"username"`
	// Password of mi user
	Password string `json:"password"`
	// AccessToken of mi
	AccessToken string `json:"accessToken"`
}

// GetMICredentials returns credentials for mi
func GetMICredentials(env string) (MiCredential, error) {

	store, err := GetDefaultCredentialStore()
	if err != nil {
		return MiCredential{}, err
	}

	if !utils.MIExistsInEnv(env, utils.MainConfigFilePath) {
		fmt.Println("MI does not exists in", env, "Add it using add env")
		os.Exit(1)
	}

	if !store.HasMI(env) {

		fmt.Println("Login to MI in", env)
		err = RunMILogin(store, env, "", "")
		if err != nil {
			return MiCredential{}, err
		}
		fmt.Println()
	}
	cred, err := store.GetMICredentials(env)
	if err != nil {
		return MiCredential{}, err
	}
	return cred, nil
}

// UpdateMIAccessToken updates the access token for mi
func UpdateMIAccessToken(env, accessToken string) error {

	store, err := GetDefaultCredentialStore()
	if err != nil {
		return err
	}
	cred, err := store.GetMICredentials(env)
	if err != nil {
		return err
	}
	return store.SetMICredentials(env, cred.Username, cred.Password, accessToken)
}

// GetOAuthAccessTokenForMI returns access token for mi
func GetOAuthAccessTokenForMI(username, password, env string) (string, error) {

	b64encodedCredentials := Base64Encode(username + ":" + password)

	tokenEndpoint := utils.GetMIManagementEndpointOfResource(utils.MiManagementMiLoginResource, env, utils.MainConfigFilePath)

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials

	resp, err := utils.InvokeGETRequest(tokenEndpoint, headers)
	utils.Logln(utils.LogPrefixInfo + "connecting to " + tokenEndpoint)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("Unable to connect to MI Token endpoint. Status: " + resp.Status())
	}
	responseDataMap := make(map[string]string)
	data := []byte(resp.Body())
	unmarshalError := json.Unmarshal(data, &responseDataMap)

	if unmarshalError != nil {
		return "", unmarshalError
	}
	if accessToken, ok := responseDataMap["AccessToken"]; ok {
		return accessToken, nil
	}
	return "", errors.New("AccessToken not found")
}

// RevokeAccessTokenForMI revokes the mi management token when the user log out from the environment
func RevokeAccessTokenForMI(env, token string) error {

	tokenRevokeEndpoint := utils.GetMIManagementEndpointOfResource(utils.MiManagementMiLogoutResource, env, utils.MainConfigFilePath)

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + token

	resp, err := utils.InvokeGETRequest(tokenRevokeEndpoint, headers)
	utils.Logln(utils.LogPrefixInfo + "connecting to " + tokenRevokeEndpoint)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("Error logging out of the MI in " + env + " Status: " + resp.Status())
	}

	return nil
}

// RunMILogin prompt user to input MI management API username and password
func RunMILogin(store Store, environment, username, password string) error {

	if username == "" {
		fmt.Print("Username:")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			username = scanner.Text()
		}
	}

	if password == "" {
		fmt.Print("Password:")
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		password = string(pass)
		fmt.Println()
	}

	accessToken, err := GetOAuthAccessTokenForMI(username, password, environment)
	if err != nil {
		return err
	}

	fmt.Println("Logged into MI in", environment, "environment")
	err = store.SetMICredentials(environment, username, password, accessToken)
	if err != nil {
		return err
	}
	return nil
}

// RunMILogout revoke mi management token and remove credentials from the store
func RunMILogout(environment string) error {
	cred, err := GetMICredentials(environment)
	if err != nil {
		return err
	}
	err = RevokeAccessTokenForMI(environment, cred.AccessToken)
	if err != nil {
		return err
	}
	store, err := GetDefaultCredentialStore()
	if err != nil {
		return err
	}
	fmt.Println("Logged out from MI in", environment, "environment")
	return store.EraseMI(environment)
}

// HandleMissingCredentials check for missing credentials and prompt to enter credentials or print error and exit
func HandleMissingCredentials(env string) {
	_, err := GetMICredentials(env)
	if err != nil {
		utils.HandleErrorAndExit("Error getting credentials", err)
	}
}
