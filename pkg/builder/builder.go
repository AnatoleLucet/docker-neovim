package builder

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Builder struct {
	// For future buildkit integration
	useLocalDocker bool
}

func New() (*Builder, error) {
	// For now, use local docker commands
	// TODO: Integrate with buildkit client when dependency issues are resolved
	return &Builder{useLocalDocker: true}, nil
}

func (b *Builder) Close() error {
	// No cleanup needed for docker commands
	return nil
}

func (b *Builder) BuildAndPush(ctx context.Context, dockerfile, imageName string, platforms []string) error {
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

	// Build for multiple platforms using docker buildx
	platformsStr := strings.Join(platforms, ",")

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