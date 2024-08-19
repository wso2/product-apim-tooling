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
	"strings"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/adminservices"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	Environments            []Environment `yaml:"environments"`
	IndexingDelay           int           `yaml:"indexing-delay"`
	MaxInvocationAttempts   int           `yaml:"max-invocation-attempts"`
	DCRVersion              string        `yaml:"dcr-version"`
	AdminRESTAPIVersion     string        `yaml:"admin-rest-api-version"`
	DevportalRESTAPIVersion string        `yaml:"devportal-rest-api-version"`
	PublisherRESTAPIVersion string        `yaml:"publisher-rest-api-version"`
	DevOpsRESTAPIVersion    string        `yaml:"devops-rest-api-version"`
	APICTLVersion           string        `yaml:"apictl-version"`
}

type Environment struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Offset int    `yaml:"offset"`
}

const (
	superAdminUser     = adminservices.AdminUsername
	superAdminPassword = adminservices.AdminPassword

	userMgtService   = "RemoteUserStoreManagerService"
	tenantMgtService = "TenantMgtAdminService"

	DEFAULT_TENANT_DOMAIN = adminservices.DefaultTenantDomain
	TENANT1               = adminservices.Tenant1
)

var (
	Users = map[string][]adminservices.User{
		"creator":    {{UserName: adminservices.CreatorUsername, Password: adminservices.Password, Roles: []string{"Internal/creator"}}},
		"publisher":  {{UserName: adminservices.PublisherUsername, Password: adminservices.Password, Roles: []string{"Internal/publisher"}}},
		"subscriber": {{UserName: adminservices.SubscriberUsername, Password: adminservices.Password, Roles: []string{"Internal/subscriber"}}},
		"devops":     {{UserName: adminservices.DevopsUsername, Password: adminservices.Password, Roles: []string{"Internal/devops"}}},
	}

	yamlConfig YamlConfig

	envs = map[string]Environment{}

	tenants = []adminservices.Tenant{
		{AdminUserName: adminservices.AdminUsername, AdminPassword: adminservices.AdminPassword, Domain: adminservices.Tenant1},
	}

	creator    = Users["creator"][0]
	subscriber = Users["subscriber"][0]
	publisher  = Users["publisher"][0]
	devops     = Users["devops"][0]

	apimClients = map[string]*apim.Client{}

	// Table driven testing user combinations
	testCaseUsers = []testutils.TestCaseUsers{
		{
			Description:   "CTL user admin Super Tenant",
			ApiCreator:    testutils.Credentials{Username: creator.UserName, Password: creator.Password},
			ApiPublisher:  testutils.Credentials{Username: publisher.UserName, Password: publisher.Password},
			ApiSubscriber: testutils.Credentials{Username: subscriber.UserName, Password: subscriber.Password},
			Admin:         testutils.Credentials{Username: adminservices.AdminUsername, Password: adminservices.AdminPassword},
			CtlUser:       testutils.Credentials{Username: adminservices.AdminUsername, Password: adminservices.AdminPassword},
		},
		{
			Description:   "CTL user admin Tenant",
			ApiCreator:    testutils.Credentials{Username: creator.UserName + "@" + TENANT1, Password: creator.Password},
			ApiPublisher:  testutils.Credentials{Username: publisher.UserName + "@" + TENANT1, Password: publisher.Password},
			ApiSubscriber: testutils.Credentials{Username: subscriber.UserName + "@" + TENANT1, Password: subscriber.Password},
			Admin:         testutils.Credentials{Username: adminservices.AdminUsername + "@" + TENANT1, Password: adminservices.AdminPassword},
			CtlUser:       testutils.Credentials{Username: adminservices.AdminUsername + "@" + TENANT1, Password: adminservices.AdminPassword},
		},
		{
			Description:   "CTL user devops Super Tenant",
			ApiCreator:    testutils.Credentials{Username: creator.UserName, Password: creator.Password},
			ApiPublisher:  testutils.Credentials{Username: publisher.UserName, Password: publisher.Password},
			ApiSubscriber: testutils.Credentials{Username: subscriber.UserName, Password: subscriber.Password},
			Admin:         testutils.Credentials{Username: adminservices.AdminUsername, Password: adminservices.AdminPassword},
			CtlUser:       testutils.Credentials{Username: devops.UserName, Password: devops.Password},
		},
		{
			Description:   "CTL user devops Tenant",
			ApiCreator:    testutils.Credentials{Username: creator.UserName + "@" + TENANT1, Password: creator.Password},
			ApiPublisher:  testutils.Credentials{Username: publisher.UserName + "@" + TENANT1, Password: publisher.Password},
			ApiSubscriber: testutils.Credentials{Username: subscriber.UserName + "@" + TENANT1, Password: subscriber.Password},
			Admin:         testutils.Credentials{Username: adminservices.AdminUsername + "@" + TENANT1, Password: adminservices.AdminPassword},
			CtlUser:       testutils.Credentials{Username: devops.UserName + "@" + TENANT1, Password: devops.Password},
		},
	}
)

func GetDevClient() *apim.Client {
	return apimClients["development"]
}

func GetProdClient() *apim.Client {
	return apimClients["production"]
}

func TestMain(m *testing.M) {
	flag.Parse()

	readConfigs()

	base.ExtractArchiveFile("../build/target/")

	for _, env := range envs {
		client := apim.Client{}
		client.Setup(env.Name, env.Host, env.Offset, yamlConfig.DCRVersion, yamlConfig.AdminRESTAPIVersion,
			yamlConfig.DevportalRESTAPIVersion, yamlConfig.PublisherRESTAPIVersion, yamlConfig.DevOpsRESTAPIVersion)
		apimClients[env.Name] = &client
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

	indexingDelayOS := os.Getenv("INDEXING_DELAY")

	if !strings.EqualFold(indexingDelayOS, "") {
		intVar, _ := strconv.Atoi(indexingDelayOS)
		yamlConfig.IndexingDelay = intVar

	}

	versionOS := os.Getenv("VERSION")

	if !strings.EqualFold(versionOS, "") {
		yamlConfig.APICTLVersion = versionOS
	}

	base.SetIndexingDelay(yamlConfig.IndexingDelay)
	base.SetMaxInvocationAttempts(yamlConfig.MaxInvocationAttempts)

	base.Log("envs:", envs)
	base.Log("indexing delay:", yamlConfig.IndexingDelay)
	base.Log("max invocation attempts", yamlConfig.MaxInvocationAttempts)
	base.Log("dcr version:", yamlConfig.DCRVersion)
	base.Log("admin rest api Version:", yamlConfig.AdminRESTAPIVersion)
	base.Log("devportal rest api Version:", yamlConfig.DevportalRESTAPIVersion)
	base.Log("publisher rest api Version:", yamlConfig.PublisherRESTAPIVersion)
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
	removeEndpointCerts()
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
		deleteAllTenantUserApps(subscriber.UserName, subscriber.Password)
	}

	devopsUsers := Users["devops"]
	for _, devops := range devopsUsers {
		deleteAllTenantUserApps(devops.UserName, devops.Password)
	}

	for _, client := range apimClients {
		deleteUserApps(client, superAdminUser, superAdminPassword)
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			deleteUserApps(client, tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
		}
	}
}

func deleteAllTenantUserApps(username string, password string) {
	for _, client := range apimClients {
		deleteUserApps(client, username, password)
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			deleteUserApps(client, username+"@"+tenant.Domain, password)
		}
	}
}

func deleteUserApps(client *apim.Client, username string, password string) {
	client.Login(username, password)
	client.DeleteAllSubscriptions()
	client.DeleteAllApplications()
}

func deleteApis() {
	devopsUsers := Users["devops"]
	for _, devops := range devopsUsers {
		deleteAllTenantUserApis(devops.UserName, devops.Password)
	}

	for _, client := range apimClients {
		deleteUserApis(client, superAdminUser, superAdminPassword)
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			deleteUserApis(client, tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
		}
	}
}

func deleteAllTenantUserApis(username string, password string) {
	for _, client := range apimClients {
		deleteUserApis(client, username, password)
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			deleteUserApis(client, username+"@"+tenant.Domain, password)
		}
	}
}

func deleteUserApis(client *apim.Client, username string, password string) {
	client.Login(username, password)
	client.DeleteAllAPIs()
}

func deleteApiProducts() {
	publishers := Users["publisher"]
	for _, publisher := range publishers {
		deleteAllTenantUserApiProducts(publisher.UserName, publisher.Password)
	}

	for _, client := range apimClients {
		deleteUserApiProducts(client, superAdminUser, superAdminPassword)
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			deleteUserApiProducts(client, tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
		}
	}
}

func deleteAllTenantUserApiProducts(username string, password string) {
	for _, client := range apimClients {
		deleteUserApiProducts(client, username, password)
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			deleteUserApiProducts(client, username+"@"+tenant.Domain, password)
		}
	}
}

func deleteUserApiProducts(client *apim.Client, username string, password string) {
	client.Login(username, password)
	client.DeleteAllAPIProducts()
}

func removeEndpointCerts() {
	creators := Users["creator"]
	for _, creator := range creators {
		removeAllTenantUserEndpointCerts(creator.UserName, creator.Password)
	}

	for _, client := range apimClients {
		removeUserEndpointCerts(client, superAdminUser, superAdminPassword)

		for _, tenant := range tenants {
			removeUserEndpointCerts(client, tenant.AdminUserName+"@"+tenant.Domain, tenant.AdminPassword)
		}
	}
}

func removeAllTenantUserEndpointCerts(username string, password string) {
	for _, client := range apimClients {
		removeUserEndpointCerts(client, username, password)
	}

	for _, tenant := range tenants {
		for _, client := range apimClients {
			removeUserEndpointCerts(client, username+"@"+tenant.Domain, password)
		}
	}
}

func removeUserEndpointCerts(client *apim.Client, username string, password string) {
	client.Login(username, password)
	client.RemoveAllEndpointCerts()
}

func getSOAPServiceURL(host string, offset int, service string) string {
	port := 9443 + offset
	return "https://" + host + ":" + strconv.Itoa(port) + "/services/" + service
}

func isTenantUser(username, tenant string) bool {
	return strings.Contains(username, "@"+tenant)
}
