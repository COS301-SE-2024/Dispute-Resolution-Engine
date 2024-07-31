package env

type Env interface {
	LoadFromFile(files ...string)
	Register(key string)
	RegisterDefault(key, fallback string)
	Get(key string) (string, error)
}