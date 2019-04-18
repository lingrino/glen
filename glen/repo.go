package glen

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
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
	var err error

	// Make sure git is installed
	_, err = exec.LookPath("git")
	if err != nil {
		return errors.Wrapf(err, "'git' is not installed and in your path")
	}

	// Get information about the GitLab remote
	remoteOUT, err := exec.Command("git", "-C", r.LocalPath, "remote", "get-url", r.RemoteName).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Failed to get info about the git repo remote: %s", remoteOUT)
	}

	// We get all needed information about the repo based on the remote url
	remote := string(remoteOUT)
	r.RemoteURL = remote

	// Parse the remote url into needed information, different paths for ssh vs http remotes
	if strings.Contains(remote, "@") {
		r.IsSSH = true
		r.IsHTTP = false

		remote = strings.TrimSpace(remote)
		remote = strings.TrimPrefix(remote, "git@")
		remote = strings.TrimSuffix(remote, ".git")
		remoteS := strings.Split(remote, ":")

		r.HTTPURL = remoteS[0] + "/" + remoteS[1]
		r.BaseURL = remoteS[0]
		r.Path = remoteS[1]
	} else if strings.Contains(remote, "://") {
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
	} else {
		return fmt.Errorf("Your remote (%s), %s, is not an SSH or HTTP remote", r.RemoteName, remote)
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
