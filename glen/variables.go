package glen

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
)

// Variables represents a set of CI/CD environment variables and
// the repo that those variables were collected from.
type Variables struct {
	Env     map[string]string
	Recurse bool
	Repo    *Repo

	apiKey string
}

// Gitlab list apis default to 20 per page with 100 being the max.
// https://docs.gitlab.com/ee/api/README.html#offset-based-pagination
// Set to max to reduce number of API calls required.
const pageSize = 100

// NewVariables takes a *Repo and returns an empty Variables struct. This
// assumes that you have a GitLab API key set as GITLAB_TOKEN. If not, make
// sure you set one with Variables.SetAPIKey(). Note that by default we do not
// 'recurse' and get the group variables. If you want group variables merged
// in make sure you set Variables.Recurse=true.
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

// Get the group variables and add them to v.Env.
func (v *Variables) getGroupVariables(glc *gitlab.Client, group string) error {
	groupVariablesOpt := &gitlab.ListGroupVariablesOptions{
		PerPage: pageSize,
		Page:    1,
	}

	for {
		gvs, response, err := glc.GroupVariables.ListVariables(group, groupVariablesOpt)
		if err != nil {
			return fmt.Errorf("failed to get variables from group %s: %w", group, err)
		}
		for _, gv := range gvs {
			v.Env[gv.Key] = gv.Value
		}

		if response.CurrentPage >= response.TotalPages {
			break
		}

		groupVariablesOpt.Page = response.NextPage
	}

	return nil
}

// Get the project variables and add them to v.Env.
func (v *Variables) getProjectVariables(glc *gitlab.Client) error {
	projectVariablesOpt := &gitlab.ListProjectVariablesOptions{
		PerPage: pageSize,
		Page:    1,
	}

	for {
		pvs, response, err := glc.ProjectVariables.ListVariables(v.Repo.Path, projectVariablesOpt)
		if err != nil {
			return fmt.Errorf("failed to get variables from project %s: %w", v.Repo.Path, err)
		}
		for _, pv := range pvs {
			v.Env[pv.Key] = pv.Value
		}

		if response.CurrentPage >= response.TotalPages {
			break
		}

		projectVariablesOpt.Page = response.NextPage
	}

	return nil
}

// Init collects GitLab variables from the repo, and optionally from the parent groups
// if Variables.Recurse=true. Variable precedence respects
// https://docs.gitlab.com/ee/ci/variables/#priority-of-environment-variables
func (v *Variables) Init() error {
	var err error

	// Initialize the GitLab client
	glURL := fmt.Sprintf("https://%s/api/v4", v.Repo.BaseURL)
	glc, err := gitlab.NewClient(v.apiKey, gitlab.WithBaseURL(glURL))
	if err != nil {
		return fmt.Errorf("failed to create gitlab client: %w", err)
	}

	// Get variables from the parent groups, if recurse
	if v.Recurse {
		for _, group := range v.Repo.Groups {
			//nolint:errcheck
			v.getGroupVariables(glc, group)
		}
	}

	//nolint:errcheck
	v.getProjectVariables(glc)

	return err
}
