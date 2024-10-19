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

/*
 * Package "mapper" contains artifacts relate to fetching APIs and
 * API related updates from the control plane event-hub.
 * This file contains functions to retrieve APIs and API updates.
 */

package mapper

import (
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	internalk8sClient "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/k8sClient"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/transformer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MapAndCreateCR will read the CRD Yaml and based on the Kind of the CR, unmarshal and maps the
// data and sends to the K8-Client for creating the respective CR inside the cluster
func MapAndCreateCR(k8sArtifact transformer.K8sArtifacts, k8sClient client.Client) *error {
	namespace, err := getDeploymentNamespace(k8sArtifact)
	if err != nil {
		return &err
	}
	k8sArtifact.API.Namespace = namespace
	
	for _, configMaps := range k8sArtifact.ConfigMaps {
		configMaps.Namespace = namespace
		internalk8sClient.DeployConfigMapCR(configMaps, k8sClient)
	}
	for _, authPolicies := range k8sArtifact.Authentication {
		authPolicies.Namespace = namespace
		internalk8sClient.DeployAuthenticationCR(authPolicies, k8sClient)
	}
	for _, interceptorServices := range k8sArtifact.InterceptorServices {
		interceptorServices.Namespace = namespace
		internalk8sClient.DeployInterceptorServicesCR(interceptorServices, k8sClient)
	}
	if k8sArtifact.BackendJWT != nil {
		k8sArtifact.BackendJWT.Namespace = namespace
		internalk8sClient.DeployBackendJWTCR(k8sArtifact.BackendJWT, k8sClient)
	}
	for _, scopes := range k8sArtifact.Scopes {
		scopes.Namespace = namespace
		internalk8sClient.DeployScopeCR(scopes, k8sClient)
	}
	for _, rateLimitPolicy := range k8sArtifact.RateLimitPolicies {
		rateLimitPolicy.Namespace = namespace
		internalk8sClient.DeployRateLimitPolicyCR(rateLimitPolicy, k8sClient)
	}
	for _, aiRateLimitPolicy := range k8sArtifact.AIRateLimitPolicies {
		aiRateLimitPolicy.Namespace = namespace
		internalk8sClient.DeployAIRateLimitPolicyCR(aiRateLimitPolicy, k8sClient)
	}
	for _, secrets := range k8sArtifact.Secrets {
		secrets.Namespace = namespace
		internalk8sClient.DeploySecretCR(secrets, k8sClient)
	}
	for _, apiPolicies := range k8sArtifact.APIPolicies {
		apiPolicies.Namespace = namespace
		internalk8sClient.DeployAPIPolicyCR(apiPolicies, k8sClient)
	}
	for _, httpRoutes := range k8sArtifact.HTTPRoutes {
		httpRoutes.Namespace = namespace
		internalk8sClient.DeployHTTPRouteCR(httpRoutes, k8sClient)
	}
	for _, gqlRoutes := range k8sArtifact.GQLRoutes {
		gqlRoutes.Namespace = namespace
		internalk8sClient.DeployGQLRouteCR(gqlRoutes, k8sClient)
	}
	for _, backends := range k8sArtifact.Backends {
		backends.Namespace = namespace
		internalk8sClient.DeployBackendCR(backends, k8sClient)
	}
	internalk8sClient.DeployAPICR(&k8sArtifact.API, k8sClient)
	return nil
}
func getDeploymentNamespace(k8sArtifact transformer.K8sArtifacts) (string, error) {
	conf, errReadConfig := config.ReadConfigs()
	if errReadConfig != nil {
		logger.LoggerMapper.Errorf("Error reading configs: %v", errReadConfig)
		return "", errReadConfig
	}
	return conf.DataPlane.Namespace, nil
}
