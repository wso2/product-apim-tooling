package cmd

import (
	"reflect"
	"testing"

	"github.com/go-resty/resty"
)

func TestExportAPI(t *testing.T) {
	type args struct {
		name        string
		version     string
		url         string
		accessToken string
	}
	tests := []struct {
		name string
		args args
		want *resty.Response
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExportAPI(tt.args.name, tt.args.version, tt.args.url, tt.args.accessToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
