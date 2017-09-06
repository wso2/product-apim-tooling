package utils

import (
	"fmt"
	"github.com/go-resty/resty"
	"bufio"
	"os"
	"golang.org/x/crypto/ssh/terminal"
)

// Invoke http-post request using go-resty
func InvokePOSTRequest(url string, headers map[string]string, body string) (*resty.Response, error) {
	resp, err := resty.R().SetHeaders(headers).SetBody(body).Post(url)

	return resp, err
}

func PromptForUsername() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	return username
}

func PromptForPassword() string {
	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(0)
	password := string(bytePassword)

	defer fmt.Println()
	return password
}


