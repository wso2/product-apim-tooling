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
	Username     string        `yaml:"username"`
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

type APIListResponse struct {
	Count int32 `json:"count"`
	List []API `json:"list"`
}
