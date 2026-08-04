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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var inputPropertiesfile string
var encryptionAlgorithm string
var outputType string
var symmetricFormat string

const secretCreateCmdLiteral = "create"
const secretCreateCmdShortDesc = "Encrypt secrets"

const secretCreateCmdLongDesc = "Create secrets based on given arguments"

var secretCreateCmdExamples = "To encrypt secret and get output on console\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + "\n" +
	"To encrypt secret using an initialized symmetric encryption key and get output on console\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " " + symmetricModeLiteral + "\n" +
	"To encrypt secret in the carbon-crypto-service ciphertext format\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " " + symmetricModeLiteral + " " + symmetricInternalModeLiteral + "\n" +
	"To encrypt secret in the cipher-tool/carbon-secvault/carbon-mediation ciphertext format\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " " + symmetricModeLiteral + " " + symmetricExternalModeLiteral + "\n" +
	"To encrypt secret and get output as a .properties file (stored in the security folder in apictl executable directory)\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -o file\n" +
	"To encrypt secret and get output as a .yaml file (stored in the security folder in apictl executable directory)\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -o k8\n" +
	"To bulk encrypt secrets defined in a properties file\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -f <file_path>\n" +
	"To bulk encrypt secrets defined in a properties file and get a .yaml file (stored in the security folder in apictl executable directory)\n" +
	"  " + utils.ProjectName + " " + secretCmdLiteral + " " + secretCreateCmdLiteral + " -o k8 -f <file_path>"

var secretCreateCmd = &cobra.Command{
	Use: secretCreateCmdLiteral + " [" + symmetricModeLiteral + " [" + symmetricInternalModeLiteral +
		"|" + symmetricExternalModeLiteral + "]]",
	Short:   secretCreateCmdShortDesc,
	Long:    secretCreateCmdLongDesc,
	Example: secretCreateCmdExamples,
	Args:    validateSymmetricModeCreateArg,
	Run: func(cmd *cobra.Command, args []string) {
		resolveCreateCipher(cmd, args)
		symmetricFormat = ""
		if len(args) == 2 {
			symmetricFormat = args[1]
		}
		err := validateFlags()
		if err != nil {
			utils.HandleErrorAndExit("Invalid flag", err)
		}
		var keyStoreConfig *utils.KeyStoreConfig
		var encryptionKeyConfig *utils.EncryptionKeyConfig
		if utils.IsAES256Encryption(encryptionAlgorithm) {
			encryptionKeyConfig, err = utils.GetEncryptionKeyConfigFromFile(utils.GetEncryptionKeyConfigFilePath())
			if err != nil {
				utils.HandleErrorAndExit("Encryption key has not been initialized.", err)
			}
			if encryptionKeyConfig == nil || !utils.IsValidSymmetricEncryptionConfig(encryptionKeyConfig) {
				utils.HandleErrorAndExit("Encryption key has not been initialized.", nil)
			}
		} else {
			keyStoreConfig, err = utils.GetKeyStoreConfigFromFile(utils.GetKeyStoreConfigFilePath())
			if err != nil {
				utils.HandleErrorAndExit("Key Store has not been initialized.", err)
			}
			if keyStoreConfig == nil || !utils.IsValidKeyStoreConfig(keyStoreConfig) {
				utils.HandleErrorAndExit("Key Store has not been initialized.", nil)
			}
		}
		initSecretInformation(keyStoreConfig, encryptionKeyConfig)
	},
}

func init() {
	SecretCmd.AddCommand(secretCreateCmd)
	secretCreateCmd.Flags().StringVarP(&inputPropertiesfile, "from-file", "f", "", "Path to the properties file which contains secrets to be encrypted")
	secretCreateCmd.Flags().StringVarP(&outputType, "output", "o", "console", "Get the output in yaml (k8) or properties (file) format. By default the output is printed to the console")
	secretCreateCmd.Flags().StringVarP(&encryptionAlgorithm, "cipher", "c", "", "Encryption algorithm. Supports RSA/ECB/OAEPWithSHA1AndMGF1Padding, RSA/ECB/PKCS1Padding, AES/GCM/NoPadding and AES256")
}

func resolveCreateCipher(cmd *cobra.Command, args []string) {
	if cmd.Flags().Changed(cipherFlagLiteral) {
		return
	}
	if len(args) >= 1 && args[0] == symmetricModeLiteral {
		encryptionAlgorithm = utils.SecretEncryptionAlgorithmAESGCM
		return
	}
	encryptionAlgorithm = utils.SecretEncryptionAlgorithmRSAOAEP
}

func initSecretInformation(keyStoreConfig *utils.KeyStoreConfig, encryptionKeyConfig *utils.EncryptionKeyConfig) {
	secretConfig := utils.SecretConfig{
		OutputType:      outputType,
		Algorithm:       encryptionAlgorithm,
		SymmetricFormat: symmetricFormat,
	}
	if utils.IsAES256Encryption(secretConfig.Algorithm) {
		encryptionKey, err := utils.GetStoredEncryptionKey(encryptionKeyConfig)
		if err != nil {
			utils.HandleErrorAndExit("Invalid stored encryption key.", err)
		}
		secretConfig.EncryptionKey = encryptionKey
	}
	if isNonEmptyString(inputPropertiesfile) {
		secretConfig.InputType = "file"
		secretConfig.InputFile = inputPropertiesfile
	} else {
		secretConfig.InputType = "console"
		startConsoleForSecretInfo(&secretConfig)
	}
	err := utils.EncryptSecrets(keyStoreConfig, secretConfig)
	if err != nil {
		utils.HandleErrorAndExit("Error encrypting secrets.", err)
	}
}

func startConsoleForSecretInfo(secretConfig *utils.SecretConfig) {
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
	if !(utils.IsOAEPEncryption(encryptionAlgorithm) || utils.IsPKCS1Encryption(encryptionAlgorithm) ||
		utils.IsAES256Encryption(encryptionAlgorithm)) {
		return errors.New("Accepts RSA/ECB/OAEPWithSHA1AndMGF1Padding, RSA/ECB/PKCS1Padding, AES/GCM/NoPadding or AES256 as encryption algorithms (-c)")
	}
	if !(utils.IsConsole(outputType) || utils.IsFile(outputType) || utils.IsK8(outputType)) {
		return errors.New("Accepts k8, file or console as output formats (-o)")
	}
	return nil
}

func validateSymmetricModeCreateArg(cmd *cobra.Command, args []string) error {
	if len(args) > 2 {
		return cobra.MaximumNArgs(2)(cmd, args)
	}
	if len(args) >= 1 && args[0] != symmetricModeLiteral {
		return errors.New("accepts only '" + symmetricModeLiteral + "' as the first optional argument")
	}
	if len(args) == 2 && args[1] != symmetricInternalModeLiteral && args[1] != symmetricExternalModeLiteral {
		return errors.New("accepts only '" + symmetricInternalModeLiteral + "' or '" + symmetricExternalModeLiteral +
			"' as the second optional argument")
	}
	if len(args) >= 1 && args[0] == symmetricModeLiteral && cmd.Flags().Changed(cipherFlagLiteral) &&
		!utils.IsAES256Encryption(encryptionAlgorithm) {
		return errors.New("the optional argument '" + symmetricModeLiteral + "' only supports AES/GCM/NoPadding or AES256 with -c")
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
