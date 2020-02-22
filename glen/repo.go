package glen

import (
	"fmt"
	"path"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

// Repo represents information about a git repo
// Repo does not represent ALL information about a repo, only the information
// needed for this package (for gathering GitLab variables)
type Repo struct {
	IsSSH  bool
	IsHTTP bool

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
// named 'origin' then make sure you set those before you Init() the repo
func NewRepo() *Repo {
	r := &Repo{}

	r.LocalPath = "."
	r.RemoteName = "origin"

	return r
}

// Init gathers information about the repo struct, populating all required fields
func (r *Repo) Init() error {
	// We get all needed information about the repo based on the remote url
	remote, err := getRemoteFromLocalRepoPath(r.LocalPath, r.RemoteName)

	// Parse the remote url into needed information, different paths for ssh vs http remotes
	switch {
	case strings.Contains(remote, "@"):
		r.IsSSH = true
		r.IsHTTP = false

		remote = strings.TrimSpace(remote)
		remote = strings.TrimPrefix(remote, "git@")
		remote = strings.TrimSuffix(remote, ".git")
		remoteS := strings.Split(remote, ":")

		r.HTTPURL = remoteS[0] + "/" + remoteS[1]
		r.BaseURL = remoteS[0]
		r.Path = remoteS[1]
	case strings.Contains(remote, "://"):
		r.IsSSH = false
		r.IsHTTP = true

		remote = strings.TrimSpace(remote)
		remote = strings.TrimPrefix(remote, "http://")
		remote = strings.TrimPrefix(remote, "https://")
		remote = strings.TrimSuffix(remote, ".git")
		remoteS := strings.SplitN(remote, "/", 2)

		r.HTTPURL = remote
		r.BaseURL = remoteS[0]
		r.Path = remoteS[1]
	default:
		return fmt.Errorf("your remote (%s), %s, is not an SSH or HTTP remote", r.RemoteName, remote)
	}

	// We create a list of gitlab groups that we can collect variables from. These are
	// just the paths before the repo name.
	groups := path.Dir(r.Path)
	groupsN := strings.Count(groups, "/") + 1
	for i := groupsN; i > 0; i-- {
		r.Groups = append(r.Groups, groups)
		groups = path.Dir(groups)
	}

	return err
}

func getRemoteFromLocalRepoPath(path string, remote string) (string, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return "", fmt.Errorf("unable to open git repository (%s) with the following error: %s", path, err)
	}
	rm, err := r.Remote(remote)
	if err != nil {
		return "", fmt.Errorf("unable to find selected remote (%s) with the following error: %s", remote, err)
	}

	firstURL := rm.Config().URLs[0]

	return firstURL, nil
}
