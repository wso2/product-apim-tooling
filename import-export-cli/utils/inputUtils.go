/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// ReadInputString reads input from user with prompting printText and validating against regex: validRegex
func ReadInputString(printText string, defaultVal string, validRegex string, retryOnInvalid bool) (string, error) {
	validate := func(value string) bool {
		reg := regexp.MustCompile(validRegex)
		return reg.MatchString(value)
	}

	return ReadInput(printText, defaultVal, validate, "", retryOnInvalid)
}

// ReadOption reads an option from user
func ReadOption(printText string, defaultVal int, maxValue int, retryOnInvalid bool) (int, error) {
	validate := func(value string) bool {
		option, _ := strconv.Atoi(value)
		return option > 0 && option <= maxValue
	}

	optionStr, err := ReadInput(printText, strconv.Itoa(defaultVal), validate, "Choose a number", retryOnInvalid)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(optionStr)
}

// ReadInput reads input from user with prompting printText
func ReadInput(printText string, defaultVal string, validate func(value string) bool, invalidText string, retryOnInvalid bool) (string, error) {
	retry := true
	value := ""
	reader := bufio.NewReader(os.Stdin)
	text := fmt.Sprintf("%s: %s: ", printText, defaultVal)

	for retry {
		fmt.Print(text)
		inputValue, err := reader.ReadString('\n')
		value = strings.TrimSpace(inputValue)

		if err != nil {
			return "", err
		}

		if value == "" {
			return defaultVal, nil
		}

		isValid := validate(value)

		if !isValid {
			fmt.Println(invalidText)
			if !retryOnInvalid {
				return value, errors.New("input validation failed")
			}
		}

		retry = retryOnInvalid && !isValid
	}

	return value, nil
}

// ReadPassword reads password from user with prompting printText
func ReadPassword(printText string) (string, error) {
	if printText == "" {
		printText = "Enter Password: "
	}
	fmt.Print(printText + ": ")

	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		return "", err
	}
	return string(password), nil
}
