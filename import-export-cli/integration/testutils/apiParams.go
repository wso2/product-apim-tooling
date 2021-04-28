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

type Params struct {
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name    string  `yaml:"name"`
	Configs Configs `yaml:"configs"`
}

type Configs struct {
	EndpointType           string                   `yaml:"endpointType,omitempty"`
	EndpointRoutingPolicy  string                   `yaml:"endpointRoutingPolicy,omitempty"`
	LoadBalanceEndpoints   LoadBalanceEndpointsData `yaml:"loadBalanceEndpoints,omitempty"`
	FailoverEndpoints      FailoverEndpointsData    `yaml:"failoverEndpoints,omitempty"`
	AWSLambdaEndpoints     AWSLambdaEndpointsData   `yaml:"awsLambdaEndpoints,omitempty"`
	Endpoints              Endpoints                `yaml:"endpoints"`
	Security               SecurityEndpoints        `yaml:"security,omitempty"`
	DeploymentEnvironments []DeploymentEnvironments `yaml:"deploymentEnvironments,omitempty"`
	Certs                  []Cert                   `yaml:"certs,omitempty"`
	MsslCerts              []MsslCert               `yaml:"mutualSslCerts,omitempty"`
	Policies               []string                 `yaml:"policies,omitempty"`
	DependentAPIs          map[string]interface{}   `yaml:"dependentAPIs,omitempty"`
}

type LoadBalanceEndpointsData struct {
	Production         []map[string]interface{} `yaml:"production,omitempty"`
	Sandbox            []map[string]interface{} `yaml:"sandbox,omitempty"`
	SessionManagement  string                   `yaml:"sessionManagement,omitempty"`
	SessionTimeout     int                      `yaml:"sessionTimeOut,omitempty"`
	AlgorithmClassName string                   `yaml:"algoClassName,omitempty"`
}

type FailoverEndpointsData struct {
	Production          Endpoint   `yaml:"production,omitempty"`
	ProductionFailovers []Endpoint `yaml:"productionFailovers,omitempty"`
	Sandbox             Endpoint   `yaml:"sandbox,omitempty"`
	SandboxFailovers    []Endpoint `yaml:"sandboxFailovers,omitempty"`
}

type AWSLambdaEndpointsData struct {
	AccessMethod  string `yaml:"accessMethod,omitempty"`
	AmznRegion    string `yaml:"amznRegion,omitempty"`
	AmznAccessKey string `yaml:"amznAccessKey,omitempty"`
	AmznSecretKey string `yaml:"amznSecretKey,omitempty"`
}
type DeploymentEnvironments struct {
	DisplayOnDevportal    bool   `yaml:"displayOnDevportal,omitempty"`
	DeploymentEnvironment string `yaml:"deploymentEnvironment,omitempty"`
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

type SecurityEndpoints struct {
	Production Security `yaml:"production,omitempty"`
	Sandbox    Security `yaml:"sandbox,omitempty"`
}

type Security struct {
	Enabled  bool   `yaml:"enabled"`
	Type     string `yaml:"type"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Cert struct {
	HostName string `yaml:"hostName"`
	Alias    string `yaml:"alias"`
	Path     string `yaml:"path"`
}

type MsslCert struct {
	TierName string `yaml:"tierName"`
	Alias    string `yaml:"alias"`
	Path     string `yaml:"path"`
}
