package cmd

import (
	"reflect"
	"testing"

	"github.com/menuka94/wso2apim-cli/utils"
)

func TestGetAPIList(t *testing.T) {
	type args struct {
		query              string
		accessToken        string
		apiManagerEndpoint string
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		want1   []utils.API
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetAPIList(tt.args.query, tt.args.accessToken, tt.args.apiManagerEndpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAPIList() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetAPIList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
