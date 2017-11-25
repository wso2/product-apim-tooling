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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvokePOSTRequestUnreachable(t *testing.T) {
	var httpStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s'\n", http.MethodPost, r.Method)
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer httpStub.Close()

	resp, err := InvokePOSTRequest(httpStub.URL, make(map[string]string), "")
	if resp.StatusCode() != http.StatusInternalServerError {
		t.Errorf("Error in InvokePOSTRequest(): %s\n", err)
	}

}

func TestInvokeGETRequestOK(t *testing.T) {
	SkipTLSVerification = true
	var httpStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected '%s', got '%s'\n", http.MethodGet, r.Method)
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer httpStub.Close()

	resp, err := InvokeGETRequest(httpStub.URL, make(map[string]string))
	if resp.StatusCode() != http.StatusInternalServerError {
		t.Errorf("Error in InvokePOSTRequest(): %s\n", err)
	}
}

func TestInvokePOSTRequestOK(t *testing.T) {
	SkipTLSVerification = true
	var httpStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST', got '%s'\n", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer httpStub.Close()

	resp, err := InvokePOSTRequest(httpStub.URL, make(map[string]string), "")
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Error in InvokePOSTRequest(): %s\n", err)
	}
}

func TestPromptForPassword(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "admin", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PromptForPassword(); got != tt.want {
				t.Errorf("PromptForPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChopPath(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{source: "/user/home", expected: "home"},
		{source: "home", expected: "home"},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			if got := chopPath(tt.source); got != tt.expected {
				t.Errorf("PromptForPassword() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestWhereAmI(t *testing.T) {
	WhereAmI(5)
	WhereAmI()
}

func TestShowHelpCommandTip(t *testing.T) {
	ShowHelpCommandTip("export-api")
}
