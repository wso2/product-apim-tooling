/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// PlainTextWarnMessage warning message
const PlainTextWarnMessage = "WARNING: credentials are stored as a plain text in %s\n"

// JsonStore is storing keys in json format
type JsonStore struct {
	// Path to file
	Path string

	// internal usage
	credentials Credentials
}

// NewJsonStore creates a new store
func NewJsonStore(path string) *JsonStore {
	return &JsonStore{Path: path}
}

// Load json store
func (s *JsonStore) Load() error {
	if info, err := os.Stat(s.Path); err == nil && !info.IsDir() {
		data, err := ioutil.ReadFile(s.Path)
		if err != nil {
			return err
		}

		var cred Credentials
		err = json.Unmarshal(data, &cred)
		if err != nil {
			return err
		}

		s.credentials = cred
		return nil
	} else if err == nil && info.IsDir() {
		return fmt.Errorf("%s is a directory", s.Path)
	}

	s.credentials = Credentials{
		Environments:   make(map[string]Environment),
		MgwAdapterEnvs: make(map[string]MgToken),
	}
	return nil
}

// saves to disk
func (s *JsonStore) persist() error {
	data, err := json.MarshalIndent(s.credentials, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.Path, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// GetAPIMCredentials returns credentials for apim from the store or an error
func (s *JsonStore) GetAPIMCredentials(env string) (Credential, error) {
	if environment, ok := s.credentials.Environments[env]; ok {
		username, err := Base64Decode(environment.APIM.Username)
		if err != nil {
			return Credential{}, err
		}
		password, err := Base64Decode(environment.APIM.Password)
		if err != nil {
			return Credential{}, err
		}
		clientID, err := Base64Decode(environment.APIM.ClientId)
		if err != nil {
			return Credential{}, err
		}
		clientSecret, err := Base64Decode(environment.APIM.ClientSecret)
		if err != nil {
			return Credential{}, err
		}
		credential := Credential{
			username, password, clientID, clientSecret,
		}
		return credential, nil
	}
	return Credential{}, fmt.Errorf("credentials not found for APIM in %s, use login", env)
}

// SetAPIMCredentials sets credentials for micro integrator using username, password, clientID and client secret
func (s *JsonStore) SetAPIMCredentials(env, username, password, clientId, clientSecret string) error {
	environment := s.credentials.Environments[env]
	environment.APIM = Credential{
		Username:     Base64Encode(username),
		Password:     Base64Encode(password),
		ClientId:     Base64Encode(clientId),
		ClientSecret: Base64Encode(clientSecret),
	}
	s.credentials.Environments[env] = environment
	err := s.persist()
	if err != nil {
		return err
	}
	fmt.Printf(PlainTextWarnMessage, s.Path)
	return nil
}

// GetMICredentials returns credentials for micro integrator from the store or an error
func (s *JsonStore) GetMICredentials(env string) (MiCredential, error) {
	if environment, ok := s.credentials.Environments[env]; ok {
		username, err := Base64Decode(environment.MI.Username)
		if err != nil {
			return MiCredential{}, err
		}
		password, err := Base64Decode(environment.MI.Password)
		if err != nil {
			return MiCredential{}, err
		}
		accessToken, err := Base64Decode(environment.MI.AccessToken)
		if err != nil {
			return MiCredential{}, err
		}
		credential := MiCredential{
			username, password, accessToken,
		}
		return credential, nil
	}
	return MiCredential{}, fmt.Errorf("credentials not found for Mi in %s, use login", env)
}

// SetMICredentials set credentials for mi using username, password, accessToken
func (s *JsonStore) SetMICredentials(env, username, password, accessToken string) error {
	environment := s.credentials.Environments[env]
	environment.MI = MiCredential{
		Username:    Base64Encode(username),
		Password:    Base64Encode(password),
		AccessToken: Base64Encode(accessToken),
	}
	s.credentials.Environments[env] = environment
	err := s.persist()
	if err != nil {
		return err
	}
	fmt.Printf(PlainTextWarnMessage, s.Path)
	return nil
}

// GetMGToken returns token for microgateway adapter from the store or an error
func (s *JsonStore) GetMGToken(env string) (MgToken, error) {
	if mgToken, ok := s.credentials.MgwAdapterEnvs[env]; ok {
		return mgToken, nil
	}
	return MgToken{}, fmt.Errorf("Tokens not found for Mgw in %s, use login", env)
}

// SetMGToken set token for microgateway adapter
func (s *JsonStore) SetMGToken(env, accessToken string) error {
	mgwAdapterEnv := s.credentials.MgwAdapterEnvs[env]
	mgwAdapterEnv.AccessToken = accessToken
	err := s.persist()
	if err != nil {
		return err
	}
	return nil
}

// EraseAPIM remove apim credentials from the store
func (s *JsonStore) EraseAPIM(env string) error {
	environment, ok := s.credentials.Environments[env]
	if !ok {
		return fmt.Errorf("%s was not found", env)
	}
	if !miCredentialsExists(environment.MI) {
		// delete the environment
		delete(s.credentials.Environments, env)
	} else {
		// remove only apim credentials
		environment.APIM = Credential{}
		s.credentials.Environments[env] = environment
	}
	return s.persist()
}

// EraseMI remove mi credentials from the store
func (s *JsonStore) EraseMI(env string) error {
	environment, ok := s.credentials.Environments[env]
	if !ok {
		return fmt.Errorf("%s was not found", env)
	}
	if !apimCredentialsExists(environment.APIM) {
		// delete the environment
		delete(s.credentials.Environments, env)
	} else {
		// remove only mi credentials
		environment.MI = MiCredential{}
		s.credentials.Environments[env] = environment
	}
	return s.persist()
}

// EraseMG remove mg tokens from the store
func (s *JsonStore) EraseMG(env string) error {
	_, ok := s.credentials.MgwAdapterEnvs[env]
	if !ok {
		return fmt.Errorf("%s was not found", env)
	} else {
		// remove only mg tokens
		delete(s.credentials.MgwAdapterEnvs, env)
	}
	return s.persist()
}

// IsKeychainEnabled returns if another store is activated
func (s *JsonStore) IsKeychainEnabled() bool {
	return s.credentials.CredStore != ""
}

// HasAPIM return the existance of apim credentials in the store for a given environment
func (s *JsonStore) HasAPIM(env string) bool {
	if environment, ok := s.credentials.Environments[env]; ok {
		return apimCredentialsExists(environment.APIM)
	}
	return false
}

// HasMI return the existance of mi credentials in the store for a given environment
func (s *JsonStore) HasMI(env string) bool {
	if environment, ok := s.credentials.Environments[env]; ok {
		return miCredentialsExists(environment.MI)
	}
	return false
}

func miCredentialsExists(miCred MiCredential) bool {
	return miCred.AccessToken != "" && miCred.Username != "" && miCred.Password != ""
}

func apimCredentialsExists(apimCred Credential) bool {
	return apimCred.ClientId != "" && apimCred.ClientSecret != "" && apimCred.Username != "" && apimCred.Password != ""
}
