package v2

const (
	EpHttp        = "http"
	EpLoadbalance = "load_balance"
	EpFailover    = "failover"
)

// APIDefinition represents an API artifact in APIM
type APIDefinition struct {
	ID                                 ID                 `json:"id,omitempty"`
	UUID                               string             `json:"uuid,omitempty"`
	Description                        string             `json:"description,omitempty"`
	Type                               string             `json:"type,omitempty"`
	Context                            string             `json:"context"`
	ContextTemplate                    string             `json:"contextTemplate,omitempty"`
	Tags                               []string           `json:"tags"`
	Documents                          []interface{}      `json:"documents,omitempty"`
	LastUpdated                        string             `json:"lastUpdated,omitempty"`
	AvailableTiers                     []AvailableTiers   `json:"availableTiers,omitempty"`
	AvailableSubscriptionLevelPolicies []interface{}      `json:"availableSubscriptionLevelPolicies,omitempty"`
	URITemplates                       []URITemplates     `json:"uriTemplates"`
	APIHeaderChanged                   bool               `json:"apiHeaderChanged,omitempty"`
	APIResourcePatternsChanged         bool               `json:"apiResourcePatternsChanged,omitempty"`
	Status                             string             `json:"status,omitempty"`
	TechnicalOwner                     string             `json:"technicalOwner,omitempty"`
	TechnicalOwnerEmail                string             `json:"technicalOwnerEmail,omitempty"`
	BusinessOwner                      string             `json:"businessOwner,omitempty"`
	BusinessOwnerEmail                 string             `json:"businessOwnerEmail,omitempty"`
	Visibility                         string             `json:"visibility,omitempty"`
	EndpointSecured                    bool               `json:"endpointSecured,omitempty"`
	EndpointAuthDigest                 bool               `json:"endpointAuthDigest,omitempty"`
	EndpointUTUsername                 string             `json:"endpointUTUsername,omitempty"`
	Transports                         string             `json:"transports,omitempty"`
	InSequence                         string             `json:"inSequence,omitempty"`
	OutSequence                        string             `json:"outSequence,omitempty"`
	FaultSequence                      string             `json:"faultSequence,omitempty"`
	AdvertiseOnly                      bool               `json:"advertiseOnly,omitempty"`
	CorsConfiguration                  *CorsConfiguration `json:"corsConfiguration,omitempty"`
	ProductionUrl                      string             `json:"productionUrl,omitempty"`
	SandboxUrl                         string             `json:"sandboxUrl,omitempty"`
	EndpointConfig                     *string            `json:"endpointConfig,omitempty"`
	ResponseCache                      string             `json:"responseCache,omitempty"`
	CacheTimeout                       int                `json:"cacheTimeout,omitempty"`
	Implementation                     string             `json:"implementation,omitempty"`
	AuthorizationHeader                string             `json:"authorizationHeader,omitempty"`
	Scopes                             []interface{}      `json:"scopes,omitempty"`
	IsDefaultVersion                   bool               `json:"isDefaultVersion,omitempty"`
	IsPublishedDefaultVersion          bool               `json:"isPublishedDefaultVersion,omitempty"`
	Environments                       []string           `json:"environments,omitempty"`
	CreatedTime                        string             `json:"createdTime,omitempty"`
	AdditionalProperties               map[string]string  `json:"additionalProperties,omitempty"`
	EnvironmentList                    []string           `json:"environmentList,omitempty"`
	APISecurity                        string             `json:"apiSecurity,omitempty"`
	AccessControl                      string             `json:"accessControl,omitempty"`
	Rating                             float64            `json:"rating,omitempty"`
	IsLatest                           bool               `json:"isLatest,omitempty"`
}
type ID struct {
	ProviderName string `json:"providerName"`
	APIName      string `json:"apiName"`
	Version      string `json:"version"`
}
type AvailableTiers struct {
	Name               string `json:"name,omitempty"`
	DisplayName        string `json:"displayName,omitempty"`
	Description        string `json:"description,omitempty"`
	RequestsPerMin     int    `json:"requestsPerMin,omitempty"`
	RequestCount       int    `json:"requestCount,omitempty"`
	UnitTime           int    `json:"unitTime,omitempty"`
	TimeUnit           string `json:"timeUnit,omitempty"`
	TierPlan           string `json:"tierPlan,omitempty"`
	StopOnQuotaReached bool   `json:"stopOnQuotaReached,omitempty"`
}
type Scopes struct {
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	Roles       string `json:"roles,omitempty"`
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
}
type MediationScripts struct {
}
type URITemplates struct {
	URITemplate          string            `json:"uriTemplate,omitempty"`
	HTTPVerb             string            `json:"httpVerb,omitempty"`
	AuthType             string            `json:"authType,omitempty"`
	HTTPVerbs            []string          `json:"httpVerbs,omitempty"`
	AuthTypes            []string          `json:"authTypes,omitempty"`
	ThrottlingConditions []interface{}     `json:"throttlingConditions,omitempty"`
	ThrottlingTier       string            `json:"throttlingTier,omitempty"`
	ThrottlingTiers      []string          `json:"throttlingTiers,omitempty"`
	MediationScript      string            `json:"mediationScript,omitempty"`
	Scopes               []*Scopes         `json:"scopes,omitempty"`
	MediationScripts     *MediationScripts `json:"mediationScripts,omitempty"`
}
type CorsConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled,omitempty"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty"`
}
