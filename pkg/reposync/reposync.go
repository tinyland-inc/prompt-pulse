// Package reposync provides historical external-repository sync verification
// and CI compatibility utilities. These helpers preserve the old monorepo
// export model for analysis and bounded migration work, but they are not the
// canonical day-to-day contribution path now that prompt-pulse is maintained
// directly in its standalone GitHub repo.
package reposync

import (
	"fmt"
	"strings"
	"time"
)

// SyncConfig describes the mapping between a source monorepo path and a
// target standalone repository.
type SyncConfig struct {
	// SourceRepo is the full repository path of the source monorepo
	// (e.g., "github.com/tinyland-inc/lab").
	SourceRepo string

	// SourcePath is the subdirectory within SourceRepo containing the project.
	SourcePath string

	// TargetRepo is the standalone repository home that receives synced files.
	TargetRepo string

	// TargetModule is the Go module path preserved or written into the target
	// repository. This can intentionally differ from TargetRepo when repo-home
	// authority and module-path authority have not converged.
	TargetModule string

	// TargetBranch is the branch in TargetRepo that receives synced commits.
	TargetBranch string

	// SyncPaths lists paths (relative to SourcePath) to include in sync.
	SyncPaths []string

	// ExcludePaths lists paths (relative to SourcePath) to exclude from sync.
	ExcludePaths []string

	// CITemplate is the path to the historical CI template in the source repo
	// that used to drive synchronization.
	CITemplate string
}

// SyncStatus captures the current state of synchronization between source
// and target repositories.
type SyncStatus struct {
	// InSync is true when source and target are aligned.
	InSync bool

	// SourceCommit is the latest commit hash on the source side.
	SourceCommit string

	// TargetCommit is the latest commit hash on the target side.
	TargetCommit string

	// DriftFiles lists paths that differ between source and target.
	DriftFiles []string

	// LastSync records when the last successful sync occurred.
	LastSync time.Time
}

// DefaultConfig returns the standard sync configuration for prompt-pulse.
func DefaultConfig() *SyncConfig {
	return &SyncConfig{
		SourceRepo:   "github.com/tinyland-inc/lab",
		SourcePath:   "cmd/prompt-pulse/",
		TargetRepo:   "github.com/tinyland-inc/prompt-pulse",
		TargetModule: "gitlab.com/tinyland/lab/prompt-pulse",
		TargetBranch: "main",
		SyncPaths: []string{
			"pkg/",
			"go.mod",
			"go.sum",
			"main.go",
			"*.go",
			"docs/",
			"vendor/",
		},
		ExcludePaths: []string{
			"display/",
			"waifu/",
			"collectors/",
			"shell/",
			"config/",
			"cache/",
			"status/",
			"internal/",
			"cmd/",
			"scripts/",
			"tests/",
		},
		CITemplate: "ci/archive/sync-external.yml",
	}
}

// ValidateConfig checks a SyncConfig for common problems and returns a slice
// of human-readable error strings. An empty return means the config is valid.
func ValidateConfig(c *SyncConfig) []string {
	var errs []string

	if c == nil {
		return []string{"config is nil"}
	}

	if strings.TrimSpace(c.SourceRepo) == "" {
		errs = append(errs, "source_repo is required")
	}
	if strings.TrimSpace(c.SourcePath) == "" {
		errs = append(errs, "source_path is required")
	}
	if strings.TrimSpace(c.TargetRepo) == "" {
		errs = append(errs, "target_repo is required")
	}
	if strings.TrimSpace(c.TargetModule) == "" {
		errs = append(errs, "target_module is required")
	}
	if strings.TrimSpace(c.TargetBranch) == "" {
		errs = append(errs, "target_branch is required")
	}
	if len(c.SyncPaths) == 0 {
		errs = append(errs, "sync_paths must not be empty")
	}

	// Warn if source and target are the same repository.
	if c.SourceRepo != "" && c.SourceRepo == c.TargetRepo {
		errs = append(errs, "source_repo and target_repo must differ")
	}

	// Validate that exclude paths don't overlap with explicit sync paths.
	for _, ep := range c.ExcludePaths {
		for _, sp := range c.SyncPaths {
			if ep == sp {
				errs = append(errs, fmt.Sprintf("path %q appears in both sync_paths and exclude_paths", ep))
			}
		}
	}

	return errs
}
