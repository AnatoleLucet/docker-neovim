package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func Login(username, password string) error {
	if password == "" {
		return nil // No password provided, skip login
	}

	cmd := exec.Command("docker", "login", "-u", username, "--password-stdin")
	cmd.Stdin = strings.NewReader(password)
	return cmd.Run()
}

func ImageExists(imageName string) bool {
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

func GetLatestNeovimVersion() (string, error) {
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