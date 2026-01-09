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
		{
			name:        "ssh:// protocol format with port",
			remoteURL:   "ssh://git@gitlab.example.com:2222/group/project.git",
			wantBaseURL: "gitlab.example.com:2222",
			wantPath:    "group/project",
			wantHTTPURL: "gitlab.example.com:2222/group/project",
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
				require.ErrorIs(t, err, ErrInvalidRemoteURL)

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

	tests := []struct {
		name        string
		remoteName  string
		setupRepo   func(t *testing.T) string // returns path to test directory
		wantURL     string
		wantErr     bool
		errContains string
	}{
		{
			name:       "reads origin remote URL",
			remoteName: "origin",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				repo, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)
				_, err = repo.CreateRemote(&config.RemoteConfig{
					Name: "origin",
					URLs: []string{"git@gitlab.com:testgroup/testproject.git"},
				})
				require.NoError(t, err)

				return tmpDir
			},
			wantURL: "git@gitlab.com:testgroup/testproject.git",
			wantErr: false,
		},
		{
			name:       "reads custom remote name",
			remoteName: "upstream",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				repo, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)
				_, err = repo.CreateRemote(&config.RemoteConfig{
					Name: "upstream",
					URLs: []string{"https://gitlab.example.com/org/repo.git"},
				})
				require.NoError(t, err)

				return tmpDir
			},
			wantURL: "https://gitlab.example.com/org/repo.git",
			wantErr: false,
		},
		{
			name:       "returns first URL when multiple URLs configured",
			remoteName: "origin",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				repo, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)
				_, err = repo.CreateRemote(&config.RemoteConfig{
					Name: "origin",
					URLs: []string{
						"git@gitlab.com:group/project.git",
						"https://gitlab.com/group/project.git",
					},
				})
				require.NoError(t, err)

				return tmpDir
			},
			wantURL: "git@gitlab.com:group/project.git",
			wantErr: false,
		},
		{
			name:       "error when remote does not exist",
			remoteName: "nonexistent",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				_, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)

				return tmpDir
			},
			wantErr:     true,
			errContains: "unable to find selected remote",
		},
		{
			name:       "error when path is not a git repo",
			remoteName: "origin",
			setupRepo: func(t *testing.T) string {
				t.Helper()

				return t.TempDir()
			},
			wantErr:     true,
			errContains: "unable to open git repository",
		},
		{
			name:       "error when path does not exist",
			remoteName: "origin",
			setupRepo: func(_ *testing.T) string {
				return "/nonexistent/path/to/repo"
			},
			wantErr:     true,
			errContains: "unable to open git repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.setupRepo(t)
			remoteURL, err := getRemoteFromLocalRepoPath(path, tt.remoteName)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantURL, remoteURL)
		})
	}
}

// TestRepoInit tests the full Init() flow with real git repositories.
func TestRepoInit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		remoteName    string
		setupRepo     func(t *testing.T) string // returns path to test directory
		wantBaseURL   string
		wantPath      string
		wantHTTPURL   string
		wantRemoteURL string
		wantGroups    []string
		wantErr       bool
		errContains   string
	}{
		{
			name:       "initializes repo with SSH remote",
			remoteName: "origin",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				repo, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)
				_, err = repo.CreateRemote(&config.RemoteConfig{
					Name: "origin",
					URLs: []string{"git@gitlab.com:myorg/myteam/myproject.git"},
				})
				require.NoError(t, err)

				return tmpDir
			},
			wantBaseURL:   "gitlab.com",
			wantPath:      "myorg/myteam/myproject",
			wantHTTPURL:   "gitlab.com/myorg/myteam/myproject",
			wantRemoteURL: "git@gitlab.com:myorg/myteam/myproject.git",
			wantGroups:    []string{"myorg/myteam", "myorg"},
		},
		{
			name:       "initializes repo with HTTPS remote",
			remoteName: "origin",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				repo, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)
				_, err = repo.CreateRemote(&config.RemoteConfig{
					Name: "origin",
					URLs: []string{"https://gitlab.example.com/company/product.git"},
				})
				require.NoError(t, err)

				return tmpDir
			},
			wantBaseURL:   "gitlab.example.com",
			wantPath:      "company/product",
			wantHTTPURL:   "gitlab.example.com/company/product",
			wantRemoteURL: "https://gitlab.example.com/company/product.git",
			wantGroups:    []string{"company"},
		},
		{
			name:       "initializes repo with custom remote name",
			remoteName: "gitlab",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				repo, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)
				_, err = repo.CreateRemote(&config.RemoteConfig{
					Name: "gitlab",
					URLs: []string{"git@gitlab.com:team/project.git"},
				})
				require.NoError(t, err)

				return tmpDir
			},
			wantBaseURL:   "gitlab.com",
			wantPath:      "team/project",
			wantHTTPURL:   "gitlab.com/team/project",
			wantRemoteURL: "git@gitlab.com:team/project.git",
			wantGroups:    []string{"team"},
		},
		{
			name:       "error when git repo does not exist",
			remoteName: "origin",
			setupRepo: func(_ *testing.T) string {
				return "/nonexistent/path"
			},
			wantErr:     true,
			errContains: "unable to open git repository",
		},
		{
			name:       "error when remote does not exist",
			remoteName: "nonexistent",
			setupRepo: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				_, err := git.PlainInit(tmpDir, false)
				require.NoError(t, err)

				return tmpDir
			},
			wantErr:     true,
			errContains: "unable to find selected remote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.setupRepo(t)
			repo := &Repo{
				LocalPath:  path,
				RemoteName: tt.remoteName,
			}

			err := repo.Init()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantBaseURL, repo.BaseURL)
			assert.Equal(t, tt.wantPath, repo.Path)
			assert.Equal(t, tt.wantHTTPURL, repo.HTTPURL)
			assert.Equal(t, tt.wantRemoteURL, repo.RemoteURL)
			assert.Equal(t, tt.wantGroups, repo.Groups)
		})
	}
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
