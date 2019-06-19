package credentials

import (
	"path/filepath"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var DefaultConfigFile = "keys.json"

type Credential struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type Credentials struct {
	Environments map[string]Credential `json:"environments"`
	CredStore    string                `json:"credStore,omitempty"`
}

func GetCredentialStore(f string) (Store, error) {
	// load as a json store first
	js := NewJsonStore(f)
	err := js.Load()
	if err != nil {
		return nil, err
	}
	return js, nil
}

func GetDefaultCredentialStore() (Store, error) {
	return GetCredentialStore(filepath.Join(utils.ConfigDirPath, DefaultConfigFile))
}
