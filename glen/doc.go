/*
Package glen provides an API to get GitLab project and group variables

Glen has a VERY simple API with two structs that combined can get GitLab environment variables from
a local repo with a GitLab remote.

Repo

The Repo struct holds information about a local repo and can be passed to a Variables struct. By
default Repo assumes that your current working directory has a remote named 'origin'. However, you
can specify custom directtories and custom origin names.

Variables

The Variables struct holds information about a GitLab project's CI/CD variables. By default
Variables collects a project's variables using your GITLAB_TOKEN environment variable API key.
However, you can set a custom API key and you can set Variables.Recurse=true to also collect the
group variables of the project's parent groups. Variables are merged according to the GitLab
specified precedence here:
https://docs.gitlab.com/ee/ci/variables/#priority-of-environment-variables
*/
package glen
