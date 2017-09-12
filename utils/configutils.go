package utils

// ------------------- Structs for YAML Config Files ----------------------------------

// For env_keys_all.yaml
// Not to be manually edited
type EnvKeysAll struct {
	Environments map[string]EnvKeys `yaml:"environments"`
}

// For env_endpoints_all.yaml
// To be manually edited by the user
type EnvEndpointsAll struct {
	Environments map[string]EnvEndpoints `yaml:"environments"`
}

type EnvKeys struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"` // to be encrypted (with the user's password) and stored
	Username     string `yaml:"username"`
}

type EnvEndpoints struct {
	APIManagerEndpoint   string `yaml:"api_manager_endpoint"`
	RegistrationEndpoint string `yaml:"registration_endpoint"`
	TokenEndpoint        string `yaml:"token_endpoint"`
}

// ---------------- End of Structs for YAML Config Files ---------------------------------

type API struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Context         string `json:"context"`
	Version         string `json:"version"`
	Provider        string `json:"provider"`
	LifeCycleStatus string `json:"lifeCycleStatus"`
	WorkflowStatus  string `json:"workflowStatus"`
}

type RegistrationResponse struct {
	ClientSecretExpiresAt string `json:"client_secret_expires_at"`
	ClientID              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	ClientName            string `json:"client_name"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int32 `json:"expires_in"`
}

type APIListResponse struct {
	Count int32 `json:"count"`
	List  []API `json:"list"`
}
