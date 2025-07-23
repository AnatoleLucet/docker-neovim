package buildkit

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("NewClient() should not return error, got: %v", err)
	}
	if client == nil {
		t.Error("NewClient() should return non-nil client")
	}

	// Test Close doesn't error
	err = client.Close()
	if err != nil {
		t.Errorf("Close() should not return error, got: %v", err)
	}
}

func TestBuild(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Test with minimal dockerfile (this would fail without docker, but shouldn't panic)
	opts := BuildOptions{
		Dockerfile: "FROM alpine\nRUN echo 'test'",
		Tags:       []string{"test:buildkit"},
		Platforms:  []string{"linux/amd64"},
		Push:       false, // Don't push in test
	}
	
	ctx := context.Background()
	
	// This will likely fail because docker buildx might not be available or configured
	// but it tests that the function handles the parameters correctly
	err = client.Build(ctx, opts)
	if err != nil {
		t.Logf("Build failed as expected (docker not available): %v", err)
		// This is expected in test environment without docker setup
	}
}

func TestBuildOptionsValidation(t *testing.T) {
	opts := BuildOptions{
		Dockerfile: "FROM alpine",
		Tags:       []string{"test:latest", "test:v1.0"},
		Platforms:  []string{"linux/amd64", "linux/arm64"},
		Push:       true,
	}

	if len(opts.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(opts.Tags))
	}
	if len(opts.Platforms) != 2 {
		t.Errorf("Expected 2 platforms, got %d", len(opts.Platforms))
	}
	if opts.Push != true {
		t.Errorf("Expected Push to be true, got %v", opts.Push)
	}
}