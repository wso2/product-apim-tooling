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

func TestImportApplication1(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s' instead\n", http.MethodPost, r.Method)
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

	name := "sampleApp.zip"
	owner := "admin"
	accessToken := "access-token"

	_, err := ImportApplication(name, owner, server.URL, accessToken, "")
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
	utils.Insecure = true
	_, err = ImportApplication(name, owner, server.URL, accessToken, "")
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}

func TestNewAppFileUploadRequest(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected '%s', got '%s' instead\n", http.MethodPut, r.Method)
		}

		if !strings.Contains(r.Header.Get(utils.HeaderAccept), utils.HeaderValueMultiPartFormData) {
			t.Errorf("Expected '%s', got '%s' instead\n", utils.HeaderValueApplicationZip,
				r.Header.Get(utils.HeaderContentType))
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
	filePath := filepath.Join("sampleApp.zip")
	accessToken := "access-token"
	_, err := NewAppFileUploadRequest(server.URL, extraParams, "file", filePath, accessToken)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}
