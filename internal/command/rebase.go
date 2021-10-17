package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rebaseCmd = &cobra.Command{
		Use:  "rebase [<branch>...]",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("rebase")
			return stacker.Rebase(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(rebaseCmd)
}
