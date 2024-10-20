package env

type MockEnv struct {
	ThrowError bool
	Error 	error
	ReturnGet string
}


func (m MockEnv) LoadFromFile(files ...string) {
	return
}

func (m MockEnv) Register(key string) {
	return
}

func (m MockEnv) RegisterDefault(key, fallback string) {
	return
}

func (m MockEnv) Get(key string) (string, error) {
	if m.ThrowError {
		return "", m.Error
	}
	return m.ReturnGet, nil
}
