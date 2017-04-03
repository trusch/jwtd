package storage

type Group struct {
	Name   string                       `json:"name"`
	Rights map[string]map[string]string `json:"rights"`
}

type User struct {
	Name         string   `json:"name"`
	PasswordHash string   `json:"passwordhash,omitempty"`
	Groups       []string `json:"groups"`
}

type ProjectConfig struct {
	Users  []*User
	Groups []*Group
}

type StorageBackend interface {
	Load() (*ProjectConfig, error)
	Save(cfg *ProjectConfig) error
}
