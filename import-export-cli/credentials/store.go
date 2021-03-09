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

type Store interface {
	// Has return the existance of apim credentials in the store for a given environment
	HasAPIM(env string) bool
	// HasMI return the existance of mi credentials in the store for a given environment
	HasMI(env string) bool
	// HasMG return the existance of mg tokens in the store for a given mgw adapter environment
	HasMG(env string) bool
	// GetAPIMCredentials returns credentials for apim from the store or an error
	GetAPIMCredentials(env string) (Credential, error)
	// GetMICredentials returns credentials for micro integrator from the store or an error
	GetMICredentials(env string) (MiCredential, error)
	// GetMgwAdapterToken returns the Access Token of the Microgateway Adapter
	GetMGToken(env string) (MgAdapterEnv, error)
	// SetAPIMCredentials sets credentials for micro integrator using username, password, clientID and client secret
	SetAPIMCredentials(env, username, password, clientID, clientSecret string) error
	// SetMICredentials sets credentials for micro integrator using username, password and access token
	SetMICredentials(env, username, password, accessToken string) error
	// SetMGToken sets the Access Token for a Microgateway Adapter env
	SetMGToken(env, accessToken string) error
	// Erase apim credentials in a given environment
	EraseAPIM(env string) error
	// Erase mi credentials in a given environment
	EraseMI(env string) error
	// Erase mg token in a given microgateway Adapter env
	EraseMG(env string) error
	// Load store
	Load() error
}
