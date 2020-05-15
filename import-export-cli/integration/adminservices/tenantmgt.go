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

// Tenant : data strcuture defining tenant
type Tenant struct {
	AdminUserName string
	AdminPassword string
	Domain        string
}

type addTenantRequest struct {
	XMLName       xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS     string   `xml:"xmlns:soapenv,attr"`
	SERNS         string   `xml:"xmlns:ser,attr"`
	XSDNS         string   `xml:"xmlns:xsd,attr"`
	Active        bool     `xml:"soapenv:Body>ser:addSkeletonTenant>ser:tenantInfoBean>xsd:active"`
	Admin         string   `xml:"soapenv:Body>ser:addSkeletonTenant>ser:tenantInfoBean>xsd:admin"`
	AdminPassword string   `xml:"soapenv:Body>ser:addSkeletonTenant>ser:tenantInfoBean>xsd:adminPassword"`
	Email         string   `xml:"soapenv:Body>ser:addSkeletonTenant>ser:tenantInfoBean>xsd:email"`
	FirstName     string   `xml:"soapenv:Body>ser:addSkeletonTenant>ser:tenantInfoBean>xsd:firstname"`
	LastName      string   `xml:"soapenv:Body>ser:addSkeletonTenant>ser:tenantInfoBean>xsd:lastname"`
	TenantDomain  string   `xml:"soapenv:Body>ser:addSkeletonTenant>ser:tenantInfoBean>xsd:tenantDomain"`
}

type getTenantRequest struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS    string   `xml:"xmlns:soapenv,attr"`
	SERNS        string   `xml:"xmlns:ser,attr"`
	TenantDomain string   `xml:"soapenv:Body>ser:getTenant>ser:tenantDomain"`
}

type getTenantResponse struct {
	RawXML   string `xml:",innerxml"`
	isActive bool
	isExists bool
}

func (instance *getTenantResponse) parse() {
	splits := strings.Split(instance.RawXML, "\n")

	isActiveFound := false
	isTenantIDFound := false

	for _, s := range splits {
		if !isActiveFound && strings.Contains(s, "active>true") {
			instance.isActive = true
			isActiveFound = true
		}

		if !isTenantIDFound && strings.Contains(s, "tenantId") {
			if !strings.Contains(s, ">0</") {
				instance.isExists = true
			}
			isTenantIDFound = true
		}
	}
}

type activateTenantRequest struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS    string   `xml:"xmlns:soapenv,attr"`
	SERNS        string   `xml:"xmlns:ser,attr"`
	TenantDomain string   `xml:"soapenv:Body>ser:activateTenant>ser:tenantDomain"`
}

type deactivateTenantRequest struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	SOAPEnvNS    string   `xml:"xmlns:soapenv,attr"`
	SERNS        string   `xml:"xmlns:ser,attr"`
	TenantDomain string   `xml:"soapenv:Body>ser:deactivateTenant>ser:tenantDomain"`
}

// InitTenant : Init new tenant in APIM
func InitTenant(tenant *Tenant, url string, username string, password string) {
	getResp := getTenant(tenant.Domain, url, username, password)

	if !getResp.isExists {
		addTenant(tenant, url, username, password)
		activateTenant(tenant.Domain, url, username, password)
	} else if !getResp.isActive {
		activateTenant(tenant.Domain, url, username, password)
	}
}

// DeactivateTenant : Deactivate Tenant in APIM
func DeactivateTenant(tenantDomain string, url string, username string, password string) {
	msg := createDeactivateTenantRequestPayload(tenantDomain)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:deactivateTenant", request)

	base.LogRequest("adminservices.DeactivateTenant()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.DeactivateTenant()", response, 200)
}

func addTenant(tenant *Tenant, url string, username string, password string) {
	msg := createAddTenantRequestPayload(tenant)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:addTenant", request)

	base.LogRequest("adminservices.addTenant()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.addTenant()", response, 200)
}

func getTenant(tenantDomain string, url string, username string, password string) getTenantResponse {
	msg := createGetTenantRequestPayload(tenantDomain)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:getTenant", request)

	base.LogRequest("adminservices.getTenant()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.getTenant()", response, 200)

	getResp := getTenantResponse{}
	xml.NewDecoder(response.Body).Decode(&getResp)
	getResp.parse()

	return getResp
}

func activateTenant(tenantDomain string, url string, username string, password string) {
	msg := createActivateTenantRequestPayload(tenantDomain)

	request := base.CreatePost(url, msg)

	base.SetSOAPHeaders(username, password, "urn:activateTenant", request)

	base.LogRequest("adminservices.activateTenant()", request)

	response := base.SendHTTPRequest(request)

	defer response.Body.Close()

	base.ValidateAndLogResponse("adminservices.activateTenant()", response, 200)
}

func createAddTenantRequestPayload(tenant *Tenant) *bytes.Buffer {
	request := addTenantRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = tenantSerNS
	request.XSDNS = xsdNS
	request.Active = true
	request.Admin = tenant.AdminUserName
	request.AdminPassword = tenant.AdminPassword
	request.Email = tenant.AdminUserName + "@" + tenant.Domain
	request.FirstName = tenant.AdminUserName
	request.LastName = tenant.AdminUserName
	request.TenantDomain = tenant.Domain

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

func createGetTenantRequestPayload(tenantDomain string) *bytes.Buffer {
	request := getTenantRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = tenantSerNS
	request.TenantDomain = tenantDomain

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

func createActivateTenantRequestPayload(tenantDomain string) *bytes.Buffer {
	request := activateTenantRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = tenantSerNS
	request.TenantDomain = tenantDomain

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

func createDeactivateTenantRequestPayload(tenantDomain string) *bytes.Buffer {
	request := deactivateTenantRequest{}

	request.SOAPEnvNS = xmlNS
	request.SERNS = tenantSerNS
	request.TenantDomain = tenantDomain

	b, err := xml.MarshalIndent(request, "", "    ")

	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}
