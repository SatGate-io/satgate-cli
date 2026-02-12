package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds CLI configuration
type Config struct {
	Surface     string `yaml:"surface"`      // gateway | cloud
	Gateway     string `yaml:"gateway"`      // gateway admin URL
	AdminToken  string `yaml:"admin_token"`  // X-Admin-Token (gateway surface)
	BearerToken string `yaml:"bearer_token"` // Bearer token (cloud surface)
	Tenant      string `yaml:"tenant"`       // tenant slug (cloud surface)
	Format      string `yaml:"format"`       // table | json | yaml
}

var current *Config

// Get returns the current config
func Get() *Config {
	if current == nil {
		current = &Config{
			Surface: "gateway",
			Gateway: "http://localhost:9090",
			Format:  "table",
		}
	}
	return current
}

// Load reads config from file, env vars, and applies defaults
func Load(cfgFile string) {
	cfg := &Config{}

	// 1. Read config file
	if cfgFile == "" {
		home, _ := os.UserHomeDir()
		cfgFile = filepath.Join(home, ".satgate", "config.yaml")
	}
	data, err := os.ReadFile(cfgFile)
	if err == nil {
		yaml.Unmarshal(data, cfg)
	}

	// 2. Env var overrides
	if v := os.Getenv("SATGATE_SURFACE"); v != "" {
		cfg.Surface = v
	}
	if v := os.Getenv("SATGATE_GATEWAY"); v != "" {
		cfg.Gateway = v
	}
	if v := os.Getenv("SATGATE_ADMIN_TOKEN"); v != "" {
		cfg.AdminToken = v
	}
	if v := os.Getenv("SATGATE_BEARER_TOKEN"); v != "" {
		cfg.BearerToken = v
	}
	if v := os.Getenv("SATGATE_TENANT"); v != "" {
		cfg.Tenant = v
	}
	if v := os.Getenv("SATGATE_FORMAT"); v != "" {
		cfg.Format = v
	}

	// 3. Auto-detect surface from gateway URL
	if cfg.Surface == "" {
		if strings.Contains(cfg.Gateway, "cloud.satgate.io") || strings.Contains(cfg.Gateway, "satgate.io/api") {
			cfg.Surface = "cloud"
		} else {
			cfg.Surface = "gateway"
		}
	}

	// 4. Defaults
	if cfg.Gateway == "" {
		cfg.Gateway = "http://localhost:9090"
	}
	if cfg.Format == "" {
		cfg.Format = "table"
	}

	current = cfg
}

// AuthHeader returns the appropriate auth header key and value
func (c *Config) AuthHeader() (string, string) {
	if c.Surface == "cloud" {
		return "Authorization", fmt.Sprintf("Bearer %s", c.BearerToken)
	}
	return "X-Admin-Token", c.AdminToken
}

// Validate checks that required config is present
func (c *Config) Validate() error {
	if c.Gateway == "" {
		return fmt.Errorf("gateway URL not configured. Run 'satgate configure' or set SATGATE_GATEWAY")
	}
	if c.Surface == "cloud" && c.BearerToken == "" {
		return fmt.Errorf("bearer token not configured for cloud surface. Set SATGATE_BEARER_TOKEN")
	}
	if c.Surface == "gateway" && c.AdminToken == "" {
		return fmt.Errorf("admin token not configured. Set SATGATE_ADMIN_TOKEN")
	}
	return nil
}
