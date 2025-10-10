package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Scraping ScrapingConfig `yaml:"scraping"`
	Features FeaturesConfig `yaml:"features"`
	Auth     AuthConfig     `yaml:"auth"`
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

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Database.Path == "" {
		config.Database.Path = "./data/scraper.db"
	}
	if config.Scraping.UserAgent == "" {
		config.Scraping.UserAgent = "WebScraper/1.0"
	}
	if config.Scraping.Timeout == 0 {
		config.Scraping.Timeout = 30
	}
	if config.Scraping.MaxRedirects == 0 {
		config.Scraping.MaxRedirects = 10
	}
	if config.Scraping.MaxLinks == 0 {
		config.Scraping.MaxLinks = 100
	}
	if config.Scraping.MaxImages == 0 {
		config.Scraping.MaxImages = 50
	}
	if config.Features.CacheDuration == 0 {
		config.Features.CacheDuration = 3600
	}

	if config.Auth.JWTSecret == "" {
		config.Auth.JWTSecret = os.Getenv("JWT_SECRET")

		if config.Auth.JWTSecret == "" {

			if os.Getenv("ENV") == "development" {
				secret, err := generateRandomSecret(32)
				if err != nil {
					return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
				}
				config.Auth.JWTSecret = secret
				fmt.Println("⚠️  WARNING: Using auto-generated JWT secret. Set JWT_SECRET env var or add jwt_secret to config.yaml for production!")
			} else {
				return nil, fmt.Errorf("JWT_SECRET must be set in config.yaml or as environment variable in production")
			}
		}
	}

	if len(config.Auth.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	if config.Auth.TokenDuration == 0 {
		config.Auth.TokenDuration = 1
	}
	if config.Auth.BCryptCost == 0 {
		config.Auth.BCryptCost = 12
	}
	if config.Auth.DefaultRole == "" {
		config.Auth.DefaultRole = "user"
	}
	return &config, nil
}

func generateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
