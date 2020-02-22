package glen

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
)

// Variables represents a set of CI/CD environment variables and
// the repo that those variables were collected from
type Variables struct {
	Env     map[string]string
	Recurse bool
	Repo    *Repo

	apiKey string
}

// NewVariables takes a *Repo and returns an empty Variables struct. This
// assumes that you have a GitLab API key set as GITLAB_TOKEN. If not, make
// sure you set one with Variables.SetAPIKey(). Note that by default we do not
// 'recurse' and get the group variables. If you want group variables merged
// in make sure you set Variables.Recurse=true
func NewVariables(r *Repo) *Variables {
	v := &Variables{}

	v.Env = make(map[string]string)
	v.Recurse = false
	v.Repo = r
	v.apiKey = os.Getenv("GITLAB_TOKEN")

	return v
}

// SetAPIKey takes a GitLab API key and adds it to the Variables struct. Note
// that if you used NewVariables() to create your struct and you had GITLAB_TOKEN
// exported then that token is already set and you don't need to call this function.
func (v *Variables) SetAPIKey(key string) {
	v.apiKey = key
}

// Init collects GitLab variables from the repo, and optionally from the parent groups
// if Variables.Recurse=true. Variable precedence respects
// https://docs.gitlab.com/ee/ci/variables/#priority-of-environment-variables
func (v *Variables) Init() error {
	var err error

	// Initialize the GitLab client
	glc := gitlab.NewClient(nil, v.apiKey)
	err = glc.SetBaseURL("https://" + v.Repo.BaseURL + "/api/v4")
	if err != nil {
		return fmt.Errorf("failed to set gitlab client base URL: %w", err)
	}

	// Get variables from the parent groups, if recurse
	if v.Recurse {
		for _, group := range v.Repo.Groups {
			gvs, _, err := glc.GroupVariables.ListVariables(group)
			if err != nil {
				return fmt.Errorf("failed to get variables from group %s: %w", group, err)
			}
			for _, gv := range gvs {
				v.Env[gv.Key] = gv.Value
			}
		}
	}

	// Get the project variables and add them to v.Env
	pvs, _, err := glc.ProjectVariables.ListVariables(v.Repo.Path)
	if err != nil {
		return fmt.Errorf("failed to get variables from project %s: %w", v.Repo.Path, err)
	}
	for _, pv := range pvs {
		v.Env[pv.Key] = pv.Value
	}

	return err
}
