# Glen

[![Go Report Card](https://goreportcard.com/badge/github.com/lingrino/glen)](https://goreportcard.com/report/github.com/lingrino/glen)
[![GoDoc](https://godoc.org/github.com/lingrino/glen/glen?status.svg)](https://godoc.org/github.com/lingrino/glen/glen)

Glen is a simple command line tool that, when run within a GitLab project, will call the
GitLab API to get all environment variables from your project's CI/CD pipeline and print
them locally, ready for exporting.

With the default flags you can run 'eval $(glen -r)' to export your project's variables
and the variables of every parent group.

Note that glen requires that you have 'git' installed locally and in your PATH.

## Installation

```text
go get github.com/lingrino/glen
```

## Usage

```text
$ glen --help
Glen is a simple command line tool that, when run within a GitLab project,
will call the GitLab API to get all environment variables from your project's
CI/CD pipeline and print them locally, ready for exporting.

With the default flags you can run 'eval $(glen -r)' to export your project's
variables and the variables of every parent group.

Note that glen requires that you have 'git' installed locally and in your PATH.

Usage:
  glen [flags]

Flags:
  -k, --api-key string       Your GitLab API key. NOTE - It's preferrable to specify your key as a GITLAB_TOKEN environment variable (default "GITLAB_TOKEN")
  -d, --directory string     The directory where you're git repo lives. Defaults to your current working directory (default ".")
  -h, --help                 help for glen
  -o, --output string        The output format. One of 'export', 'json', 'table'. Defaults to 'export', which can be executed to export all variables. (default "export")
  -r, --recurse              Set recurse to true if you want to include the variables of the parent groups
  -n, --remote-name string   The name of the GitLab remote in your git repo. Defaults to 'origin'. Check with 'git remote -v' (default "origin")
```

## Improvements

Glen is working and ready to use. Glen will be tagged `v1.0.0` after the following are complete:

- [ ] Tests (coverage badge in README)
- [ ] Detailed README/Docs
- [ ] Automated Releases
- [ ] Mirroring to GitLab
- [ ] CI/CD on GitLab

Issues or PRs are welcome for any functionality you would like to see!
