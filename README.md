# Glen

[![Glen](assets/logo-md.png?raw=true)](assets/logo-md.png "Glen")

[![PkgGoDev](https://pkg.go.dev/badge/github.com/lingrino/glen/glen)](https://pkg.go.dev/github.com/lingrino/glen/glen)
[![goreportcard](https://goreportcard.com/badge/github.com/lingrino/glen)](https://goreportcard.com/report/github.com/lingrino/glen)

Glen is a simple command line tool that, when run within a local GitLab project, will call the GitLab API to get all environment variables from your project's CI/CD pipeline and print them locally, ready for exporting.

With the default flags you can run `eval $(glen -r)` to export the variables of your project and the variables of every parent group.

## Installation

The easiest way to install glen is with [homebrew][]

```console
brew install lingrino/tap/glen
```

Glen can also be installed using [asdf][]:

```console
asdf plugin-add glen
asdf install glen latest
asdf global glen latest
```

Glen can also be installed by downloading the latest binary from the releases page and adding it to your path.

Alternatively you can install glen using `go get`, assuming you have `$GOPATH/bin` in your path.

```console
go install github.com/lingrino/glen@latest
```

## Usage

By default glen assumes that you have a GitLab API key exported as `GITLAB_TOKEN` and that you are calling glen from within a git repo where the GitLab remote is named `origin` (see `git remote -v`).

You can override all of these settings, specifying the API key, git directory, or GitLab remote name as flags on the command line (see `glen --help`).

By default glen will only get the variables from the current GitLab project. If you would also like glen to merge in variables from all of the project's parent groups then you can use `glen -r`

Lastly, the default output for glen is called `export`, meaning that the output is ready to be read into your shell and will export all variables. This lets you call glen as `eval $(glen)` as a one line command to export all variables locally. You can also specify a `json` or `table` output for more machine or human friendly outputs.

```console
$ glen --help
Glen is a simple command line tool that, when run within a GitLab project,
will call the GitLab API to get all environment variables from your project's
CI/CD pipeline and print them locally, ready for exporting.

With the default flags you can run 'eval $(glen -r)' to export the variables of
your project and the variables of every parent group.

Usage:
  glen [flags]
  glen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Returns the current glen version

Flags:
  -k, --api-key string       Your GitLab API key, if not set as a GITLAB_TOKEN environment variable (default "GITLAB_TOKEN")
  -d, --directory string     The directory where your git repo lives. Defaults to your current working directory (default ".")
  -g, --group-only           Set group to true to get only variables from the parent groups.
  -h, --help                 Help for glen
  -o, --output string        One of 'export', 'json', 'table'. Default 'export', which can be executed to export variables (default "export")
  -r, --recurse              Set recurse to true if you want to include the variables of the parent groups
  -n, --remote-name string   Name of the GitLab remote in your git repo. Defaults to 'origin' (default "origin")

Use "glen [command] --help" for more information about a command.
```

## Contributing

Glen does one thing (reads variables from GitLab projects) and should do that one thing well. If you notice a bug with glen please file an issue or submit a PR.

Also, all contributions and ideas are welcome! Please submit an issue or a PR with anything that you think could be improved.

In particular, this project could benefit from the following:

- [ ] Tests that mock git repos

|                Contributors                |
| :----------------------------------------: |
| [@solidnerd](https://github.com/solidnerd) |
|    [@bradym](https://github.com/bradym)    |
|   [@mgonnav](https://github.com/mgonnav)   |

[homebrew]: https://brew.sh/
[asdf]: https://asdf-vm.com/
