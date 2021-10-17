package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func overrideRepoValidation(*cobra.Command, []string) error { return nil }

var (
	// used by some subcommands
	forceFlag bool

	rootCmd = &cobra.Command{
		Use: "stacker <command>",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("validate stacker refs")
			return nil
		},
		DisableFlagsInUseLine: true,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
