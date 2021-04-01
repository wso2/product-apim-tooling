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

package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
	"github.com/pavel-v-chernykh/keystore-go/v4"
	"gopkg.in/yaml.v2"
)

const keystoreDirName = "keystore"
const keyStoreConfigFileName = "keystore_info.yaml"
const encryptedSecretsPropertiesFileName = "wso2-secrets.properties"
const encryptedSecretsYamlFileName = "wso2-secrets.yaml"

type k8sSecretConfig struct {
	APIVerion  string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	MetaData   metaData          `yaml:"metadata"`
	StringData map[string]string `yaml:"stringData"`
	Type       string            `yaml:"type"`
}

type metaData struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type SecretConfig struct {
	OutputType          string
	Algorithm           string
	InputType           string
	InputFile           string
	PlainTextAlias      string
	PlainTextSecretText string
}

type KeyStoreConfig struct {
	KeyStorePath     string `yaml:"keyStorePath"`
	KeyStorePassword string `yaml:"keyStorePassword"`
	KeyAlias         string `yaml:"keyAlias"`
	KeyPassword      string `yaml:"keyPassword"`
}

type encryptFunc func(key *rsa.PublicKey, plainText string) (string, error)

// IsValidKeyStoreConfig return true if the KeyStoreConfig is valid
func IsValidKeyStoreConfig(config *KeyStoreConfig) bool {
	if IsNonEmptyString(config.KeyStorePath) && IsNonEmptyString(config.KeyStorePassword) &&
		IsNonEmptyString(config.KeyAlias) && IsNonEmptyString(config.KeyPassword) {
		return true
	}
	return false
}

// EncryptSecrets encrypts the secrets using the keystore and write them to a file or console depending on the config map argument
func EncryptSecrets(keyStoreConfig *KeyStoreConfig, secretConfig SecretConfig) error {
	encryptionKey, err := getEncryptionKey(keyStoreConfig)
	if err != nil {
		return err
	}
	var encryptedSecrets map[string]string
	plainTextSecrets := getPlainTextSecrets(secretConfig)

	if IsPKCS1Encryption(secretConfig.Algorithm) {
		encryptedSecrets, err = encrypt(encryptionKey, plainTextSecrets, encryptPKCS1v15)
	} else {
		encryptedSecrets, err = encrypt(encryptionKey, plainTextSecrets, encryptOAEP)
	}
	if err != nil {
		return err
	}
	if IsK8(secretConfig.OutputType) {
		printSecretsToYamlFile(encryptedSecrets)
	} else if IsFile(secretConfig.OutputType) {
		printSecretsToPropertiesFile(encryptedSecrets)
	} else {
		printSecretsToConsole(encryptedSecrets)
	}
	return nil
}

// WritePropertiesToFile write a map to a .properties file
func WritePropertiesToFile(variables map[string]string, fileName string) {
	props := properties.LoadMap(variables)
	writer, err := os.Create(fileName)
	if err != nil {
		HandleErrorAndExit("Unable to create file.", err)
	}
	_, err = props.Write(writer, properties.UTF8)
	if err != nil {
		HandleErrorAndExit("Unable to write properties to file.", err)
	}
	writer.Close()
}

func readPropertiesFromFile(fileName string) map[string]string {
	props := properties.MustLoadFile(fileName, properties.UTF8)
	return props.Map()
}

// GetKeyStoreDirectoryPath join keystore with the config directory path
func GetKeyStoreDirectoryPath() string {
	return filepath.Join(ConfigDirPath, keystoreDirName)
}

// GetKeyStoreConfigFilePath join keystore-info.yaml with the keystore path
func GetKeyStoreConfigFilePath() string {
	return filepath.Join(GetKeyStoreDirectoryPath(), keyStoreConfigFileName)
}

// GetKeyStoreConfigFromFile read and return KeyStoreConfig
func GetKeyStoreConfigFromFile(filePath string) (*KeyStoreConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("Config file not found.\nExecute 'apictl secret init --help' for more information")
	}
	config := &KeyStoreConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.New("Parsing error.\nExecute 'apictl secret init --help' for more information")
	}
	if !IsValidKeyStoreConfig(config) {
		return nil, errors.New("Missing required fields.\nExecute 'apictl secret init --help' for more information")
	}
	return config, nil
}

func getEncryptionKey(keyStoreConfig *KeyStoreConfig) (*rsa.PublicKey, error) {
	keyStorePath := keyStoreConfig.KeyStorePath
	keyStorePassword, _ := base64.StdEncoding.DecodeString(keyStoreConfig.KeyStorePassword)
	keyStore, err := readKeyStore(keyStorePath, keyStorePassword)
	if err != nil {
		return nil, errors.New("Reading Key Store: " + err.Error())
	}
	keyAlias := keyStoreConfig.KeyAlias
	keyPassword, _ := base64.StdEncoding.DecodeString(keyStoreConfig.KeyPassword)
	pke, err := keyStore.GetPrivateKeyEntry(keyAlias, keyPassword)
	if err != nil {
		return nil, errors.New("Reading Key Entry: " + err.Error())
	}
	key, err := x509.ParsePKCS8PrivateKey(pke.PrivateKey)
	rsaKey := key.(*rsa.PrivateKey)
	if err != nil {
		return nil, errors.New("Parsing Key Entry: " + err.Error())
	}
	return &rsaKey.PublicKey, nil
}

func encrypt(encryptionKey *rsa.PublicKey, plainTextSecrets map[string]string, encryptFunction encryptFunc) (map[string]string, error) {
	var encryptedSecrets = make(map[string]string)
	for alias, plainText := range plainTextSecrets {
		encryptedSecret, err := encryptFunction(encryptionKey, plainText)
		if err != nil {
			return nil, err
		}
		encryptedSecrets[alias] = encryptedSecret
	}
	return encryptedSecrets, nil
}

func getPlainTextSecrets(secretConfig SecretConfig) map[string]string {
	var plainTexts = make(map[string]string)
	if IsFile(secretConfig.InputType) {
		plainTexts = readPropertiesFromFile(secretConfig.InputFile)
	} else {
		plainTexts[secretConfig.PlainTextAlias] = secretConfig.PlainTextSecretText
	}
	return plainTexts
}

func printSecretsToConsole(secrets map[string]string) {
	for alias, secret := range secrets {
		fmt.Println(alias, ":", secret)
	}
}

func printSecretsToPropertiesFile(secrets map[string]string) {
	secretFilePath := getSecretFilePath(encryptedSecretsPropertiesFileName)
	WritePropertiesToFile(secrets, secretFilePath)
	fmt.Println("Secret properties file created in", secretFilePath)
}

func printSecretsToYamlFile(secrets map[string]string) {
	secretConfig := k8sSecretConfig{
		APIVerion:  "v1",
		Kind:       "Secret",
		StringData: secrets,
		Type:       "Opaque",
		MetaData: metaData{
			Name:      "wso2secret",
		},
	}
	secretFilePath := getSecretFilePath(encryptedSecretsYamlFileName)
	WriteConfigFile(secretConfig, secretFilePath)
	fmt.Println("Kubernetes secret file created in", secretFilePath, "with default name and namespace")
	fmt.Println("You can change the default values as required before applying.")
}

func getSecretFilePath(fileName string) string {
	currentDir, _ := os.Getwd()
	secretDirPath := filepath.Join(currentDir, "security")
	CreateDirIfNotExist(secretDirPath)
	return filepath.Join(secretDirPath, fileName)
}

func encryptOAEP(key *rsa.PublicKey, plainText string) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, key, []byte(plainText), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func encryptPKCS1v15(key *rsa.PublicKey, plainText string) (string, error) {
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func readKeyStore(filename string, password []byte) (*keystore.KeyStore, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
	}()
	keyStore := keystore.New()
	if err := keyStore.Load(f, password); err != nil {
		return nil, err
	}
	return &keyStore, nil
}

// IsConsole return true if outputType is console
func IsConsole(outputType string) bool {
	return strings.EqualFold(outputType, "console")
}

// IsFile return true if outputType is file
func IsFile(outputType string) bool {
	return strings.EqualFold(outputType, "file")
}

// IsK8 return true if outputType is k8
func IsK8(outputType string) bool {
	return strings.EqualFold(outputType, "k8")
}

// IsPKCS1Encryption return true if the encryption algorithm is RSA/ECB/PKCS1Padding
func IsPKCS1Encryption(algorithm string) bool {
	return strings.EqualFold(algorithm, "RSA/ECB/PKCS1Padding")
}

// IsOAEPEncryption return true if the encryption algorithm is RSA/ECB/OAEPWithSHA1AndMGF1Padding
func IsOAEPEncryption(algorithm string) bool {
	return strings.EqualFold(algorithm, "RSA/ECB/OAEPWithSHA1AndMGF1Padding")
}

// IsNonEmptyString return true if the passed string is non empty
func IsNonEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) > 0
}
