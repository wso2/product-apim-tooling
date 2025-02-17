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

const (
	// Multipart form fields for apk CRD generation request
	apiTypeMultipartField          = "apiType"
	apkConfigurationMultipartField = "apkConfiguration"
	definitionFileMultipartField   = "definitionFile"
	restType                       = "REST"

	// Http protocol related constants
	postHTTPMethod    = "POST"
	contentTypeHeader = "Content-Type"
	internalKeyHeader = "internal-key"
	apiKeyHeader      = "apikey"

	// K8s CRD fields
	k8sKindField                = "kind"
	k8sSpecField                = "spec"
	k8sMetadataField            = "metadata"
	k8sNamespaceField           = "namespace"
	k8sOrganizationField        = "organization"
	k8sHostnamesField           = "hostnames"
	k8sLabelsField              = "labels"
	k8RevisionField             = "revisionID"
	k8APIUuidField              = "apiUUID"
	k8sRateLimitPolicyNameField = "rateLimitPolicyName"

	// K8s CRD values
	k8sKindHTTPRoute   = "HTTPRoute"
	k8sKindAPI         = "API"
	k8sKindTokenIssuer = "TokenIssuer"
	apkCRDAPIVersion   = "dp.wso2.com/v1alpha1"

	// Auth Types
	mTLS   = "mTLS"
	jwt    = "JWT"
	oAuth2 = "OAuth2"
	apiKey = "APIKey"

	// Security Scheme values
	oAuth2SecScheme              = "oauth2"
	applicationSecurityMandatory = "oauth_basic_auth_api_key_mandatory"
	applicationSecurityOptional  = "oauth_basic_auth_api_key_optional"
	mutualSSL                    = "mutualssl"
	mutualSSLMandatory           = "mutualssl_mandatory"
	apiKeySecScheme              = "api_key"

	// Optionality constants
	mandatory = "mandatory"
	optional  = "optional"

	// Interceptor constants
	requestHeader                 = "request_header"
	requestBody                   = "request_body"
	requestTrailers               = "request_trailers"
	requestContext                = "request_context"
	includes                      = "includes"
	interceptorServiceURL         = "interceptorServiceURL"
	https                         = "https"
	requestInterceptorSecretName  = "request-interceptor-tls-secret"
	responseInterceptorSecretName = "response-interceptor-tls-secret"
	tlsKey                        = "tls.crt"

	// BackendJWT constants
	encoding         = "encoding"
	header           = "header"
	signingAlgorithm = "signingAlgorithm"
	tokenTTL         = "tokenTTL"
	base64Url        = "Base64Url"

	// APK Operation Policy constants
	interceptorPolicy     = "Interceptor"
	backendJWTPolicy      = "BackendJwt"
	addHeaderPolicy       = "AddHeader"
	removeHeaderPolicy    = "RemoveHeader"
	requestRedirectPolicy = "RequestRedirect"
	requestMirrorPolicy   = "RequestMirror"
	modelBasedRoundRobin  = "ModelBasedRoundRobin"

	// APK BackendJWT parameter constants
	base64url = "Base64url"

	// APK header modification parameter constants
	url         = "url"
	statusCode  = "statusCode"
	headerName  = "headerName"
	headerValue = "headerValue"

	// Version constants
	v1 = "v1"
	v2 = "v2"
)
