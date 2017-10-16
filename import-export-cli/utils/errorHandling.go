/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package utils

import (
	"fmt"
	"os"
)

func HandleErrorAndExit(msg string, err error) {
	fmt.Println("\n=======  DEBUG LOG ==================")
	// TODO:: Remove debug log in production
	for i := 1; i <= 6; i++ {
		fmt.Println(WhereAmI(i))
	}
	fmt.Println("=======  END OF DEBUG LOG ===========\n")
	if err == nil {
		fmt.Fprintf(os.Stderr, "wso2apim: %v\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "wso2apim: %v Reason: %v\n", msg, err.Error())
		Logln(LogPrefixError + msg + ": " + err.Error())
	}
	defer printAndExit()
}

func printAndExit() {
	fmt.Println("Exiting...")
	os.Exit(1)
}
