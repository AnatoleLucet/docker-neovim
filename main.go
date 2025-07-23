package main

import (
	"context"
	"flag"
	"log"

	"github.com/AnatoleLucet/docker-neovim/pkg/builder"
	"github.com/AnatoleLucet/docker-neovim/pkg/config"
	"github.com/AnatoleLucet/docker-neovim/pkg/dockerfile"
	"github.com/AnatoleLucet/docker-neovim/pkg/registry"
)

func main() {
	var (
		configPath = flag.String("config", "build-config.json", "Path to build configuration file")
		version    = flag.String("version", "", "Neovim version to build")
		buildType  = flag.String("type", "package", "Build type: 'nightly' or 'package'")
		baseImage  = flag.String("base", "alpine", "Base image: 'alpine', 'bookworm', or 'bullseye'")
		tag        = flag.String("tag", "", "Custom tag for the image")
	)
	flag.Parse()

	cfg := config.Load(*configPath)

	// Override config with command line args if provided
	if *version != "" {
		cfg.Version = *version
	}
	if *buildType != "" {
		cfg.BuildType = *buildType
	}
	if *baseImage != "" {
		cfg.BaseImages = []string{*baseImage}
	}

	// Login to Docker registry
	err := registry.Login(cfg.Username, cfg.Password)
	if err != nil {
		log.Fatalf("Failed to login to Docker registry: %v", err)
	}

	// Create builder
	b, err := builder.New()
	if err != nil {
		log.Fatalf("Failed to create builder: %v", err)
	}
	defer b.Close()

	ctx := context.Background()

	for _, base := range cfg.BaseImages {
		imageName := cfg.GenerateImageName(base, *tag)

		if !cfg.AllowOverride && registry.ImageExists(imageName) {
			log.Printf("Image %s already exists, skipping", imageName)
			continue
		}

		log.Printf("Building image: %s", imageName)

		dockerfileContent := dockerfile.Generate(cfg.BuildType, base, cfg.Version)

		err := b.BuildAndPush(ctx, dockerfileContent, imageName, cfg.Platforms)
		if err != nil {
			log.Fatalf("Failed to build image %s: %v", imageName, err)
		}

		log.Printf("Successfully built and pushed: %s", imageName)
	}
}

