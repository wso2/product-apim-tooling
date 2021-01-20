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
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var keyStorePropertiesFile = miUtils.GetkeyStorePropertiesFilePath()
var inputPropertiesfile string
var encryptionAlgorithm string
var outputType string

const secretCreateCmdLiteral = "create"
const secretCreateCmdShortDesc = "Encrypt secrets"

const secretCreateCmdLongDesc = "Create secrets based on given arguments"

var secretCreateCmdExamples = "To encrypt secret and get output on console\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + "\n" +
	"To encrypt secret and get output as a .properties file (stored in the security folder in apictl executable directory)\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -o file\n" +
	"To encrypt secret and get output as a .yaml file (stored in the security folder in apictl executable directory)\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -o k8\n" +
	"To bulk encrypt secrets defined in a properties file\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -f <file_path>\n" +
	"To bulk encrypt secrets defined in a properties file and get a .yaml file (stored in the security folder in apictl executable directory)\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -o k8 -f <file_path>"

var secretCreateCmd = &cobra.Command{
	Use:     secretCreateCmdLiteral,
	Short:   secretCreateCmdShortDesc,
	Long:    secretCreateCmdLongDesc,
	Example: secretCreateCmdExamples,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.IsFileExist(keyStorePropertiesFile) {
			utils.HandleErrorAndExit("Key Store has not been initialized.\nExecute 'apictl mi secret init --help' for more information", nil)
		}
		err := validateFlags()
		if err != nil {
			utils.HandleErrorAndExit("Invalid flag", err)
		}
		initSecretInformation()
	},
}

func init() {
	SecretCmd.AddCommand(secretCreateCmd)
	secretCreateCmd.Flags().StringVarP(&inputPropertiesfile, "from-file", "f", "", "Path to the properties file which contain secrets to be encrypted")
	secretCreateCmd.Flags().StringVarP(&outputType, "output", "o", "console", "Get the output in yaml(k8) or properties(file) format. By default the output is printed to the console")
	secretCreateCmd.Flags().StringVarP(&encryptionAlgorithm, "cipher", "c", "RSA/ECB/OAEPWithSHA1AndMGF1Padding", "Encryption algorithm")
}

func initSecretInformation() {
	secretConfig := miUtils.SecretConfig{
		OutputType: outputType,
		Algorithm:  encryptionAlgorithm,
	}
	if isNonEmptyString(inputPropertiesfile) {
		secretConfig.InputType = "file"
		secretConfig.InputFile = inputPropertiesfile
	} else {
		secretConfig.InputType = "console"
		startConsoleForSecretInfo(&secretConfig)
	}
	err := miUtils.EncryptSecrets(keyStorePropertiesFile, secretConfig)
	if err != nil {
		utils.HandleErrorAndExit("Error encrypting secrets.", err)
	}
}

func startConsoleForSecretInfo(secretConfig *miUtils.SecretConfig) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter plain alias for secret:")
	alias, _ := reader.ReadString('\n')
	secretConfig.PlainTextAlias = strings.TrimSpace(alias)

	fmt.Printf("Enter plain text secret:")
	byteSecret, _ := terminal.ReadPassword(int(syscall.Stdin))
	secret := string(byteSecret)
	fmt.Println()

	fmt.Printf("Repeat plain text secret:")
	byteRepeatSecret, _ := terminal.ReadPassword(int(syscall.Stdin))
	repeatSecret := string(byteRepeatSecret)
	fmt.Println()

	if isMatchingSecrets(secret, repeatSecret) {
		secretConfig.PlainTextSecretText = strings.TrimSpace(secret)
	} else {
		fmt.Println("Entered secret values did not match.")
		startConsoleForSecretInfo(secretConfig)
	}
}

func validateFlags() error {
	if !(miUtils.IsOAEPEncryption(encryptionAlgorithm) || miUtils.IsPKCS1Encryption(encryptionAlgorithm)) {
		return errors.New("Accepts RSA/ECB/OAEPWithSHA1AndMGF1Padding or RSA/ECB/PKCS1Padding as encryption algorithms (-c)")
	}
	if !(miUtils.IsConsole(outputType) || miUtils.IsFile(outputType) || miUtils.IsK8(outputType)) {
		return errors.New("Accepts k8, file or console as output formats (-o)")
	}
	return nil
}

func isMatchingSecrets(secret, repeatSecret string) bool {
	if secret == repeatSecret {
		return true
	}
	return false
}

func isNonEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) > 0
}
