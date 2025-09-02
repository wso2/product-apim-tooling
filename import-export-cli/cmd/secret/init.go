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

package secret

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

const secretInitCmdLiteral = "init"
const secretInitCmdShortDesc = "Initialize Key Store"

const secretInitCmdLongDesc = "Initialize the Key Store information required for secret encryption"

var secretInitCmdExamples = "To initialize a Key Store information\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + secretCmdLiteral + " " + secretInitCmdLiteral + "\n" +
	"NOTE: Secret encryption supports JKS and PKCS12 Key Stores (.jks, .p12, .pfx)"

var secretInitCmd = &cobra.Command{
	Use:     secretInitCmdLiteral,
	Short:   secretInitCmdShortDesc,
	Long:    secretInitCmdLongDesc,
	Example: secretInitCmdExamples,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		startConsoleForKeyStore()
	},
}

func init() {
	SecretCmd.AddCommand(secretInitCmd)
}

func startConsoleForKeyStore() {
	reader := bufio.NewReader(os.Stdin)
	keyStoreConfig := &utils.KeyStoreConfig{}

	fmt.Printf("Enter Key Store location: ")
	path, _ := reader.ReadString('\n')
	if !isValidKeyStore(path) {
		utils.HandleErrorAndExit("Invalid Key Store Type. Supports only JKS and PKCS12 Key Stores (.jks, .p12, .pfx)", nil)
	}
	keyStoreConfig.KeyStorePath = strings.TrimSpace(path)

	fmt.Printf("Enter Key Store password: ")
	byteStorePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	keyStorePassword := string(byteStorePassword)
	fmt.Println()
	keyStoreConfig.KeyStorePassword = base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(keyStorePassword)))

	fmt.Printf("Enter Key alias: ")
	alias, _ := reader.ReadString('\n')
	keyStoreConfig.KeyAlias = strings.TrimSpace(alias)

	fmt.Printf("Enter Key password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	keyPassword := string(bytePassword)
	fmt.Println()
	keyStoreConfig.KeyPassword = base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(keyPassword)))

	if utils.IsValidKeyStoreConfig(keyStoreConfig) {
		utils.CreateDirIfNotExist(utils.GetKeyStoreDirectoryPath())
		keyStoreConfigFilePath := utils.GetKeyStoreConfigFilePath()
		utils.WriteConfigFile(keyStoreConfig, keyStoreConfigFilePath)
		fmt.Println("Key Store initialization completed.")
	} else {
		fmt.Println("Key Store initialization failed.")
	}
}

func updateMap(params map[string]string, key, value string) {
	params[key] = strings.TrimSpace(value)
}

func isJKSKeyStore(path string) bool {
	return filepath.Ext(strings.TrimSpace(path)) == ".jks"
}

func isPKCS12KeyStore(path string) bool {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(path)))
	return ext == ".p12" || ext == ".pfx"
}

func isValidKeyStore(path string) bool {
	return isJKSKeyStore(path) || isPKCS12KeyStore(path)
}
