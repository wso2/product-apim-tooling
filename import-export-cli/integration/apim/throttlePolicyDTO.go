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

// ApplicationThrottlePolicy : Application Throttle Policy DTO
type ApplicationThrottlePolicy struct {
	PolicyName   string       `json:"policyName"`
	DisplayName  string       `json:"displayName"`
	Description  string       `json:"description"`
	IsDeployed   bool         `json:"isDeployed"`
	Type         string       `json:"type"`
	DefaultLimit DefaultLimit `json:"defaultLimit"`
}

// CustomThrottlePolicy : Custom Throttle Policy DTO
type CustomThrottlePolicy struct {
	PolicyName  string `json:"policyName"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IsDeployed  bool   `json:"isDeployed"`
	Type        string `json:"type"`
	SiddhiQuery string `json:"siddhiQuery"`
	KeyTemplate string `json:"keyTemplate"`
}

// AdvancedThrottlePolicy : Advanced Throttle Policy DTO
type AdvancedThrottlePolicy struct {
	PolicyName        string                           `json:"policyName"`
	DisplayName       string                           `json:"displayName"`
	Description       string                           `json:"description"`
	IsDeployed        bool                             `json:"isDeployed"`
	Type              string                           `json:"type"`
	DefaultLimit      DefaultLimit                     `json:"defaultLimit"`
	ConditionalGroups []AdvancedPolicyConditionalGroup `json:"conditionalGroups"`
}

// AdvancedPolicyConditionalGroup : Conditional Groups in Advanced Throttle Policy DTO
type AdvancedPolicyConditionalGroup struct {
	Description string                    `json:"description"`
	Conditions  []AdvancedPolicyCondition `json:"conditions"`
	Limit       DefaultLimit              `json:"limit"`
}

// AdvancedPolicyCondition : Advanced Throttle Policy Condition in Conditional Groups
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

// SubscriptionThrottlePolicy : Subscription Throttle Policy DTO
type SubscriptionThrottlePolicy struct {
	PolicyName           string       `json:"policyName"`
	DisplayName          string       `json:"displayName"`
	Description          string       `json:"description"`
	IsDeployed           bool         `json:"isDeployed"`
	Type                 string       `json:"type"`
	GraphQLMaxComplexity int          `json:"graphQLMaxComplexity"`
	GraphQLMaxDepth      int          `json:"graphQLMaxDepth"`
	DefaultLimit         DefaultLimit `json:"defaultLimit"`
	Monetization         struct {
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

type DefaultLimit struct {
	Type         string       `json:"type"`
	RequestCount RequestCount `json:"requestCount"`
	Bandwidth    Bandwidth    `json:"bandwidth"`
	EventCount   EventCount   `json:"eventCount"`
	AiApiQuota   AiApiQuota   `json:"aiApiQuota"`
}

type RequestCount struct {
	TimeUnit     string `json:"timeUnit"`
	UnitTime     int    `json:"unitTime"`
	RequestCount int    `json:"requestCount"`
}

type Bandwidth struct {
	TimeUnit   string `json:"timeUnit"`
	UnitTime   int    `json:"unitTime"`
	DataAmount int    `json:"dataAmount"`
	DataUnit   string `json:"dataUnit"`
}

type EventCount struct {
	TimeUnit   string `json:"timeUnit"`
	UnitTime   int    `json:"unitTime"`
	EventCount int    `json:"eventCount"`
}

type AiApiQuota struct {
	TimeUnit             string `json:"timeUnit"`
	UnitTime             int    `json:"unitTime"`
	RequestCount         int    `json:"requestCount"`
	TotalTokenCount      int    `json:"totalTokenCount"`
	PromptTokenCount     int    `json:"promptTokenCount"`
	CompletionTokenCount int    `json:"completionTokenCount"`
}
