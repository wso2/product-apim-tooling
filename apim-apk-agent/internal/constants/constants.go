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

package constants

// Gateway related constants
const (
	GatewayName  = "wso2-apk-default"
	GatewayGroup = "gateway.networking.k8s.io"
	GatewayKind  = "Gateway"
)

// TokenIssuer related constants
const (
	ConsumerKeyClaim           = "azp"
	ScopesClaim                = "scope"
	InternalKeyTokenIssuerName = "Internal Key TokenIssuer"
	InternalKeySecretName      = "apim-apk-issuer-cert"
	InternalKeySecretKey       = "wso2.crt"
	InternalKeySuffix          = "-internal-key-issuer"
)

// APIM Mediation constants
const (
	InterceptorService      = "CallInterceptorService"
	BackendJWT              = "backEndJWT"
	AddHeader               = "apkAddHeader"
	RemoveHeader            = "apkRemoveHeader"
	MirrorRequest           = "apkMirrorRequest"
	RedirectRequest         = "apkRedirectRequest"
	ModelWeightedRoundRobin = "modelWeightedRoundRobin"
	ModelRoundRobin         = "modelRoundRobin"

	// Version constants
	V1 = "v1"
	V2 = "v2"

	// Policy Types
	CommonType = "common"
)
