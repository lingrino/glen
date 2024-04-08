package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lingrino/glen/glen"
	"github.com/spf13/cobra"
)

const (
	flagRecurseDesc      = "Set recurse to true if you want to include the variables of the parent groups"
	flagAPIKeyDesc       = "Your GitLab API key, if not set as a GITLAB_TOKEN environment variable" //nolint:gosec
	flagDirectoryDesc    = "The directory where you're git repo lives. Defaults to your current working directory"
	flagRemoteNameDesc   = "Name of the GitLab remote in your git repo. Defaults to 'origin'"
	flagOutputFormatDesc = "One of 'export', 'json', 'table'. Default 'export', which can be executed to export variables"
	flagGroupDesc        = "Set group to true to get only variables from the parent groups."
)

func glenCmd() *cobra.Command {
	var (
		recurse      bool   // recurse determines if glen with also get variables from the project's parent groups
		apiKey       string // apiKey is the GitLab key that we should use when calling the API
		directory    string // directory is the path to the git repo that we should run glen on
		remoteName   string // remoteName is the name of the GitLab remote in your git repo
		outputFormat string // outputFormat is the text format that we should use to print our results to stdout
		groupOnly    bool   // groupOnly determines if glen only gets variables from the project's parent groups
	)

	cmd := &cobra.Command{
		Use:   "glen",
		Short: "glen prints variables for a GitLab project and it's parent groups",
		Long: `Glen is a simple command line tool that, when run within a GitLab project,
will call the GitLab API to get all environment variables from your project's
CI/CD pipeline and print them locally, ready for exporting.

With the default flags you can run 'eval $(glen -r)' to export your project's
variables and the variables of every parent group.`,
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			var err error

			repo := glen.NewRepo()
			repo.LocalPath = directory
			repo.RemoteName = remoteName

			err = repo.Init()
			if err != nil {
				slog.Error("failed to initialize the repository", "error", err)
				os.Exit(1)
			}

			vars := glen.NewVariables(repo)
			vars.GroupOnly = groupOnly
			vars.Recurse = recurse
			if apiKey != "GITLAB_TOKEN" {
				vars.SetAPIKey(apiKey)
			}

			if !vars.IsAPIKeySet() {
				fmt.Println("GitLab API key not set. Please use --api-key/-k flag or set GITLAB_TOKEN environment variable.")
				os.Exit(1)
			}

			err = vars.Init()
			if err != nil {
				slog.Error("failed to initialize variables", "error", err)
				os.Exit(1)
			}

			output(vars.Env, outputFormat)
		},
	}

	cmd.Flags().BoolVarP(&recurse, "recurse", "r", false, flagRecurseDesc)
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "GITLAB_TOKEN", flagAPIKeyDesc)
	cmd.Flags().StringVarP(&directory, "directory", "d", ".", flagDirectoryDesc)
	cmd.Flags().StringVarP(&remoteName, "remote-name", "n", "origin", flagRemoteNameDesc)
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "export", flagOutputFormatDesc)
	cmd.Flags().BoolVarP(&groupOnly, "group-only", "g", false, flagGroupDesc)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(v string) error {
	glen := glenCmd()
	glen.AddCommand(versionCmd(v))

	err := glen.Execute()
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	return nil
}
