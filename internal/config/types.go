package config

type Config struct {
	IOMode          string                   `hcl:"provider"`
	CommandResource []*CommandResourceConfig `hcl:"resource,block"`
}

type CommandResourceConfig struct {
	Resource    string               `hcl:",label"`
	Name        string               `hcl:",label"`
	Environment []*EnvironmentConfig `hcl:"environment,block"`
}

type EnvironmentConfig struct {
	Type    string     `hcl:"type,label"`
	Command string     `hcl:"command"`
	Args    string     `hcl:"args"`
	EnvVar  *EnvMapper `hcl:"env"`
}

type EnvMapper struct {
	ID       string `cty:"id"`
	Login    string `cty:"login"`
	Password string `cty:"password"`
}

type KeyPath struct {
	Env   string `yaml:"env,omitempty"`
	Path  string `yaml:"path"`
	Field string `yaml:"field,omitempty"`
}

type EnvEntry struct {
	Key           string
	Field         string
	ValuePassword string
	ProviderName  string
	IsFound       bool
	ValueLogin    string
}

type Provider interface {
	Get(p KeyPath) (*EnvEntry, error)
}

func (k *KeyPath) Found(id string, v string) EnvEntry {
	return EnvEntry{
		IsFound:       true,
		Key:           k.Env,
		Field:         k.Field,
		ValuePassword: v,
		ValueLogin:    id,
	}
}

type MetaInfo struct {
	Description    string
	Name           string
	Authentication string
}
type RegisteredProvider struct {
	Meta    MetaInfo
	Builder func() Provider
}

type CommandDetails struct {
	Command                    string
	CommandArg                 string
	CommandEnvVariableID       string
	CommandEnvVariableLogin    string
	CommandEnvVariablePassword string
}

type FileSettings struct {
	PasswordManager string
	SessionTimeout  int
	VaultLocation   string
	AuthSecret      string
}
