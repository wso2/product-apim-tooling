/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import "testing"

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
