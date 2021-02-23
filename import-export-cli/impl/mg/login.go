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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

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

	accessToken, err := credentials.GetOAuthAccessTokenForMGAdapter(loginUsername, loginPassword,
		mgwAdapterEndpoints.AdapterEndpoint)
	if err != nil {
		utils.HandleErrorAndExit("Error getting access token from adapter", err)
	}

	err = store.SetMGToken(environment, accessToken)
	if err != nil {
		return err
	}
	fmt.Println("Successfully logged into Microgateway Adapter in environment: ", environment)
	return nil
}
