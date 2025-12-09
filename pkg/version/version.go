// Package version provides version information for the Groundflow API.
// It embeds the VERSION file and provides utilities for accessing version data
// throughout the application.
package version

import (
	"context"
	_ "embed"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//go:embed VERSION
var versionFile string

const unknownValue = "unknown"

// BuildInfo contains version and build information
type BuildInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

// These variables can be set at build time using -ldflags
var (
	GitCommit = unknownValue
	BuildTime = unknownValue
)

// GetVersion reads the version from the embedded VERSION file
func GetVersion() string {
	return strings.TrimSpace(versionFile)
}

// GetGitCommit gets the current git commit hash
func GetGitCommit() string {
	if GitCommit != unknownValue {
		return GitCommit // Use build-time injected value if available
	}

	cmd := exec.CommandContext(context.Background(), "git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return unknownValue
	}
	return strings.TrimSpace(string(output))
}

// GetBuildTime returns the build time
func GetBuildTime() string {
	if BuildTime != unknownValue {
		return BuildTime // Use build-time injected value if available
	}
	return time.Now().UTC().Format(time.RFC3339)
}

// GetBuildInfo returns complete build information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   GetVersion(),
		GitCommit: GetGitCommit(),
		BuildTime: GetBuildTime(),
		GoVersion: runtime.Version(),
	}
}

// GetSemanticVersion returns a semantic version with git info for development builds
func GetSemanticVersion() string {
	version := GetVersion()

	// Check if we're on a tag
	cmd := exec.CommandContext(context.Background(), "git", "describe", "--tags", "--exact-match", "HEAD")
	if _, err := cmd.Output(); err == nil {
		return version // Clean version on tag
	}

	// Add commit info for non-tag builds
	gitCommit := GetGitCommit()
	if gitCommit != unknownValue {
		return fmt.Sprintf("%s-dev+%s", version, gitCommit)
	}

	return fmt.Sprintf("%s-dev", version)
}

// PrintVersion prints version information to stdout
func PrintVersion() {
	info := GetBuildInfo()
	fmt.Printf("Version:    %s\n", info.Version)
	fmt.Printf("Git Commit: %s\n", info.GitCommit)
	fmt.Printf("Build Time: %s\n", info.BuildTime)
	fmt.Printf("Go Version: %s\n", info.GoVersion)
}
