package utils

import (
	"testing"
	"fmt"
)

func TestGetClientIDSecret(t *testing.T) {
	one, two, three := GetClientIDSecret("admin", "admin", "")
	fmt.Println("one:", one)
	fmt.Println("two:", two)
	fmt.Println("three:", three)
}

func TestGetOAuthTokens(t *testing.T) {

}

func TestGetBase64EncodedCredentials(t *testing.T) {
	usernames := []string{"admin", "user", "admin"}
	passwords := []string{"admin", "password", "123456"}
	encodedPairs := []string{"YWRtaW46YWRtaW4=", "dXNlcjpwYXNzd29yZA==", "YWRtaW46MTIzNDU2"}

	for i, s := range encodedPairs {
		if s != GetBase64EncodedCredentials(usernames[i], passwords[i])	{
			t.Errorf("Error in Base64 Encoding. Base64(" + usernames[i] + ":" + passwords[i] + ") = " + encodedPairs[i])
		}
	}
}

