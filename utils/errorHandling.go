package utils

import (
	"fmt"
	"os"
)


func HandleErrorAndExit(msg string, err error) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "wso2apim: %v\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "wso2apim: %v Reason: %v\n", msg, err.Error())
		Logln(LogPrefixError + msg + ": " + err.Error())
	}
	defer printAndExit()
}

func printAndExit(){
	fmt.Println("Exiting...")
	os.Exit(1)
}
