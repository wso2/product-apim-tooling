package cmd

import (
	"testing"
	"net/http/httptest"
	"github.com/renstrom/dedent"
	"net/http"
	"github.com/menuka94/wso2apim-cli/utils"
	"strings"
)

func TestImportAPI(t *testing.T) {
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

	_, err := ImportAPI(name, server.URL, accessToken)
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
	filePath := utils.ApplicationRoot + "/" + "sampleapi.zip"
	accessToken := "access-token"
	_, err := NewFileUploadRequest(server.URL, extraParams, "file", filePath, accessToken)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}


