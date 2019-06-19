package credentials

import "encoding/base64"

// Base64Encode encodes given str to Base64 string
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode decodes a given str to original or returns an error
func Base64Decode(str string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(enc), nil
}
