package main

type AppConfig struct {
	Runtime struct {
		EnvPrefix    string            `toml:"env_prefix"`
		EnvVariables map[string]string `toml:"env_variables"`
	} `toml:"runtime"`
	PackageName string
}
