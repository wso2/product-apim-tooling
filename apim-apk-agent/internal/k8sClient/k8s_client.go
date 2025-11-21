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
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"

	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/constants"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/logging"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	corev1 "k8s.io/api/core/v1"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// DeployAPICR applies the given API struct to the Kubernetes cluster.
func DeployAPICR(api *dpv1alpha3.API, k8sClient client.Client) {
	crAPI := &dpv1alpha3.API{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: api.ObjectMeta.Namespace, Name: api.Name}, crAPI); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get API CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), api); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create API CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("API CR created: " + api.Name)
		}
	} else {
		crAPI.Spec = api.Spec
		crAPI.ObjectMeta.Labels = api.ObjectMeta.Labels
		if err := k8sClient.Update(context.Background(), crAPI); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update API CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("API CR updated: " + api.Name)
		}
	}
}

// UndeployK8sAPICR removes the API Custom Resource from the Kubernetes cluster based on API ID label.
func UndeployK8sAPICR(k8sClient client.Client, k8sAPI dpv1alpha3.API) error {
	err := k8sClient.Delete(context.Background(), &k8sAPI, &client.DeleteOptions{})
	if err != nil {
		loggers.LoggerK8sClient.Errorf("Unable to delete API CR: %v", err)
		return err
	}
	loggers.LoggerK8sClient.Infof("Deleted API CR: %s", k8sAPI.Name)
	return nil
}

// UndeployAPICR removes the API Custom Resource from the Kubernetes cluster based on API ID label.
func UndeployAPICR(apiID string, k8sClient client.Client) {
	conf, errReadConfig := config.ReadConfigs()
	if errReadConfig != nil {
		loggers.LoggerK8sClient.Errorf("Error reading configurations: %v", errReadConfig)
	}
	apiList := &dpv1alpha3.APIList{}
	err := k8sClient.List(context.Background(), apiList, &client.ListOptions{Namespace: conf.DataPlane.Namespace, LabelSelector: labels.SelectorFromSet(map[string]string{"apiUUID": apiID})})
	// Retrieve all API CRs from the Kubernetes cluster
	if err != nil {
		loggers.LoggerK8sClient.Errorf("Unable to list API CRs: %v", err)
	}
	for _, api := range apiList.Items {
		if err := UndeployK8sAPICR(k8sClient, api); err != nil {
			loggers.LoggerK8sClient.Errorf("Unable to delete API CR: %v", err)
		}
		loggers.LoggerK8sClient.Infof("Deleted API CR: %s", api.Name)
	}
}

// DeployConfigMapCR applies the given ConfigMap struct to the Kubernetes cluster.
func DeployConfigMapCR(configMap *corev1.ConfigMap, k8sClient client.Client) {
	crConfigMap := &corev1.ConfigMap{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: configMap.ObjectMeta.Namespace, Name: configMap.Name}, crConfigMap); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get ConfigMap CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), configMap); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create ConfigMap CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("ConfigMap CR created: " + configMap.Name)
		}
	} else {
		crConfigMap.Data = configMap.Data
		if err := k8sClient.Update(context.Background(), crConfigMap); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update ConfigMap CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("ConfigMap CR updated: " + configMap.Name)
		}
	}
}

// DeployHTTPRouteCR applies the given HttpRoute struct to the Kubernetes cluster.
func DeployHTTPRouteCR(httpRoute *gwapiv1.HTTPRoute, k8sClient client.Client) {
	crHTTPRoute := &gwapiv1.HTTPRoute{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: httpRoute.ObjectMeta.Namespace, Name: httpRoute.Name}, crHTTPRoute); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get HTTPRoute CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), httpRoute); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create HTTPRoute CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("HTTPRoute CR created: " + httpRoute.Name)
		}
	} else {
		crHTTPRoute.Spec = httpRoute.Spec
		if err := k8sClient.Update(context.Background(), crHTTPRoute); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update HTTPRoute CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("HTTPRoute CR updated: " + httpRoute.Name)
		}
	}
}

// DeployGQLRouteCR applies the given GqlRoute struct to the Kubernetes cluster.
func DeployGQLRouteCR(gqlRoute *dpv1alpha2.GQLRoute, k8sClient client.Client) {
	crGQLRoute := &dpv1alpha2.GQLRoute{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: gqlRoute.ObjectMeta.Namespace, Name: gqlRoute.Name}, crGQLRoute); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get GQLRoute CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), gqlRoute); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create GQLRoute CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("GQLRoute CR created: " + gqlRoute.Name)
		}
	} else {
		crGQLRoute.Spec = gqlRoute.Spec
		if err := k8sClient.Update(context.Background(), crGQLRoute); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update GQLRoute CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("GQLRoute CR updated: " + gqlRoute.Name)
		}
	}
}

// DeploySecretCR applies the given Secret struct to the Kubernetes cluster.
func DeploySecretCR(secret *corev1.Secret, k8sClient client.Client) {
	crSecret := &corev1.Secret{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: secret.ObjectMeta.Namespace, Name: secret.Name}, crSecret); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get Secret CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), secret); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create Secret CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Secret CR created: " + secret.Name)
		}
	} else {
		crSecret.Data = secret.Data
		if err := k8sClient.Update(context.Background(), crSecret); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update Secret CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Secret CR updated: " + secret.Name)
		}
	}
}

// DeployAuthenticationCR applies the given Authentication struct to the Kubernetes cluster.
func DeployAuthenticationCR(authPolicy *dpv1alpha2.Authentication, k8sClient client.Client) {
	crAuthPolicy := &dpv1alpha2.Authentication{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: authPolicy.ObjectMeta.Namespace, Name: authPolicy.Name}, crAuthPolicy); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get Authentication CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), authPolicy); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create Authentication CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Authentication CR created: " + authPolicy.Name)
		}
	} else {
		crAuthPolicy.Spec = authPolicy.Spec
		if err := k8sClient.Update(context.Background(), crAuthPolicy); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update Authentication CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Authentication CR updated: " + authPolicy.Name)
		}
	}
}

// DeployBackendJWTCR applies the given BackendJWT struct to the Kubernetes cluster.
func DeployBackendJWTCR(backendJWT *dpv1alpha1.BackendJWT, k8sClient client.Client) {
	crBackendJWT := &dpv1alpha1.BackendJWT{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: backendJWT.ObjectMeta.Namespace, Name: backendJWT.Name}, crBackendJWT); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get BackendJWT CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), backendJWT); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create BackendJWT CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("BackendJWT CR created: " + backendJWT.Name)
		}
	} else {
		crBackendJWT.Spec = backendJWT.Spec
		if err := k8sClient.Update(context.Background(), crBackendJWT); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update BackendJWT CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("BackendJWT CR updated: " + backendJWT.Name)
		}
	}
}

// DeployAPIPolicyCR applies the given APIPolicies struct to the Kubernetes cluster.
func DeployAPIPolicyCR(apiPolicies *dpv1alpha4.APIPolicy, k8sClient client.Client) {
	crAPIPolicies := &dpv1alpha4.APIPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: apiPolicies.ObjectMeta.Namespace, Name: apiPolicies.Name}, crAPIPolicies); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get APIPolicies CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), apiPolicies); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create APIPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("APIPolicies CR created: " + apiPolicies.Name)
		}
	} else {
		crAPIPolicies.Spec = apiPolicies.Spec
		if err := k8sClient.Update(context.Background(), crAPIPolicies); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update APIPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("APIPolicies CR updated: " + apiPolicies.Name)
		}
	}
}

// DeployInterceptorServicesCR applies the given InterceptorServices struct to the Kubernetes cluster.
func DeployInterceptorServicesCR(interceptorServices *dpv1alpha1.InterceptorService, k8sClient client.Client) {
	crInterceptorServices := &dpv1alpha1.InterceptorService{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: interceptorServices.ObjectMeta.Namespace, Name: interceptorServices.Name}, crInterceptorServices); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get InterceptorServices CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), interceptorServices); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create InterceptorServices CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("InterceptorServices CR created: " + interceptorServices.Name)
		}
	} else {
		crInterceptorServices.Spec = interceptorServices.Spec
		if err := k8sClient.Update(context.Background(), crInterceptorServices); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update InterceptorServices CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("InterceptorServices CR updated: " + interceptorServices.Name)
		}
	}
}

// DeployScopeCR applies the given Scope struct to the Kubernetes cluster.
func DeployScopeCR(scope *dpv1alpha1.Scope, k8sClient client.Client) {
	crScope := &dpv1alpha1.Scope{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: scope.ObjectMeta.Namespace, Name: scope.Name}, crScope); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get Scope CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), scope); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create Scope CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Scope CR created: " + scope.Name)
		}
	} else {
		crScope.Spec = scope.Spec
		if err := k8sClient.Update(context.Background(), crScope); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update Scope CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Scope CR updated: " + scope.Name)
		}
	}
}

// DeployAIProviderCR applies the given AIProvider struct to the Kubernetes cluster.
func DeployAIProviderCR(aiProvider *dpv1alpha4.AIProvider, k8sClient client.Client) {
	crAIProvider := &dpv1alpha4.AIProvider{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: aiProvider.ObjectMeta.Namespace, Name: aiProvider.Name}, crAIProvider); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get AIProvider CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), aiProvider); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create AIProvider CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("AIProvider CR created: " + aiProvider.Name)
		}
	} else {
		crAIProvider.Spec = aiProvider.Spec
		if err := k8sClient.Update(context.Background(), crAIProvider); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update AIProvider CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("AIProvider CR updated: " + aiProvider.Name)
		}
	}
}

// DeleteAIProviderCR removes the AIProvider Custom Resource from the Kubernetes cluster based on CR name
func DeleteAIProviderCR(aiProviderName string, k8sClient client.Client) {
	conf, errReadConfig := config.ReadConfigs()
	if errReadConfig != nil {
		loggers.LoggerK8sClient.Errorf("Error reading configurations: %v", errReadConfig)
		return
	}

	crAIProvider := &dpv1alpha4.AIProvider{}
	err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: conf.DataPlane.Namespace, Name: aiProviderName}, crAIProvider)
	if err != nil {
		if k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Infof("AI Provider CR not found: %s", aiProviderName)
		} else {
			loggers.LoggerK8sClient.Error("Unable to get AIProvider CR: " + err.Error())
		}
		return
	}

	// Proceed to delete the CR if it was successfully retrieved
	err = k8sClient.Delete(context.Background(), crAIProvider, &client.DeleteOptions{})
	if err != nil {
		loggers.LoggerK8sClient.Errorf("Unable to delete AI Provider CR: %v", err)
	} else {
		loggers.LoggerK8sClient.Infof("Deleted AI Provider CR: %s Successfully", aiProviderName)
	}
}

// DeleteAIRatelimitPolicy removes the AIRatelimitPolicy Custom Resource from the Kubernetes cluster based on CR name
func DeleteAIRatelimitPolicy(airlName string, k8sClient client.Client) {
	conf, errReadConfig := config.ReadConfigs()
	if errReadConfig != nil {
		loggers.LoggerK8sClient.Errorf("Error reading configurations: %v", errReadConfig)
		return
	}

	crAIRatelimitPolicy := &dpv1alpha3.AIRateLimitPolicy{}
	err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: conf.DataPlane.Namespace, Name: airlName}, crAIRatelimitPolicy)
	if err != nil {
		if k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Infof("AIRatelimitPolicy CR not found: %s", airlName)
		} else {
			loggers.LoggerK8sClient.Error("Unable to get AIRatelimitPolicy CR: " + err.Error())
		}
		return
	}

	// Proceed to delete the CR if it was successfully retrieved
	err = k8sClient.Delete(context.Background(), crAIRatelimitPolicy, &client.DeleteOptions{})
	if err != nil {
		loggers.LoggerK8sClient.Errorf("Unable to delete AIRatelimitPolicy CR: %v", err)
	} else {
		loggers.LoggerK8sClient.Infof("Deleted AIRatelimitPolicy CR: %s Successfully", airlName)
	}
}

// DeployRateLimitPolicyCR applies the given RateLimitPolicies struct to the Kubernetes cluster.
func DeployRateLimitPolicyCR(rateLimitPolicies *dpv1alpha1.RateLimitPolicy, k8sClient client.Client) {
	crRateLimitPolicies := &dpv1alpha1.RateLimitPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: rateLimitPolicies.ObjectMeta.Namespace, Name: rateLimitPolicies.Name}, crRateLimitPolicies); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get RateLimitPolicies CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), rateLimitPolicies); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("RateLimitPolicies CR created: " + rateLimitPolicies.Name)
		}
	} else {
		crRateLimitPolicies.Spec = rateLimitPolicies.Spec
		crRateLimitPolicies.ObjectMeta.Labels = rateLimitPolicies.ObjectMeta.Labels
		if err := k8sClient.Update(context.Background(), crRateLimitPolicies); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("RateLimitPolicies CR updated: " + rateLimitPolicies.Name)
		}
	}
}

// DeployAIRateLimitPolicyCR applies the given AIRateLimitPolicies struct to the Kubernetes cluster.
func DeployAIRateLimitPolicyCR(aiRateLimitPolicies *dpv1alpha3.AIRateLimitPolicy, k8sClient client.Client) {
	crAIRateLimitPolicies := &dpv1alpha3.AIRateLimitPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: aiRateLimitPolicies.ObjectMeta.Namespace, Name: aiRateLimitPolicies.Name}, crAIRateLimitPolicies); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get RateLimitPolicies CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), aiRateLimitPolicies); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("RateLimitPolicies CR created: " + aiRateLimitPolicies.Name)
		}
	} else {
		crAIRateLimitPolicies.Spec = aiRateLimitPolicies.Spec
		crAIRateLimitPolicies.ObjectMeta.Labels = aiRateLimitPolicies.ObjectMeta.Labels
		if err := k8sClient.Update(context.Background(), crAIRateLimitPolicies); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("RateLimitPolicies CR updated: " + aiRateLimitPolicies.Name)
		}
	}
}

// UpdateRateLimitPolicyCR applies the updated policy details to all the RateLimitPolicies struct which has the provided label to the Kubernetes cluster.
func UpdateRateLimitPolicyCR(policy eventhubTypes.RateLimitPolicy, k8sClient client.Client) {
	conf, _ := config.ReadConfigs()
	policyName := getSha1Value(policy.Name)
	policyOrganization := getSha1Value(policy.TenantDomain)

	// retrieve all RateLimitPolicies from the Kubernetes cluster with the provided label selector "rateLimitPolicyName"
	rateLimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	labelMap := map[string]string{"rateLimitPolicyName": policyName, "organization": policyOrganization}
	// Create a list option with the label selector
	listOption := &client.ListOptions{
		Namespace:     conf.DataPlane.Namespace,
		LabelSelector: labels.SelectorFromSet(labelMap),
	}
	err := k8sClient.List(context.Background(), rateLimitPolicyList, listOption)
	if err != nil {
		loggers.LoggerK8sClient.Errorf("Unable to list RateLimitPolicies CR: %v", err)
	}
	loggers.LoggerK8sClient.Infof("RateLimitPolicies CR list retrieved: %v", rateLimitPolicyList.Items)
	for _, rateLimitPolicy := range rateLimitPolicyList.Items {
		rateLimitPolicy.Spec.Default.API.RequestsPerUnit = uint32(policy.DefaultLimit.RequestCount.RequestCount)
		rateLimitPolicy.Spec.Default.API.Unit = policy.DefaultLimit.RequestCount.TimeUnit
		loggers.LoggerK8sClient.Infof("RateLimitPolicy CR updated: %v", rateLimitPolicy)
		if err := k8sClient.Update(context.Background(), &rateLimitPolicy); err != nil {
			loggers.LoggerK8sClient.Errorf("Unable to update RateLimitPolicies CR: %v", err)
		} else {
			loggers.LoggerK8sClient.Infof("RateLimitPolicies CR updated: %v", rateLimitPolicy.Name)
		}
	}
}

// DeploySubscriptionRateLimitPolicyCR applies the given RateLimitPolicies struct to the Kubernetes cluster.
func DeploySubscriptionRateLimitPolicyCR(policy eventhubTypes.SubscriptionPolicy, k8sClient client.Client) {
	conf, _ := config.ReadConfigs()
	crRateLimitPolicy := dpv1alpha3.RateLimitPolicy{}
	crName := PrepareSubscritionPolicyCRName(policy.Name, policy.TenantDomain)
	labelMap := map[string]string{
		"InitiateFrom": "CP",
		"CPName":       policy.Name,
	}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: conf.DataPlane.Namespace, Name: crName}, &crRateLimitPolicy); err != nil {
		crRateLimitPolicy = dpv1alpha3.RateLimitPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName,
				Namespace: conf.DataPlane.Namespace,
				Labels:    labelMap,
			},
			Spec: dpv1alpha3.RateLimitPolicySpec{
				Override: &dpv1alpha3.RateLimitAPIPolicy{
					Subscription: &dpv1alpha3.SubscriptionRateLimitPolicy{
						StopOnQuotaReach: policy.StopOnQuotaReach,
						Organization:     policy.TenantDomain,
						RequestCount: &dpv1alpha3.RequestCount{
							RequestsPerUnit: uint32(policy.DefaultLimit.RequestCount.RequestCount),
							Unit:            policy.DefaultLimit.RequestCount.TimeUnit,
						},
					},
				},
				TargetRef: gwapiv1b1.NamespacedPolicyTargetReference{Group: constants.GatewayGroup, Kind: "Subscription", Name: "default"},
			},
		}
		if err := k8sClient.Create(context.Background(), &crRateLimitPolicy); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("RateLimitPolicies CR created: " + crRateLimitPolicy.Name)
		}
	} else {
		crRateLimitPolicy.Spec.Override.Subscription.StopOnQuotaReach = policy.StopOnQuotaReach
		crRateLimitPolicy.Spec.Override.Subscription.Organization = policy.TenantDomain
		crRateLimitPolicy.Spec.Override.Subscription.RequestCount.RequestsPerUnit = uint32(policy.DefaultLimit.RequestCount.RequestCount)
		crRateLimitPolicy.Spec.Override.Subscription.RequestCount.Unit = policy.DefaultLimit.RequestCount.TimeUnit
		if err := k8sClient.Update(context.Background(), &crRateLimitPolicy); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update RateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("RateLimitPolicies CR updated: " + crRateLimitPolicy.Name)
		}
	}

}

// DeployAIRateLimitPolicyFromCPPolicy applies the given AIRateLimitPolicies struct to the Kubernetes cluster.
func DeployAIRateLimitPolicyFromCPPolicy(policy eventhubTypes.SubscriptionPolicy, k8sClient client.Client) {
	conf, _ := config.ReadConfigs()
	tokenCount := &dpv1alpha3.TokenCount{}
	requestCount := &dpv1alpha3.RequestCount{}
	if policy.DefaultLimit.AiAPIQuota.PromptTokenCount != nil &&
		policy.DefaultLimit.AiAPIQuota.CompletionTokenCount != nil &&
		policy.DefaultLimit.AiAPIQuota.TotalTokenCount != nil {
		tokenCount = &dpv1alpha3.TokenCount{
			Unit:               policy.DefaultLimit.AiAPIQuota.TimeUnit,
			RequestTokenCount:  uint32(*policy.DefaultLimit.AiAPIQuota.PromptTokenCount),
			ResponseTokenCount: uint32(*policy.DefaultLimit.AiAPIQuota.CompletionTokenCount),
			TotalTokenCount:    uint32(*policy.DefaultLimit.AiAPIQuota.TotalTokenCount),
		}
	} else {
		tokenCount = nil
	}
	if policy.DefaultLimit.AiAPIQuota.RequestCount != nil {
		requestCount = &dpv1alpha3.RequestCount{
			RequestsPerUnit: uint32(*policy.DefaultLimit.AiAPIQuota.RequestCount),
			Unit:            policy.DefaultLimit.AiAPIQuota.TimeUnit,
		}
	} else {
		requestCount = nil
	}
	labelMap := map[string]string{
		"InitiateFrom": "CP",
		"CPName":       policy.Name,
	}

	crRateLimitPolicies := dpv1alpha3.AIRateLimitPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrepareSubscritionPolicyCRName(policy.Name, policy.TenantDomain),
			Namespace: conf.DataPlane.Namespace,
			Labels:    labelMap,
		},
		Spec: dpv1alpha3.AIRateLimitPolicySpec{
			Override: &dpv1alpha3.AIRateLimit{
				Organization: policy.TenantDomain,
				TokenCount:   tokenCount,
				RequestCount: requestCount,
			},
			TargetRef: gwapiv1b1.NamespacedPolicyTargetReference{Group: constants.GatewayGroup, Kind: "Subscription", Name: "default"},
		},
	}
	crRateLimitPolicyFetched := &dpv1alpha3.AIRateLimitPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: crRateLimitPolicies.ObjectMeta.Namespace, Name: crRateLimitPolicies.Name}, crRateLimitPolicyFetched); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get AiratelimitPolicy CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), &crRateLimitPolicies); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create AIRateLimitPolicies CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("AIRateLimitPolicies CR created: " + crRateLimitPolicies.Name)
		}
	} else {
		crRateLimitPolicyFetched.Spec = crRateLimitPolicies.Spec
		crRateLimitPolicyFetched.ObjectMeta.Labels = crRateLimitPolicies.ObjectMeta.Labels
		if err := k8sClient.Update(context.Background(), crRateLimitPolicyFetched); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update AiRatelimitPolicy CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("AiRatelimitPolicy CR updated: " + crRateLimitPolicyFetched.Name)
		}
	}
}

// UnDeploySubscriptionRateLimitPolicyCR applies the given RateLimitPolicies struct to the Kubernetes cluster.
func UnDeploySubscriptionRateLimitPolicyCR(crName string, k8sClient client.Client) {
	conf, _ := config.ReadConfigs()
	crRateLimitPolicies := &dpv1alpha1.RateLimitPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: conf.DataPlane.Namespace, Name: crName}, crRateLimitPolicies); err != nil {
		loggers.LoggerK8sClient.Error("Unable to get RateLimitPolicies CR: " + err.Error())
	}
	err := k8sClient.Delete(context.Background(), crRateLimitPolicies, &client.DeleteOptions{})
	if err != nil {
		loggers.LoggerK8sClient.Error("Unable to delete RateLimitPolicies CR: " + err.Error())
	}
	loggers.LoggerK8sClient.Debug("RateLimitPolicies CR deleted: " + crRateLimitPolicies.Name)
}

// UndeploySubscriptionAIRateLimitPolicyCR applies the given AIRateLimitPolicies struct to the Kubernetes cluster.
func UndeploySubscriptionAIRateLimitPolicyCR(crName string, k8sClient client.Client) {
	conf, _ := config.ReadConfigs()
	crAIRateLimitPolicies := &dpv1alpha3.AIRateLimitPolicy{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: conf.DataPlane.Namespace, Name: crName}, crAIRateLimitPolicies); err != nil {
		loggers.LoggerK8sClient.Error("Unable to get AIRateLimitPolicies CR: " + err.Error())
	}
	err := k8sClient.Delete(context.Background(), crAIRateLimitPolicies, &client.DeleteOptions{})
	if err != nil {
		loggers.LoggerK8sClient.Error("Unable to delete AIRateLimitPolicies CR: " + err.Error())
	}
	loggers.LoggerK8sClient.Debug("AIRateLimitPolicies CR deleted: " + crAIRateLimitPolicies.Name)
}

// DeployBackendCR applies the given Backends struct to the Kubernetes cluster.
func DeployBackendCR(backends *dpv1alpha2.Backend, k8sClient client.Client) {
	crBackends := &dpv1alpha2.Backend{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: backends.ObjectMeta.Namespace, Name: backends.Name}, crBackends); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get Backends CR: " + err.Error())
		}
		if err := k8sClient.Create(context.Background(), backends); err != nil {
			loggers.LoggerK8sClient.Error("Unable to create Backends CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Backends CR created: " + backends.Name)
		}
	} else {
		crBackends.Spec = backends.Spec
		if err := k8sClient.Update(context.Background(), crBackends); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update Backends CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("Backends CR updated: " + backends.Name)
		}
	}
}

// CreateAndUpdateTokenIssuersCR applies the given TokenIssuers struct to the Kubernetes cluster.
func CreateAndUpdateTokenIssuersCR(keyManager eventhubTypes.ResolvedKeyManager, k8sClient client.Client) error {
	conf, _ := config.ReadConfigs()
	sha1ValueofKmName := getSha1Value(keyManager.Name)
	sha1ValueOfOrganization := getSha1Value(keyManager.Organization)
	labelMap := map[string]string{"name": sha1ValueofKmName,
		"organization": sha1ValueOfOrganization,
		"InitiateFrom": "CP",
	}

	tokenIssuer := dpv1alpha2.TokenIssuer{
		ObjectMeta: metav1.ObjectMeta{Name: keyManager.UUID,
			Namespace: conf.DataPlane.Namespace,
			Labels:    labelMap,
		},
		Spec: dpv1alpha2.TokenIssuerSpec{
			Name:          keyManager.Name,
			Organization:  keyManager.Organization,
			Issuer:        keyManager.KeyManagerConfig.Issuer,
			ClaimMappings: marshalClaimMappings(keyManager.KeyManagerConfig.ClaimMappings),
			TargetRef:     &v1alpha2.NamespacedPolicyTargetReference{Group: constants.GatewayGroup, Kind: constants.GatewayKind, Name: constants.GatewayName},
		},
	}
	signatureValidation, err := marshalSignatureValidation(keyManager.KeyManagerConfig)
	if err != nil {
		loggers.LoggerK8sClient.Errorf("Failed to marshal signature validation: %v", err)
		return err
	}
	tokenIssuer.Spec.SignatureValidation = signatureValidation
	tokenIssuer.Spec.ConsumerKeyClaim = constants.ConsumerKeyClaim
	if keyManager.KeyManagerConfig.ConsumerKeyClaim != "" {
		tokenIssuer.Spec.ConsumerKeyClaim = keyManager.KeyManagerConfig.ConsumerKeyClaim
	}
	keyManager.KeyManagerConfig.ScopesClaim = constants.ScopesClaim
	if keyManager.KeyManagerConfig.ScopesClaim != "" {
		tokenIssuer.Spec.ScopesClaim = keyManager.KeyManagerConfig.ScopesClaim
	}
	crTokenIssuer := &dpv1alpha2.TokenIssuer{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: tokenIssuer.ObjectMeta.Namespace, Name: tokenIssuer.Name}, crTokenIssuer); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get TokenIssuer CR: " + err.Error())
		}
		err := k8sClient.Create(context.Background(), &tokenIssuer)
		if err != nil {
			loggers.LoggerK8sClient.Error("Unable to create TokenIssuer CR: " + err.Error())
			return err
		}
		loggers.LoggerK8sClient.Infof("TokenIssuer CR created: " + tokenIssuer.Name)
	} else {
		crTokenIssuer.Spec = tokenIssuer.Spec
		if err := k8sClient.Update(context.Background(), crTokenIssuer); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update TokenIssuer CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("TokenIssuer CR updated: " + tokenIssuer.Name)
		}
	}

	internalKeyTokenIssuer := dpv1alpha2.TokenIssuer{
		ObjectMeta: metav1.ObjectMeta{Name: keyManager.Organization + constants.InternalKeySuffix,
			Namespace: conf.DataPlane.Namespace,
			Labels:    labelMap,
		},
		Spec: dpv1alpha2.TokenIssuerSpec{
			Name:          constants.InternalKeyTokenIssuerName,
			Organization:  keyManager.Organization,
			Issuer:        conf.ControlPlane.InternalKeyIssuer,
			ClaimMappings: marshalClaimMappings(keyManager.KeyManagerConfig.ClaimMappings),
			SignatureValidation: &dpv1alpha2.SignatureValidation{
				Certificate: &dpv1alpha2.CERTConfig{
					SecretRef: &dpv1alpha2.RefConfig{
						Name: constants.InternalKeySecretName,
						Key:  constants.InternalKeySecretKey,
					},
				},
			},
			TargetRef: &v1alpha2.NamespacedPolicyTargetReference{Group: constants.GatewayGroup, Kind: constants.GatewayKind, Name: constants.GatewayName},
		},
	}
	internalKeyTokenIssuer.Spec.ConsumerKeyClaim = constants.ConsumerKeyClaim
	internalKeyTokenIssuer.Spec.ScopesClaim = constants.ScopesClaim
	crInternalTokenIssuer := &dpv1alpha2.TokenIssuer{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: internalKeyTokenIssuer.ObjectMeta.Namespace, Name: internalKeyTokenIssuer.Name}, crInternalTokenIssuer); err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerK8sClient.Error("Unable to get Internal TokenIssuer CR: " + err.Error())
		}
		err = k8sClient.Create(context.Background(), &internalKeyTokenIssuer)
		if err != nil {
			loggers.LoggerK8sClient.Error("Unable to create Internal TokenIssuer CR: " + err.Error())
			return err
		}
		loggers.LoggerK8sClient.Infof("Internal TokenIssuer CR created: " + internalKeyTokenIssuer.Name)
	} else {
		crInternalTokenIssuer.Spec = internalKeyTokenIssuer.Spec
		if err := k8sClient.Update(context.Background(), crInternalTokenIssuer); err != nil {
			loggers.LoggerK8sClient.Error("Unable to update Internal TokenIssuer CR: " + err.Error())
		} else {
			loggers.LoggerK8sClient.Info("TokenIssuer CR updated: " + internalKeyTokenIssuer.Name)
		}
	}
	return nil
}

// DeleteTokenIssuerCR deletes the TokenIssuer struct from the Kubernetes cluster.
func DeleteTokenIssuerCR(k8sClient client.Client, tokenIssuer dpv1alpha2.TokenIssuer) error {
	// Skip the deletion if the token issuer is for internal keys
	if !strings.Contains(tokenIssuer.Name, constants.InternalKeySuffix) {
		err := k8sClient.Delete(context.Background(), &tokenIssuer, &client.DeleteOptions{})
		if err != nil {
			loggers.LoggerK8sClient.Error("Unable to delete TokenIssuer CR: " + err.Error())
			return err
		}
		loggers.LoggerK8sClient.Debug("TokenIssuer CR deleted: " + tokenIssuer.Name)
	}
	return nil
}

// DeleteTokenIssuersCR deletes the TokenIssuers struct from the Kubernetes cluster.
func DeleteTokenIssuersCR(k8sClient client.Client, keymanagerName string, tenantDomain string) error {
	conf, _ := config.ReadConfigs()
	sha1ValueofKmName := getSha1Value(keymanagerName)
	sha1ValueOfOrganization := getSha1Value(tenantDomain)
	labelMap := map[string]string{"name": sha1ValueofKmName, "organization": sha1ValueOfOrganization}
	// Create a list option with the label selector
	listOption := &client.ListOptions{
		Namespace:     conf.DataPlane.Namespace,
		LabelSelector: labels.SelectorFromSet(labelMap),
	}

	tokenIssuerList := &dpv1alpha2.TokenIssuerList{}
	err := k8sClient.List(context.Background(), tokenIssuerList, listOption)
	if err != nil {
		loggers.LoggerK8sClient.Error("Unable to list TokenIssuer CR: " + err.Error())
	}
	if len(tokenIssuerList.Items) == 0 {
		loggers.LoggerK8sClient.Debug("No TokenIssuer CR found for deletion")
	}
	for _, tokenIssuer := range tokenIssuerList.Items {
		err := DeleteTokenIssuerCR(k8sClient, tokenIssuer)
		if err != nil {
			loggers.LoggerK8sClient.Error("Unable to delete TokenIssuer CR: " + err.Error())
			return err
		}
		loggers.LoggerK8sClient.Debug("TokenIssuer CR deleted: " + tokenIssuer.Name)
	}
	return nil
}

// UpdateTokenIssuersCR applies the given TokenIssuers struct to the Kubernetes cluster.
func UpdateTokenIssuersCR(keyManager eventhubTypes.ResolvedKeyManager, k8sClient client.Client) error {
	conf, _ := config.ReadConfigs()
	sha1ValueofKmName := getSha1Value(keyManager.Name)
	sha1ValueOfOrganization := getSha1Value(keyManager.Organization)
	labelMap := map[string]string{"name": sha1ValueofKmName, "organization": sha1ValueOfOrganization}
	tokenIssuer := &dpv1alpha2.TokenIssuer{}
	err := k8sClient.Get(context.Background(), client.ObjectKey{Name: keyManager.UUID, Namespace: conf.DataPlane.Namespace}, tokenIssuer)
	if err != nil {
		loggers.LoggerK8sClient.Error("Unable to get TokenIssuer CR: " + err.Error())
		return err
	}
	tokenIssuer.ObjectMeta.Labels = labelMap
	tokenIssuer.Spec.Name = keyManager.Name
	tokenIssuer.Spec.Organization = keyManager.Organization
	tokenIssuer.Spec.Issuer = keyManager.KeyManagerConfig.Issuer
	tokenIssuer.Spec.ClaimMappings = marshalClaimMappings(keyManager.KeyManagerConfig.ClaimMappings)
	signatureValidation, err := marshalSignatureValidation(keyManager.KeyManagerConfig)
	if err != nil {
		loggers.LoggerK8sClient.Errorf("Failed to marshal signature validation: %v", err)
		return err
	}
	tokenIssuer.Spec.SignatureValidation = signatureValidation
	tokenIssuer.Spec.TargetRef = &v1alpha2.NamespacedPolicyTargetReference{Group: constants.GatewayGroup, Kind: constants.GatewayKind, Name: constants.GatewayName}
	if keyManager.KeyManagerConfig.ConsumerKeyClaim != "" {
		tokenIssuer.Spec.ConsumerKeyClaim = keyManager.KeyManagerConfig.ConsumerKeyClaim
	}
	if keyManager.KeyManagerConfig.ScopesClaim != "" {
		tokenIssuer.Spec.ScopesClaim = keyManager.KeyManagerConfig.ScopesClaim
	}
	err = k8sClient.Update(context.Background(), tokenIssuer)
	if err != nil {
		loggers.LoggerK8sClient.Error("Unable to update TokenIssuer CR: " + err.Error())
		return err
	}
	loggers.LoggerK8sClient.Debug("TokenIssuer CR updated: " + tokenIssuer.Name)
	return nil
}

func marshalSignatureValidation(keyManagerConfig eventhubTypes.KeyManagerConfig) (*dpv1alpha2.SignatureValidation, error) {
	if keyManagerConfig.CertificateType != "" && keyManagerConfig.CertificateValue != "" {
		if keyManagerConfig.CertificateType == "JWKS" {
			loggers.LoggerK8sClient.Debugf("Using JWKS for signature validation")
			return &dpv1alpha2.SignatureValidation{JWKS: &dpv1alpha2.JWKS{URL: keyManagerConfig.CertificateValue}}, nil
		}
		loggers.LoggerK8sClient.Debugf("Using Certificate for signature validation %s", keyManagerConfig.CertificateValue)
		certValue := keyManagerConfig.CertificateValue
		// Check if the certificate value is base64 encoded and decode it
		if decodedCert, err := base64.StdEncoding.DecodeString(certValue); err == nil {
			// Successfully decoded, use the decoded value
			decodedCertStr := string(decodedCert)
			loggers.LoggerK8sClient.Debugf("Certificate value was base64 encoded, using decoded value")
			certValue = decodedCertStr
		}
		// Validate that the certificate is in proper PEM format
		block, _ := pem.Decode([]byte(certValue))
		if block == nil {
			loggers.LoggerK8sClient.Errorf("Failed to decode PEM block from certificate")
			return nil, fmt.Errorf("certificate is not in valid PEM format")
		}
		// Validate that it's a valid X.509 certificate
		if _, err := x509.ParseCertificate(block.Bytes); err != nil {
			loggers.LoggerK8sClient.Errorf("Failed to parse X.509 certificate: %v", err)
			return nil, fmt.Errorf("invalid X.509 certificate: %w", err)
		}
		loggers.LoggerK8sClient.Debugf("Certificate validated successfully")
		return &dpv1alpha2.SignatureValidation{Certificate: &dpv1alpha2.CERTConfig{CertificateInline: &certValue}}, nil
	}
	return nil, nil
}

func marshalClaimMappings(claimMappings []eventhubTypes.Claim) *[]dpv1alpha2.ClaimMapping {
	resolvedClaimMappings := make([]dpv1alpha2.ClaimMapping, 0)
	for _, claim := range claimMappings {
		resolvedClaimMappings = append(resolvedClaimMappings, dpv1alpha2.ClaimMapping{RemoteClaim: claim.RemoteClaim, LocalClaim: claim.LocalClaim})
	}
	return &resolvedClaimMappings
}
func getSha1Value(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

// RetrieveAllAPISFromK8s retrieves all the API CRs from the Kubernetes cluster
func RetrieveAllAPISFromK8s(k8sClient client.Client, nextToken string) ([]dpv1alpha3.API, string, error) {
	conf, _ := config.ReadConfigs()
	apiList := dpv1alpha3.APIList{}
	resolvedAPIList := make([]dpv1alpha3.API, 0)
	var err error
	if nextToken == "" {
		err = k8sClient.List(context.Background(), &apiList, &client.ListOptions{Namespace: conf.DataPlane.Namespace})
	} else {
		err = k8sClient.List(context.Background(), &apiList, &client.ListOptions{Namespace: conf.DataPlane.Namespace, Continue: nextToken})
	}
	if err != nil {
		loggers.LoggerK8sClient.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedAPIList = append(resolvedAPIList, apiList.Items...)
	if apiList.Continue != "" {
		tempAPIList, _, err := RetrieveAllAPISFromK8s(k8sClient, apiList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedAPIList = append(resolvedAPIList, tempAPIList...)
	}
	return resolvedAPIList, apiList.Continue, nil
}

// RetrieveAllAIProvidersFromK8s retrieves all the API CRs from the Kubernetes cluster
func RetrieveAllAIProvidersFromK8s(k8sClient client.Client, nextToken string) ([]dpv1alpha4.AIProvider, string, error) {
	conf, _ := config.ReadConfigs()
	aiProviderList := dpv1alpha4.AIProviderList{}
	resolvedAIProviderList := make([]dpv1alpha4.AIProvider, 0)
	var err error
	if nextToken == "" {
		err = k8sClient.List(context.Background(), &aiProviderList, &client.ListOptions{Namespace: conf.DataPlane.Namespace})
	} else {
		err = k8sClient.List(context.Background(), &aiProviderList, &client.ListOptions{Namespace: conf.DataPlane.Namespace, Continue: nextToken})
	}
	if err != nil {
		loggers.LoggerK8sClient.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get ai provider from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedAIProviderList = append(resolvedAIProviderList, aiProviderList.Items...)
	if aiProviderList.Continue != "" {
		tempAIProviderList, _, err := RetrieveAllAIProvidersFromK8s(k8sClient, aiProviderList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedAIProviderList = append(resolvedAIProviderList, tempAIProviderList...)
	}
	return resolvedAIProviderList, aiProviderList.Continue, nil
}

// RetrieveAllRatelimitPoliciesSFromK8s retrieves all the API CRs from the Kubernetes cluster
func RetrieveAllRatelimitPoliciesSFromK8s(k8sClient client.Client, nextToken string) ([]dpv1alpha3.RateLimitPolicy, string, error) {
	conf, _ := config.ReadConfigs()
	rlList := dpv1alpha3.RateLimitPolicyList{}
	resolvedRLList := make([]dpv1alpha3.RateLimitPolicy, 0)
	var err error
	if nextToken == "" {
		err = k8sClient.List(context.Background(), &rlList, &client.ListOptions{Namespace: conf.DataPlane.Namespace})
	} else {
		err = k8sClient.List(context.Background(), &rlList, &client.ListOptions{Namespace: conf.DataPlane.Namespace, Continue: nextToken})
	}
	if err != nil {
		loggers.LoggerK8sClient.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get ratelimitpolicies from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedRLList = append(resolvedRLList, rlList.Items...)
	if rlList.Continue != "" {
		tempRLList, _, err := RetrieveAllRatelimitPoliciesSFromK8s(k8sClient, rlList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedRLList = append(resolvedRLList, tempRLList...)
	}
	return resolvedRLList, rlList.Continue, nil
}

// RetrieveAllAIRatelimitPoliciesSFromK8s retrieves all the API CRs from the Kubernetes cluster
func RetrieveAllAIRatelimitPoliciesSFromK8s(k8sClient client.Client, nextToken string) ([]dpv1alpha3.AIRateLimitPolicy, string, error) {
	conf, _ := config.ReadConfigs()
	airlList := dpv1alpha3.AIRateLimitPolicyList{}
	resolvedAIRLList := make([]dpv1alpha3.AIRateLimitPolicy, 0)
	var err error
	if nextToken == "" {
		err = k8sClient.List(context.Background(), &airlList, &client.ListOptions{Namespace: conf.DataPlane.Namespace})
	} else {
		err = k8sClient.List(context.Background(), &airlList, &client.ListOptions{Namespace: conf.DataPlane.Namespace, Continue: nextToken})
	}
	if err != nil {
		loggers.LoggerK8sClient.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get airatelimitpolicies from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedAIRLList = append(resolvedAIRLList, airlList.Items...)
	if airlList.Continue != "" {
		tempAIRLList, _, err := RetrieveAllAIRatelimitPoliciesSFromK8s(k8sClient, airlList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedAIRLList = append(resolvedAIRLList, tempAIRLList...)
	}
	return resolvedAIRLList, airlList.Continue, nil
}

// PrepareSubscritionPolicyCRName prepare the cr name for a given policy name and organization pair
func PrepareSubscritionPolicyCRName(name, org string) string {
	return getSha1Value(fmt.Sprintf("%s-%s", name, org))
}
