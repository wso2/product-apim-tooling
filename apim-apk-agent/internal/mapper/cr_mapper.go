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
	"archive/zip"
	"io"

	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	internalk8sClient "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/k8sClient"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	k8Yaml "sigs.k8s.io/yaml"
)

// MapAndCreateCR will read the CRD Yaml and based on the Kind of the CR, unmarshal and maps the
// data and sends to the K8-Client for creating the respective CR inside the cluster
func MapAndCreateCR(zipFile *zip.File, k8sClient client.Client, conf *config.Config) (string, interface{}) {
	fileReader, err := zipFile.Open()
	if err != nil {
		logger.LoggerTransformer.Errorf("Failed to open YAML file inside zip: %v", err)
		return "", nil
	}
	defer fileReader.Close()

	yamlData, err := io.ReadAll(fileReader)
	if err != nil {
		logger.LoggerTransformer.Errorf("Failed to read YAML file inside zip: %v", err)
		return "", nil
	}

	var crdData map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &crdData); err != nil {
		logger.LoggerTransformer.Errorf("Failed to unmarshal YAML data to parse the Kind: %v", err)
		return "", nil
	}

	kind, ok := crdData["kind"].(string)
	if !ok {
		logger.LoggerTransformer.Errorf("Kind attribute not found in the given yaml file.")
		return "", nil
	}

	switch kind {
	case "APIPolicy":
		var apiPolicy dpv1alpha2.APIPolicy
		err = k8Yaml.Unmarshal(yamlData, &apiPolicy)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling APIPolicy YAML: %v", err)
		}
		apiPolicy.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateAPIPolicyCR(&apiPolicy, k8sClient)
	case "HTTPRoute":
		var httpRoute gwapiv1b1.HTTPRoute
		err = k8Yaml.Unmarshal(yamlData, &httpRoute)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling HTTPRoute YAML: %v", err)
		}
		httpRoute.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateHTTPRouteCR(&httpRoute, k8sClient)
	case "Backend":
		var backend dpv1alpha1.Backend
		err = k8Yaml.Unmarshal(yamlData, &backend)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling Backend YAML: %v", err)
		}
		backend.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateBackendCR(&backend, k8sClient)
	case "ConfigMap":
		var configMap corev1.ConfigMap
		err = k8Yaml.Unmarshal(yamlData, &configMap)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling ConfigMap YAML: %v", err)
		}
		configMap.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateConfigMapCR(&configMap, k8sClient)
	case "Authentication":
		var authPolicy dpv1alpha2.Authentication
		err = k8Yaml.Unmarshal(yamlData, &authPolicy)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling Authentication YAML: %v", err)
		}
		authPolicy.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateAuthenticationCR(&authPolicy, k8sClient)
	case "API":
		var api dpv1alpha2.API
		err = k8Yaml.Unmarshal(yamlData, &api)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling API YAML: %v", err)
		}
		api.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateAPICR(&api, k8sClient)
	case "InterceptorService":
		var interceptorService dpv1alpha1.InterceptorService
		err = k8Yaml.Unmarshal(yamlData, &interceptorService)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling InterceptorService YAML: %v", err)
		}
		interceptorService.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateInterceptorServicesCR(&interceptorService, k8sClient)
	case "BackendJWT":
		var backendJWT dpv1alpha1.BackendJWT
		err = k8Yaml.Unmarshal(yamlData, &backendJWT)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling BackendJWT YAML: %v", err)
		}
		backendJWT.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateBackendJWTCR(&backendJWT, k8sClient)
	case "Scope":
		var scope dpv1alpha1.Scope
		err = k8Yaml.Unmarshal(yamlData, &scope)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling Scope YAML: %v", err)
		}
		scope.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateScopeCR(&scope, k8sClient)
	case "RateLimitPolicy":
		var rateLimitPolicy dpv1alpha1.RateLimitPolicy
		err = k8Yaml.Unmarshal(yamlData, &rateLimitPolicy)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling RateLimitPolicy YAML: %v", err)
		}
		rateLimitPolicy.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateRateLimitPolicyCR(&rateLimitPolicy, k8sClient)
	case "Secret":
		var secret corev1.Secret
		err = k8Yaml.Unmarshal(yamlData, &secret)
		if err != nil {
			logger.LoggerSync.Errorf("Error unmarshaling Secret YAML: %v", err)
		}
		secret.ObjectMeta.Namespace = conf.DataPlane.Namespace
		internalk8sClient.CreateSecretCR(&secret, k8sClient)
	default:
		logger.LoggerSync.Errorf("[!]Unknown Kind parsed from the YAML File: %v", kind)
	}
	return kind, crdData
}
