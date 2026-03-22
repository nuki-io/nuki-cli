package bleflows

// AuthStore loads and persists authorization contexts for Nuki devices.
// Implementations are responsible for the storage backend (config file, OS keychain, etc.).
type AuthStore interface {
	Load(deviceId string) (*AuthorizeContext, error)
	Store(deviceId string, ctx *AuthorizeContext) error
}
