/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package transformer

// APIIdentifier  holds information about an API associated for a given client certificate
type APIIdentifier struct {
	ProviderName string `json:"providerName"`
	APIName      string `json:"apiName"`
	Version      string `json:"version"`
	UUID         string `json:"uuid"`
	ID           int    `json:"id"`
}

// ClientCert holds the data belongs to a single client certificate configuration
type ClientCert struct {
	Alias         string        `json:"alias"`
	Certificate   string        `json:"certificate"`
	TierName      string        `json:"tierName"`
	APIIdentifier APIIdentifier `json:"apiIdentifier"`
}

// CertDescriptor contains data related to one or more client certificates for an API
type CertDescriptor struct {
	CertData []ClientCert `json:"data"`
}

// MetadataBlock will hold the metadata information belong to the particular configmap
type MetadataBlock struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

// CertConfigYaml holds the data belong to the main configmap structure created upon enabling mTLS
type CertConfigYaml struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   MetadataBlock     `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}
