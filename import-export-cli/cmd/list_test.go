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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/renstrom/dedent"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func TestGetAPIListOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodGet {
			t.Errorf("Expected method '%s', got '%s'\n", http.MethodGet, r.Method)
		}
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)

		if !strings.Contains(r.Header.Get(utils.HeaderAuthorization), utils.HeaderValueAuthBearerPrefix) {
			t.Errorf("Error in Authorization Header. Got '%s'\n", w.Header().Get(utils.HeaderAuthorization))
		}

		body := dedent.Dedent(`
			{
	"count": 3,
	"list": [{
			"id": "17e0f83c-dce5-4e9b-aa6a-db49b55591c5",
			"name": "test1",
			"context": "/test1",
			"version": "1.0.0",
			"provider": "admin",
			"status": "Created"
		},
		{
			"id": "9c740e42-309e-44aa-a8e1-6b8830aa7146",
			"name": "test2",
			"context": "/test2",
			"version": "1.0.0",
			"provider": "admin",
			"status": "Created"
		},
		{
			"id": "39899b8c-5893-4864-a935-9c149bc7461d",
			"name": "test3",
			"context": "/test3",
			"version": "1.0",
			"provider": "admin",
			"status": "Created"
		}
	]
}`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	count, apiList, err := GetAPIList("", "", "access_token", server.URL)
	fmt.Println("Count:", count)
	fmt.Println("List:", apiList)

	if count != 3 {
		t.Errorf("Incorrect count. Exptected %d, got %d\n", 3, count)
	}

	if err != nil {
		t.Error("Error" + err.Error())
	}
}

func TestGetAPIListUnreachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		if r.Method != http.MethodGet {
			t.Errorf("Expected method '%s', got '%s'\n", http.MethodGet, r.Method)
		}
		//if !strings.Contains(w.Header().Get(utils.HeaderAuthorization), "Bearer") {
		//	t.Error("Error in Authorization Header")
		//}
	}))
	defer server.Close()

	count, list, err := GetAPIList("", "", "access_token", server.URL)
	if count != 0 {
		t.Errorf("Incorrect Count. Expected %d, got %d\n", 0, count)
	}
	if list != nil {
		t.Errorf("")
	}

	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestGetApplicationListOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodGet {
			t.Errorf("Expected method '%s', got '%s'\n", http.MethodGet, r.Method)
		}
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)

		if !strings.Contains(r.Header.Get(utils.HeaderAuthorization), utils.HeaderValueAuthBearerPrefix) {
			t.Errorf("Error in Authorization Header. Got '%s'\n", w.Header().Get(utils.HeaderAuthorization))
		}

		body := dedent.Dedent(`
			{
    "count": 2,
    "list": [
        {
            "applicationId": "0e09806c-65bb-4114-b483-3f7521e51a70",
            "name": "testApp1",
            "owner": "admin",
            "status": "APPROVED",
            "groupId": ""
        },
        {
            "applicationId": "d2b2a966-97e6-40da-9f73-7202d6c2bf9b",
            "name": "testApp2",
            "owner": "admin",
            "status": "APPROVED",
            "groupId": "testGrp"
        }
		
    ]
}`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	count, appList, err := GetApplicationList("admin", "access_token", server.URL, "")
	fmt.Println("Count:", count)
	fmt.Println("List:", appList)

	if count != 2 {
		t.Errorf("Incorrect count. Exptected %d, got %d\n", 3, count)
	}

	if err != nil {
		t.Error("Error" + err.Error())
	}
}

func TestGetApplicationListUnreachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		if r.Method != http.MethodGet {
			t.Errorf("Expected method '%s', got '%s'\n", http.MethodGet, r.Method)
		}
	}))
	defer server.Close()

	count, list, err := GetAPIList("", "", "access_token", server.URL)
	if count != 0 {
		t.Errorf("Incorrect Count. Expected %d, got %d\n", 0, count)
	}
	if list != nil {
		t.Errorf("")
	}

	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestGetAPIProductListOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodGet {
			t.Errorf("Expected method '%s', got '%s'\n", http.MethodGet, r.Method)
		}
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)

		if !strings.Contains(r.Header.Get(utils.HeaderAuthorization), utils.HeaderValueAuthBearerPrefix) {
			t.Errorf("Error in Authorization Header. Got '%s'\n", w.Header().Get(utils.HeaderAuthorization))
		}

		body := dedent.Dedent(`
			{
	"count": 3,
	"list": [{
			"id": "17e0f83c-dce5-4e9b-aa6a-db49b55591c5",
			"name": "testproduct1",
			"context": "/testproduct1",
			"version": "1.0.0",
			"provider": "admin",
			"status": "Created"
		},
		{
			"id": "9c740e42-309e-44aa-a8e1-6b8830aa7146",
			"name": "testproduct2",
			"context": "/testproduct2",
			"version": "1.0.0",
			"provider": "admin",
			"status": "Created"
		},
		{
			"id": "39899b8c-5893-4864-a935-9c149bc7461d",
			"name": "testproduct3",
			"context": "/testproduct3",
			"version": "1.0",
			"provider": "admin@test.com",
			"status": "Created"
		}
	]
}`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	count, apiList, err := GetAPIProductList("", " ", "access_token", server.URL)
	fmt.Println("Count:", count)
	fmt.Println("List:", apiList)

	if count != 3 {
		t.Errorf("Incorrect count. Exptected %d, got %d\n", 3, count)
	}

	if err != nil {
		t.Error("Error" + err.Error())
	}
}

func TestGetAPIProductListUnreachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		if r.Method != http.MethodGet {
			t.Errorf("Expected method '%s', got '%s'\n", http.MethodGet, r.Method)
		}
		//if !strings.Contains(w.Header().Get(utils.HeaderAuthorization), "Bearer") {
		//	t.Error("Error in Authorization Header")
		//}
	}))
	defer server.Close()

	count, list, err := GetAPIProductList("", " ", "access_token", server.URL)
	if count != 0 {
		t.Errorf("Incorrect Count. Expected %d, got %d\n", 0, count)
	}
	if list != nil {
		t.Errorf("")
	}

	if err == nil {
		t.Error("Error should not be nil")
	}
}
