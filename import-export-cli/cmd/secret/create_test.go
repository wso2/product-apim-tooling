/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
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
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/magiconair/properties"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func TestSecretCreateSymmetricUsesInitializedEncryptionKey(t *testing.T) {
	tempDir := t.TempDir()
	oldWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() returned an error: %v", err)
	}
	if err = os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir() returned an error: %v", err)
	}
	defer os.Chdir(oldWorkingDir)

	oldConfigDirPath := utils.ConfigDirPath
	defer func() {
		utils.ConfigDirPath = oldConfigDirPath
	}()
	utils.ConfigDirPath = filepath.Join(tempDir, ".wso2apictl")

	oldInputPropertiesFile := inputPropertiesfile
	oldEncryptionAlgorithm := encryptionAlgorithm
	oldOutputType := outputType
	oldSymmetricFormat := symmetricFormat

	cipherFlag := secretCreateCmd.Flags().Lookup(cipherFlagLiteral)
	if cipherFlag == nil {
		t.Fatalf("expected %q flag to be registered on secretCreateCmd", cipherFlagLiteral)
	}
	oldCipherFlagChanged := cipherFlag.Changed
	defer func() {
		inputPropertiesfile = oldInputPropertiesFile
		encryptionAlgorithm = oldEncryptionAlgorithm
		outputType = oldOutputType
		symmetricFormat = oldSymmetricFormat
		cipherFlag.Changed = oldCipherFlagChanged
	}()

	plainTextKey := "12345678901234567890123456789012"
	encryptionKeyConfig := &utils.EncryptionKeyConfig{
		Algorithm:     utils.SecretEncryptionAlgorithmAESGCM,
		EncryptionKey: base64.StdEncoding.EncodeToString([]byte(plainTextKey)),
	}
	utils.CreateDirIfNotExist(utils.GetEncryptionKeyDirectoryPath())
	utils.WriteConfigFile(encryptionKeyConfig, utils.GetEncryptionKeyConfigFilePath())

	propertiesFilePath := filepath.Join(tempDir, "secrets.properties")
	if err = os.WriteFile(propertiesFilePath, []byte("sample.alias=super-secret\n"), 0644); err != nil {
		t.Fatalf("WriteFile() returned an error: %v", err)
	}

	inputPropertiesfile = propertiesFilePath
	outputType = "file"
	encryptionAlgorithm = ""
	cipherFlag.Changed = false
	if err = validateSymmetricModeCreateArg(secretCreateCmd, []string{symmetricModeLiteral}); err != nil {
		t.Fatalf("validateSymmetricModeCreateArg() returned an error: %v", err)
	}
	secretCreateCmd.Run(secretCreateCmd, []string{symmetricModeLiteral})

	if encryptionAlgorithm != utils.SecretEncryptionAlgorithmAESGCM {
		t.Fatalf("expected symmetric shortcut to default to %s, got %s", utils.SecretEncryptionAlgorithmAESGCM, encryptionAlgorithm)
	}

	encryptedSecretsFilePath := filepath.Join(tempDir, "security", "wso2-secrets.properties")
	encryptedSecrets := properties.MustLoadFile(encryptedSecretsFilePath, properties.UTF8)
	encryptedValue, ok := encryptedSecrets.Get("sample.alias")
	if !ok {
		t.Fatalf("expected encrypted value for alias %q in %s", "sample.alias", encryptedSecretsFilePath)
	}

	keyBytes, err := utils.ResolveAES256Key(plainTextKey)
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}
	decryptedValue, err := utils.DecryptAES256(keyBytes, encryptedValue)
	if err != nil {
		t.Fatalf("DecryptAES256() returned an error: %v", err)
	}
	if decryptedValue != "super-secret" {
		t.Fatalf("expected decrypted value %q, got %q", "super-secret", decryptedValue)
	}
}
