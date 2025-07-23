package builder

import (
	"context"

	"github.com/AnatoleLucet/docker-neovim/pkg/buildkit"
)

type Builder struct {
	buildkitClient *buildkit.BuildkitClient
}

func New() (*Builder, error) {
	client, err := buildkit.NewClient()
	if err != nil {
		return nil, err
	}

	return &Builder{buildkitClient: client}, nil
}

func (b *Builder) Close() error {
	return b.buildkitClient.Close()
}

func (b *Builder) BuildAndPush(ctx context.Context, dockerfile, imageName string, platforms []string) error {
	opts := buildkit.BuildOptions{
		Dockerfile: dockerfile,
		Tags:       []string{imageName},
		Platforms:  platforms,
		Push:       true,
	}

	return b.buildkitClient.Build(ctx, opts)
}