package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test loading with no config file (defaults)
	cfg := Load("nonexistent.json")
	if cfg.Registry != "docker.io" {
		t.Errorf("Expected default registry 'docker.io', got %s", cfg.Registry)
	}
	if cfg.Username != "anatolelucet" {
		t.Errorf("Expected default username 'anatolelucet', got %s", cfg.Username)
	}
	if cfg.AllowOverride != false {
		t.Errorf("Expected default AllowOverride false, got %v", cfg.AllowOverride)
	}
}

func TestGenerateImageName(t *testing.T) {
	cfg := Config{
		Username:   "test",
		Repository: "neovim",
	}

	tests := []struct {
		buildType string
		baseImage string
		customTag string
		expected  string
	}{
		{"nightly", "alpine", "", "test/neovim:nightly-alpine"},
		{"nightly", "bookworm", "", "test/neovim:nightly-bookworm"},
		{"latest", "alpine", "", "test/neovim:latest"},
		{"latest", "bookworm", "", "test/neovim:latest-bookworm"},
		{"package", "alpine", "", "test/neovim:unknown-alpine"},
		{"", "alpine", "custom", "test/neovim:custom"},
	}

	for _, tt := range tests {
		cfg.BuildType = tt.buildType
		result := cfg.GenerateImageName(tt.baseImage, tt.customTag)
		if result != tt.expected {
			t.Errorf("GenerateImageName(%s, %s, %s) = %s, want %s", 
				tt.buildType, tt.baseImage, tt.customTag, result, tt.expected)
		}
	}
}

func TestLoadWithFile(t *testing.T) {
	// Create a temporary config file
	configContent := `{
		"registry": "custom.registry.com",
		"username": "testuser",
		"repository": "testapp",
		"allow_override": true
	}`
	
	tmpFile, err := ioutil.TempFile("", "test-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	
	cfg := Load(tmpFile.Name())
	if cfg.Registry != "custom.registry.com" {
		t.Errorf("Expected registry 'custom.registry.com', got %s", cfg.Registry)
	}
	if cfg.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", cfg.Username)
	}
	if cfg.AllowOverride != true {
		t.Errorf("Expected AllowOverride true, got %v", cfg.AllowOverride)
	}
}

func TestEnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("DOCKER_PASSWORD", "secret")
	os.Setenv("VERSION", "v1.0.0")
	os.Setenv("BUILD_TYPE", "nightly")
	os.Setenv("ALLOW_OVERRIDE", "true")
	
	defer func() {
		os.Unsetenv("DOCKER_PASSWORD")
		os.Unsetenv("VERSION")
		os.Unsetenv("BUILD_TYPE")
		os.Unsetenv("ALLOW_OVERRIDE")
	}()
	
	cfg := Load("nonexistent.json")
	if cfg.Password != "secret" {
		t.Errorf("Expected password from env, got %s", cfg.Password)
	}
	if cfg.Version != "v1.0.0" {
		t.Errorf("Expected version from env, got %s", cfg.Version)
	}
	if cfg.BuildType != "nightly" {
		t.Errorf("Expected build type from env, got %s", cfg.BuildType)
	}
	if cfg.AllowOverride != true {
		t.Errorf("Expected AllowOverride from env, got %v", cfg.AllowOverride)
	}
}