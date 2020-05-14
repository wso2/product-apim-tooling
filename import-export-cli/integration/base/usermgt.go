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

package base

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

const (
	xmlNS  = "http://schemas.xmlsoap.org/soap/envelope/"
	xsdNS  = "http://org.apache.axis2/xsd"
	xsdNS1 = "http://common.mgt.user.carbon.wso2.org/xsd"
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

type addUserRequest struct {
	XMLName   xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS string   `xml:"xmlns:soapenv,attr"`
	XSDNS     string   `xml:"xmlns:xsd,attr"`
	XSD1NS    string   `xml:"xmlns:xsd1,attr"`
	UserName  string   `xml:"soapenv:Body>xsd:addUser>xsd:userName"`
	Password  string   `xml:"soapenv:Body>xsd:addUser>xsd:password"`
	Roles     []string `xml:"soapenv:Body>xsd:addUser>xsd:roles"`
}

type addRoleRequest struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS    string   `xml:"xmlns:soapenv,attr"`
	XSDNS        string   `xml:"xmlns:xsd,attr"`
	XSD1NS       string   `xml:"xmlns:xsd1,attr"`
	RoleName     string   `xml:"soapenv:Body>xsd:addRole>xsd:roleName"`
	IsSharedRole bool     `xml:"soapenv:Body>xsd:addRole>xsd:isSharedRole"`
}

// AddRole : Add new role in APIM
func AddRole(role *Role, url string, username string, password string, action string) {
	msg := createAddRoleRequest(role)

	request := CreatePost(url, msg)

	SetSOAPHeaders(username, password, action, request)

	SendHTTPRequest(request)
}

// AddUser : Add new user in APIM
func AddUser(user *User, url string, username string, password string, action string) {
	msg := createAddUserRequest(user)

	request := CreatePost(url, msg)

	SetSOAPHeaders(username, password, action, request)

	SendHTTPRequest(request)
}

func createAddUserRequest(user *User) *bytes.Buffer {
	request := addUserRequest{}

	request.SOAPEnvNS = xmlNS
	request.XSDNS = xsdNS
	request.XSD1NS = xsdNS1
	request.UserName = user.UserName
	request.Password = user.Password
	request.Roles = user.Roles

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	fmt.Printf("xml output : %s", b)

	return bytes.NewBuffer(b)
}

func createAddRoleRequest(role *Role) *bytes.Buffer {
	request := addRoleRequest{}

	request.SOAPEnvNS = xmlNS
	request.XSDNS = xsdNS
	request.XSD1NS = xsdNS1
	request.RoleName = role.RoleName
	request.IsSharedRole = role.IsSharedRole

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	fmt.Printf("xml output : %s", b)

	return bytes.NewBuffer(b)
}
