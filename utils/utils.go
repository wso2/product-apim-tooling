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
