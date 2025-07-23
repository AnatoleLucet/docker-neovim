package dockerfile

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Test nightly dockerfile generation
	nightlyDockerfile := Generate("nightly", "alpine", "master")
	if !strings.Contains(nightlyDockerfile, "FROM alpine AS builder") {
		t.Error("Nightly dockerfile should use multi-stage build")
	}
	if !strings.Contains(nightlyDockerfile, "git clone https://github.com/neovim/neovim.git") {
		t.Error("Nightly dockerfile should clone neovim source")
	}

	// Test package dockerfile generation
	packageDockerfile := Generate("package", "alpine", "")
	if !strings.Contains(packageDockerfile, "apk add --no-cache neovim") {
		t.Error("Package dockerfile should install neovim package")
	}
	if strings.Contains(packageDockerfile, "git clone") {
		t.Error("Package dockerfile should not clone source code")
	}
}

func TestGenerateNightly(t *testing.T) {
	tests := []struct {
		baseImage string
		expected  []string
	}{
		{"alpine", []string{"FROM alpine AS builder", "apk add --no-cache"}},
		{"bookworm", []string{"FROM debian:bookworm AS builder", "apt update"}},
		{"bullseye", []string{"FROM debian:bullseye AS builder", "apt update"}},
	}

	for _, tt := range tests {
		dockerfile := generateNightly(tt.baseImage, "master")
		for _, expected := range tt.expected {
			if !strings.Contains(dockerfile, expected) {
				t.Errorf("Nightly dockerfile for %s should contain '%s'", tt.baseImage, expected)
			}
		}
	}
}

func TestGeneratePackage(t *testing.T) {
	tests := []struct {
		baseImage string
		expected  []string
	}{
		{"alpine", []string{"FROM alpine", "apk add --no-cache neovim"}},
		{"bookworm", []string{"FROM debian:bookworm", "apt install -y neovim"}},
		{"bullseye", []string{"FROM debian:bullseye", "apt install -y neovim"}},
	}

	for _, tt := range tests {
		dockerfile := generatePackage(tt.baseImage, "")
		for _, expected := range tt.expected {
			if !strings.Contains(dockerfile, expected) {
				t.Errorf("Package dockerfile for %s should contain '%s'", tt.baseImage, expected)
			}
		}
	}
}

func TestDefaultBaseImage(t *testing.T) {
	// Test that unknown base images default to alpine
	nightlyDockerfile := generateNightly("unknown", "master")
	if !strings.Contains(nightlyDockerfile, "FROM alpine AS builder") {
		t.Error("Unknown base image should default to alpine for nightly builds")
	}

	packageDockerfile := generatePackage("unknown", "")
	if !strings.Contains(packageDockerfile, "FROM alpine") {
		t.Error("Unknown base image should default to alpine for package builds")
	}
}