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

package testutils

type APIParams struct {
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name                string    `yaml:"name"`
	Endpoints           Endpoints `yaml:"endpoints"`
	Security            Security  `yaml:"security,omitempty"`
	GatewayEnvironments []string  `yaml:"gatewayEnvironments,omitempty"`
	Certs               []Cert    `yaml:"certs,omitempty"`
}

type Endpoints struct {
	Production map[string]interface{} `yaml:"production,omitempty"`
	Sandbox    map[string]interface{} `yaml:"sandbox,omitempty"`
}

type Endpoint struct {
	URL    string  `yaml:"url"`
	Config *Config `yaml:"config,omitempty"`
}

type Config struct {
	RetryTimeOut int `yaml:"retryTimeOut"`
	RetryDelay   int `yaml:"retryDelay"`
	Factor       int `yaml:"factor"`
}

type Security struct {
	Production OAuthEndpointSecurity `yaml:"production"`
	Sandbox    OAuthEndpointSecurity `yaml:"sandbox"`
	Enabled    bool                  `yaml:"enabled"`
	Type       string                `yaml:"type"`
	Username   string                `yaml:"username"`
	Password   string                `yaml:"password"`
}

// OAuthEndpointSecurity contains details about the OAuth 2.0 endpoint security
type OAuthEndpointSecurity struct {
	Password          string            `yaml:"password"`
	Username          string            `yaml:"username"`
	TokenUrl          string            `yaml:"tokenUrl"`
	ClientId          string            `yaml:"clientId"`
	ClientSecret      string            `yaml:"clientSecret"`
	CustomParameters  map[string]string `yaml:"customParameters"`
	Type              string            `yaml:"type"`
	GrantType         string            `yaml:"grantType"`
	Enabled           bool              `yaml:"enabled"`
	IsSecretEncrypted bool              `yaml:"isSecretEncrypted"`
}

type Cert struct {
	HostName string `yaml:"hostName"`
	Alias    string `yaml:"alias"`
	Path     string `yaml:"path"`
}
