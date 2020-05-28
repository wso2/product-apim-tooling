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

package integration

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/adminservices"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	Environments   []Environment `yaml:"environments"`
	DCRVersion     string        `yaml:"dcr-version"`
	RESTAPIVersion string        `yaml:"rest-api-version"`
	APICTLVersion  string        `yaml:"apictl-version"`
}

type Environment struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Offset int    `yaml:"offset"`
}

const (
	superAdminUser     = "admin"
	superAdminPassword = "admin"

	userMgtService   = "RemoteUserStoreManagerService"
	tenantMgtService = "TenantMgtAdminService"

	TENANT1 = "test.com"
)

var (
	Users = map[string][]adminservices.User{
		"creator":    {{UserName: "creator", Password: "password", Roles: []string{"Internal/creator"}}},
		"publisher":  {{UserName: "publisher", Password: "password", Roles: []string{"Internal/publisher"}}},
		"subscriber": {{UserName: "subscriber", Password: "password", Roles: []string{"Internal/subscriber"}}},
	}

	yamlConfig YamlConfig

	envs = map[string]Environment{}

	tenants = []adminservices.Tenant{
		{AdminUserName: "admin", AdminPassword: "admin", Domain: TENANT1},
	}

	creator    = Users["creator"][0]
	subscriber = Users["subscriber"][0]
	publisher  = Users["publisher"][0]

	apimClients []*apim.Client
)

func TestMain(m *testing.M) {
	flag.Parse()

	readConfigs()

	base.ExtractArchiveFile()

	for _, env := range envs {
		client := apim.Client{}
		client.Setup(env.Name, env.Host, env.Offset, yamlConfig.DCRVersion, yamlConfig.RESTAPIVersion)
		apimClients = append(apimClients, &client)
	}

	cleanupUsersAndTenants()
	addUsersAndTenants()
	cleanupAPIM()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func readConfigs() {
	reader, err := os.Open("config.yaml")

	if err != nil {
		base.Fatal(err)
	}
	defer reader.Close()

	yamlConfig = YamlConfig{}
	yaml.NewDecoder(reader).Decode(&yamlConfig)

	for _, env := range yamlConfig.Environments {
		envs[env.Name] = env
	}

	base.Log("envs:", envs)
	base.Log("dcr version:", yamlConfig.DCRVersion)
	base.Log("rest api Version:", yamlConfig.RESTAPIVersion)
	base.Log("apictl version:", yamlConfig.APICTLVersion)

	if len(envs) != 2 {
		base.Fatal("Expected number of Envs have not been configured for intergration tests")
	}
}

func cleanupUsersAndTenants() {
	for _, env := range envs {
		initTenants(env.Host, env.Offset)

		deleteUsers(env.Host, env.Offset)
		deactivateTenants(env.Host, env.Offset)
	}
}

func addUsersAndTenants() {
	for _, env := range envs {
		initTenants(env.Host, env.Offset)
		addUsers(env.Host, env.Offset)
	}
}

func cleanupAPIM() {
	deleteApps()
	deleteApiProducts()
	deleteApis()
}

func initTenants(host string, offset int) {
	url := getSOAPServiceURL(host, offset, tenantMgtService)

	for _, tenant := range tenants {
		adminservices.InitTenant(&tenant, url, superAdminUser, superAdminPassword)
	}

}

func addUsers(host string, offset int) {
	url := getSOAPServiceURL(host, offset, userMgtService)

	for _, userCategory := range Users {
		for _, user := range userCategory {
			// Add super tenant user
			adminservices.AddUser(&user, url, superAdminUser, superAdminPassword)

			for _, tenant := range tenants {
				// Add tenant user
				adminservices.AddUser(&user, url, tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
			}
		}
	}
}

func deleteUsers(host string, offset int) {
	url := getSOAPServiceURL(host, offset, userMgtService)

	for _, userCategory := range Users {
		for _, user := range userCategory {
			// Delete super tenant user
			if adminservices.IsUserExists(user.UserName, url, superAdminUser, superAdminPassword) {
				adminservices.DeleteUser(user.UserName, url, superAdminUser, superAdminPassword)
			}

			for _, tenant := range tenants {
				// Delete tenant user
				if adminservices.IsUserExists(user.UserName, url, tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword) {
					adminservices.DeleteUser(user.UserName, url, tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
				}
			}
		}
	}
}

func deactivateTenants(host string, offset int) {
	url := getSOAPServiceURL(host, offset, tenantMgtService)

	for _, tenant := range tenants {
		adminservices.DeactivateTenant(tenant.Domain, url, superAdminUser, superAdminPassword)
	}
}

func deleteApps() {
	subscribers := Users["subscriber"]
	for _, subscriber := range subscribers {
		for _, client := range apimClients {
			client.Login(subscriber.UserName, subscriber.Password)
			client.DeleteAllSubscriptions()
			client.DeleteAllApplications()
		}

		for _, tenant := range tenants {
			for _, client := range apimClients {
				client.Login(subscriber.UserName+"@"+tenant.Domain, subscriber.Password)
				client.DeleteAllSubscriptions()
				client.DeleteAllApplications()
			}
		}
	}

	for _, client := range apimClients {
		client.Login(superAdminUser, superAdminPassword)
		client.DeleteAllSubscriptions()
		client.DeleteAllApplications()
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			client.Login(tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
			client.DeleteAllSubscriptions()
			client.DeleteAllApplications()
		}
	}
}

func deleteApis() {
	creators := Users["creator"]
	for _, creator := range creators {
		for _, client := range apimClients {
			client.Login(creator.UserName, creator.Password)
			client.DeleteAllAPIs()
		}

		for _, tenant := range tenants {
			for _, client := range apimClients {
				client.Login(creator.UserName+"@"+tenant.Domain, creator.Password)
				client.DeleteAllAPIs()
			}
		}
	}

	for _, client := range apimClients {
		client.Login(superAdminUser, superAdminPassword)
		client.DeleteAllAPIs()
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			client.Login(tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
			client.DeleteAllAPIs()
		}
	}
}

func deleteApiProducts() {
	publishers := Users["publisher"]
	for _, publisher := range publishers {
		for _, client := range apimClients {
			client.Login(publisher.UserName, publisher.Password)
			client.DeleteAllAPIProducts()
		}

		for _, tenant := range tenants {
			for _, client := range apimClients {
				client.Login(publisher.UserName+"@"+tenant.Domain, publisher.Password)
				client.DeleteAllAPIProducts()
			}
		}
	}

	for _, client := range apimClients {
		client.Login(superAdminUser, superAdminPassword)
		client.DeleteAllAPIProducts()
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			client.Login(tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
			client.DeleteAllAPIProducts()
		}
	}
}

func getSOAPServiceURL(host string, offset int, service string) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/services/" + service
}
