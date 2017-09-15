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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvokePOSTRequestUnreachable(t *testing.T) {
	var httpStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST', got '%s'\n", r.Method)
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer httpStub.Close()

	resp, err := InvokePOSTRequest(httpStub.URL, make(map[string]string), "")
	if resp.StatusCode() != http.StatusInternalServerError {
		t.Errorf("Error in InvokePOSTRequest(): %s\n", err)
	}

}

func TestInvokePOSTRequestOK(t *testing.T) {
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

func TestPromptForUsername(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PromptForUsername(); got != tt.want {
				t.Errorf("PromptForUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPromptForPassword(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PromptForPassword(); got != tt.want {
				t.Errorf("PromptForPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
