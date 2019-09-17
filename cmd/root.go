package cmd

import (
	"log"
	"os"

	"github.com/lingrino/glen/glen"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
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

// init is where we set all flags
func init() {
	rootCmd.Flags().BoolVarP(&recurse, "recurse", "r", false, "Set recurse to true if you want to include the variables of the parent groups")
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "k", "GITLAB_TOKEN", "Your GitLab API key. NOTE - It's preferrable to specify your key as a GITLAB_TOKEN environment variable")
	rootCmd.Flags().StringVarP(&directory, "directory", "d", ".", "The directory where you're git repo lives. Defaults to your current working directory")
	rootCmd.Flags().StringVarP(&remoteName, "remote-name", "n", "origin", "The name of the GitLab remote in your git repo. Defaults to 'origin'. Check with 'git remote -v'")
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "export", "The output format. One of 'export', 'json', 'table'. Defaults to 'export', which can be executed to export all variables.")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(vers string) {
	version = vers

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
