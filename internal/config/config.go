// Package config loads OMEN's runtime configuration from a YAML file, with
// environment variables overriding secret-bearing fields.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration wraps time.Duration so it can be expressed in config.yaml as a
// plain string like "30m" or "24h".
type Duration struct {
	time.Duration
}

// UnmarshalYAML implements yaml.Unmarshaler for Duration.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	parsed, err := time.ParseDuration(value.Value)
	if err != nil {
		return fmt.Errorf("parse duration %q: %w", value.Value, err)
	}
	d.Duration = parsed
	return nil
}

// Config is OMEN's full runtime configuration.
type Config struct {
	Server struct {
		ListenAddr string `yaml:"listen_addr"`
	} `yaml:"server"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`

	Poll struct {
		VictimsInterval Duration `yaml:"victims_interval"`
		GroupsInterval  Duration `yaml:"groups_interval"`
		KEVInterval     Duration `yaml:"kev_interval"`
		// BackfillMonths is a pointer so an explicit "0" (disable backfill)
		// can be distinguished from "unset" (apply the default).
		BackfillMonths *int `yaml:"backfill_months"`
	} `yaml:"poll"`

	RansomwareLive struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"ransomwarelive"`

	Shodan struct {
		APIKey            string   `yaml:"api_key"`
		RateLimitInterval Duration `yaml:"rate_limit_interval"`
	} `yaml:"shodan"`

	KEV struct {
		CachePath string `yaml:"cache_path"`
		FeedURL   string `yaml:"feed_url"`
	} `yaml:"kev"`

	Watchlist struct {
		DefaultNotifyVia string `yaml:"default_notify_via"`
	} `yaml:"watchlist"`
}

// Load reads and parses the YAML config file at path, then applies
// environment variable overrides for secret fields (SHODAN_API_KEY).
func Load(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	applyDefaults(&cfg)

	if key := os.Getenv("SHODAN_API_KEY"); key != "" {
		cfg.Shodan.APIKey = key
	}

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Server.ListenAddr == "" {
		cfg.Server.ListenAddr = ":8080"
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "./data/omen.db"
	}
	if cfg.Poll.VictimsInterval.Duration == 0 {
		cfg.Poll.VictimsInterval.Duration = 30 * time.Minute
	}
	if cfg.Poll.GroupsInterval.Duration == 0 {
		cfg.Poll.GroupsInterval.Duration = 24 * time.Hour
	}
	if cfg.Poll.KEVInterval.Duration == 0 {
		cfg.Poll.KEVInterval.Duration = 24 * time.Hour
	}
	if cfg.Poll.BackfillMonths == nil {
		defaultMonths := 3
		cfg.Poll.BackfillMonths = &defaultMonths
	}
	if cfg.RansomwareLive.BaseURL == "" {
		cfg.RansomwareLive.BaseURL = "https://api.ransomware.live/v2"
	}
	if cfg.Shodan.RateLimitInterval.Duration == 0 {
		cfg.Shodan.RateLimitInterval.Duration = time.Second
	}
	if cfg.KEV.CachePath == "" {
		cfg.KEV.CachePath = "./data/kev-cache.json"
	}
	if cfg.KEV.FeedURL == "" {
		cfg.KEV.FeedURL = "https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json"
	}
}
