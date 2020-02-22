package cmd

import (
	"log"

	"github.com/lingrino/glen/glen"
	"github.com/spf13/cobra"
)

const (
	flagRecurseDesc      = "Set recurse to true if you want to include the variables of the parent groups"
	flagAPIKeyDesc       = "Your GitLab API key, if not set as a GITLAB_TOKEN environment variable"
	flagDirectoryDesc    = "The directory where you're git repo lives. Defaults to your current working directory"
	flagRemoteNameDesc   = "Name of the GitLab remote in your git repo. Defaults to 'origin'"
	flagOutputFormatDesc = "One of 'export', 'json', 'table'. Default 'export', which can be executed to export variables"
)

func glenCmd() *cobra.Command {
	// recurse determines if glen with also get variables from the project's parent groups
	var recurse bool

	// apiKey is the GitLab key that we should use when calling the API
	var apiKey string

	// directory is the path to the git repo that we should run glen on
	var directory string

	// remoteName is the name of the GitLab remote in your git repo
	var remoteName string

	// outputFormat is the text format that we should use to print our results to stdout
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "glen",
		Short: "glen prints variables for a GitLab project and it's parent groups",
		Long: `Glen is a simple command line tool that, when run within a GitLab project,
will call the GitLab API to get all environment variables from your project's
CI/CD pipeline and print them locally, ready for exporting.

With the default flags you can run 'eval $(glen -r)' to export your project's
variables and the variables of every parent group.`,
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			repo := glen.NewRepo()
			repo.LocalPath = directory
			repo.RemoteName = remoteName

			err = repo.Init()
			if err != nil {
				log.Fatal(err)
			}

			vars := glen.NewVariables(repo)
			vars.Recurse = recurse
			if apiKey != "GITLAB_TOKEN" {
				vars.SetAPIKey(apiKey)
			}

			err = vars.Init()
			if err != nil {
				log.Fatal(err)
			}

			print(vars.Env, outputFormat)
		},
	}

	cmd.Flags().BoolVarP(&recurse, "recurse", "r", false, flagRecurseDesc)
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "GITLAB_TOKEN", flagAPIKeyDesc)
	cmd.Flags().StringVarP(&directory, "directory", "d", ".", flagDirectoryDesc)
	cmd.Flags().StringVarP(&remoteName, "remote-name", "n", "origin", flagRemoteNameDesc)
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "export", flagOutputFormatDesc)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(v string) error {
	glen := glenCmd()
	glen.AddCommand(versionCmd(v))

	err := glen.Execute()
	if err != nil {
		return err
	}

	return err
}
