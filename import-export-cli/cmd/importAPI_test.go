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

package cmd

import (
	"github.com/renstrom/dedent"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

// TestImportAPISuccessful - 200 OK
func TestImportAPISuccessful(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected 'PUT', got '%s'\n", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)
		w.Header().Set(utils.HeaderContentEncoding, utils.HeaderValueGZIP)
		w.Header().Set(utils.HeaderTransferEncoding, utils.HeaderValueChunked)

		body := dedent.Dedent(`
		`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	name := "sampleapi.zip"
	accessToken := "access-token"

	utils.SkipTLSVerification = true

	_, err := ImportAPI(name, server.URL, accessToken, utils.CurrentDir)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}

// TestImportAPIError - 404 Not Found
func TestImportAPIError(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected 'PUT', got '%s'\n", r.Method)
		}

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)
		w.Header().Set(utils.HeaderContentEncoding, utils.HeaderValueGZIP)
		w.Header().Set(utils.HeaderTransferEncoding, utils.HeaderValueChunked)

		body := dedent.Dedent(`
		`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	name := "sampleapi.zip"
	accessToken := "access-token"

	utils.SkipTLSVerification = true

	_, err := ImportAPI(name, server.URL, accessToken, utils.CurrentDir)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}


func TestNewFileUploadRequest(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected 'PUT', got '%s'\n", r.Method)
		}

		if !strings.Contains(r.Header.Get(utils.HeaderAccept), utils.HeaderValueMultiPartFormData) {
			t.Errorf("Expected '"+utils.HeaderValueApplicationZip+"', got '%s'\n", r.Header.Get(utils.HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)
		w.Header().Set(utils.HeaderContentEncoding, utils.HeaderValueGZIP)
		w.Header().Set(utils.HeaderTransferEncoding, utils.HeaderValueChunked)

		body := dedent.Dedent(`
		`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	extraParams := map[string]string{}
	filePath := filepath.Join(utils.CurrentDir, "sampleapi.zip")
	accessToken := "access-token"
	_, err := NewFileUploadRequest(server.URL, extraParams, "file", filePath, accessToken)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}

func TestPrintAPIS(t *testing.T){
	var apis = []utils.API{
		{Context:"context", ID:"id", LifeCycleStatus:"created", Name:"test-api", Provider:"admin", Version:"1.0.0",
		WorkflowStatus:"work-flow-status"},
	}
	printAPIs(apis)
}
