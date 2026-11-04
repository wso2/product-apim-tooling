/*
 *  Copyright (c) 2026, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package k8sclient

import (
	"context"
	"testing"

	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newFakeClient(t *testing.T) client.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := dpv1alpha3.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add dp/v1alpha3 to scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).Build()
}

// newSubscriptionPolicy returns a subscription policy with time units already resolved to the CR format.
func newSubscriptionPolicy(requestCount int, timeUnit string, burstCount int32, burstTimeUnit string) eventhubTypes.SubscriptionPolicy {
	policy := eventhubTypes.SubscriptionPolicy{
		Name:              "BurstPolicy",
		TenantDomain:      "carbon.super",
		QuotaType:         "requestCount",
		RateLimitCount:    burstCount,
		RateLimitTimeUnit: burstTimeUnit,
		StopOnQuotaReach:  true,
	}
	policy.DefaultLimit.RequestCount.RequestCount = requestCount
	policy.DefaultLimit.RequestCount.TimeUnit = timeUnit
	return policy
}

func getSubscriptionRateLimitPolicyCR(t *testing.T, k8sClient client.Client,
	policy eventhubTypes.SubscriptionPolicy) *dpv1alpha3.SubscriptionRateLimitPolicy {
	t.Helper()
	conf, _ := config.ReadConfigs()
	crName := PrepareRateLimitPolicyCRName(policy.Name, policy.TenantDomain)
	rateLimitPolicyCR := &dpv1alpha3.RateLimitPolicy{}
	if err := k8sClient.Get(context.Background(),
		client.ObjectKey{Namespace: conf.DataPlane.Namespace, Name: crName}, rateLimitPolicyCR); err != nil {
		t.Fatalf("RateLimitPolicy CR %s was not found: %v", crName, err)
	}
	if rateLimitPolicyCR.Spec.Override == nil || rateLimitPolicyCR.Spec.Override.Subscription == nil {
		t.Fatalf("RateLimitPolicy CR %s does not have a subscription override", crName)
	}
	if rateLimitPolicyCR.Labels["CPName"] != policy.Name {
		t.Errorf("expected CPName label %q, got %q", policy.Name, rateLimitPolicyCR.Labels["CPName"])
	}
	return rateLimitPolicyCR.Spec.Override.Subscription
}

func assertRequestCount(t *testing.T, subscription *dpv1alpha3.SubscriptionRateLimitPolicy, requestsPerUnit uint32, unit string) {
	t.Helper()
	if subscription.RequestCount == nil {
		t.Fatalf("expected request count %d/%s, got nil", requestsPerUnit, unit)
	}
	if subscription.RequestCount.RequestsPerUnit != requestsPerUnit || subscription.RequestCount.Unit != unit {
		t.Errorf("expected request count %d/%s, got %d/%s", requestsPerUnit, unit,
			subscription.RequestCount.RequestsPerUnit, subscription.RequestCount.Unit)
	}
}

func assertBurstControl(t *testing.T, subscription *dpv1alpha3.SubscriptionRateLimitPolicy, requestsPerUnit uint32, unit string) {
	t.Helper()
	if subscription.BurstControl == nil {
		t.Fatalf("expected burst control %d/%s, got nil", requestsPerUnit, unit)
	}
	if subscription.BurstControl.RequestsPerUnit != requestsPerUnit || subscription.BurstControl.Unit != unit {
		t.Errorf("expected burst control %d/%s, got %d/%s", requestsPerUnit, unit,
			subscription.BurstControl.RequestsPerUnit, subscription.BurstControl.Unit)
	}
}

func assertNoBurstControl(t *testing.T, subscription *dpv1alpha3.SubscriptionRateLimitPolicy) {
	t.Helper()
	if subscription.BurstControl != nil {
		t.Errorf("expected no burst control, got %+v", *subscription.BurstControl)
	}
}

func TestDeploySubscriptionRateLimitPolicyCRCreateWithBurstControl(t *testing.T) {
	k8sClient := newFakeClient(t)
	policy := newSubscriptionPolicy(5000, "Day", 5, "Minute")

	DeploySubscriptionRateLimitPolicyCR(policy, k8sClient)

	subscription := getSubscriptionRateLimitPolicyCR(t, k8sClient, policy)
	assertRequestCount(t, subscription, 5000, "Day")
	assertBurstControl(t, subscription, 5, "Minute")
	if !subscription.StopOnQuotaReach || subscription.Organization != "carbon.super" {
		t.Errorf("unexpected stopOnQuotaReach/organization: %v/%s", subscription.StopOnQuotaReach, subscription.Organization)
	}
}

func TestDeploySubscriptionRateLimitPolicyCRCreateWithoutBurstControl(t *testing.T) {
	k8sClient := newFakeClient(t)
	policy := newSubscriptionPolicy(1000, "Day", 0, "")

	DeploySubscriptionRateLimitPolicyCR(policy, k8sClient)

	subscription := getSubscriptionRateLimitPolicyCR(t, k8sClient, policy)
	assertRequestCount(t, subscription, 1000, "Day")
	assertNoBurstControl(t, subscription)
}

func TestDeploySubscriptionRateLimitPolicyCRUpdateBurstControl(t *testing.T) {
	k8sClient := newFakeClient(t)

	// Burst control added to a policy that did not have it.
	DeploySubscriptionRateLimitPolicyCR(newSubscriptionPolicy(1000, "Day", 0, ""), k8sClient)
	policy := newSubscriptionPolicy(1000, "Day", 10, "Minute")
	DeploySubscriptionRateLimitPolicyCR(policy, k8sClient)
	subscription := getSubscriptionRateLimitPolicyCR(t, k8sClient, policy)
	assertRequestCount(t, subscription, 1000, "Day")
	assertBurstControl(t, subscription, 10, "Minute")

	// Burst control and request count changed.
	policy = newSubscriptionPolicy(200, "Hour", 3, "Second")
	DeploySubscriptionRateLimitPolicyCR(policy, k8sClient)
	subscription = getSubscriptionRateLimitPolicyCR(t, k8sClient, policy)
	assertRequestCount(t, subscription, 200, "Hour")
	assertBurstControl(t, subscription, 3, "Second")

	// Burst control removed from the policy. The CR must not keep an empty or stale burst control.
	policy = newSubscriptionPolicy(200, "Hour", 0, "")
	DeploySubscriptionRateLimitPolicyCR(policy, k8sClient)
	subscription = getSubscriptionRateLimitPolicyCR(t, k8sClient, policy)
	assertRequestCount(t, subscription, 200, "Hour")
	assertNoBurstControl(t, subscription)
}

func TestSubscriptionBurstControl(t *testing.T) {
	if burstControl := subscriptionBurstControl(newSubscriptionPolicy(100, "Minute", 0, "")); burstControl != nil {
		t.Errorf("expected nil burst control when not configured, got %+v", *burstControl)
	}
	if burstControl := subscriptionBurstControl(newSubscriptionPolicy(100, "Minute", 5, "")); burstControl != nil {
		t.Errorf("expected nil burst control without a time unit, got %+v", *burstControl)
	}
	if burstControl := subscriptionBurstControl(newSubscriptionPolicy(100, "Minute", 0, "Minute")); burstControl != nil {
		t.Errorf("expected nil burst control without a count, got %+v", *burstControl)
	}
	burstControl := subscriptionBurstControl(newSubscriptionPolicy(100, "Minute", 5, "Second"))
	if burstControl == nil || burstControl.RequestsPerUnit != 5 || burstControl.Unit != "Second" {
		t.Errorf("expected burst control 5/Second, got %+v", burstControl)
	}
}
