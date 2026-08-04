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
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	AES256KeySize    = 32
	AES256HexKeySize = 64
	GCMIVSize        = 12
	GCMTagSize       = 128
	hexCharacters    = "0123456789abcdefABCDEF"
)

type cipherInitializationVectorHolder struct {
	Cipher               string `json:"cipher"`
	InitializationVector string `json:"initializationVector"`
	KeyId                string `json:"keyId,omitempty"`
}

type cipherMetaDataHolder struct {
	C  string `json:"c"`
	T  string `json:"t"`
	Iv string `json:"iv"`
}

// legacyCipherHolder is the flat ciphertext shape used by cipher-tool, carbon-secvault, and
// carbon-mediation's current secure vault.
type legacyCipherHolder struct {
	CipherText string `json:"cipherText"`
	Iv         string `json:"iv"`
}

// Returns md5 hash of a given string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		HandleErrorAndExit("Error in encryption", err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		HandleErrorAndExit("Error in encryption", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// Decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		HandleErrorAndExit("Error in decryption", err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		HandleErrorAndExit("Error in Decryption: Ciphertext too short", nil)
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

// ResolveAES256Key validates and converts a direct encryption key into raw key bytes.
func ResolveAES256Key(encryptionKey string) ([]byte, error) {
	trimmedKey := strings.TrimSpace(encryptionKey)
	if trimmedKey == "" {
		return nil, errors.New("Encryption key cannot be empty")
	}
	if isHexAES256Key(trimmedKey) {
		keyBytes, err := hex.DecodeString(trimmedKey)
		if err != nil {
			return nil, errors.New("Invalid hexadecimal characters found in encryption key")
		}
		return keyBytes, nil
	}
	keyBytes := []byte(trimmedKey)
	if len(keyBytes) != AES256KeySize {
		return nil, fmt.Errorf("Invalid AES key length: %d bytes. AES-256 requires a 32-byte (256-bit) key", len(keyBytes))
	}
	return keyBytes, nil
}

// newAES256GCMWithRandomIV builds an AES-256/GCM cipher for key and generates a fresh random IV.
func newAES256GCMWithRandomIV(key []byte) (cipher.AEAD, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, GCMIVSize)
	if err != nil {
		return nil, nil, err
	}
	iv := make([]byte, GCMIVSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}
	return aesGCM, iv, nil
}

// EncryptAES256 encrypts plain text using AES-256 GCM and returns a self-contained ciphertext in
// the nested {"c","t","iv"} shape used by carbon-crypto-service.
func EncryptAES256(key []byte, text string) (string, error) {
	aesGCM, iv, err := newAES256GCMWithRandomIV(key)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, iv, []byte(text), nil)

	inner := cipherInitializationVectorHolder{
		Cipher:               base64.StdEncoding.EncodeToString(ciphertext),
		InitializationVector: base64.StdEncoding.EncodeToString(iv),
		KeyId:                aes256KeyId(key),
	}
	innerJSON, err := json.Marshal(inner)
	if err != nil {
		return "", err
	}
	outer := cipherMetaDataHolder{
		C:  base64.StdEncoding.EncodeToString(innerJSON),
		T:  SecretEncryptionAlgorithmAESGCM,
		Iv: base64.StdEncoding.EncodeToString(iv),
	}
	outerJSON, err := json.Marshal(outer)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(outerJSON), nil
}

// EncryptAES256External encrypts plain text using AES-256 GCM and returns ciphertext in the flat
// {"cipherText","iv"} shape used by cipher-tool, carbon-secvault, and carbon-mediation's current
// ciphertext format.
func EncryptAES256External(key []byte, text string) (string, error) {
	aesGCM, iv, err := newAES256GCMWithRandomIV(key)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, iv, []byte(text), nil)

	holder := legacyCipherHolder{
		CipherText: base64.StdEncoding.EncodeToString(ciphertext),
		Iv:         base64.StdEncoding.EncodeToString(iv),
	}
	holderJSON, err := json.Marshal(holder)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(holderJSON), nil
}

// EncryptAES256Plain encrypts plain text using AES-256 GCM and returns a single opaque base64
// value (the IV prepended to the ciphertext) - no JSON envelope, just the encrypted value itself.
func EncryptAES256Plain(key []byte, text string) (string, error) {
	aesGCM, iv, err := newAES256GCMWithRandomIV(key)
	if err != nil {
		return "", err
	}

	sealed := aesGCM.Seal(iv, iv, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func aes256KeyId(key []byte) string {
	hash := sha256.Sum256(key)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// DecryptAES256 decrypts AES-256 GCM ciphertext produced by EncryptAES256, EncryptAES256External
// or EncryptAES256Plain, auto-detecting which of the three shapes cryptoText is in.
func DecryptAES256(key []byte, cryptoText string) (string, error) {
	outerBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(cryptoText))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, GCMIVSize)
	if err != nil {
		return "", err
	}

	var outer map[string]interface{}
	if err = json.Unmarshal(outerBytes, &outer); err != nil {
		// Not a JSON envelope - treat it as the plain iv||ciphertext blob.
		if len(outerBytes) < GCMIVSize {
			return "", errors.New("Ciphertext too short")
		}
		iv := outerBytes[:GCMIVSize]
		ciphertext := outerBytes[GCMIVSize:]
		plainText, err := aesGCM.Open(nil, iv, ciphertext, nil)
		if err != nil {
			return "", err
		}
		return string(plainText), nil
	}

	var cipherB64, ivB64 string
	if c, ok := outer["c"]; ok {
		// Current format: {"c": <base64 inner json>, "t": ..., "iv": ...}
		cStr, ok := c.(string)
		if !ok {
			return "", errors.New("Invalid ciphertext format: \"c\" field is not a string")
		}
		innerJSONBytes, err := base64.StdEncoding.DecodeString(cStr)
		if err != nil {
			return "", errors.New("Invalid base64 in \"c\" field")
		}
		var inner cipherInitializationVectorHolder
		if err = json.Unmarshal(innerJSONBytes, &inner); err != nil {
			return "", errors.New("Invalid inner ciphertext JSON")
		}
		cipherB64 = inner.Cipher
		ivB64 = inner.InitializationVector
	} else if ct, ok := outer["cipherText"]; ok {
		// Flat format: {"cipherText": ..., "iv": ...}
		cipherStr, ok := ct.(string)
		if !ok {
			return "", errors.New("Invalid ciphertext format: \"cipherText\" field is not a string")
		}
		cipherB64 = cipherStr
		if iv, ok := outer["iv"]; ok {
			ivStr, ok := iv.(string)
			if !ok {
				return "", errors.New("Invalid ciphertext format: \"iv\" field is not a string")
			}
			ivB64 = ivStr
		}
	} else {
		return "", errors.New("Unrecognized ciphertext format")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(ivB64)
	if err != nil {
		return "", err
	}

	plainText, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func isHexAES256Key(encryptionKey string) bool {
	if len(encryptionKey) != AES256HexKeySize {
		return false
	}
	for _, char := range encryptionKey {
		if !strings.ContainsRune(hexCharacters, char) {
			return false
		}
	}
	return true
}
