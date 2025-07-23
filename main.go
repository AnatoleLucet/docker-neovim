package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func main() {
	var (
		configPath = flag.String("config", "build-config.json", "Path to build configuration file")
		version    = flag.String("version", "", "Neovim version to build")
		buildType  = flag.String("type", "package", "Build type: 'nightly' or 'package'")
		baseImage  = flag.String("base", "alpine", "Base image: 'alpine', 'bookworm', or 'bullseye'")
		tag        = flag.String("tag", "", "Custom tag for the image")
	)
	flag.Parse()

	config := loadConfig(*configPath)

	// Override config with command line args if provided
	if *version != "" {
		config.Version = *version
	}
	if *buildType != "" {
		config.BuildType = *buildType
	}
	if *baseImage != "" {
		config.BaseImages = []string{*baseImage}
	}

	// Login to Docker registry
	if config.Password != "" {
		err := dockerLogin(config.Username, config.Password)
		if err != nil {
			log.Fatalf("Failed to login to Docker registry: %v", err)
		}
	}

	for _, base := range config.BaseImages {
		imageName := generateImageName(config, base, *tag)

		if !config.AllowOverride && imageExists(imageName) {
			log.Printf("Image %s already exists, skipping", imageName)
			continue
		}

		log.Printf("Building image: %s", imageName)

		dockerfile := generateDockerfile(config.BuildType, base, config.Version)

		err := buildAndPushImage(dockerfile, imageName, config)
		if err != nil {
			log.Fatalf("Failed to build image %s: %v", imageName, err)
		}

		log.Printf("Successfully built and pushed: %s", imageName)
	}
}

func loadConfig(path string) Config {
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

func generateImageName(config Config, baseImage, customTag string) string {
	repo := fmt.Sprintf("%s/%s", config.Username, config.Repository)

	if customTag != "" {
		return fmt.Sprintf("%s:%s", repo, customTag)
	}

	switch config.BuildType {
	case "nightly":
		if baseImage == "alpine" {
			return fmt.Sprintf("%s:nightly-alpine", repo)
		}
		return fmt.Sprintf("%s:nightly-%s", repo, baseImage)
	case "latest":
		if baseImage == "alpine" {
			return fmt.Sprintf("%s:latest", repo)
		}
		return fmt.Sprintf("%s:latest-%s", repo, baseImage)
	default:
		// Version tag
		if baseImage == "alpine" {
			return fmt.Sprintf("%s:%s-alpine", repo, config.Version)
		}
		return fmt.Sprintf("%s:%s-%s", repo, config.Version, baseImage)
	}
}

func generateDockerfile(buildType, baseImage, version string) string {
	switch buildType {
	case "nightly":
		return generateNightlyDockerfile(baseImage, version)
	default:
		return generatePackageDockerfile(baseImage, version)
	}
}

func generateNightlyDockerfile(baseImage, version string) string {
	switch baseImage {
	case "alpine":
		return fmt.Sprintf(`FROM alpine AS builder

LABEL maintainer="AnatoleLucet"

ARG BUILD_DEPS="autoconf automake cmake curl g++ gettext gettext-dev git libtool make ninja openssl pkgconfig unzip binutils wget"
ARG VERSION=%s

RUN apk add --no-cache ${BUILD_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM alpine
COPY --from=builder /usr/local /usr/local/
# Required shared libraries
COPY --from=builder /lib/ld-musl-*.so.1 /lib/
COPY --from=builder /usr/lib/libgcc_s.so.1 /usr/lib/
COPY --from=builder /usr/lib/libintl.so.8 /usr/lib/

CMD ["/usr/local/bin/nvim"]
`, version)

	case "bookworm":
		return fmt.Sprintf(`FROM debian:bookworm AS builder

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive
ARG BUILD_DEPS="ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip git binutils wget"
ARG VERSION=%s

RUN apt update && apt upgrade -y && \
  apt install -y ${BUILD_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM debian:bookworm
COPY --from=builder /usr/local /usr/local/

CMD ["/usr/local/bin/nvim"]
`, version)

	case "bullseye":
		return fmt.Sprintf(`FROM debian:bullseye AS builder

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive
ARG BUILD_DEPS="ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip git binutils wget"
ARG VERSION=%s

RUN apt update && apt upgrade -y && \
  apt install -y ${BUILD_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM debian:bullseye
COPY --from=builder /usr/local /usr/local/

CMD ["/usr/local/bin/nvim"]
`, version)

	default:
		return generateNightlyDockerfile("alpine", version)
	}
}

func generatePackageDockerfile(baseImage, version string) string {
	switch baseImage {
	case "alpine":
		return `FROM alpine

LABEL maintainer="AnatoleLucet"

RUN apk add --no-cache neovim

CMD ["/usr/bin/nvim"]
`

	case "bookworm":
		return `FROM debian:bookworm

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && apt upgrade -y && \
  apt install -y neovim && \
  apt clean && \
  rm -rf /var/lib/apt/lists/*

CMD ["/usr/bin/nvim"]
`

	case "bullseye":
		return `FROM debian:bullseye

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && apt upgrade -y && \
  apt install -y neovim && \
  apt clean && \
  rm -rf /var/lib/apt/lists/*

CMD ["/usr/bin/nvim"]
`

	default:
		return generatePackageDockerfile("alpine", version)
	}
}

func dockerLogin(username, password string) error {
	cmd := exec.Command("docker", "login", "-u", username, "--password-stdin")
	cmd.Stdin = strings.NewReader(password)
	return cmd.Run()
}

func buildAndPushImage(dockerfile, imageName string, config Config) error {
	// Create a temporary directory for the build context
	tmpDir, err := ioutil.TempDir("", "docker-build")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write Dockerfile to temp directory
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	err = ioutil.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Dockerfile: %v", err)
	}

	// Build for multiple platforms
	platformsStr := strings.Join(config.Platforms, ",")

	// Use docker buildx to build and push
	cmd := exec.Command("docker", "buildx", "build",
		"--push",
		"--platform", platformsStr,
		"-t", imageName,
		"-f", dockerfilePath,
		tmpDir,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func imageExists(imageName string) bool {
	// Extract registry, repository and tag from image name
	parts := strings.Split(imageName, ":")
	if len(parts) != 2 {
		return false
	}

	repository := parts[0]
	tag := parts[1]

	// Check if image exists on Docker Hub
	url := fmt.Sprintf("https://index.docker.io/v1/repositories/%s/tags/%s", repository, tag)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func getLatestNeovimVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/neovim/neovim/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}