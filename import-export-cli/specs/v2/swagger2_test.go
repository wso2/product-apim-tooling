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
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
)

func Test_swagger2WSO2Cors(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	cors, ok, err := swagger2XWSO2Cors(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, []string{"GET", "PUT", "POST"}, cors.AccessControlAllowMethods, "should have same elements for access control")
	assert.ElementsMatch(t, []string{"test.com", "example.com"}, cors.AccessControlAllowOrigins, "should have same elements for origins")
}

func Test_swagger2Tags(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	tags := swagger2Tags(doc)
	assert.ElementsMatch(t, []string{"pet", "user", "store"}, tags, "should have same elements")
}

func Test_swagger2WSO2ProductionEndpoints(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := swagger2XWSO2ProductionEndpoints(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_swagger2WSO2SandboxEndpoints(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := swagger2XWSO2SandboxEndpoints(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func TestSwagger2Populate(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "Swagger Petstore", def.Name, "Should return correct api name")
	assert.Equal(t, "/petstore/v1/1.0.0", def.Context)
}

func findOperation(t *testing.T, operations []interface{}, target, verb string) map[string]interface{} {
	t.Helper()
	for _, op := range operations {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}
		if opMap["target"] == target && opMap["verb"] == verb {
			return opMap
		}
	}
	t.Fatalf("operation not found: %s %s", verb, target)
	return nil
}

func findScope(t *testing.T, scopes []interface{}, name string) map[string]interface{} {
	t.Helper()
	for _, s := range scopes {
		scopeEntry, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		scope, ok := scopeEntry["scope"].(map[string]interface{})
		if !ok {
			continue
		}
		if scope["name"] == name {
			return scope
		}
	}
	t.Fatalf("scope not found: %s", name)
	return nil
}

func TestSwagger2PopulateAddsScopesAndOperations(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 2, "expected oauth scopes to be populated")
	assert.NotEmpty(t, def.Operations, "expected operations to be populated")

	postOp := findOperation(t, def.Operations, "/pet", "POST")
	assert.Equal(t, "POST", postOp["verb"])
	assert.Equal(t, "/pet", postOp["target"])
	assert.ElementsMatch(t, []string{"read:pets", "write:pets"}, postOp["scopes"].([]string))
	// petstore_swagger2.yaml carries no x-auth-type/x-throttling-tier, so the exported-project defaults apply
	assert.Equal(t, defaultOperationAuthType, postOp["authType"])
	assert.Equal(t, defaultOperationThrottlingTier, postOp["throttlingPolicy"])

	getByID := findOperation(t, def.Operations, "/pet/{petId}", "GET")
	assert.ElementsMatch(t, []string{}, getByID["scopes"].([]string))
}

func TestSwagger2PopulateOas3ScopesWithBindingsAndDefaults(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/oas3_scopes.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 1, "expected one oauth2 scope from components.securitySchemes")
	scope := findScope(t, def.Scopes, "test4_scope")
	assert.Equal(t, "Scope to protect menu resource", scope["description"])
	assert.ElementsMatch(t, []string{"role4"}, scope["bindings"])

	// /menu has no x-auth-type/x-throttling-tier -> falls back to exported-project defaults
	menuOp := findOperation(t, def.Operations, "/menu", "GET")
	assert.ElementsMatch(t, []string{"test4_scope"}, menuOp["scopes"].([]string))
	assert.Equal(t, defaultOperationAuthType, menuOp["authType"])
	assert.Equal(t, defaultOperationThrottlingTier, menuOp["throttlingPolicy"])

	// /order POST sets explicit x-auth-type/x-throttling-tier -> must be preserved, not overridden
	orderPost := findOperation(t, def.Operations, "/order", "POST")
	assert.Equal(t, "Application", orderPost["authType"])
	assert.Equal(t, "Gold", orderPost["throttlingPolicy"])

	// /order GET has no security requirement anywhere -> empty scopes, still gets default authType/tier
	orderGet := findOperation(t, def.Operations, "/order", "GET")
	assert.ElementsMatch(t, []string{}, orderGet["scopes"].([]string))
	assert.Equal(t, defaultOperationAuthType, orderGet["authType"])
	assert.Equal(t, defaultOperationThrottlingTier, orderGet["throttlingPolicy"])
}

func TestSwagger2PopulateAwsGatewayExportNoSpuriousScopes(t *testing.T) {
	// AWS Gateway security schemes (apiKey, sigv4, Cognito) are never oauth2-with-flows.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/aws_gateway_export.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "should not error on AWS Gateway apiKey/sigv4 security schemes")
	assert.Empty(t, def.Scopes, "apiKey/sigv4 schemes have no oauth2 flows, so no scopes should be generated")
	op := findOperation(t, def.Operations, "/pets", "GET")
	assert.Equal(t, []string{}, op["scopes"].([]string), "sigv4 is not oauth2, so the operation itself must have no scopes either")
}

func TestSwagger2PopulateNonStringScopeDescriptionBecomesEmpty(t *testing.T) {
	// Only reachable via OAS3: Swagger 2.0's typed model rejects a non-string
	// description at load time, but OAS3's raw-JSON path doesn't enforce that.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/oas3_non_string_scope_description.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	nullScope := findScope(t, def.Scopes, "null_scope")
	assert.Equal(t, "", nullScope["description"], "a null description must become empty, not the literal string '<nil>'")

	objectScope := findScope(t, def.Scopes, "object_scope")
	assert.Equal(t, "", objectScope["description"], "a non-string (object) description must become empty, not its Go representation")
}

func TestSwagger2PopulateSkipsEmptyScopeName(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/empty_scope_name.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 1, "an empty-string scope name must be skipped, not turned into a nameless scope entry")
	findScope(t, def.Scopes, "real_scope")
}

func TestSwagger2PopulateMatchesHttpMethodsCaseInsensitively(t *testing.T) {
	// "GET" instead of "get" loads fine (go-openapi doesn't enforce casing here) - a
	// case-sensitive-only lookup would silently drop the operation and its scopes.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/uppercase_method.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Operations, 1, "an uppercase 'GET' method key must still be captured as an operation")
	op := findOperation(t, def.Operations, "/test", "GET")
	assert.Equal(t, "/test", op["target"])
}

func TestSwagger2PopulateTrimsWhitespaceConsistently(t *testing.T) {
	// Trimming only one side of a name/key comparison would silently break the match.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/whitespace_padded_names.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 1, "the padded scope name must resolve to exactly one scope")
	scope := findScope(t, def.Scopes, "padded_scope")
	assert.Equal(t, "a description with padding", scope["description"])
	assert.ElementsMatch(t, []string{"padded_role"}, scope["bindings"])

	op := findOperation(t, def.Operations, "/test", "GET")
	assert.ElementsMatch(t, []string{"padded_scope"}, op["scopes"].([]string),
		"the operation must reference the same trimmed scope name as the declared scope, not a mismatched padded variant")
}

func TestSwagger2PopulateMultiSchemeScopeMergeAndPrecedence(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/multi_scheme_scopes.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 4, "expected scopes merged across both security schemes: shared_scope, order_write, api_level_scope, order_read")

	// shared_scope: scheme_a's empty binding must not block scheme_b's from filling in.
	shared := findScope(t, def.Scopes, "shared_scope")
	assert.Equal(t, "description from scheme_a", shared["description"])
	assert.ElementsMatch(t, []string{"role_shared_reader"}, shared["bindings"])

	// api_level_scope has a multi-role (array) binding
	apiLevel := findScope(t, def.Scopes, "api_level_scope")
	assert.ElementsMatch(t, []string{"role_admin", "role_superadmin"}, apiLevel["bindings"])

	// order_write has a single-role (string, not array) binding
	orderWrite := findScope(t, def.Scopes, "order_write")
	assert.ElementsMatch(t, []string{"role_writer"}, orderWrite["bindings"])

	// Operation-level security must be used verbatim, not merged with the API-level default
	ordersPost := findOperation(t, def.Operations, "/orders", "POST")
	assert.ElementsMatch(t, []string{"order_write"}, ordersPost["scopes"].([]string))

	ordersGet := findOperation(t, def.Operations, "/orders", "GET")
	assert.ElementsMatch(t, []string{"order_read"}, ordersGet["scopes"].([]string))

	// /admin PUT has no security requirement of its own, so it must fall back to the API-level (global) security
	adminPut := findOperation(t, def.Operations, "/admin", "PUT")
	assert.ElementsMatch(t, []string{"api_level_scope"}, adminPut["scopes"].([]string))
}

func TestSwagger2PopulateWithBasePath(t *testing.T) {
	var def1, def2 APIDTODefinition

	// Basepath without {version}
	doc1, err1 := loads.Spec("testdata/petstore_with_basepath1.yaml")
	assert.Nil(t, err1, "err should be nil")
	err1 = Swagger2Populate(&def1, doc1)
	assert.Nil(t, err1, "err should be nil")

	assert.Equal(t, "/petstore/v1/1.0.0", def1.Context)
	assert.Equal(t, true, def1.IsDefaultVersion)

	// Basepath with {version}
	doc2, err2 := loads.Spec("testdata/petstore_with_basepath2.yaml")
	assert.Nil(t, err2, "err should be nil")
	err2 = Swagger2Populate(&def2, doc2)
	assert.Nil(t, err2, "err should be nil")

	assert.Equal(t, "/petstore/v1/1.0.0", def2.Context)
	assert.Equal(t, false, def2.IsDefaultVersion)
}

func TestSwagger2PopulateNoSecurityAnywhereYieldsNilScopes(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/no_security_anywhere.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Nil(t, def.Scopes, "a spec with no security definitions anywhere must yield nil scopes, not an empty-but-non-nil slice")
	op := findOperation(t, def.Operations, "/open", "GET")
	assert.ElementsMatch(t, []string{}, op["scopes"].([]string))
	assert.Equal(t, defaultOperationAuthType, op["authType"])
	assert.Equal(t, defaultOperationThrottlingTier, op["throttlingPolicy"])
}

func TestSwagger2PopulateOas3MultipleFlowsSameSchemeMergeDeterministically(t *testing.T) {
	// Same scope declared under two flows of one scheme must merge deterministically.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/oas3_multi_flow_same_scheme.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 1)
	scope := findScope(t, def.Scopes, "shared_flow_scope")
	assert.Equal(t, "description from authorizationCode flow", scope["description"],
		"'authorizationCode' sorts before 'password', so its description must win")
	assert.ElementsMatch(t, []string{"role_from_auth_code_flow"}, scope["bindings"])
}

func TestSwagger2PopulateWhitespaceOnlyScopeNameSkipped(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/whitespace_only_scope_name.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 1, "a scope name that is nothing but whitespace must be trimmed to empty and skipped, same as a literal empty string")
	findScope(t, def.Scopes, "real_scope")
}

func TestSwagger2PopulateBindingEdgeCases(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/binding_edge_cases.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	dup := findScope(t, def.Scopes, "dup_binding_scope")
	assert.ElementsMatch(t, []string{"role_a", "role_b"}, dup["bindings"], "a role repeated in the bindings array must be deduplicated")

	noBindings := findScope(t, def.Scopes, "no_bindings_scope")
	assert.Equal(t, []string{}, noBindings["bindings"], "a scope with no x-scopes-bindings entry at all must get an empty (not nil) bindings list")

	allInvalid := findScope(t, def.Scopes, "all_invalid_bindings_scope")
	assert.Equal(t, []string{}, allInvalid["bindings"], "a bindings array containing only non-string elements must yield an empty list, not garbage string conversions")

	commaBound := findScope(t, def.Scopes, "comma_binding_scope")
	assert.ElementsMatch(t, []string{"role_a", "role_b"}, commaBound["bindings"], "a comma-separated string binding - the real format APIM's own import path understands - must split into separate roles, not survive as one literal 'role_a, role_b' binding")
}

func TestBindingsFromValueSplitsCommaSeparatedString(t *testing.T) {
	// The real server (OAS2Parser/OAS3Parser.getScopes()) splits on "," itself; an
	// actual array crashes the import instead - must split the same way here.
	assert.Equal(t, []string{"role_a", "role_b"}, bindingsFromValue("role_a,role_b"))
	assert.Equal(t, []string{"role_a", "role_b"}, bindingsFromValue("role_a, role_b"))
	assert.Equal(t, []string{"solo_role"}, bindingsFromValue("solo_role"))
	assert.Equal(t, []string{"role_a", "role_b"}, bindingsFromValue("role_a,role_a,role_b"))
	assert.Equal(t, []string{}, bindingsFromValue(""))
	assert.Equal(t, []string{}, bindingsFromValue(",,"))
}

func TestSwagger2PopulateMultipleSchemesInOneSecurityRequirement(t *testing.T) {
	// One requirement entry naming two schemes together - both must be captured.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/multi_scheme_single_requirement.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	op := findOperation(t, def.Operations, "/both-required", "GET")
	assert.ElementsMatch(t, []string{"scope_from_a", "scope_from_b"}, op["scopes"].([]string))
}

func TestSwagger2PopulateExplicitEmptySecurityOverridesInsteadOfInheriting(t *testing.T) {
	// "security: []" (present but empty) must override the API-level default, not
	// inherit it - only a genuinely absent "security" key inherits. See securityScopes.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/explicit_empty_security_vs_inherited.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	public := findOperation(t, def.Operations, "/public", "GET")
	assert.ElementsMatch(t, []string{}, public["scopes"].([]string),
		"an operation with an explicit empty security array must NOT inherit the API-level scope")

	inherits := findOperation(t, def.Operations, "/inherits", "GET")
	assert.ElementsMatch(t, []string{"api_level_scope"}, inherits["scopes"].([]string),
		"an operation with no security key at all must still inherit the API-level scope")
}

func TestSwagger2PopulateScopeEntryHasCorrectShape(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/whitespace_padded_names.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 1)
	entry := def.Scopes[0].(map[string]interface{})
	assert.Equal(t, false, entry["shared"], "apictl init always generates local (non-shared) scopes")
	scope := entry["scope"].(map[string]interface{})
	assert.Equal(t, "padded_scope", scope["name"], "name must be trimmed")
	assert.Equal(t, "padded_scope", scope["displayName"], "displayName always mirrors name - apictl has no way to author a distinct display name from a plain OpenAPI spec")
}

func TestSwagger2PopulateOperationEntryHasCorrectShape(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/whitespace_padded_names.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	op := findOperation(t, def.Operations, "/test", "GET")
	assert.Equal(t, "", op["id"])
	assert.Equal(t, []string{}, op["usedProductIds"])
	policies := op["operationPolicies"].(map[string]interface{})
	assert.Equal(t, []string{}, policies["request"])
	assert.Equal(t, []string{}, policies["response"])
	assert.Equal(t, []string{}, policies["fault"])
}

func TestSwagger2PopulateMixedCaseMethodKeysResolveDeterministically(t *testing.T) {
	// Duplicate case-variants ("GET" and "Get") get different scopes so the winner
	// is distinguishable - "exactly one operation" alone can't prove determinism.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/mixed_case_method_keys.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	op := findOperation(t, def.Operations, "/test", "GET")
	assert.Equal(t, []string{"get_scope"}, op["scopes"].([]string),
		"the 'GET' key must win over 'Get' (sortedKeys orders 'GET' first), and alt_get_scope from the losing variant must not leak in")
}

func TestInterfaceMapCaseInsensitiveSkipsNonMapMatchToFindAValidOne(t *testing.T) {
	// A non-map case-insensitive match must not abort the scan before a later,
	// valid map match is reached.
	m := map[string]interface{}{
		"GET": "not a map - malformed",
		"Get": map[string]interface{}{"summary": "the real operation"},
	}
	for i := 0; i < 20; i++ {
		result, ok := interfaceMapCaseInsensitive(m, "get")
		assert.True(t, ok, "must find the valid map variant, not give up on the non-map one")
		assert.Equal(t, "the real operation", result["summary"])
	}
}

func TestSwagger2PopulateOrphanScopeStillListedEvenIfUnreferenced(t *testing.T) {
	// A scope declared but never referenced by any operation must still appear in
	// def.Scopes - swaggerScopes and swaggerOperations are independent passes.
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/orphan_scope_unreferenced.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Len(t, def.Scopes, 2)
	orphan := findScope(t, def.Scopes, "orphan_scope")
	assert.ElementsMatch(t, []string{"role_orphan"}, orphan["bindings"])

	op := findOperation(t, def.Operations, "/test", "GET")
	assert.ElementsMatch(t, []string{"used_scope"}, op["scopes"].([]string),
		"the orphan scope must not leak into an operation that never referenced it")
}

func TestSwagger2PopulateScalarBindingValueYieldsEmpty(t *testing.T) {
	// Neither a string nor an array (e.g. a bare number/boolean) must yield [].
	assert.Equal(t, []string{}, bindingsFromValue(float64(123)))
	assert.Equal(t, []string{}, bindingsFromValue(true))
}
