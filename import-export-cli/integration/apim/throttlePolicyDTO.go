/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package apim

type ApplicationThrottlePolicy struct {
	PolicyName   string `json:"policyName"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	IsDeployed   bool   `json:"isDeployed"`
	Type         string `json:"type"`
	DefaultLimit struct {
		Type         string `json:"type"`
		RequestCount struct {
			TimeUnit     string `json:"timeUnit"`
			UnitTime     int    `json:"unitTime"`
			RequestCount int    `json:"requestCount"`
		} `json:"requestCount"`
		Bandwidth struct {
			TimeUnit   string `json:"timeUnit"`
			UnitTime   int    `json:"unitTime"`
			DataAmount int    `json:"dataAmount"`
			DataUnit   string `json:"dataUnit"`
		} `json:"bandwidth"`
		EventCount struct {
			TimeUnit   string `json:"timeUnit"`
			UnitTime   int    `json:"unitTime"`
			EventCount int    `json:"eventCount"`
		} `json:"eventCount"`
	} `json:"defaultLimit"`
}

type CustomThrottlePolicy struct {
	PolicyName  string `json:"policyName"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IsDeployed  bool   `json:"isDeployed"`
	Type        string `json:"type"`
	SiddhiQuery string `json:"siddhiQuery"`
	KeyTemplate string `json:"keyTemplate"`
}

type AdvancedThrottlePolicy struct {
	PolicyName   string `json:"policyName"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	IsDeployed   bool   `json:"isDeployed"`
	Type         string `json:"type"`
	DefaultLimit struct {
		Type         string `json:"type"`
		RequestCount struct {
			TimeUnit     string `json:"timeUnit"`
			UnitTime     int    `json:"unitTime"`
			RequestCount int    `json:"requestCount"`
		} `json:"requestCount"`
		Bandwidth struct {
			TimeUnit   string `json:"timeUnit"`
			UnitTime   int    `json:"unitTime"`
			DataAmount int    `json:"dataAmount"`
			DataUnit   string `json:"dataUnit"`
		} `json:"bandwidth"`
		EventCount struct {
			TimeUnit   string `json:"timeUnit"`
			UnitTime   int    `json:"unitTime"`
			EventCount int    `json:"eventCount"`
		} `json:"eventCount"`
	} `json:"defaultLimit"`
	ConditionalGroups []AdvancedPolicyConditionalGroup `json:"conditionalGroups"`
}

type AdvancedPolicyConditionalGroup struct {
	Description string                    `json:"description"`
	Conditions  []AdvancedPolicyCondition `json:"conditions"`
	Limit       struct {
		Type         string `json:"type"`
		RequestCount struct {
			TimeUnit     string `json:"timeUnit"`
			UnitTime     int    `json:"unitTime"`
			RequestCount int    `json:"requestCount"`
		} `json:"requestCount"`
		Bandwidth  interface{} `json:"bandwidth"`
		EventCount interface{} `json:"eventCount"`
	} `json:"limit"`
}

type AdvancedPolicyCondition struct {
	Type            string `json:"type"`
	InvertCondition bool   `json:"invertCondition"`
	HeaderCondition struct {
		HeaderName  string `json:"headerName"`
		HeaderValue string `json:"headerValue"`
	} `json:"headerCondition"`
	IpCondition struct {
		IpConditionType string      `json:"ipConditionType"`
		SpecificIP      string      `json:"specificIP"`
		StartingIP      interface{} `json:"startingIP"`
		EndingIP        interface{} `json:"endingIP"`
	} `json:"ipCondition"`
	JwtClaimsCondition      interface{} `json:"jwtClaimsCondition"`
	QueryParameterCondition *struct {
		ParameterName  string `json:"parameterName"`
		ParameterValue string `json:"parameterValue"`
	} `json:"queryParameterCondition"`
}
type SubscriptionThrottlePolicy struct {
	PolicyName           string `json:"policyName"`
	DisplayName          string `json:"displayName"`
	Description          string `json:"description"`
	IsDeployed           bool   `json:"isDeployed"`
	Type                 string `json:"type"`
	GraphQLMaxComplexity int    `json:"graphQLMaxComplexity"`
	GraphQLMaxDepth      int    `json:"graphQLMaxDepth"`
	DefaultLimit         struct {
		Type         string `json:"type"`
		RequestCount struct {
			TimeUnit     string `json:"timeUnit"`
			UnitTime     int    `json:"unitTime"`
			RequestCount int    `json:"requestCount"`
		} `json:"requestCount"`
		Bandwidth struct {
			TimeUnit   string `json:"timeUnit"`
			UnitTime   int    `json:"unitTime"`
			DataAmount int    `json:"dataAmount"`
			DataUnit   string `json:"dataUnit"`
		} `json:"bandwidth"`
		EventCount struct {
			TimeUnit   string `json:"timeUnit"`
			UnitTime   int    `json:"unitTime"`
			EventCount int    `json:"eventCount"`
		} `json:"eventCount"`
	} `json:"defaultLimit"`
	Monetization struct {
		MonetizationPlan string `json:"monetizationPlan"`
		Properties       struct {
			Property1 string `json:"property1"`
			Property2 string `json:"property2"`
		} `json:"properties"`
	} `json:"monetization"`
	RateLimitCount    int           `json:"rateLimitCount"`
	RateLimitTimeUnit string        `json:"rateLimitTimeUnit"`
	SubscriberCount   int           `json:"subscriberCount"`
	CustomAttributes  []interface{} `json:"customAttributes"`
	StopOnQuotaReach  bool          `json:"stopOnQuotaReach"`
	BillingPlan       string        `json:"billingPlan"`
	Permissions       struct {
		PermissionType string   `json:"permissionType"`
		Roles          []string `json:"roles"`
	} `json:"permissions"`
}
