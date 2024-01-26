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

// Package k8sclient contains the common implementation methods to invoke k8s APIs in the agent
package k8sclient

import (
	"context"

	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// CreateAPICR applies the given API struct to the Kubernetes cluster.
func CreateAPICR(api *dpv1alpha2.API, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), api); err != nil {
		loggers.LoggerXds.Error("Unable to create API CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("API CR created: " + api.Name)
	}
}

// CreateConfigMapCR applies the given ConfigMap struct to the Kubernetes cluster.
func CreateConfigMapCR(configMap *corev1.ConfigMap, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), configMap); err != nil {
		loggers.LoggerXds.Error("Unable to create ConfigMap CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("ConfigMap CR created: " + configMap.Name)
	}
}

// CreateHTTPRouteCR applies the given HttpRoute struct to the Kubernetes cluster.
func CreateHTTPRouteCR(httpRoute *gwapiv1b1.HTTPRoute, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), httpRoute); err != nil {
		loggers.LoggerXds.Error("Unable to create HttpRoute CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("HttpRoute CR created: " + httpRoute.Name)
	}
}

// CreateSecretCR applies the given Secret struct to the Kubernetes cluster.
func CreateSecretCR(secret *corev1.Secret, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), secret); err != nil {
		loggers.LoggerXds.Error("Unable to create Secret CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("Secret CR created: " + secret.Name)
	}
}

// CreateAuthenticationCR applies the given Authentication struct to the Kubernetes cluster.
func CreateAuthenticationCR(authPolicy *dpv1alpha2.Authentication, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), authPolicy); err != nil {
		loggers.LoggerXds.Error("Unable to create Authentication CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("Authentication CR created: " + authPolicy.Name)
	}
}

// CreateBackendJWTCR applies the given BackendJWT struct to the Kubernetes cluster.
func CreateBackendJWTCR(backendJWT *dpv1alpha1.BackendJWT, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), backendJWT); err != nil {
		loggers.LoggerXds.Error("Unable to create BackendJWT CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("BackendJWT CR created: " + backendJWT.Name)
	}
}

// CreateAPIPolicyCR applies the given APIPolicies struct to the Kubernetes cluster.
func CreateAPIPolicyCR(apiPolicies *dpv1alpha2.APIPolicy, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), apiPolicies); err != nil {
		loggers.LoggerXds.Error("Unable to create APIPolicies CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("APIPolicies CR created: " + apiPolicies.Name)
	}
}

// CreateInterceptorServicesCR applies the given InterceptorServices struct to the Kubernetes cluster.
func CreateInterceptorServicesCR(interceptorServices *dpv1alpha1.InterceptorService, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), interceptorServices); err != nil {
		loggers.LoggerXds.Error("Unable to create InterceptorServices CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("InterceptorServices CR created: " + interceptorServices.Name)
	}
}

// CreateScopeCR applies the given Scope struct to the Kubernetes cluster.
func CreateScopeCR(scope *dpv1alpha1.Scope, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), scope); err != nil {
		loggers.LoggerXds.Error("Unable to create Scope CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("Scope CR created: " + scope.Name)
	}
}

// CreateRateLimitPolicyCR applies the given RateLimitPolicies struct to the Kubernetes cluster.
func CreateRateLimitPolicyCR(rateLimitPolicies *dpv1alpha1.RateLimitPolicy, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), rateLimitPolicies); err != nil {
		loggers.LoggerXds.Error("Unable to create RateLimitPolicies CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("RateLimitPolicies CR created: " + rateLimitPolicies.Name)
	}
}

// CreateBackendCR applies the given Backends struct to the Kubernetes cluster.
func CreateBackendCR(backends *dpv1alpha1.Backend, k8sClient client.Client) {
	if err := k8sClient.Create(context.Background(), backends); err != nil {
		loggers.LoggerXds.Error("Unable to create Backends CR: " + err.Error())
	} else {
		loggers.LoggerXds.Info("Backends CR created: " + backends.Name)
	}
}
