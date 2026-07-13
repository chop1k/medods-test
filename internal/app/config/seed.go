package config

import "flag"

type SeedConfig struct {
	Path string
}

func RegisterSeedFlags(fs *flag.FlagSet) *SeedConfig {
	cfg := &SeedConfig{}

	fs.StringVar(&cfg.Path, "seeds-path", envOrDefault("SEEDS_PATH", "seeds"),
		"Filesystem path to the seed files (env: SEEDS_PATH)")

	return cfg
}
