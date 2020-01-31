package utils

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"strings"
	"syscall"
)

const UsernameValidationRegex = `^[\w\d\-]+$`
const UrlValidationRegex = `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`

// ReadInputString reads input from user with prompting printText and validating against regex: validRegex
func ReadInputString(printText string, defaultVal string, validRegex string, retryOnInvalid bool) (string, error) {
	retry := true
	value := ""
	reader := bufio.NewReader(os.Stdin)
	text := printText + ": "
	if defaultVal != "" {
		text += defaultVal + ": "
	}

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

		reg := regexp.MustCompile(validRegex)
		isValid := reg.MatchString(value)
		if !retryOnInvalid && !isValid {
			return value, errors.New("input validation failed")
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
