package glen

import (
	"os"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	t.Parallel()

	expected := &Repo{
		LocalPath:  ".",
		RemoteName: "origin",
	}
	repo := NewRepo()

	assert.Equal(t, expected, repo)
}

func TestParseRemoteURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		remoteURL   string
		wantBaseURL string
		wantPath    string
		wantHTTPURL string
		wantErr     bool
	}{
		// HTTPS URLs
		{
			name:        "https with .git suffix",
			remoteURL:   "https://gitlab.com/group/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},
		{
			name:        "https without .git suffix",
			remoteURL:   "https://gitlab.com/group/project",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},
		{
			name:        "https with subgroups",
			remoteURL:   "https://gitlab.com/group/subgroup/subsubgroup/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/subgroup/subsubgroup/project",
			wantHTTPURL: "gitlab.com/group/subgroup/subsubgroup/project",
			wantErr:     false,
		},
		{
			name:        "https self-hosted gitlab",
			remoteURL:   "https://gitlab.example.com/org/repo.git",
			wantBaseURL: "gitlab.example.com",
			wantPath:    "org/repo",
			wantHTTPURL: "gitlab.example.com/org/repo",
			wantErr:     false,
		},
		{
			name:        "https with port",
			remoteURL:   "https://gitlab.example.com:8443/group/project.git",
			wantBaseURL: "gitlab.example.com:8443",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.example.com:8443/group/project",
			wantErr:     false,
		},

		// HTTP URLs
		{
			name:        "http with .git suffix",
			remoteURL:   "http://gitlab.com/group/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},

		// SSH URLs (git@host:path format)
		{
			name:        "ssh git@ format with .git suffix",
			remoteURL:   "git@gitlab.com:group/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},
		{
			name:        "ssh git@ format without .git suffix",
			remoteURL:   "git@gitlab.com:group/project",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},
		{
			name:        "ssh git@ format with subgroups",
			remoteURL:   "git@gitlab.com:group/subgroup/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/subgroup/project",
			wantHTTPURL: "gitlab.com/group/subgroup/project",
			wantErr:     false,
		},
		{
			name:        "ssh git@ format self-hosted",
			remoteURL:   "git@gitlab.example.com:org/repo.git",
			wantBaseURL: "gitlab.example.com",
			wantPath:    "org/repo",
			wantHTTPURL: "gitlab.example.com/org/repo",
			wantErr:     false,
		},

		// SSH URLs (ssh:// protocol format)
		{
			name:        "ssh:// protocol format",
			remoteURL:   "ssh://git@gitlab.com/group/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},
		{
			name:        "ssh:// protocol format with subgroups",
			remoteURL:   "ssh://git@gitlab.com/org/team/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "org/team/project",
			wantHTTPURL: "gitlab.com/org/team/project",
			wantErr:     false,
		},

		// Edge cases - whitespace handling
		{
			name:        "url with leading whitespace",
			remoteURL:   "  https://gitlab.com/group/project.git",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},
		{
			name:        "url with trailing whitespace",
			remoteURL:   "git@gitlab.com:group/project.git  ",
			wantBaseURL: "gitlab.com",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.com/group/project",
			wantErr:     false,
		},

		// Error cases
		{
			name:      "empty string",
			remoteURL: "",
			wantErr:   true,
		},
		{
			name:      "whitespace only",
			remoteURL: "   ",
			wantErr:   true,
		},
		{
			name:      "invalid url - no protocol or @",
			remoteURL: "gitlab.com/group/project",
			wantErr:   true,
		},
		{
			name:      "https url without path",
			remoteURL: "https://gitlab.com/",
			wantErr:   true,
		},
		{
			name:      "https url with only domain",
			remoteURL: "https://gitlab.com",
			wantErr:   true,
		},
		{
			name:      "ssh url without path",
			remoteURL: "git@gitlab.com:",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			baseURL, repoPath, httpURL, err := ParseRemoteURL(tt.remoteURL)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidRemoteURL)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantBaseURL, baseURL, "baseURL mismatch")
			assert.Equal(t, tt.wantPath, repoPath, "path mismatch")
			assert.Equal(t, tt.wantHTTPURL, httpURL, "httpURL mismatch")
		})
	}
}

func TestExtractGroups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		projectPath string
		wantGroups  []string
	}{
		{
			name:        "single group",
			projectPath: "group/project",
			wantGroups:  []string{"group"},
		},
		{
			name:        "two levels of groups",
			projectPath: "group/subgroup/project",
			wantGroups:  []string{"group/subgroup", "group"},
		},
		{
			name:        "three levels of groups",
			projectPath: "org/team/subteam/project",
			wantGroups:  []string{"org/team/subteam", "org/team", "org"},
		},
		{
			name:        "deeply nested groups",
			projectPath: "a/b/c/d/e/project",
			wantGroups:  []string{"a/b/c/d/e", "a/b/c/d", "a/b/c", "a/b", "a"},
		},
		{
			name:        "no groups - just project name",
			projectPath: "project",
			wantGroups:  nil,
		},
		{
			name:        "empty path",
			projectPath: "",
			wantGroups:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			groups := ExtractGroups(tt.projectPath)
			assert.Equal(t, tt.wantGroups, groups)
		})
	}
}

// TestGetRemoteFromLocalRepoPath tests reading remote URLs from git repositories.
// These tests use go-git to create real git repos in temp directories.
// No system git binary is required - go-git is a pure Go implementation.
func TestGetRemoteFromLocalRepoPath(t *testing.T) {
	t.Parallel()

	t.Run("reads origin remote URL", func(t *testing.T) {
		t.Parallel()

		// Create temp directory for test repo
		tmpDir := t.TempDir()

		// Initialize git repo using go-git (pure Go, no git binary needed)
		repo, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err, "failed to init git repo")

		// Add a remote
		expectedURL := "git@gitlab.com:testgroup/testproject.git"
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{expectedURL},
		})
		require.NoError(t, err, "failed to create remote")

		// Test getRemoteFromLocalRepoPath
		remoteURL, err := getRemoteFromLocalRepoPath(tmpDir, "origin")
		require.NoError(t, err)
		assert.Equal(t, expectedURL, remoteURL)
	})

	t.Run("reads custom remote name", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		repo, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err)

		expectedURL := "https://gitlab.example.com/org/repo.git"
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "upstream",
			URLs: []string{expectedURL},
		})
		require.NoError(t, err)

		remoteURL, err := getRemoteFromLocalRepoPath(tmpDir, "upstream")
		require.NoError(t, err)
		assert.Equal(t, expectedURL, remoteURL)
	})

	t.Run("returns first URL when multiple URLs configured", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		repo, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err)

		firstURL := "git@gitlab.com:group/project.git"
		secondURL := "https://gitlab.com/group/project.git"
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{firstURL, secondURL},
		})
		require.NoError(t, err)

		remoteURL, err := getRemoteFromLocalRepoPath(tmpDir, "origin")
		require.NoError(t, err)
		assert.Equal(t, firstURL, remoteURL)
	})

	t.Run("error when remote does not exist", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		_, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err)

		_, err = getRemoteFromLocalRepoPath(tmpDir, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to find selected remote")
	})

	t.Run("error when path is not a git repo", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		_, err := getRemoteFromLocalRepoPath(tmpDir, "origin")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to open git repository")
	})

	t.Run("error when path does not exist", func(t *testing.T) {
		t.Parallel()

		_, err := getRemoteFromLocalRepoPath("/nonexistent/path/to/repo", "origin")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to open git repository")
	})
}

// TestRepoInit tests the full Init() flow with real git repositories.
func TestRepoInit(t *testing.T) {
	t.Parallel()

	t.Run("initializes repo with SSH remote", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		gitRepo, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err)

		_, err = gitRepo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{"git@gitlab.com:myorg/myteam/myproject.git"},
		})
		require.NoError(t, err)

		repo := &Repo{
			LocalPath:  tmpDir,
			RemoteName: "origin",
		}

		err = repo.Init()
		require.NoError(t, err)

		assert.Equal(t, "gitlab.com", repo.BaseURL)
		assert.Equal(t, "myorg/myteam/myproject", repo.Path)
		assert.Equal(t, "gitlab.com/myorg/myteam/myproject", repo.HTTPURL)
		assert.Equal(t, "git@gitlab.com:myorg/myteam/myproject.git", repo.RemoteURL)
		assert.Equal(t, []string{"myorg/myteam", "myorg"}, repo.Groups)
	})

	t.Run("initializes repo with HTTPS remote", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		gitRepo, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err)

		_, err = gitRepo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{"https://gitlab.example.com/company/product.git"},
		})
		require.NoError(t, err)

		repo := &Repo{
			LocalPath:  tmpDir,
			RemoteName: "origin",
		}

		err = repo.Init()
		require.NoError(t, err)

		assert.Equal(t, "gitlab.example.com", repo.BaseURL)
		assert.Equal(t, "company/product", repo.Path)
		assert.Equal(t, "gitlab.example.com/company/product", repo.HTTPURL)
		assert.Equal(t, []string{"company"}, repo.Groups)
	})

	t.Run("initializes repo with custom remote name", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		gitRepo, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err)

		_, err = gitRepo.CreateRemote(&config.RemoteConfig{
			Name: "gitlab",
			URLs: []string{"git@gitlab.com:team/project.git"},
		})
		require.NoError(t, err)

		repo := &Repo{
			LocalPath:  tmpDir,
			RemoteName: "gitlab",
		}

		err = repo.Init()
		require.NoError(t, err)

		assert.Equal(t, "gitlab.com", repo.BaseURL)
		assert.Equal(t, "team/project", repo.Path)
	})

	t.Run("error when git repo does not exist", func(t *testing.T) {
		t.Parallel()

		repo := &Repo{
			LocalPath:  "/nonexistent/path",
			RemoteName: "origin",
		}

		err := repo.Init()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to open git repository")
	})

	t.Run("error when remote does not exist", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		_, err := git.PlainInit(tmpDir, false)
		require.NoError(t, err)

		repo := &Repo{
			LocalPath:  tmpDir,
			RemoteName: "nonexistent",
		}

		err = repo.Init()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to find selected remote")
	})
}

// TestRepoInitFromSubdirectory verifies that Init() requires the repo root path.
// go-git's PlainOpen does not walk up directories like the git CLI does.
func TestRepoInitFromSubdirectory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	// Initialize git repo in the root
	gitRepo, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)

	_, err = gitRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"git@gitlab.com:group/project.git"},
	})
	require.NoError(t, err)

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "src", "pkg")
	err = os.MkdirAll(subDir, 0o750)
	require.NoError(t, err)

	// Init from subdirectory fails - go-git requires the repo root path
	repo := &Repo{
		LocalPath:  subDir,
		RemoteName: "origin",
	}

	err = repo.Init()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to open git repository")

	// But init from the repo root works
	repo = &Repo{
		LocalPath:  tmpDir,
		RemoteName: "origin",
	}

	err = repo.Init()
	require.NoError(t, err)
	assert.Equal(t, "gitlab.com", repo.BaseURL)
	assert.Equal(t, "group/project", repo.Path)
}
