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
	"strings"

	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
)

const (
	crTimeUnitSecond = "Second"
	crTimeUnitMinute = "Minute"
	crTimeUnitHour   = "Hour"
	crTimeUnitDay    = "Day"
)

// resolveCRTimeUnit maps a time unit received from the control plane to the time unit used in the
// data plane rate limit CRs. The control plane internal data endpoints send "sec", "min", "hour" and
// "day", while older code paths expected "hours" and "days", so both singular and plural forms are
// accepted. The second return value is false when the time unit is empty or not supported by the
// data plane.
func resolveCRTimeUnit(cpTimeUnit string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(cpTimeUnit)) {
	case "sec", "second", "seconds":
		return crTimeUnitSecond, true
	case "min", "minute", "minutes":
		return crTimeUnitMinute, true
	case "hour", "hours":
		return crTimeUnitHour, true
	case "day", "days":
		return crTimeUnitDay, true
	default:
		return "", false
	}
}

// resolveCRTimeUnitWithoutSeconds is the same as resolveCRTimeUnit, but for CR fields that only accept
// "Minute", "Hour" and "Day", such as API level rate limits and AI rate limits. A time unit of seconds is
// reported as not supported.
func resolveCRTimeUnitWithoutSeconds(cpTimeUnit string) (string, bool) {
	crTimeUnit, ok := resolveCRTimeUnit(cpTimeUnit)
	if !ok || crTimeUnit == crTimeUnitSecond {
		return "", false
	}
	return crTimeUnit, true
}

// normalizeSubscriptionRateLimitPolicy converts the request count and burst control time units of a
// request count based subscription policy to the format used in the data plane CRs.
//
// It returns false when the policy cannot be deployed because its request count time unit is missing
// or not supported. Burst control is optional: when it is not configured in the control plane (the
// burst time unit is null), or when its time unit is not supported, the burst control fields are
// cleared so that the request count quota is still deployed without burst control.
func normalizeSubscriptionRateLimitPolicy(policy *eventhubTypes.SubscriptionPolicy) bool {
	requestCountUnit, ok := resolveCRTimeUnit(policy.DefaultLimit.RequestCount.TimeUnit)
	if !ok {
		if policy.DefaultLimit.RequestCount.TimeUnit == "" {
			logger.LoggerSynchronizer.Debugf("Empty RequestCount timeunit for subscription policy %s, skipping", policy.Name)
		} else {
			logger.LoggerSynchronizer.Errorf("Unsupported RequestCount timeunit '%s' for subscription policy %s, skipping",
				policy.DefaultLimit.RequestCount.TimeUnit, policy.Name)
		}
		return false
	}
	policy.DefaultLimit.RequestCount.TimeUnit = requestCountUnit

	if policy.RateLimitCount <= 0 || policy.RateLimitTimeUnit == "" {
		logger.LoggerSynchronizer.Debugf("Burst control is not configured for subscription policy %s", policy.Name)
		policy.RateLimitCount = 0
		policy.RateLimitTimeUnit = ""
		return true
	}
	burstUnit, ok := resolveCRTimeUnit(policy.RateLimitTimeUnit)
	if !ok {
		logger.LoggerSynchronizer.Errorf("Unsupported RateLimitTimeUnit '%s' for subscription policy %s. "+
			"The policy is deployed without burst control", policy.RateLimitTimeUnit, policy.Name)
		policy.RateLimitCount = 0
		policy.RateLimitTimeUnit = ""
		return true
	}
	policy.RateLimitTimeUnit = burstUnit
	return true
}
