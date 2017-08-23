package utils

import "testing"

func TestMD5DigestLength(t *testing.T){
	passwords := []string{"admin", "1234", "!@#$"}

	for _, p := range passwords {
		md5Digest := GetMD5Hash(p)
		if len(md5Digest) != 32 {
			t.Errorf("MD5 Digest of %s does not have 32 characters.", p)
		}
	}
}

func TestEncryptDecrypt(t *testing.T){
	data := []string{"123412", "jfal;dsjf 3214134", "a&8S4#"}
	key := []byte(GetMD5Hash("password"))
	encryptedData := make([]string, len(data))
	for i, s := range data {
		encryptedData[i] = Encrypt(key, s)
		if s != Decrypt(key, encryptedData[i]) {
			t.Errorf("Encryption/Decryption does not work for '" + s + "'" )
		}
	}
}
