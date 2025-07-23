# Docker NeoVim

An up-to-date, ready-to-use NeoVim image built with Go and Docker Buildx.

## Features

- Supports multiple build types (nightly from source, latest/versioned from packages)
- Multi-platform builds (linux/amd64, linux/arm64)
- Multiple base images (Alpine, Debian Bookworm, Debian Bullseye)
- Automatically updated via GitHub Actions
- Built with Go for better maintainability

## Image Tags

### Nightly Builds (from source)
- `anatolelucet/neovim:nightly-alpine` - Built from master branch on Alpine
- `anatolelucet/neovim:nightly-bookworm` - Built from master branch on Debian Bookworm
- `anatolelucet/neovim:nightly-bullseye` - Built from master branch on Debian Bullseye

### Latest Stable (from package managers)
- `anatolelucet/neovim:latest` - Latest stable from Alpine packages (alias for latest-alpine)
- `anatolelucet/neovim:latest-alpine` - Latest stable from Alpine packages
- `anatolelucet/neovim:latest-bookworm` - Latest stable from Debian Bookworm packages
- `anatolelucet/neovim:latest-bullseye` - Latest stable from Debian Bullseye packages

### Versioned Releases (from package managers)
- `anatolelucet/neovim:X.Y.Z-alpine` - Specific version from Alpine packages
- `anatolelucet/neovim:X.Y.Z-bookworm` - Specific version from Debian Bookworm packages
- `anatolelucet/neovim:X.Y.Z-bullseye` - Specific version from Debian Bullseye packages

## How to use

You can use this image anywhere you want with something like:
```bash
# Use latest stable version (Alpine-based)
docker run -it -v `pwd`:/mnt/volume --workdir=/mnt/volume anatolelucet/neovim:latest

# Use nightly build (Alpine-based)
docker run -it -v `pwd`:/mnt/volume --workdir=/mnt/volume anatolelucet/neovim:nightly-alpine

# Use specific version on Debian Bookworm
docker run -it -v `pwd`:/mnt/volume --workdir=/mnt/volume anatolelucet/neovim:0.10.0-bookworm
```

You can also extend this image in a Dockerfile to make your own (possibly containing your personal setup).

## Build System

This repository uses a Go-based build system that:
- Generates appropriate Dockerfiles for each base image and build type
- Uses Docker Buildx for multi-platform builds
- Automatically determines the correct Neovim installation method (source vs package)
- Manages image tagging according to the new naming convention

### Local Building

To build images locally:

```bash
# Build the Go program
go build -o docker-neovim-builder main.go

# Build nightly image for Alpine
./docker-neovim-builder -type nightly -base alpine -version master

# Build latest stable for all base images
./docker-neovim-builder -type latest -base alpine
./docker-neovim-builder -type latest -base bookworm
./docker-neovim-builder -type latest -base bullseye

# Build specific version
./docker-neovim-builder -type package -base alpine -version 0.10.0
```

## Releases

You can find every release here: https://hub.docker.com/r/anatolelucet/neovim/tags
