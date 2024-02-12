package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionCmd(v string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Returns the current glen version",

		Args: cobra.NoArgs,

		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(v)
		},
	}
}
