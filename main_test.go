package main

import (
	"strings"
	"testing"

	"github.com/AnatoleLucet/docker-neovim/pkg/config"
	"github.com/AnatoleLucet/docker-neovim/pkg/dockerfile"
)

func TestGenerateImageName(t *testing.T) {
	cfg := config.Config{
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
		{"nightly", "bullseye", "", "test/neovim:nightly-bullseye"},
		{"latest", "alpine", "", "test/neovim:latest"},
		{"latest", "bookworm", "", "test/neovim:latest-bookworm"},
		{"latest", "bullseye", "", "test/neovim:latest-bullseye"},
		{"package", "alpine", "", "test/neovim:unknown-alpine"},
		{"package", "bookworm", "", "test/neovim:unknown-bookworm"},
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

func TestGenerateDockerfile(t *testing.T) {
	// Test nightly dockerfile generation
	nightlyDockerfile := dockerfile.Generate("nightly", "alpine", "master")
	if !strings.Contains(nightlyDockerfile, "FROM alpine AS builder") {
		t.Error("Nightly dockerfile should use multi-stage build")
	}
	if !strings.Contains(nightlyDockerfile, "git clone https://github.com/neovim/neovim.git") {
		t.Error("Nightly dockerfile should clone neovim source")
	}

	// Test package dockerfile generation
	packageDockerfile := dockerfile.Generate("package", "alpine", "")
	if !strings.Contains(packageDockerfile, "apk add --no-cache neovim") {
		t.Error("Package dockerfile should install neovim package")
	}
	if strings.Contains(packageDockerfile, "git clone") {
		t.Error("Package dockerfile should not clone source code")
	}

	// Test different base images
	bookwormDockerfile := dockerfile.Generate("package", "bookworm", "")
	if !strings.Contains(bookwormDockerfile, "FROM debian:bookworm") {
		t.Error("Bookworm dockerfile should use debian:bookworm base")
	}

	bullseyeDockerfile := dockerfile.Generate("package", "bullseye", "")
	if !strings.Contains(bullseyeDockerfile, "FROM debian:bullseye") {
		t.Error("Bullseye dockerfile should use debian:bullseye base")
	}
}