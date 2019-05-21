package utils

var DefaultAPISpecs = []byte(`id:
  providerName: admin
  version: ''
  apiName: ''
type: HTTP
tags: []
documents: []
availableTiers:
- name: Bronze
  displayName: Bronze
  description: Allows 1000 requests per minute
  requestsPerMin: 1000
  requestCount: 1000
  unitTime: 1
  timeUnit: min
  tierPlan: FREE
  stopOnQuotaReached: true
- name: Gold
  displayName: Gold
  description: Allows 5000 requests per minute
  requestsPerMin: 5000
  requestCount: 5000
  unitTime: 1
  timeUnit: min
  tierPlan: FREE
  stopOnQuotaReached: true
- name: Silver
  displayName: Silver
  description: Allows 2000 requests per minute
  requestsPerMin: 2000
  requestCount: 2000
  unitTime: 1
  timeUnit: min
  tierPlan: FREE
  stopOnQuotaReached: true
- name: Unlimited
  displayName: Unlimited
  description: Allows unlimited requests
  requestsPerMin: 2147483647
  requestCount: 2147483647
  unitTime: 0
  timeUnit: ms
  tierPlan: FREE
  stopOnQuotaReached: true
availableSubscriptionLevelPolicies: []
uriTemplates: 
apiHeaderChanged: false
apiResourcePatternsChanged: false
status: CREATED
visibility: public
endpointSecured: false
endpointAuthDigest: false
transports: http,https
advertiseOnly: false
corsConfiguration:
  corsConfigurationEnabled: false
  accessControlAllowOrigins:
  - "*"
  accessControlAllowCredentials: false
  accessControlAllowHeaders:
  - authorization
  - Access-Control-Allow-Origin
  - Content-Type
  - SOAPAction
  accessControlAllowMethods:
  - GET
  - PUT
  - POST
  - DELETE
  - PATCH
  - OPTIONS
endpointConfig: 
responseCache: Disabled
cacheTimeout: 300
implementation: ENDPOINT
scopes: 
isDefaultVersion: false
isPublishedDefaultVersion: false
environments:
- Production and Sandbox
environmentList:
- SANDBOX
- PRODUCTION
apiSecurity: oauth2
accessControl: all
rating: 0
isLatest: false
`)
