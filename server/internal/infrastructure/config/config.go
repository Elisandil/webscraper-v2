package config

import (
	"fmt"
	"os"
	"webscraper-v2/pkg/crypto"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Scraping ScrapingConfig `yaml:"scraping"`
	Features FeaturesConfig `yaml:"features"`
	Auth     AuthConfig     `yaml:"auth"`
	Chat     *ChatConfig    `yaml:"chat,omitempty"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type ScrapingConfig struct {
	UserAgent      string `yaml:"user_agent"`
	Timeout        int    `yaml:"timeout"`
	MaxRedirects   int    `yaml:"max_redirects"`
	ExtractImages  bool   `yaml:"extract_images"`
	ExtractFavicon bool   `yaml:"extract_favicon"`
	ExtractHeaders bool   `yaml:"extract_headers"`
	MaxLinks       int    `yaml:"max_links"`
	MaxImages      int    `yaml:"max_images"`
}

type FeaturesConfig struct {
	EnableAnalytics bool `yaml:"enable_analytics"`
	EnableCaching   bool `yaml:"enable_caching"`
	CacheDuration   int  `yaml:"cache_duration"`
}

type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	TokenDuration int    `yaml:"token_duration_hours"`
	BCryptCost    int    `yaml:"bcrypt_cost"`
	RequireAuth   bool   `yaml:"require_auth"`
	DefaultRole   string `yaml:"default_role"`
}

type ChatConfig struct {
	HFAPIToken string `yaml:"hf_api_token"`
	HFModelID  string `yaml:"hf_model_id"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.setDefaults()

	if err := config.setupJWTSecret(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) setDefaults() {

	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Database.Path == "" {
		c.Database.Path = "./data/scraper.db"
	}
	if c.Scraping.UserAgent == "" {
		c.Scraping.UserAgent = "WebScraper/1.0"
	}
	if c.Scraping.Timeout == 0 {
		c.Scraping.Timeout = 30
	}
	if c.Scraping.MaxRedirects == 0 {
		c.Scraping.MaxRedirects = 10
	}
	if c.Scraping.MaxLinks == 0 {
		c.Scraping.MaxLinks = 100
	}
	if c.Scraping.MaxImages == 0 {
		c.Scraping.MaxImages = 50
	}
	if c.Features.CacheDuration == 0 {
		c.Features.CacheDuration = 3600
	}
	if c.Auth.TokenDuration == 0 {
		c.Auth.TokenDuration = 1
	}
	if c.Auth.BCryptCost == 0 {
		c.Auth.BCryptCost = 12
	}
	if c.Auth.DefaultRole == "" {
		c.Auth.DefaultRole = "user"
	}
}

func (c *Config) setupJWTSecret() error {

	if c.Auth.JWTSecret == "" {

		c.Auth.JWTSecret = os.Getenv("JWT_SECRET")
		if c.Auth.JWTSecret == "" {
			if os.Getenv("ENV") == "development" {
				secret, err := crypto.GenerateRandomSecret(32)
				if err != nil {
					return fmt.Errorf("failed to generate JWT secret: %w", err)
				}
				c.Auth.JWTSecret = secret
				fmt.Println("⚠️  WARNING: Using auto-generated JWT secret. Set JWT_SECRET env var or add jwt_secret to config.yaml for production!")
			} else {
				return fmt.Errorf("JWT_SECRET must be set in config.yaml or as environment variable in production")
			}
		}
	}

	if err := crypto.ValidateSecretLength(c.Auth.JWTSecret, 32); err != nil {
		return err
	}

	return nil
}
