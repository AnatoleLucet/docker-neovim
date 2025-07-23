package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Registry      string   `json:"registry"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	Repository    string   `json:"repository"`
	Version       string   `json:"version"`
	BuildType     string   `json:"build_type"` // "nightly" or "package"
	BaseImages    []string `json:"base_images"`
	Platforms     []string `json:"platforms"`
	AllowOverride bool     `json:"allow_override"`
}

func Load(path string) Config {
	// Default configuration
	config := Config{
		Registry:      "docker.io",
		Username:      "anatolelucet",
		Repository:    "neovim",
		BaseImages:    []string{"alpine"},
		Platforms:     []string{"linux/amd64", "linux/arm64"},
		AllowOverride: false,
	}

	// Try to load from file if it exists
	if _, err := os.Stat(path); err == nil {
		data, err := ioutil.ReadFile(path)
		if err == nil {
			json.Unmarshal(data, &config)
		}
	}

	// Override with environment variables
	if password := os.Getenv("DOCKER_PASSWORD"); password != "" {
		config.Password = password
	}
	if version := os.Getenv("VERSION"); version != "" {
		config.Version = version
	}
	if buildType := os.Getenv("BUILD_TYPE"); buildType != "" {
		config.BuildType = buildType
	}
	if allowOverride := os.Getenv("ALLOW_OVERRIDE"); allowOverride == "true" {
		config.AllowOverride = true
	}

	return config
}

func (c *Config) GenerateImageName(baseImage, customTag string) string {
	repo := c.Username + "/" + c.Repository

	if customTag != "" {
		return repo + ":" + customTag
	}

	switch c.BuildType {
	case "nightly":
		if baseImage == "alpine" {
			return repo + ":nightly-alpine"
		}
		return repo + ":nightly-" + baseImage
	case "latest":
		if baseImage == "alpine" {
			return repo + ":latest"
		}
		return repo + ":latest-" + baseImage
	default:
		// Version tag
		version := c.Version
		if version == "" {
			version = "unknown"
		}
		if baseImage == "alpine" {
			return repo + ":" + version + "-alpine"
		}
		return repo + ":" + version + "-" + baseImage
	}
}