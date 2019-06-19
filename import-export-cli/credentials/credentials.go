package credentials

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// DefaultConfigFile name
var DefaultConfigFile = "keys.json"

// Credential for storing user details
type Credential struct {
	// Username of user
	Username string `json:"username"`
	// Password of user
	Password string `json:"password"`
	// ClientId for cli
	ClientId string `json:"clientId"`
	// ClientSecret for cli
	ClientSecret string `json:"clientSecret"`
}

// Credentials of cli
type Credentials struct {
	// Environments specific credentials
	Environments map[string]Credential `json:"environments"`
	// CredStore represent type of store to be used
	CredStore string `json:"credStore,omitempty"`
}

// GetCredentialStore from file
// Note to set a different store please use credStore variable
func GetCredentialStore(f string) (Store, error) {
	// load as a json store first
	js := NewJsonStore(f)
	err := js.Load()
	if err != nil {
		return nil, err
	}
	return js, nil
}

// GetDefaultCredentialStore returns store from default path
func GetDefaultCredentialStore() (Store, error) {
	return GetCredentialStore(filepath.Join(utils.ConfigDirPath, DefaultConfigFile))
}

// GetOAuthAccessToken generates an accesstoken for CLI
func GetOAuthAccessToken(credential Credential, env string) (string, error) {
	tokenEndpoint := utils.GetTokenEndpointOfEnv(env, utils.MainConfigFilePath)
	data, err := utils.GetOAuthTokens(credential.Username, credential.Password,
		Base64Encode(credential.ClientId+":"+credential.ClientSecret),
		tokenEndpoint)
	if err != nil {
		return "", err
	}
	if accessToken, ok := data["access_token"]; ok {
		return accessToken, nil
	}
	return "", errors.New("access_token not found")
}

// GetBasicAuth returns basic auth username:password encoded in base64
func GetBasicAuth(credential Credential) string {
	return Base64Encode(fmt.Sprintf("%s:%s", credential.Username, credential.Password))
}
