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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	wso2Api "github.com/wso2/k8s-apim-operator/apim-operator/pkg/controller/api"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagApiOperatorFile string
var isLocalInstallation bool

// installOperatorCmd represents the installOperator command
var installOperatorCmd = &cobra.Command{
	Use:   "install-operator",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		isLocalInstallation = flagApiOperatorFile != ""
		if !isLocalInstallation {
			installOLM("0.13.0")
		}
		installApiOperator()

		createDockerConfigSecret("https://registry-1.docker.io/", "renuka", "hello this is a dummy :D")
	},
}

// installOLM installs Operator Lifecycle Manager (OLM) with the given version
func installOLM(version string) {
	cmdString := fmt.Sprintf("curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/%s/install.sh | bash -s %s", version, version)
	cmd := exec.Command("bash", "-c", cmdString) //TODO: renuka: remove bash ?

	//print curl and install.sh errors
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

	if err := cmd.Start(); err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}
}

// installApiOperator installs WSO2 api-operator
func installApiOperator() {
	operatorFile := flagApiOperatorFile
	if isLocalInstallation {
		operatorFile = "https://operatorhub.io/install/api-operator.yaml" //TODO: renuka move to const
	}

	// Install the operator by running the following command
	cmd := exec.Command(
		utils.Kubectl,
		"create",
		"-f",
		operatorFile,
	)

	//print kubernetes errors
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

	if err := cmd.Run(); err != nil {
		utils.HandleErrorAndExit("Error installing WSO2 api-operator", err)
	}
}

func createDockerConfigSecret(registryUrl string, username string, password string) {
	encodedCredential := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	auth := wso2Api.Auth{Auths: map[string]wso2Api.Credential{registryUrl: wso2Api.Credential{Auth: encodedCredential}}}
	authJsonByte, err := json.Marshal(auth)
	if err != nil {
		utils.HandleErrorAndExit("Error marshalling docker secret credentials ", err)
	}

	//write configmap to a temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "docker-config-secret-*.json")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err = tmpFile.Write(authJsonByte); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}
	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
	// execute kubernetes command to create secret for accessing registry
	cmd := exec.Command(
		utils.Kubectl,
		utils.Create,
		"secret",
		"docker-config", //TODO: renuka make a constant
		"--from-file",
		tmpFile.Name(),
	)
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	errAddApi := cmd.Run()
	if errAddApi != nil {
		fmt.Println(errAddApi)
	}
}

// init using Cobra
func init() {
	RootCmd.AddCommand(installOperatorCmd)

	installOperatorCmd.Flags().StringVarP(&flagApiOperatorFile, "from-file", "f", "", "Path to API Operator directory")
}
