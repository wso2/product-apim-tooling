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
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestMD5DigestLength(t *testing.T) {
	passwords := []string{"admin", "1234", "!@#$"}

	for _, p := range passwords {
		md5Digest := GetMD5Hash(p)
		if len(md5Digest) != 32 {
			t.Errorf("MD5 Digest of %s does not have 32 characters.", p)
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {
	data := []string{"123412", "jfal;dsjf 3214134", "a&8S4#"}
	key := []byte(GetMD5Hash("password"))
	encryptedData := make([]string, len(data))
	for i, s := range data {
		encryptedData[i] = Encrypt(key, s)
		if s != Decrypt(key, encryptedData[i]) {
			t.Errorf("Encryption/Decryption does not work for '" + s + "'")
		}
	}
}

func TestResolveAES256Key(t *testing.T) {
	plainTextKey := "12345678901234567890123456789012"
	keyBytes, err := ResolveAES256Key(plainTextKey)
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error for a valid plain text key: %v", err)
	}
	if string(keyBytes) != plainTextKey {
		t.Fatal("ResolveAES256Key() did not preserve the plain text key bytes")
	}

	hexKey := strings.Repeat("ab", 32)
	keyBytes, err = ResolveAES256Key(hexKey)
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error for a valid hex key: %v", err)
	}
	if len(keyBytes) != AES256KeySize {
		t.Fatalf("ResolveAES256Key() returned an invalid key size: %d", len(keyBytes))
	}
}

func TestResolveAES256KeyInvalidLength(t *testing.T) {
	_, err := ResolveAES256Key("short-key")
	if err == nil {
		t.Fatal("ResolveAES256Key() did not fail for an invalid key length")
	}
}

func TestEncryptDecryptAES256(t *testing.T) {
	data := []string{"123412", "jfal;dsjf 3214134", "a&8S4#"}
	key, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	for _, s := range data {
		encryptedData, encryptErr := EncryptAES256(key, s)
		if encryptErr != nil {
			t.Fatalf("EncryptAES256() returned an error: %v", encryptErr)
		}
		if s == encryptedData {
			t.Fatal("EncryptAES256() returned plain text without encryption")
		}
		decryptedData, decryptErr := DecryptAES256(key, encryptedData)
		if decryptErr != nil {
			t.Fatalf("DecryptAES256() returned an error: %v", decryptErr)
		}
		if s != decryptedData {
			t.Errorf("EncryptAES256()/DecryptAES256() does not work for '%s'", s)
		}
	}
}

// TestEncryptAES256OutputFormat asserts EncryptAES256 produces a self-contained ciphertext with
// the outer {c,t,iv} / inner {cipher,initializationVector,keyId} shape, and a 12-byte (96-bit) IV.
func TestEncryptAES256OutputFormat(t *testing.T) {
	keyBytes, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	stored, err := EncryptAES256(keyBytes, "AnotherSecret#123")
	if err != nil {
		t.Fatalf("unexpected error encrypting: %v", err)
	}

	outerBytes, err := base64.StdEncoding.DecodeString(stored)
	if err != nil {
		t.Fatalf("stored value is not valid base64: %v", err)
	}

	var outer map[string]interface{}
	if err := json.Unmarshal(outerBytes, &outer); err != nil {
		t.Fatalf("outer envelope is not valid JSON: %v", err)
	}
	for _, field := range []string{"c", "t", "iv"} {
		if _, ok := outer[field]; !ok {
			t.Fatalf("outer envelope missing required field %q: %v", field, outer)
		}
	}
	if outer["t"] != SecretEncryptionAlgorithmAESGCM {
		t.Fatalf(`expected "t" = %q, got %v`, SecretEncryptionAlgorithmAESGCM, outer["t"])
	}

	innerBytes, err := base64.StdEncoding.DecodeString(outer["c"].(string))
	if err != nil {
		t.Fatalf(`"c" field is not valid base64: %v`, err)
	}
	var inner map[string]interface{}
	if err := json.Unmarshal(innerBytes, &inner); err != nil {
		t.Fatalf("inner envelope is not valid JSON: %v", err)
	}
	for _, field := range []string{"cipher", "initializationVector", "keyId"} {
		if _, ok := inner[field]; !ok {
			t.Fatalf("inner envelope missing required field %q: %v", field, inner)
		}
	}

	// The IV must be 12 bytes (96 bits), not the 128-byte value an earlier draft used.
	ivBytes, err := base64.StdEncoding.DecodeString(outer["iv"].(string))
	if err != nil {
		t.Fatalf(`"iv" field is not valid base64: %v`, err)
	}
	if len(ivBytes) != GCMIVSize {
		t.Fatalf("expected a %d-byte GCM IV, got %d bytes", GCMIVSize, len(ivBytes))
	}
}

func TestEncryptAES256KeyIdIsStableForSameKey(t *testing.T) {
	keyBytes, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	storedA, err := EncryptAES256(keyBytes, "valueA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	storedB, err := EncryptAES256(keyBytes, "valueB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	kidA := extractKeyId(t, storedA)
	kidB := extractKeyId(t, storedB)
	if kidA != kidB {
		t.Fatalf("expected the same key to produce the same keyId across calls, got %q and %q", kidA, kidB)
	}
}

func extractKeyId(t *testing.T, stored string) string {
	t.Helper()
	outerBytes, err := base64.StdEncoding.DecodeString(stored)
	if err != nil {
		t.Fatalf("stored value is not valid base64: %v", err)
	}
	var outer map[string]interface{}
	if err := json.Unmarshal(outerBytes, &outer); err != nil {
		t.Fatalf("outer envelope is not valid JSON: %v", err)
	}
	innerBytes, err := base64.StdEncoding.DecodeString(outer["c"].(string))
	if err != nil {
		t.Fatalf(`"c" field is not valid base64: %v`, err)
	}
	var inner cipherInitializationVectorHolder
	if err := json.Unmarshal(innerBytes, &inner); err != nil {
		t.Fatalf("inner envelope is not valid JSON: %v", err)
	}
	return inner.KeyId
}

// TestDecryptAES256LegacyFlatFormat confirms DecryptAES256 still reads the older
// {"cipherText", "iv"} flat shape, for backward compatibility with values encrypted before
// the self-contained {c,t,iv,keyId} format was introduced.
func TestDecryptAES256LegacyFlatFormat(t *testing.T) {
	keyBytes, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	plainText := "legacyValue"
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	iv := make([]byte, GCMIVSize)
	cipherBytes := gcm.Seal(nil, iv, []byte(plainText), nil)

	legacyJSON := `{"cipherText":"` + base64.StdEncoding.EncodeToString(cipherBytes) +
		`","iv":"` + base64.StdEncoding.EncodeToString(iv) + `"}`
	legacyStored := base64.StdEncoding.EncodeToString([]byte(legacyJSON))

	decrypted, err := DecryptAES256(keyBytes, legacyStored)
	if err != nil {
		t.Fatalf("unexpected error decrypting legacy format: %v", err)
	}
	if decrypted != plainText {
		t.Fatalf("expected %q, got %q", plainText, decrypted)
	}
}

// TestDecryptAES256LegacyNonceSize confirms DecryptAES256 still reads ciphertext written with the
// former 128-byte GCM nonce size (the flat {"cipherText","iv"} shape produced by the AES-256
// implementation that shipped before GCMIVSize was corrected to 12), not just the current
// 12-byte nonce.
func TestDecryptAES256LegacyNonceSize(t *testing.T) {
	keyBytes, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	const formerNonceSize = 128
	plainText := "legacyNonceSizeValue"
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, formerNonceSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	iv := make([]byte, formerNonceSize)
	cipherBytes := gcm.Seal(nil, iv, []byte(plainText), nil)

	legacyJSON := `{"cipherText":"` + base64.StdEncoding.EncodeToString(cipherBytes) +
		`","iv":"` + base64.StdEncoding.EncodeToString(iv) + `"}`
	legacyStored := base64.StdEncoding.EncodeToString([]byte(legacyJSON))

	decrypted, err := DecryptAES256(keyBytes, legacyStored)
	if err != nil {
		t.Fatalf("unexpected error decrypting former-nonce-size ciphertext: %v", err)
	}
	if decrypted != plainText {
		t.Fatalf("expected %q, got %q", plainText, decrypted)
	}
}

// TestEncryptAES256ExternalOutputFormat asserts EncryptAES256External produces the flat
// {"cipherText","iv"} shape (no "c"/"t" envelope), and that it round-trips through DecryptAES256.
func TestEncryptAES256ExternalOutputFormat(t *testing.T) {
	keyBytes, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	stored, err := EncryptAES256External(keyBytes, "ExternalSecret#456")
	if err != nil {
		t.Fatalf("unexpected error encrypting: %v", err)
	}

	outerBytes, err := base64.StdEncoding.DecodeString(stored)
	if err != nil {
		t.Fatalf("stored value is not valid base64: %v", err)
	}

	var outer map[string]interface{}
	if err := json.Unmarshal(outerBytes, &outer); err != nil {
		t.Fatalf("value is not valid JSON: %v", err)
	}
	for _, field := range []string{"cipherText", "iv"} {
		if _, ok := outer[field]; !ok {
			t.Fatalf("envelope missing required field %q: %v", field, outer)
		}
	}
	if _, ok := outer["c"]; ok {
		t.Fatalf("external format should not have a \"c\" field: %v", outer)
	}
	if _, ok := outer["t"]; ok {
		t.Fatalf("external format should not have a \"t\" field: %v", outer)
	}

	ivBytes, err := base64.StdEncoding.DecodeString(outer["iv"].(string))
	if err != nil {
		t.Fatalf(`"iv" field is not valid base64: %v`, err)
	}
	if len(ivBytes) != GCMIVSize {
		t.Fatalf("expected a %d-byte GCM IV, got %d bytes", GCMIVSize, len(ivBytes))
	}

	decrypted, err := DecryptAES256(keyBytes, stored)
	if err != nil {
		t.Fatalf("DecryptAES256() returned an error: %v", err)
	}
	if decrypted != "ExternalSecret#456" {
		t.Fatalf("expected %q, got %q", "ExternalSecret#456", decrypted)
	}
}

// TestEncryptAES256PlainOutputFormat asserts EncryptAES256Plain produces a single opaque
// base64 blob (no JSON envelope at all), and that it round-trips through DecryptAES256.
func TestEncryptAES256PlainOutputFormat(t *testing.T) {
	keyBytes, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	plainText := "PlainSecret#789"
	stored, err := EncryptAES256Plain(keyBytes, plainText)
	if err != nil {
		t.Fatalf("unexpected error encrypting: %v", err)
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(stored)
	if err != nil {
		t.Fatalf("stored value is not valid base64: %v", err)
	}
	if json.Valid(decodedBytes) {
		t.Fatalf("plain format should not be a JSON envelope, got %s", decodedBytes)
	}

	// GCMIVSize bytes of IV, plus ciphertext, plus a 16-byte GCM auth tag.
	const gcmTagSizeBytes = 16
	expectedLen := GCMIVSize + len(plainText) + gcmTagSizeBytes
	if len(decodedBytes) != expectedLen {
		t.Fatalf("expected %d decoded bytes (iv+ciphertext+tag), got %d", expectedLen, len(decodedBytes))
	}

	decrypted, err := DecryptAES256(keyBytes, stored)
	if err != nil {
		t.Fatalf("DecryptAES256() returned an error: %v", err)
	}
	if decrypted != plainText {
		t.Fatalf("expected %q, got %q", plainText, decrypted)
	}
}

// TestDecryptAES256RejectsNonStringFields confirms DecryptAES256 returns an error, rather than
// panicking, when a recognized envelope field holds a non-string JSON value.
func TestDecryptAES256RejectsNonStringFields(t *testing.T) {
	keyBytes, err := ResolveAES256Key("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("ResolveAES256Key() returned an error: %v", err)
	}

	cases := map[string]string{
		"c is a number":          `{"c":123,"t":"AES/GCM/NoPadding","iv":"aXY="}`,
		"cipherText is an array": `{"cipherText":["a"],"iv":"aXY="}`,
		"iv is an object":        `{"cipherText":"Y2lwaGVy","iv":{}}`,
	}
	for name, rawJSON := range cases {
		stored := base64.StdEncoding.EncodeToString([]byte(rawJSON))
		if _, err := DecryptAES256(keyBytes, stored); err == nil {
			t.Errorf("%s: expected an error, got nil", name)
		}
	}
}
