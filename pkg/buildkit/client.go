package buildkit

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BuildkitClient provides an interface for building Docker images
// This is a simplified implementation that uses buildkit CLI commands
// TODO: Replace with native buildkit Go client when dependency issues are resolved
type BuildkitClient struct {
	// Future: will contain actual buildkit client
	useDockerBuildx bool
}

// NewClient creates a new buildkit client
func NewClient() (*BuildkitClient, error) {
	// Check if buildkit is available, fall back to docker buildx
	return &BuildkitClient{
		useDockerBuildx: true, // For now, use docker buildx as buildkit frontend
	}, nil
}

// Close cleans up the client
func (c *BuildkitClient) Close() error {
	return nil
}

// Build builds a Docker image using buildkit
func (c *BuildkitClient) Build(ctx context.Context, opts BuildOptions) error {
	// Create a temporary directory for the build context
	tmpDir, err := ioutil.TempDir("", "buildkit-build")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write Dockerfile to temp directory
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	err = ioutil.WriteFile(dockerfilePath, []byte(opts.Dockerfile), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Dockerfile: %v", err)
	}

	// Use buildkit via docker buildx (which uses buildkit as backend)
	args := []string{"buildx", "build"}
	
	// Add platforms if specified
	if len(opts.Platforms) > 0 {
		args = append(args, "--platform", strings.Join(opts.Platforms, ","))
	}

	// Add output configuration
	if opts.Push {
		args = append(args, "--push")
	} else {
		args = append(args, "--load")
	}

	// Add tags
	for _, tag := range opts.Tags {
		args = append(args, "-t", tag)
	}

	// Add dockerfile path
	args = append(args, "-f", dockerfilePath)

	// Add build context
	args = append(args, tmpDir)

	// Execute the build command
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("buildkit build failed: %v", err)
	}

	return nil
}

// BuildOptions contains options for building an image
type BuildOptions struct {
	Dockerfile string   // Content of the Dockerfile
	Tags       []string // Image tags
	Platforms  []string // Target platforms
	Push       bool     // Whether to push the image
}