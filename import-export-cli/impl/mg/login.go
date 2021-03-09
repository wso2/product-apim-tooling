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

package mg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

type MgwAdapterInfo struct {
	Endpoint    string
	AccessToken string
}

func RunLogin(environment, loginUsername, loginPassword string, loginPasswordStdin bool) error {
	store, err := credentials.GetDefaultCredentialStore()
	if err != nil {
		utils.HandleErrorAndExit("Error occurred while loading credential store : ", err)
	}
	mgwAdapterEndpoints, err := utils.GetEndpointsOfMgwAdapterEnv(environment, utils.MainConfigFilePath)
	if err != nil {
		return errors.New("Env " + environment + " does not exists. Add it using `apictl mg add env`")
	}

	if loginUsername == "" {
		fmt.Print("Username: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			loginUsername = scanner.Text()
		}
	}
	if loginPassword != "" {
		fmt.Println("Warning: Using --password in CLI is not secure. Use --password-stdin")
		if loginPasswordStdin {
			return errors.New("--password and --password-stdin are mutually exclusive")
		}
	}
	if loginPasswordStdin {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.New("Error reading password. Cause: " + err.Error())
		}
		loginPassword = strings.TrimRight(strings.TrimSuffix(string(data), "\n"), "\r")
	}
	if loginPassword == "" {
		fmt.Print("Enter Password: ")
		loginPasswordB, err := terminal.ReadPassword(0)
		loginPassword = string(loginPasswordB)
		fmt.Println()
		if err != nil {
			return errors.New("Error reading password. Cause: " + err.Error())
		}
	}

	if loginUsername == "" || loginPassword == "" {
		return errors.New("username or password not entered")
	}

	tokenEndpoint := deriveTokenEndpointForMGAdapter(mgwAdapterEndpoints.AdapterEndpoint)
	accessToken, err := getAccessTokenFromMGAdapter(loginUsername, loginPassword, tokenEndpoint)
	if err != nil {
		utils.HandleErrorAndExit("Error getting access token from adapter endpoint: "+
			tokenEndpoint, err)
	}

	if err = store.SetMGToken(environment, accessToken); err != nil {
		return err
	}
	fmt.Println("Successfully logged into Microgateway Adapter in environment: ", environment)
	return nil
}

func getAccessTokenFromMGAdapter(username, password, tokenEndpoint string) (string, error) {
	body := make(map[string]string)
	body["username"] = username
	body["password"] = password

	headers := make(map[string]string)
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON

	resp, err := utils.InvokePOSTRequest(tokenEndpoint, headers, body)
	if err != nil {
		return "", errors.New("Unable to connect to Microgateway Token endpoint. " + err.Error())
	}
	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("Error response from Microgateway Token endpoint. Status: " + resp.Status())
	}
	accessToken, err := getAccessTokenFromResponse(resp.Body())
	if err != nil {
		return "", err
	}
	return accessToken, err
}

func getAccessTokenFromResponse(responseBody []byte) (string, error) {
	responseDataMap := make(map[string]string)
	data := []byte(responseBody)
	unmarshalError := json.Unmarshal(data, &responseDataMap)
	if unmarshalError != nil {
		return "", unmarshalError
	}
	accessToken, exists := responseDataMap["accessToken"]
	if !exists {
		return "", errors.New("accessToken not found in the response")
	}
	return accessToken, nil
}
