package model

type Config struct {
	Environments []Env
}

type Env struct {
	Platform   string `toml:"Platform"`
	Service    string `toml:"Service"`
	Account    string `toml:"Account"`
	SecretName string `toml:"SecretName"`
	ExportName string `toml:"ExportName"`
	Version    string `toml:"Version"`
	Key        string `toml:"Key"`
}
