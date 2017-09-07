package utils

import (
	"fmt"
	"os"
	"github.com/wso2/wum-client/utils"
)


func HandleErrorAndExit(msg string, err error) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "wso2apim: %v\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "wso2apim: %v Reason: %v\n", msg, err.Error())
		utils.Logln(LogPrefixError + msg + ": " + err.Error())
	}
	os.Exit(1)
}
