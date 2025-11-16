package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	PollIntervalSeconds int `json:"poll_interval_second"`
	DefaultCooldownMin 	int `json: "default_cooldown_min"`
	OutboxMaxRetry 		int `json: "outbox_max_retry"`
	StoreKind			string `json: "store_kind"`
	DBPath				string `json: "db_path"`
}

func Default() Config {
	return Config{
		PollIntervalSeconds: 20,
		DefaultCooldownMin:  30,
		OutboxMaxRetry: 	 3,
		StoreKind:			 "memory",
		DBPath:              "./alerts.db",
	}
}

func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}