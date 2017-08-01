package utils

import (
	"fmt"
	"os"
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
