package utils

import (
	"testing"
	"net/http/httptest"
	"net/http"
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

