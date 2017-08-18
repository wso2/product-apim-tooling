package utils

import (
	"fmt"
	"net/http"
	"time"
	"log"
	"github.com/go-resty/resty"
	"bufio"
	"os"
	"golang.org/x/crypto/ssh/terminal"
)

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

func Authenticate() {

}

// Send the HTTP request to the server. Default timeout is one minute
func SendRequest(req *http.Request, timeout time.Duration) *http.Response {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		// TODO: Handle accordingly
	}
	return resp
}

// invoke the client request and handle error scenarios
func invokeRequest(req *http.Request, timeout time.Duration) *http.Response {
	httpResp := SendRequest(req, timeout)
	log.Printf("Status Code: %v\n", httpResp.Status)

	// TODO:: handle error responses

	return httpResp
}
