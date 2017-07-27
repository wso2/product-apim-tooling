package utils

import (
	"fmt"
	"os"
	"net/http"
	"time"
	"log"
)

func PrintUsageErrorAndExit(msg, commandName string) {
	fmt.Fprintf(os.Stderr, "wso2apim: %v\n", msg)
	fmt.Fprintf(os.Stderr, "Try wso2apim %v --help for more information.\n", commandName)
	os.Exit(1)
}

func PrintErrorMessageAndExit(errorMsg string, err error){
	fmt.Fprintf(os.Stderr, "wso2apim: %v\n", errorMsg)
	println(err)
	os.Exit(1)
}

func HandleUnableToConnectErrorAndExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "wso2apim: %v\n", err.Error())
	}
}

func HandleErrorAndExit(msg string, err error) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "wso2apim: %v\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "wso2apim: %v Reason: %v\n", msg, err.Error())
	}
	os.Exit(1)
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

/*
func makeAPICall(req *http.Request) *http.Response {
	// Invoke request
	timeout := time.Duration(APICallTimeout * time.Minute)
	httpResp := invokeRequest(req, timeout)

	// 404, 200, or 401
	if httpResp.StatusCode == http.StatusUnauthorized {
		// Expired token. Renew Access Token. If the refresh token in invalid
		// then Authenticate() will notify and exit
		Authenticate()

		envConfig :=  GetEnvConfig()
	}
}
*/