package credentials

type Store interface {
	// Has an env in the store
	Has(env string) bool
	// Get an env from store returns intended Credential or an error
	Get(env string) (Credential, error)
	// Set credentials for env using given username,password,clientId,clientSecret
	Set(env, username, password, clientId, clientSecret string) error
	// Erase credentials from given env
	Erase(env string) error
	// Load store
	Load() error
}
