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

	s.credentials = Credentials{Environments: make(map[string]Credential)}
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

// Get credential for env
func (s *JsonStore) Get(env string) (Credential, error) {
	if cred, ok := s.credentials.Environments[env]; ok {
		username, err := Base64Decode(cred.Username)
		if err != nil {
			return Credential{}, err
		}
		password, err := Base64Decode(cred.Password)
		if err != nil {
			return Credential{}, err
		}
		clientId, err := Base64Decode(cred.ClientId)
		if err != nil {
			return Credential{}, err
		}
		clientSecret, err := Base64Decode(cred.ClientSecret)
		if err != nil {
			return Credential{}, err
		}
		credential := Credential{
			username, password, clientId, clientSecret,
		}
		return credential, nil
	}
	return Credential{}, fmt.Errorf("credentials not found for %s, use login", env)
}

// Set credentials for env using username, password, clientId, clientSecret
func (s *JsonStore) Set(env, username, password, clientId, clientSecret string) error {
	s.credentials.Environments[env] = Credential{
		Username:     Base64Encode(username),
		Password:     Base64Encode(password),
		ClientId:     Base64Encode(clientId),
		ClientSecret: Base64Encode(clientSecret),
	}
	err := s.persist()
	if err != nil {
		return err
	}
	fmt.Printf(PlainTextWarnMessage, s.Path)
	return nil
}

// Erase an env
func (s *JsonStore) Erase(env string) error {
	if _, ok := s.credentials.Environments[env]; !ok {
		return fmt.Errorf("%s was not found", env)
	}
	delete(s.credentials.Environments, env)
	return s.persist()
}

// IsKeychainEnabled returns if another store is activated
func (s *JsonStore) IsKeychainEnabled() bool {
	return s.credentials.CredStore != ""
}

// Has env in the store
func (s *JsonStore) Has(env string) bool {
	_, ok := s.credentials.Environments[env]
	return ok
}
