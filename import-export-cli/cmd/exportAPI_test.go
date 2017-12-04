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
	"github.com/go-resty/resty"
	"os"
	"path/filepath"
	"testing"

	"fmt"
	"github.com/renstrom/dedent"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"net/http/httptest"
)

func TestExportAPI(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET', got '%s'\n", r.Method)
		}

		if r.Header.Get(utils.HeaderAccept) != utils.HeaderValueApplicationZip {
			t.Errorf("Expected '"+utils.HeaderValueApplicationZip+"', got '%s'\n",
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

	resp := ExportAPI("test", "1.0", "admin", server.URL, "")
	fmt.Println(resp)
}

func TestWriteToZip(t *testing.T) {
	name := "sampleapi"
	version := "1.0.0"
	environment := "dev"
	response := new(resty.Response)
	exportDirectory := utils.CurrentDir
	WriteToZip(name, version, environment, exportDirectory, response)
	defer os.RemoveAll(filepath.Join(exportDirectory, "dev"))
}
