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

package v2

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"

	"io/ioutil"
	"os"

	"github.com/wso2/product-apim-tooling/import-export-cli/specs/params"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/Jeffail/gabs"
	"github.com/go-openapi/loads"
	"github.com/mitchellh/mapstructure"
)

// Servers represent servers of an AWS API
type Servers struct {
	Servers []struct {
		Url       string `json:"url"`
		Variables struct {
			BasePath struct {
				Default string `json:"default"`
			} `json:"basePath"`
		} `json:"variables"`
	} `json:"servers"`
}

// EndpointConfig represent EndpointConfigs of an AWS API
type EndpointConfig struct {
	EndpointType        string `yaml:"endpoint_type" json:"endpoint_type"`
	ProductionEndpoints struct {
		URL string `yaml:"url" json:"url"`
	} `yaml:"production_endpoints" json:"production_endpoints"`
	SandboxEndpoints struct {
		URL string `yaml:"url" json:"url"`
	} `yaml:"sandbox_endpoints" json:"sandbox_endpoints"`
}

// SecuritySchemes represent security schemes of an AWS API
type SecuritySchemes struct {
	Components struct {
		SecuritySchemes struct {
			CognitoAuthorizer struct {
				AuthType string `json:"x-amazon-apigateway-authtype"`
			} `json:"CognitoAuthorizer"`
			APIKey struct {
				Type string `json:"type"`
			} `json:"api_key"`
			Sigv4 struct {
				AuthType string `json:"x-amazon-apigateway-authtype"`
			} `json:"sigv4"`
		} `json:"securitySchemes"`
	} `json:"components"`
	ResourcePolicy struct {
		Version string `json:"Version"`
	} `json:"x-amazon-apigateway-policy"`
}

func swagger2XWO2BasePath(document *loads.Document) (string, bool) {
	if v, ok := document.Spec().Extensions["x-wso2-basePath"]; ok {
		str, ok := v.(string)
		return str, ok
	}
	return "", false
}

func swagger2XWSO2Cors(document *loads.Document) (*CorsConfiguration, bool, error) {
	if v, ok := document.Spec().Extensions["x-wso2-cors"]; ok {
		var cors CorsConfiguration
		err := mapstructure.Decode(v, &cors)
		if err != nil {
			return nil, true, err
		}
		cors.CorsConfigurationEnabled = true
		return &cors, true, nil
	}
	return nil, false, nil
}

func swagger2Tags(document *loads.Document) []string {
	tags := make([]string, len(document.Spec().Tags))
	for i, v := range document.Spec().Tags {
		tags[i] = v.Name
	}
	return tags
}

func swagger2XWSO2ProductionEndpoints(document *loads.Document) (*Endpoints, bool, error) {
	if v, ok := document.Spec().Extensions["x-wso2-production-endpoints"]; ok {
		var prodEp Endpoints
		err := mapstructure.Decode(v, &prodEp)
		if err != nil {
			return nil, true, err
		}
		return &prodEp, true, nil
	}
	return &Endpoints{}, false, nil
}

func swagger2XWSO2SandboxEndpoints(document *loads.Document) (*Endpoints, bool, error) {
	if v, ok := document.Spec().Extensions["x-wso2-sandbox-endpoints"]; ok {
		var sandboxEp Endpoints
		err := mapstructure.Decode(v, &sandboxEp)
		if err != nil {
			return nil, true, err
		}
		return &sandboxEp, true, nil
	}
	return &Endpoints{}, false, nil
}

// BuildAPIMEndpoints builds endpointConfig for given config
func BuildAPIMEndpoints(production, sandbox *Endpoints) (string, error) {
	epType := EpHttp
	if len(production.Urls) > 1 {
		epType = EpLoadbalance
		if production.Type == EpFailover {
			epType = EpFailover
		}
	}

	if len(production.Urls) == 0 {
		if len(sandbox.Urls) > 1 {
			epType = EpLoadbalance
		}
		if sandbox.Type == EpFailover {
			epType = EpFailover
		}
	}

	switch epType {
	case EpHttp:
		endpoint := buildHttpEndpoint(production, sandbox)
		return endpoint, nil
	case EpLoadbalance:
		endpoint := buildLoadBalancedEndpoints(production, sandbox)
		return endpoint, nil
	case EpFailover:
		endpoint := buildFailOverEndpoints(production, sandbox)
		return endpoint, nil
	default:
		return "", fmt.Errorf("unknown endpoint type")
	}
}

func buildFailOverEndpoints(production *Endpoints, sandbox *Endpoints) string {
	jsonObj, _ := gabs.ParseJSON([]byte(`
					{
						"endpoint_type": "failover",
		    			"algoCombo": "org.apache.synapse.endpoints.algorithms.RoundRobin",
		    			"algoClassName": "",
						"sessionManagement": "",
		    			"sessionTimeOut": "",
		    			"failOver": "True"
					}
				`))
	if len(production.Urls) > 0 {
		buildFailOverUrls(jsonObj, production, "production")
	}
	if len(sandbox.Urls) > 0 {
		buildFailOverUrls(jsonObj, sandbox, "sandbox")
	}
	return jsonObj.String()
}

func buildFailOverUrls(jsonObj *gabs.Container, endpoints *Endpoints, eptype string) {
	_, _ = jsonObj.Set(params.Endpoint{Url: &endpoints.Urls[0]}, fmt.Sprintf("%s_endpoints", eptype))
	rest := endpoints.Urls[1:]
	if len(rest) > 0 {
		fo := make([]params.Endpoint, len(rest))
		for i := 0; i < len(fo); i++ {
			fo[i] = params.Endpoint{Url: &rest[i]}
		}
		if len(fo) > 0 {
			_, _ = jsonObj.Set(fo, fmt.Sprintf("%s_failovers", eptype))
		}
	}
}

func buildLoadBalancedEndpoints(production *Endpoints, sandbox *Endpoints) string {
	jsonObj, _ := gabs.ParseJSON([]byte(`
		{
			"endpoint_type": "load_balance",
		    "algoCombo": "org.apache.synapse.endpoints.algorithms.RoundRobin",
		    "algoClassName": "org.apache.synapse.endpoints.algorithms.RoundRobin",
		    "sessionManagement": "",
			"sessionTimeOut": "",
			"failover" : "False"
		}
	`))
	prodEps := make([]params.Endpoint, len(production.Urls))
	for i := 0; i < len(prodEps); i++ {
		prodEps[i] = params.Endpoint{Url: &production.Urls[i]}
	}
	if len(prodEps) > 0 {
		_, _ = jsonObj.Set(prodEps, "production_endpoints")
	}

	sandboxEps := make([]params.Endpoint, len(sandbox.Urls))
	for i := 0; i < len(sandboxEps); i++ {
		sandboxEps[i] = params.Endpoint{Url: &sandbox.Urls[i]}
	}
	if len(sandboxEps) > 0 {
		_, _ = jsonObj.Set(sandboxEps, "sandbox_endpoints")
	}

	return jsonObj.String()
}

func buildHttpEndpoint(production *Endpoints, sandbox *Endpoints) string {
	jsonObj := gabs.New()
	_, _ = jsonObj.Set(EpHttp, "endpoint_type")
	if len(production.Urls) > 0 {
		var ep params.Endpoint
		ep.Url = &production.Urls[0]
		if production.AdvanceEndpointConfig != nil && production.AdvanceEndpointConfig.TimeOutInMillis != nil {
			ep.AdvanceEndpointConfig = &params.AdvanceEndpointConfiguration{
				TimeOutInMillis: production.AdvanceEndpointConfig.TimeOutInMillis,
			}
		}
		_, _ = jsonObj.SetP(ep, "production_endpoints")
	}
	if len(sandbox.Urls) > 0 {
		var ep params.Endpoint
		ep.Url = &sandbox.Urls[0]
		if sandbox.AdvanceEndpointConfig != nil && sandbox.AdvanceEndpointConfig.TimeOutInMillis != nil {
			ep.AdvanceEndpointConfig = &params.AdvanceEndpointConfiguration{
				TimeOutInMillis: sandbox.AdvanceEndpointConfig.TimeOutInMillis,
			}
		}
		_, _ = jsonObj.SetP(ep, "sandbox_endpoints")
	}
	return jsonObj.String()
}

// generateFieldsFromSwagger3 using swagger
func Swagger2Populate(def *APIDTODefinition, document *loads.Document) error {
	def.Name = document.Spec().Info.Title
	def.Version = document.Spec().Info.Version
	def.Provider = "admin"
	def.Description = document.Spec().Info.Description
	def.Context = fmt.Sprintf("/%s", def.Name)
	def.Tags = swagger2Tags(document)

	// fill basepath from swagger
	if document.BasePath() != "" {
		def.Context = path.Clean(fmt.Sprintf("/%s", document.BasePath()))
	}

	// override basepath if wso2 extension provided
	if basepath, ok := swagger2XWO2BasePath(document); ok {
		def.Context = path.Clean(basepath)
		if !strings.Contains(basepath, "{version}") {
			if strings.Contains(basepath, def.Version) {
				def.Context = path.Clean(strings.Replace(basepath, def.Version, "",
					strings.LastIndex(basepath, def.Version)))
			} else {
				def.Context = path.Clean(basepath)
			}
			def.IsDefaultVersion = true
		} else {
			def.Context = path.Clean(strings.ReplaceAll(basepath, "{version}", def.Version))
		}
	}

	// trim spaces if available
	def.Name = strings.ReplaceAll(def.Name, " ", "")
	def.Version = strings.ReplaceAll(def.Version, " ", "")
	def.Context = strings.ReplaceAll(def.Context, " ", "")

	cors, ok, err := swagger2XWSO2Cors(document)
	if err != nil && ok {
		return err
	}
	if ok {
		def.CorsConfiguration = cors
	}
	prodEp, foundProdEp, err := swagger2XWSO2ProductionEndpoints(document)
	if err != nil && foundProdEp {
		return err
	}
	sandboxEp, foundSandboxEp, err := swagger2XWSO2SandboxEndpoints(document)
	if err != nil && foundSandboxEp {
		return err
	}
	if foundProdEp || foundSandboxEp {
		ep, err := BuildAPIMEndpoints(prodEp, sandboxEp)
		if err != nil {
			return err
		}
		var endpointConfig map[string]interface{}
		err = json.Unmarshal([]byte(ep), &endpointConfig)
		if err != nil {
			return err
		}
		def.EndpointConfig = &endpointConfig
	}

	def.Scopes = swaggerScopes(document)
	def.Operations = swaggerOperations(document)
	return nil
}

func swaggerScopes(document *loads.Document) []interface{} {
	raw := map[string]interface{}{}
	if err := json.Unmarshal(document.Raw(), &raw); err != nil {
		return nil
	}

	type scopeDef struct {
		Name        string
		Description string
		Bindings    []string
	}

	scopeMap := map[string]scopeDef{}
	mergeScopes := func(scopes map[string]interface{}, bindings map[string]interface{}) {
		for scopeName, v := range scopes {
			desc := fmt.Sprintf("%v", v)
			if existing, found := scopeMap[scopeName]; found {
				if existing.Description == "" && desc != "" {
					existing.Description = desc
				}
				if len(existing.Bindings) == 0 {
					existing.Bindings = bindingsFromValue(bindings[scopeName])
				}
				scopeMap[scopeName] = existing
				continue
			}

			scopeMap[scopeName] = scopeDef{
				Name:        scopeName,
				Description: desc,
				Bindings:    bindingsFromValue(bindings[scopeName]),
			}
		}
	}

	// Swagger 2: securityDefinitions
	if securityDefinitions, ok := interfaceMap(raw["securityDefinitions"]); ok {
		for _, definition := range securityDefinitions {
			definitionMap, ok := interfaceMap(definition)
			if !ok {
				continue
			}
			scopes, _ := interfaceMap(definitionMap["scopes"])
			bindings, _ := interfaceMap(definitionMap["x-scopes-bindings"])
			mergeScopes(scopes, bindings)
		}
	}

	// OpenAPI 3: components.securitySchemes.*.flows.*
	if components, ok := interfaceMap(raw["components"]); ok {
		if securitySchemes, ok := interfaceMap(components["securitySchemes"]); ok {
			for _, scheme := range securitySchemes {
				schemeMap, ok := interfaceMap(scheme)
				if !ok {
					continue
				}

				fallbackBindings, _ := interfaceMap(schemeMap["x-scopes-bindings"])
				if flows, ok := interfaceMap(schemeMap["flows"]); ok {
					for _, flow := range flows {
						flowMap, ok := interfaceMap(flow)
						if !ok {
							continue
						}
						scopes, _ := interfaceMap(flowMap["scopes"])
						bindings, hasBindings := interfaceMap(flowMap["x-scopes-bindings"])
						if !hasBindings {
							bindings = fallbackBindings
						}
						mergeScopes(scopes, bindings)
					}
				}
			}
		}
	}

	if len(scopeMap) == 0 {
		return nil
	}

	names := make([]string, 0, len(scopeMap))
	for name := range scopeMap {
		names = append(names, name)
	}
	sort.Strings(names)

	scopes := make([]interface{}, 0, len(names))
	for _, name := range names {
		scope := scopeMap[name]
		scopes = append(scopes, map[string]interface{}{
			"scope": map[string]interface{}{
				"name":        scope.Name,
				"displayName": scope.Name,
				"description": scope.Description,
				"bindings":    scope.Bindings,
			},
			"shared": false,
		})
	}

	return scopes
}

func swaggerOperations(document *loads.Document) []interface{} {
	raw := map[string]interface{}{}
	if err := json.Unmarshal(document.Raw(), &raw); err != nil {
		return nil
	}

	paths, ok := interfaceMap(raw["paths"])
	if !ok || len(paths) == 0 {
		return nil
	}

	pathKeys := make([]string, 0, len(paths))
	for p := range paths {
		pathKeys = append(pathKeys, p)
	}
	sort.Strings(pathKeys)

	methodOrder := []string{"get", "put", "post", "delete", "patch", "options", "head", "trace"}

	operations := make([]interface{}, 0)
	for _, target := range pathKeys {
		pathObj, ok := interfaceMap(paths[target])
		if !ok {
			continue
		}

		pathLevelScopes := extractScopesFromSecurity(pathObj["security"])
		apiLevelScopes := extractScopesFromSecurity(raw["security"])

		for _, method := range methodOrder {
			opObj, ok := interfaceMap(pathObj[method])
			if !ok {
				continue
			}

			authType := stringValue(opObj["x-auth-type"])
			throttlingTier := stringValue(opObj["x-throttling-tier"])

			scopes := extractScopesFromSecurity(opObj["security"])
			if len(scopes) == 0 {
				scopes = pathLevelScopes
			}
			if len(scopes) == 0 {
				scopes = apiLevelScopes
			}

			operations = append(operations, map[string]interface{}{
				"id":               "",
				"target":           target,
				"verb":             strings.ToUpper(method),
				"authType":         authType,
				"throttlingPolicy": throttlingTier,
				"scopes":           scopes,
				"usedProductIds":   []string{},
				"operationPolicies": map[string]interface{}{
					"request":  []string{},
					"response": []string{},
					"fault":    []string{},
				},
			})
		}
	}

	if len(operations) == 0 {
		return nil
	}

	return operations
}

func interfaceMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}

func stringValue(v interface{}) string {
	str, ok := v.(string)
	if !ok {
		return ""
	}
	return str
}

func bindingsFromValue(v interface{}) []string {
	if v == nil {
		return []string{}
	}

	if value, ok := v.(string); ok {
		if value == "" {
			return []string{}
		}
		return []string{value}
	}

	arr, ok := v.([]interface{})
	if !ok {
		return []string{}
	}

	bindings := make([]string, 0, len(arr))
	seen := map[string]bool{}
	for _, item := range arr {
		value := stringValue(item)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		bindings = append(bindings, value)
	}
	sort.Strings(bindings)
	return bindings
}

func extractScopesFromSecurity(v interface{}) []string {
	requirements, ok := v.([]interface{})
	if !ok {
		return []string{}
	}

	scopes := make([]string, 0)
	seen := map[string]bool{}
	for _, requirement := range requirements {
		requirementMap, ok := interfaceMap(requirement)
		if !ok {
			continue
		}

		for _, rawScopes := range requirementMap {
			scopeArray, ok := rawScopes.([]interface{})
			if !ok {
				continue
			}
			for _, scope := range scopeArray {
				scopeName := stringValue(scope)
				if scopeName == "" || seen[scopeName] {
					continue
				}
				seen[scopeName] = true
				scopes = append(scopes, scopeName)
			}
		}
	}

	sort.Strings(scopes)
	return scopes
}

func AddAwsTag(def *APIDTODefinition) {
	def.Tags = append(def.Tags, "AWS") //adding the "aws" tag to all APIs imported using the "aws init" command
}

func GetServerUrlAndSecuritySchemes(pathToSwagger string) (string, string, []byte) {
	oas3File, err := os.Open(pathToSwagger)
	if err != nil {
		utils.HandleErrorAndExit("", err)
	}
	defer oas3File.Close()

	byteValue, _ := ioutil.ReadAll(oas3File)

	var servers Servers
	json.Unmarshal(byteValue, &servers)

	url := servers.Servers[0].Url
	stage := servers.Servers[0].Variables.BasePath.Default
	productionUrl := strings.ReplaceAll(url, "/{basePath}", stage)
	sandboxUrl := strings.ReplaceAll(url, "/{basePath}", stage)
	return productionUrl, sandboxUrl, byteValue
}

func CreateEpConfigForAwsAPIs(def *APIDTODefinition, pathToSwagger string) []byte {
	var endpointConfig EndpointConfig
	var productionEp, sandboxEp, byteValue = GetServerUrlAndSecuritySchemes(pathToSwagger)
	endpointConfig.EndpointType = "http"
	endpointConfig.ProductionEndpoints.URL = productionEp
	endpointConfig.SandboxEndpoints.URL = sandboxEp
	def.EndpointConfig = &endpointConfig
	var advertiseInfo AdvertiseInfo
	advertiseInfo.Advertised = def.AdvertiseInformation.Advertised
	advertiseInfo.ApiOwner = def.AdvertiseInformation.ApiOwner
	advertiseInfo.Vendor = def.AdvertiseInformation.Vendor
	advertiseInfo.OriginalDevPortalUrl = def.AdvertiseInformation.OriginalDevPortalUrl
	advertiseInfo.ApiExternalProductionEndpoint = productionEp
	advertiseInfo.ApiExternalSandboxEndpoint = sandboxEp
	def.AdvertiseInformation = advertiseInfo
	return byteValue
}
