package environment

// Env is a custom type to determine environment we run in.
type Env string

// Available environments we know.
const (
	Local      Env = "local"
	Production Env = "production"

	ServiceName = "driveo-pa"
)

// IsLocal returns true if we are running in a local environment.
func (e Env) IsLocal() bool {
	return e == Local
}

// IsProduction returns true if we are running in a production environment.
func (e Env) IsProduction() bool {
	return e == Production
}

// String returns string representation of the Env.
func (e Env) String() string {
	return string(e)
}

// IsLocal returns true if value is a local environment.
func IsLocal(value string) bool {
	return Env(value).IsLocal()
}

// IsProduction returns true if value is a production environment.
func IsProduction(value string) bool {
	return Env(value).IsProduction()
}
