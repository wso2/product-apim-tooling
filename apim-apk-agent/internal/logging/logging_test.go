/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package logging

import (
	"fmt"
	"testing"

	pkgLogging "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/logging"
)

func TestGetErrorMessageByCode(t *testing.T) {
	existingCode := Error1100 // Existing error code
	existingMessage := "Failed to listen on port."
	if message := GetErrorMessageByCode(existingCode); message != existingMessage {
		t.Errorf("Expected message: %s, but got: %s", existingMessage, message)
	}

	nonExistingCode := 9999 // Non-existing error code
	expectedErrorMessage := fmt.Sprintf("No error message found for error code: %v", nonExistingCode)
	if message := GetErrorMessageByCode(nonExistingCode); message != expectedErrorMessage {
		t.Errorf("Expected message: %s, but got: %s", expectedErrorMessage, message)
	}
}

func TestPrintError(t *testing.T) {
	code := Error1100
	severity := BLOCKER
	message := "Failed to listen on port."
	expectedError := pkgLogging.ErrorDetails{
		ErrorCode: code,
		Message:   message,
		Severity:  severity,
	}

	errorLog := PrintError(code, severity, message)
	if errorLog != expectedError {
		t.Errorf("PrintError returned unexpected result. Expected: %v, Got: %v", expectedError, errorLog)
	}
}

func TestPrintErrorWithDefaultMessage(t *testing.T) {
	code := Error1100
	severity := BLOCKER
	expectedError := pkgLogging.ErrorDetails{
		ErrorCode: code,
		Message:   "Failed to listen on port.",
		Severity:  severity,
	}

	errorLog := PrintErrorWithDefaultMessage(code, severity)
	if errorLog != expectedError {
		t.Errorf("PrintErrorWithDefaultMessage returned unexpected result. Expected: %v, Got: %v", expectedError, errorLog)
	}
}
