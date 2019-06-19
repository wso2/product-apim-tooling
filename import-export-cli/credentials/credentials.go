package credentials

import (
	"errors"
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

type CredentialNotFound struct {
	Env string
}

func (c CredentialNotFound) Error() string {
	return "credential not found for " + c.Env
}

func IsCredentialNotFoundError(err error) bool {
	switch err.(type) {
	case *CredentialNotFound:
		return true
	default:
		return false
	}
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
