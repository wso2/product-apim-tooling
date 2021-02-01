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

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// AdminUserName default admin username
const AdminUserName = "admin"

// AdminPassword default admin password
const AdminPassword = "admin"

const miContainer = "micontainer:/home/wso2carbon/wso2mi/repository/deployment/server/carbonapps"
const deploymentDelay = 10000

var cappList = []string{"./testdata/capps/HealthCareCompositeExporter_1.0.0.car",
	"./testdata/capps/RESTDataServiceCompositeExporter_1.0.0.car"}

// MiRESTClient Enables interacting with the Management API of MI
type MiRESTClient struct {
	portOffset  int
	host        string
	EnvName     string
	miURL       string
	accessToken string
}

// Environment store environment details of the MI
type Environment struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Offset int    `yaml:"offset"`
}

// MiConfig store credentials and REST Client of the MI
type MiConfig struct {
	Username string
	Password string
	MIClient MiRESTClient
}

// GetEnvName : Get environment name
func (instance *MiRESTClient) GetEnvName() string {
	return instance.EnvName
}

// GetMiURL : Get MI URL
func (instance *MiRESTClient) GetMiURL() string {
	return instance.miURL
}

// SetupMI : Setup MI Client config
func (instance *MiRESTClient) SetupMI(username, password, envName, host string, offset int) {
	base.Log("apim.SetupMI() - envName:", envName, ",host:", host, ",offset:", offset)
	instance.miURL = getMiURL(host, offset)
	instance.host = host
	instance.portOffset = offset
	instance.EnvName = envName
	instance.accessToken = instance.getToken(username, password)
}

func getMiURL(host string, offset int) string {
	port := 9164 + offset
	return "https://" + host + ":" + strconv.Itoa(port)
}

func (instance *MiRESTClient) getToken(username string, password string) string {
	tokenURL := getResourceURL(instance.GetMiURL(), utils.MiManagementMiLoginResource)
	request := base.CreateGet(tokenURL)
	request.SetBasicAuth(username, password)
	base.LogRequest("mi.getToken()", request)

	response := base.SendHTTPRequest(request)
	defer response.Body.Close()
	base.ValidateAndLogResponse("mi.getToken()", response, 200)

	var jsonResp map[string]string
	json.NewDecoder(response.Body).Decode(&jsonResp)
	return jsonResp["AccessToken"]
}

func getResourceURL(miURL, resource string) string {
	return miURL + "/" + utils.MiManagementAPIContext + "/" + resource
}

func getResourceURLWithQueryParam(miURL, resource string, params map[string]string) string {
	var path strings.Builder
	path.WriteString(getResourceURL(miURL, resource))
	if params != nil {
		path.WriteString("?")
		for key, value := range params {
			path.WriteString(key)
			path.WriteString("=")
			path.WriteString(value)
		}
	}
	return path.String()
}

// SetupAndLoginToMI setup Mi instance and login to it
func SetupAndLoginToMI(t *testing.T, config *MiConfig) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Password)
}

func getParamMap(key, value string) map[string]string {
	paramMap := make(map[string]string)
	paramMap[key] = value
	return paramMap
}

// DeployCApps deploy capps to the micro integrator
func DeployCApps() {
	for _, capp := range cappList {
		args := []string{"cp", capp, miContainer}
		execDockerCmd(args...)
	}
}

func execDockerCmd(args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	output, err := cmd.Output()
	return string(output), err
}

// WaitForDeployment wait for specified interval to allow MI to deploy the artifacts
func WaitForDeployment() {
	time.Sleep(time.Duration(deploymentDelay) * time.Millisecond)
}
