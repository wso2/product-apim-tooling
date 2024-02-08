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
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// DeployAPICR applies the given API struct to the Kubernetes cluster.
func DeployAPICR(api *dpv1alpha2.API, k8sClient client.Client) {
	crAPI := &dpv1alpha2.API{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: api.ObjectMeta.Namespace, Name: api.Name}, crAPI); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get API CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), api); err != nil {
			loggers.LoggerXds.Error("Unable to create API CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("API CR created: " + api.Name)
		}
	} else {
		crAPI.Spec = api.Spec
		if err := k8sClient.Update(context.Background(), crAPI); err != nil {
			loggers.LoggerXds.Error("Unable to update API CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("API CR updated: " + api.Name)
		}
	}
}

// DeployConfigMapCR applies the given ConfigMap struct to the Kubernetes cluster.
func DeployConfigMapCR(configMap *corev1.ConfigMap, k8sClient client.Client) {
	crConfigMap := &corev1.ConfigMap{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: configMap.ObjectMeta.Namespace, Name: configMap.Name}, crConfigMap); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get ConfigMap CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), configMap); err != nil {
			loggers.LoggerXds.Error("Unable to create ConfigMap CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("ConfigMap CR created: " + configMap.Name)
		}
	} else {
		crConfigMap.Data = configMap.Data
		if err := k8sClient.Update(context.Background(), crConfigMap); err != nil {
			loggers.LoggerXds.Error("Unable to update ConfigMap CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("ConfigMap CR updated: " + configMap.Name)
		}
	}
}

// DeployHTTPRouteCR applies the given HttpRoute struct to the Kubernetes cluster.
func DeployHTTPRouteCR(httpRoute *gwapiv1b1.HTTPRoute, k8sClient client.Client) {
	crHTTPRoute := &gwapiv1b1.HTTPRoute{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: httpRoute.ObjectMeta.Namespace, Name: httpRoute.Name}, crHTTPRoute); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get HTTPRoute CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), httpRoute); err != nil {
			loggers.LoggerXds.Error("Unable to create HTTPRoute CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("HTTPRoute CR created: " + httpRoute.Name)
		}
	} else {
		crHTTPRoute.Spec = httpRoute.Spec
		if err := k8sClient.Update(context.Background(), crHTTPRoute); err != nil {
			loggers.LoggerXds.Error("Unable to update HTTPRoute CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("HTTPRoute CR updated: " + httpRoute.Name)
		}
	}
}

// DeploySecretCR applies the given Secret struct to the Kubernetes cluster.
func DeploySecretCR(secret *corev1.Secret, k8sClient client.Client) {
	crSecret := &corev1.Secret{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: secret.ObjectMeta.Namespace, Name: secret.Name}, crSecret); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get Secret CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), secret); err != nil {
			loggers.LoggerXds.Error("Unable to create Secret CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Secret CR created: " + secret.Name)
		}
	} else {
		crSecret.Data = secret.Data
		if err := k8sClient.Update(context.Background(), crSecret); err != nil {
			loggers.LoggerXds.Error("Unable to update Secret CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Secret CR updated: " + secret.Name)
		}
	}
}

// DeployAuthenticationCR applies the given Authentication struct to the Kubernetes cluster.
func DeployAuthenticationCR(authPolicy *dpv1alpha2.Authentication, k8sClient client.Client) {
	crAuthPolicy := &dpv1alpha2.Authentication{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: authPolicy.ObjectMeta.Namespace, Name: authPolicy.Name}, crAuthPolicy); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get Authentication CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), authPolicy); err != nil {
			loggers.LoggerXds.Error("Unable to create Authentication CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Authentication CR created: " + authPolicy.Name)
		}
	} else {
		crAuthPolicy.Spec = authPolicy.Spec
		if err := k8sClient.Update(context.Background(), crAuthPolicy); err != nil {
			loggers.LoggerXds.Error("Unable to update Authentication CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Authentication CR updated: " + authPolicy.Name)
		}
	}
}

// DeployBackendJWTCR applies the given BackendJWT struct to the Kubernetes cluster.
func DeployBackendJWTCR(backendJWT *dpv1alpha1.BackendJWT, k8sClient client.Client) {
	crBackendJWT := &dpv1alpha1.BackendJWT{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: backendJWT.ObjectMeta.Namespace, Name: backendJWT.Name}, crBackendJWT); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get BackendJWT CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), backendJWT); err != nil {
			loggers.LoggerXds.Error("Unable to create BackendJWT CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("BackendJWT CR created: " + backendJWT.Name)
		}
	} else {
		crBackendJWT.Spec = backendJWT.Spec
		if err := k8sClient.Update(context.Background(), crBackendJWT); err != nil {
			loggers.LoggerXds.Error("Unable to update BackendJWT CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("BackendJWT CR updated: " + backendJWT.Name)
		}
	}
}

// DeployAPIPolicyCR applies the given APIPolicies struct to the Kubernetes cluster.
func DeployAPIPolicyCR(apiPolicies *dpv1alpha2.APIPolicy, k8sClient client.Client) {
	crAPIPolicies := &dpv1alpha2.APIPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: apiPolicies.ObjectMeta.Namespace, Name: apiPolicies.Name}, crAPIPolicies); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get APIPolicies CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), apiPolicies); err != nil {
			loggers.LoggerXds.Error("Unable to create APIPolicies CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("APIPolicies CR created: " + apiPolicies.Name)
		}
	} else {
		crAPIPolicies.Spec = apiPolicies.Spec
		if err := k8sClient.Update(context.Background(), crAPIPolicies); err != nil {
			loggers.LoggerXds.Error("Unable to update APIPolicies CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("APIPolicies CR updated: " + apiPolicies.Name)
		}
	}
}

// DeployInterceptorServicesCR applies the given InterceptorServices struct to the Kubernetes cluster.
func DeployInterceptorServicesCR(interceptorServices *dpv1alpha1.InterceptorService, k8sClient client.Client) {
	crInterceptorServices := &dpv1alpha1.InterceptorService{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: interceptorServices.ObjectMeta.Namespace, Name: interceptorServices.Name}, crInterceptorServices); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get InterceptorServices CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), interceptorServices); err != nil {
			loggers.LoggerXds.Error("Unable to create InterceptorServices CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("InterceptorServices CR created: " + interceptorServices.Name)
		}
	} else {
		crInterceptorServices.Spec = interceptorServices.Spec
		if err := k8sClient.Update(context.Background(), crInterceptorServices); err != nil {
			loggers.LoggerXds.Error("Unable to update InterceptorServices CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("InterceptorServices CR updated: " + interceptorServices.Name)
		}
	}
}

// DeployScopeCR applies the given Scope struct to the Kubernetes cluster.
func DeployScopeCR(scope *dpv1alpha1.Scope, k8sClient client.Client) {
	crScope := &dpv1alpha1.Scope{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: scope.ObjectMeta.Namespace, Name: scope.Name}, crScope); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get Scope CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), scope); err != nil {
			loggers.LoggerXds.Error("Unable to create Scope CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Scope CR created: " + scope.Name)
		}
	} else {
		crScope.Spec = scope.Spec
		if err := k8sClient.Update(context.Background(), crScope); err != nil {
			loggers.LoggerXds.Error("Unable to update Scope CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Scope CR updated: " + scope.Name)
		}
	}
}

// DeployRateLimitPolicyCR applies the given RateLimitPolicies struct to the Kubernetes cluster.
func DeployRateLimitPolicyCR(rateLimitPolicies *dpv1alpha1.RateLimitPolicy, k8sClient client.Client) {
	crRateLimitPolicies := &dpv1alpha1.RateLimitPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: rateLimitPolicies.ObjectMeta.Namespace, Name: rateLimitPolicies.Name}, crRateLimitPolicies); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get RateLimitPolicies CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), rateLimitPolicies); err != nil {
			loggers.LoggerXds.Error("Unable to create RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("RateLimitPolicies CR created: " + rateLimitPolicies.Name)
		}
	} else {
		crRateLimitPolicies.Spec = rateLimitPolicies.Spec
		if err := k8sClient.Update(context.Background(), crRateLimitPolicies); err != nil {
			loggers.LoggerXds.Error("Unable to update RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("RateLimitPolicies CR updated: " + rateLimitPolicies.Name)
		}
	}
}

// DeployBackendCR applies the given Backends struct to the Kubernetes cluster.
func DeployBackendCR(backends *dpv1alpha1.Backend, k8sClient client.Client) {
	crBackends := &dpv1alpha1.Backend{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: backends.ObjectMeta.Namespace, Name: backends.Name}, crBackends); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerXds.Error("Unable to get Backends CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), backends); err != nil {
			loggers.LoggerXds.Error("Unable to create Backends CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Backends CR created: " + backends.Name)
		}
	} else {
		crBackends.Spec = backends.Spec
		if err := k8sClient.Update(context.Background(), crBackends); err != nil {
			loggers.LoggerXds.Error("Unable to update Backends CR: " + err.Error())
		} else {
			loggers.LoggerXds.Info("Backends CR updated: " + backends.Name)
		}
	}
}
