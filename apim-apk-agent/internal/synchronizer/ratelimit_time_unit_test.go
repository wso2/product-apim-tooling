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

package synchronizer

import (
	"testing"

	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
)

func TestResolveCRTimeUnit(t *testing.T) {
	tests := []struct {
		cpTimeUnit string
		expected   string
		supported  bool
	}{
		// Units sent by the control plane internal data endpoints.
		{"sec", "Second", true},
		{"min", "Minute", true},
		{"hour", "Hour", true},
		{"day", "Day", true},
		// Plural and long forms expected by earlier agent versions.
		{"hours", "Hour", true},
		{"days", "Day", true},
		{"minute", "Minute", true},
		{"minutes", "Minute", true},
		{"second", "Second", true},
		{"seconds", "Second", true},
		// Case and whitespace insensitive.
		{"Day", "Day", true},
		{" MIN ", "Minute", true},
		// Missing or not supported by the data plane.
		{"", "", false},
		{"week", "", false},
		{"month", "", false},
		{"year", "", false},
	}
	for _, test := range tests {
		actual, supported := resolveCRTimeUnit(test.cpTimeUnit)
		if actual != test.expected || supported != test.supported {
			t.Errorf("resolveCRTimeUnit(%q) = (%q, %v), expected (%q, %v)",
				test.cpTimeUnit, actual, supported, test.expected, test.supported)
		}
	}
}

func TestResolveCRTimeUnitWithoutSeconds(t *testing.T) {
	tests := []struct {
		cpTimeUnit string
		expected   string
		supported  bool
	}{
		{"min", "Minute", true},
		{"hour", "Hour", true},
		{"hours", "Hour", true},
		{"day", "Day", true},
		{"days", "Day", true},
		// API level and AI rate limit CRs only accept Minute, Hour and Day.
		{"sec", "", false},
		{"seconds", "", false},
		{"", "", false},
		{"month", "", false},
	}
	for _, test := range tests {
		actual, supported := resolveCRTimeUnitWithoutSeconds(test.cpTimeUnit)
		if actual != test.expected || supported != test.supported {
			t.Errorf("resolveCRTimeUnitWithoutSeconds(%q) = (%q, %v), expected (%q, %v)",
				test.cpTimeUnit, actual, supported, test.expected, test.supported)
		}
	}
}

func newRequestCountSubscriptionPolicy(name string, requestCount int, timeUnit string,
	burstCount int32, burstTimeUnit string) eventhubTypes.SubscriptionPolicy {
	policy := eventhubTypes.SubscriptionPolicy{
		Name:              name,
		TenantDomain:      "carbon.super",
		QuotaType:         "requestCount",
		RateLimitCount:    burstCount,
		RateLimitTimeUnit: burstTimeUnit,
		StopOnQuotaReach:  true,
	}
	policy.DefaultLimit.QuotaType = "requestCount"
	policy.DefaultLimit.RequestCount.RequestCount = requestCount
	policy.DefaultLimit.RequestCount.TimeUnit = timeUnit
	policy.DefaultLimit.RequestCount.UnitTime = 1
	return policy
}

func TestNormalizeSubscriptionRateLimitPolicy(t *testing.T) {
	tests := []struct {
		name                string
		policy              eventhubTypes.SubscriptionPolicy
		expectDeploy        bool
		expectedRequestUnit string
		expectedBurstCount  int32
		expectedBurstUnit   string
	}{
		{
			name:                "day quota with minute burst",
			policy:              newRequestCountSubscriptionPolicy("BurstDay5000Min5", 5000, "day", 5, "min"),
			expectDeploy:        true,
			expectedRequestUnit: "Day",
			expectedBurstCount:  5,
			expectedBurstUnit:   "Minute",
		},
		{
			name:                "hour quota with second burst",
			policy:              newRequestCountSubscriptionPolicy("BurstHour100Sec2", 100, "hour", 2, "sec"),
			expectDeploy:        true,
			expectedRequestUnit: "Hour",
			expectedBurstCount:  2,
			expectedBurstUnit:   "Second",
		},
		{
			name:                "minute quota with minute burst",
			policy:              newRequestCountSubscriptionPolicy("BurstMin20Min5", 20, "min", 5, "min"),
			expectDeploy:        true,
			expectedRequestUnit: "Minute",
			expectedBurstCount:  5,
			expectedBurstUnit:   "Minute",
		},
		{
			// The control plane sends a null burst time unit and a zero burst count when burst control is not configured.
			name:                "default tier without burst control",
			policy:              newRequestCountSubscriptionPolicy("Gold", 5000, "min", 0, ""),
			expectDeploy:        true,
			expectedRequestUnit: "Minute",
		},
		{
			name:                "day quota without burst control",
			policy:              newRequestCountSubscriptionPolicy("NoBurstDay1000", 1000, "day", 0, ""),
			expectDeploy:        true,
			expectedRequestUnit: "Day",
		},
		{
			name:                "plural quota unit",
			policy:              newRequestCountSubscriptionPolicy("PluralDays", 1000, "days", 10, "min"),
			expectDeploy:        true,
			expectedRequestUnit: "Day",
			expectedBurstCount:  10,
			expectedBurstUnit:   "Minute",
		},
		{
			name:                "burst time unit without burst count",
			policy:              newRequestCountSubscriptionPolicy("ZeroBurstCount", 100, "min", 0, "min"),
			expectDeploy:        true,
			expectedRequestUnit: "Minute",
		},
		{
			name:                "unsupported burst time unit is dropped",
			policy:              newRequestCountSubscriptionPolicy("WeekBurst", 100, "day", 5, "week"),
			expectDeploy:        true,
			expectedRequestUnit: "Day",
		},
		{
			name:         "unsupported quota time unit",
			policy:       newRequestCountSubscriptionPolicy("MonthQuota", 100, "month", 5, "min"),
			expectDeploy: false,
		},
		{
			// Event count based (async) policies do not have a request count time unit.
			name:         "empty quota time unit",
			policy:       newRequestCountSubscriptionPolicy("AsyncGold", 0, "", 0, ""),
			expectDeploy: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			policy := test.policy
			deploy := normalizeSubscriptionRateLimitPolicy(&policy)
			if deploy != test.expectDeploy {
				t.Fatalf("expected deploy=%v, got %v", test.expectDeploy, deploy)
			}
			if !deploy {
				return
			}
			if policy.DefaultLimit.RequestCount.TimeUnit != test.expectedRequestUnit {
				t.Errorf("expected request count unit %q, got %q", test.expectedRequestUnit, policy.DefaultLimit.RequestCount.TimeUnit)
			}
			if policy.DefaultLimit.RequestCount.RequestCount != test.policy.DefaultLimit.RequestCount.RequestCount {
				t.Errorf("request count changed from %d to %d", test.policy.DefaultLimit.RequestCount.RequestCount,
					policy.DefaultLimit.RequestCount.RequestCount)
			}
			if policy.RateLimitCount != test.expectedBurstCount || policy.RateLimitTimeUnit != test.expectedBurstUnit {
				t.Errorf("expected burst %d/%q, got %d/%q", test.expectedBurstCount, test.expectedBurstUnit,
					policy.RateLimitCount, policy.RateLimitTimeUnit)
			}
		})
	}
}
