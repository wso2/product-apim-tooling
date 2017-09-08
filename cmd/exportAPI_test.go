package cmd

import (
	"testing"

	"net/http/httptest"
	"net/http"
	"github.com/renstrom/dedent"
	"github.com/menuka94/wso2apim-cli/utils"
	"fmt"
)

func TestExportAPI(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET', got '%s'\n", r.Method)
		}

		if r.Header.Get(utils.HeaderAccept) != utils.HeaderValueApplicationZip {
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

	resp := ExportAPI("test", "1.0", server.URL, "")
	fmt.Println(resp)
}
