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

// Defaults applied to operations generated from a swagger/OAS document when the
// spec itself does not carry the WSO2 x-auth-type / x-throttling-tier vendor extensions.
// These match the values APIM assigns when an API is exported, so that a project
// created via 'apictl init --oas' produces an api.yaml equivalent to an exported one.
const (
	defaultOperationAuthType       = "Application & Application User"
	defaultOperationThrottlingTier = "Unlimited"
	scopesBindingsExtension        = "x-scopes-bindings"
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

// scopeDef holds the merged scope name/description/role-bindings gathered from
// Swagger 2 securityDefinitions or OpenAPI 3 components.securitySchemes.
type scopeDef struct {
	Name        string
	Description string
	Bindings    []string
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

	raw, err := rawSwaggerMap(document)
	if err != nil {
		return err
	}
	def.Scopes = swaggerScopes(raw)
	def.Operations = swaggerOperations(raw)
	return nil
}

// rawSwaggerMap unmarshals the underlying swagger/OAS document once so that scope and
// operation extraction (which need to inspect vendor extensions not modeled by go-openapi's
// typed Swagger2 structs) can share a single parse of the raw JSON.
func rawSwaggerMap(document *loads.Document) (map[string]interface{}, error) {
	raw := map[string]interface{}{}
	if err := json.Unmarshal(document.Raw(), &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// swaggerScopes collects OAuth2 scopes declared in the spec (Swagger 2 securityDefinitions,
// or OpenAPI 3 components.securitySchemes flows), along with role bindings from the WSO2
// x-scopes-bindings vendor extension, and returns them shaped like an exported project's
// api.yaml "scopes" section.
func swaggerScopes(raw map[string]interface{}) []interface{} {
	scopeMap := map[string]scopeDef{}
	collectSwagger2Scopes(raw, scopeMap)
	collectOas3Scopes(raw, scopeMap)

	if len(scopeMap) == 0 {
		return nil
	}
	return scopesToInterfaceSlice(scopeMap)
}

// sortedKeys returns m's keys in sorted order, so callers merging data keyed by these
// maps get deterministic results regardless of Go's randomized map iteration order.
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func collectSwagger2Scopes(raw map[string]interface{}, scopeMap map[string]scopeDef) {
	securityDefinitions, ok := interfaceMap(raw["securityDefinitions"])
	if !ok {
		return
	}
	for _, name := range sortedKeys(securityDefinitions) {
		definitionMap, ok := interfaceMap(securityDefinitions[name])
		if !ok {
			continue
		}
		scopes, _ := interfaceMap(definitionMap["scopes"])
		bindings, _ := interfaceMap(definitionMap[scopesBindingsExtension])
		mergeScopes(scopeMap, scopes, bindings)
	}
}

func collectOas3Scopes(raw map[string]interface{}, scopeMap map[string]scopeDef) {
	components, ok := interfaceMap(raw["components"])
	if !ok {
		return
	}
	securitySchemes, ok := interfaceMap(components["securitySchemes"])
	if !ok {
		return
	}
	for _, schemeName := range sortedKeys(securitySchemes) {
		schemeMap, ok := interfaceMap(securitySchemes[schemeName])
		if !ok {
			continue
		}
		fallbackBindings, _ := interfaceMap(schemeMap[scopesBindingsExtension])
		flows, ok := interfaceMap(schemeMap["flows"])
		if !ok {
			continue
		}
		for _, flowName := range sortedKeys(flows) {
			flowMap, ok := interfaceMap(flows[flowName])
			if !ok {
				continue
			}
			scopes, _ := interfaceMap(flowMap["scopes"])
			bindings, hasBindings := interfaceMap(flowMap[scopesBindingsExtension])
			if !hasBindings {
				bindings = fallbackBindings
			}
			mergeScopes(scopeMap, scopes, bindings)
		}
	}
}

// mergeScopes adds scopes found in a single security definition/scheme into scopeMap,
// filling in description/bindings for scopes already seen (e.g. the same scope declared
// under multiple OAuth2 flows) only when not already populated. Descriptions go through
// stringValue rather than fmt.Sprintf since OAS3's raw-JSON path (unlike Swagger 2.0's
// typed model) doesn't enforce "scopes" values being strings. Scope names and binding
// keys are both trimmed so a scope with stray whitespace (e.g. " admin ") still matches
// its own binding entry and its operation-level reference.
func mergeScopes(scopeMap map[string]scopeDef, scopes map[string]interface{}, bindings map[string]interface{}) {
	bindings = trimMapKeys(bindings)
	for _, rawScopeName := range sortedKeys(scopes) {
		v := scopes[rawScopeName]
		scopeName := strings.TrimSpace(rawScopeName)
		if scopeName == "" {
			continue
		}
		desc := stringValue(v)
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

func scopesToInterfaceSlice(scopeMap map[string]scopeDef) []interface{} {
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

// swaggerOperations builds the "operations" section of api.yaml from the spec's paths,
// resolving each operation's scopes from its own security requirement, falling back to
// API-level (global) security when absent, and applying WSO2's default auth type /
// throttling tier when the spec has no x-auth-type / x-throttling-tier vendor extension.
func swaggerOperations(raw map[string]interface{}) []interface{} {
	paths, ok := interfaceMap(raw["paths"])
	if !ok || len(paths) == 0 {
		return nil
	}

	pathKeys := sortedKeys(paths)

	apiLevelScopes, _ := securityScopes(raw)

	operations := make([]interface{}, 0)
	for _, target := range pathKeys {
		pathObj, ok := interfaceMap(paths[target])
		if !ok {
			continue
		}
		operations = append(operations, operationsForPath(target, pathObj, apiLevelScopes)...)
	}

	if len(operations) == 0 {
		return nil
	}
	return operations
}

var swaggerMethodOrder = []string{"get", "put", "post", "delete", "patch", "options", "head", "trace"}

func operationsForPath(target string, pathObj map[string]interface{}, apiLevelScopes []string) []interface{} {
	// The OpenAPI/Swagger Path Item Object requires lowercase HTTP method keys, but
	// look them up case-insensitively anyway: a spec with e.g. "GET" instead of "get"
	// loads without error (go-openapi's JSON-based parsing doesn't enforce this), and
	// a case-sensitive-only lookup would silently drop the operation - and with it,
	// any scope requirement it declares, which is exactly the class of silent data
	// loss this whole fix exists to prevent.
	operations := make([]interface{}, 0, len(swaggerMethodOrder))
	for _, method := range swaggerMethodOrder {
		opObj, ok := interfaceMapCaseInsensitive(pathObj, method)
		if !ok {
			continue
		}
		operations = append(operations, buildOperation(target, method, opObj, apiLevelScopes))
	}
	return operations
}

// securityScopes reports whether m ("security" is either an Operation Object or the
// root document) declares its own "security" requirement, since a present-but-empty
// array ("no scopes needed here") must be distinguished from an absent key ("inherit
// the API-level default") - collapsing both to an empty list would make an operation
// that explicitly opts out of a global scope incorrectly inherit it anyway.
func securityScopes(m map[string]interface{}) ([]string, bool) {
	v, present := m["security"]
	if !present {
		return nil, false
	}
	return extractScopesFromSecurity(v), true
}

// interfaceMapCaseInsensitive looks up key in m case-insensitively, trying the exact
// key first. The scan visits candidates in sorted order and keeps looking past a
// non-map match rather than stopping, so a spec with duplicate case-variants of the
// same key (e.g. both "GET" and "Get") resolves deterministically.
func interfaceMapCaseInsensitive(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	if v, ok := interfaceMap(m[key]); ok {
		return v, true
	}
	for _, k := range sortedKeys(m) {
		if strings.EqualFold(k, key) {
			if v, ok := interfaceMap(m[k]); ok {
				return v, true
			}
		}
	}
	return nil, false
}

func buildOperation(target, method string, opObj map[string]interface{}, apiLevelScopes []string) map[string]interface{} {
	authType := stringValue(opObj["x-auth-type"])
	if authType == "" {
		authType = defaultOperationAuthType
	}
	throttlingTier := stringValue(opObj["x-throttling-tier"])
	if throttlingTier == "" {
		throttlingTier = defaultOperationThrottlingTier
	}

	// An operation's own "security" - even an explicitly empty array, meaning "no
	// scopes needed here" - always wins outright; only a genuinely absent "security"
	// key falls through to the API level. See securityScopes. There is no path-level
	// tier: neither Swagger 2.0 nor OpenAPI 3.x defines a "security" field on the Path
	// Item Object (confirmed against go-openapi/spec's own PathItemProps, which has no
	// Security field at all) - only the root document and individual operations can
	// declare it.
	scopes, opHasSecurity := securityScopes(opObj)
	if !opHasSecurity {
		scopes = apiLevelScopes
	}

	return map[string]interface{}{
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
	}
}

func interfaceMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}

// trimMapKeys returns a copy of m with every key trimmed of whitespace, so a lookup
// by a trimmed key (see stringValue) still finds it. Keys are copied in sorted order
// so two keys colliding on the same trimmed value resolve deterministically.
func trimMapKeys(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	trimmed := make(map[string]interface{}, len(m))
	for _, k := range sortedKeys(m) {
		trimmed[strings.TrimSpace(k)] = m[k]
	}
	return trimmed
}

// stringValue type-asserts v as a string, trimmed of surrounding whitespace, or
// returns "" if v isn't a string - centralized here so scope/role name comparisons
// aren't broken by stray whitespace in a hand-edited spec.
func stringValue(v interface{}) string {
	str, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(str)
}

func bindingsFromValue(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	// A single string is the real-world multi-role syntax APIM's own import
	// understands (OAS2Parser/OAS3Parser.getScopes() split it on "," server-side) - an
	// actual array crashes the import instead. Split here so a multi-role binding
	// ("role_a,role_b") doesn't survive as one garbled binding literally named that.
	if str, isString := v.(string); isString {
		return splitAndDedupeBindings(strings.Split(str, ","))
	}

	arr, ok := v.([]interface{})
	if !ok {
		return []string{}
	}
	rawValues := make([]string, 0, len(arr))
	for _, item := range arr {
		rawValues = append(rawValues, stringValue(item))
	}
	return splitAndDedupeBindings(rawValues)
}

// splitAndDedupeBindings trims, drops empties, deduplicates, and sorts a set of
// candidate role names shared by both the string (comma-separated) and array forms
// of x-scopes-bindings.
func splitAndDedupeBindings(rawValues []string) []string {
	bindings := make([]string, 0, len(rawValues))
	seen := map[string]bool{}
	for _, raw := range rawValues {
		value := strings.TrimSpace(raw)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		bindings = append(bindings, value)
	}
	sort.Strings(bindings)
	return bindings
}

// extractScopesFromSecurity flattens every scope named across all alternatives of an
// OpenAPI/Swagger "security" requirement array into a single list. This matches APIM's
// own gateway-side scope check, which already treats an operation's scope list as OR
// (any one listed scope satisfies it), so flattening OR-alternatives together here
// doesn't change the enforced semantics.
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
		scopes = appendScopesFromRequirement(requirementMap, scopes, seen)
	}

	sort.Strings(scopes)
	return scopes
}

// appendScopesFromRequirement appends the scope names listed under a single security
// requirement entry (e.g. {"default": ["read:pets", "write:pets"]}) to scopes, skipping
// duplicates already recorded in seen.
func appendScopesFromRequirement(requirementMap map[string]interface{}, scopes []string, seen map[string]bool) []string {
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
