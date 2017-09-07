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
