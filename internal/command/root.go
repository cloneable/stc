package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func overrideRepoValidation(*cobra.Command, []string) error { return nil }

var rootCmd = &cobra.Command{
	Use:   "stacker",
	Short: "",
	Long:  "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("validate stacker refs")
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
