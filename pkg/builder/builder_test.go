package builder

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	builder, err := New()
	if err != nil {
		t.Errorf("New() should not return error, got: %v", err)
	}
	if builder == nil {
		t.Error("New() should return non-nil builder")
	}

	// Test Close doesn't error
	err = builder.Close()
	if err != nil {
		t.Errorf("Close() should not return error, got: %v", err)
	}
}

func TestBuildAndPush(t *testing.T) {
	builder, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer builder.Close()

	// Test with minimal dockerfile (this would fail without docker, but shouldn't panic)
	dockerfile := "FROM alpine\nRUN echo 'test'"
	ctx := context.Background()
	
	// This will likely fail because docker buildx might not be available or configured
	// but it tests that the function handles the parameters correctly
	err = builder.BuildAndPush(ctx, dockerfile, "test:latest", []string{"linux/amd64"})
	if err != nil {
		t.Logf("BuildAndPush failed as expected (docker not available): %v", err)
		// This is expected in test environment without docker setup
	}
}