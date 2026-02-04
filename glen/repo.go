package glen

import (
	"errors"
	"fmt"
	"path"
	"strings"

	git "github.com/go-git/go-git/v5"
)

// ErrInvalidRemoteURL is returned when a remote URL cannot be parsed.
var ErrInvalidRemoteURL = errors.New("invalid remote URL")

// Repo represents information about a git repo.
// Repo does not represent ALL information about a repo, only the information
// needed for this package (for gathering GitLab variables).
type Repo struct {
	LocalPath  string
	RemoteName string

	Path      string
	BaseURL   string
	HTTPURL   string
	RemoteURL string

	Groups []string
}

// NewRepo creates a new repo struct with defaults that assume you have a remote named
// 'origin' and that you are calling this function while your current directory is the
// repo you're interested in. If you have a custom local path or your remote is not
// named 'origin' then make sure you set those before you Init() the repo.
func NewRepo() *Repo {
	r := &Repo{}

	r.LocalPath = "."
	r.RemoteName = "origin"

	return r
}

// Init gathers information about the repo struct, populating all required fields.
func (r *Repo) Init() error {
	// We get all needed information about the repo based on the remote url
	remoteURL, err := getRemoteFromLocalRepoPath(r.LocalPath, r.RemoteName)
	if err != nil {
		return err
	}

	r.RemoteURL = remoteURL

	// Parse the remote url into needed information
	baseURL, repoPath, httpURL, err := ParseRemoteURL(remoteURL)
	if err != nil {
		return fmt.Errorf("your remote (%s), %s, is not an SSH or HTTP remote: %w", r.RemoteName, remoteURL, err)
	}

	r.BaseURL = baseURL
	r.Path = repoPath
	r.HTTPURL = httpURL

	// Build the list of parent groups
	r.Groups = ExtractGroups(repoPath)

	return nil
}

// ParseRemoteURL parses a git remote URL and extracts GitLab components.
// It supports HTTPS, HTTP, and SSH URL formats.
// Returns baseURL (e.g., "gitlab.com"), repoPath (e.g., "group/project"),
// and httpURL which is the host+path without protocol (e.g., "gitlab.com/group/project").
func ParseRemoteURL(remoteURL string) (baseURL, repoPath, httpURL string, err error) {
	remote := strings.TrimSpace(remoteURL)
	if remote == "" {
		return "", "", "", ErrInvalidRemoteURL
	}

	switch {
	case strings.Contains(remote, "://"):
		// Handle http://, https://, ssh:// URLs
		remote = strings.TrimPrefix(remote, "http://")
		remote = strings.TrimPrefix(remote, "https://")
		remote = strings.TrimPrefix(remote, "ssh://git@")
		remote = strings.TrimSuffix(remote, ".git")
		parts := strings.SplitN(remote, "/", 2) //nolint:mnd
		if len(parts) < 2 || parts[1] == "" {
			return "", "", "", ErrInvalidRemoteURL
		}

		baseURL = parts[0]
		repoPath = parts[1]
		httpURL = remote

	case strings.Contains(remote, "@"):
		// Handle git@host:path format
		remote = strings.TrimPrefix(remote, "git@")
		remote = strings.TrimSuffix(remote, ".git")
		parts := strings.SplitN(remote, ":", 2) //nolint:mnd
		if len(parts) < 2 || parts[1] == "" {
			return "", "", "", ErrInvalidRemoteURL
		}

		baseURL = parts[0]
		repoPath = parts[1]
		httpURL = parts[0] + "/" + parts[1]

	default:
		return "", "", "", ErrInvalidRemoteURL
	}

	return baseURL, repoPath, httpURL, nil
}

// ExtractGroups returns the list of parent group paths for a given project path.
// Groups are ordered from the immediate parent to the root group.
// For example, "group1/subgroup/project" returns ["group1/subgroup", "group1"].
func ExtractGroups(projectPath string) []string {
	var groups []string

	groupPath := path.Dir(projectPath)
	if groupPath == "." || groupPath == "" {
		return groups
	}

	groupCount := strings.Count(groupPath, "/") + 1
	for i := groupCount; i > 0; i-- {
		groups = append(groups, groupPath)
		groupPath = path.Dir(groupPath)
	}

	return groups
}

func getRemoteFromLocalRepoPath(path string, remote string) (string, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return "", fmt.Errorf("unable to open git repository (%s) with the following error: %w", path, err)
	}
	rm, err := r.Remote(remote)
	if err != nil {
		return "", fmt.Errorf("unable to find selected remote (%s) with the following error: %w", remote, err)
	}

	urls := rm.Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("remote (%s) has no configured URLs", remote)
	}

	return urls[0], nil
}
