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

package adminservices

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

// User : data strcuture defining user
type User struct {
	UserName string
	Password string
	Roles    []string
}

// Role : data strcuture defining role
type Role struct {
	RoleName     string
	IsSharedRole bool
}

type isExistingUserRequest struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	SOAPEnvNS string   `xml:"xmlns:soap,attr"`
	SERNS     string   `xml:"xmlns:ser,attr"`
	UserName  string   `xml:"soap:Body>ser:isExistingUser>ser:userName"`
}

type addUserRequest struct {
	XMLName               xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS             string   `xml:"xmlns:soapenv,attr"`
	SERNS                 string   `xml:"xmlns:ser,attr"`
	XSDNS                 string   `xml:"xmlns:xsd,attr"`
	UserName              string   `xml:"soapenv:Body>ser:addUser>ser:userName"`
	Credential            string   `xml:"soapenv:Body>ser:addUser>ser:credential"`
	RoleList              []string `xml:"soapenv:Body>ser:addUser>ser:roleList"`
	RequirePasswordChange bool     `xml:"soapenv:Body>ser:addUser>ser:requirePasswordChange"`
}

type deleteUserRequest struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	SOAPEnvNS string   `xml:"xmlns:soap,attr"`
	SERNS     string   `xml:"xmlns:ser,attr"`
	UserName  string   `xml:"soap:Body>ser:deleteUser>ser:userName"`
}

type addRoleRequest struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS    string   `xml:"xmlns:soapenv,attr"`
	XSDNS        string   `xml:"xmlns:xsd,attr"`
	XSD1NS       string   `xml:"xmlns:xsd1,attr"`
	RoleName     string   `xml:"soapenv:Body>xsd:addRole>xsd:roleName"`
	IsSharedRole bool     `xml:"soapenv:Body>xsd:addRole>xsd:isSharedRole"`
}

type deleteRoleRequest struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	SOAPEnvNS string   `xml:"xmlns:soap,attr"`
	SERNS     string   `xml:"xmlns:ser,attr"`
	RoleName  string   `xml:"soap:Body>ser:deleteRole>ser:roleName"`
}

type isExistingUserResponse struct {
	RawXML   string `xml:",innerxml"`
	isExists bool
}

func (instance *isExistingUserResponse) parse() {
	splits := strings.Split(instance.RawXML, "\n")

	for _, s := range splits {
		if strings.Contains(s, "return>true") {
			instance.isExists = true
			break
		}
	}
}

// AddRole : Add new role in APIM
func AddRole(role *Role, url string, username string, password string) {
	msg := createAddRoleRequestPayload(role)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:addRole", request)

	base.LogRequest("adminservices.AddRole()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.AddRole()", response, 200)
}

// DeleteRole : Delete role in APIM
func DeleteRole(roleName string, url string, username string, password string) {
	msg := createDeleteRoleRequestPayload(roleName)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:deleteRole", request)

	base.LogRequest("adminservices.DeleteRole()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.DeleteRole()", response, 200)
}

// AddUser : Add new user in APIM
func AddUser(user *User, url string, username string, password string) {
	msg := createAddUserRequestPayload(user)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:addUser", request)

	base.LogRequest("adminservices.AddUser()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.AddUser()", response, 202)
}

// IsUserExists : Is user exists in APIM
func IsUserExists(name string, url string, username string, password string) bool {
	msg := createIsExistingUserRequestPayload(name)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:isExistingUser", request)

	base.LogRequest("adminservices.IsUserExists()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.IsUserExists()", response, 200)

	isExistingResp := isExistingUserResponse{}
	xml.NewDecoder(response.Body).Decode(&isExistingResp)
	isExistingResp.parse()

	return isExistingResp.isExists
}

// DeleteUser : Delete user in APIM
func DeleteUser(name string, url string, username string, password string) {
	msg := createDeleteUserRequestPayload(name)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:deleteUser", request)

	base.LogRequest("adminservices.DeleteUser()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.DeleteUser()", response, 202)
}

func createAddUserRequestPayload(user *User) *bytes.Buffer {
	request := addUserRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = umSerNS
	request.XSDNS = xsdNS
	request.UserName = user.UserName
	request.Credential = user.Password
	request.RoleList = user.Roles
	request.RequirePasswordChange = false

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

func createIsExistingUserRequestPayload(userName string) *bytes.Buffer {
	request := isExistingUserRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = umSerNS
	request.UserName = userName

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

func createDeleteUserRequestPayload(userName string) *bytes.Buffer {
	request := deleteUserRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = umSerNS
	request.UserName = userName

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

func createAddRoleRequestPayload(role *Role) *bytes.Buffer {
	request := addRoleRequest{}

	request.SOAPEnvNS = xmlNS
	request.XSDNS = xsdNS
	request.RoleName = role.RoleName
	request.IsSharedRole = role.IsSharedRole

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

func createDeleteRoleRequestPayload(roleName string) *bytes.Buffer {
	request := deleteRoleRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = umSerNS
	request.RoleName = roleName

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}
