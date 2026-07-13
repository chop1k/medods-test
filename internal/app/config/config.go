package config

import "flag"

type Config struct {
	DB         *DatabaseConfig
	HTTP       *ServerConfig
	Migrations *MigrationConfig
	Seeds      *SeedConfig
}

func RegisterConfigFlags(fs *flag.FlagSet) *Config {
	return &Config{
		DB:         RegisterDatabaseFlags(fs),
		HTTP:       RegisterServerFlags(fs),
		Migrations: RegisterMigrationFlags(fs),
		Seeds:      RegisterSeedFlags(fs),
	}
}
