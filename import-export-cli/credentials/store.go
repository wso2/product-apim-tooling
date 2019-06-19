package credentials

type Store interface {
	Has(env string) bool
	Get(env string) (Credential, error)
	Set(env, username, password, clientId, clientSecret string) error
	Erase(env string) error
	Load() error
}
